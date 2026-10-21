package adcode

import (
    "encoding/json"
    "errors"
    "os"
    "path/filepath"
    "strings"
    "sync"
)

// Entry 城市条目（从 JSON 加载）
type Entry struct {
	Name   string `json:"name"`
	Adcode string `json:"adcode"`
}

// Store 内存索引
// - 名称和 adcode 都建立索引，名称用原样与小写两份，便于无大小写/空格差异的查找
type Store struct {
	byName   map[string]Entry
	byAdcode map[string]Entry
	mu       sync.RWMutex
}

// Load 从 JSON 文件加载
func Load(path string) (*Store, error) {
	fp, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	var arr []Entry
	if err := json.NewDecoder(fp).Decode(&arr); err != nil {
		return nil, err
	}
	if len(arr) == 0 {
		return nil, errors.New("adcode empty dataset")
	}
	st := &Store{byName: make(map[string]Entry, len(arr)*2), byAdcode: make(map[string]Entry, len(arr))}
	for _, e := range arr {
		if e.Name == "" || e.Adcode == "" {
			continue
		}
		st.byAdcode[e.Adcode] = e
		st.byName[e.Name] = e
		st.byName[strings.ToLower(strings.TrimSpace(e.Name))] = e
	}
	return st, nil
}

// GetByAdcode 返回条目
func (s *Store) GetByAdcode(code string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.byAdcode[code]
	return e, ok
}

// GetByName 名称查找（尝试原样与小写）
func (s *Store) GetByName(name string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if e, ok := s.byName[name]; ok {
		return e, true
	}
	if e, ok := s.byName[strings.ToLower(strings.TrimSpace(name))]; ok {
		return e, true
	}
	return Entry{}, false
}

// GetByFuzzy 模糊查找：
// 1) 精确/小写匹配；
// 2) 尝试添加/移除常见行政区后缀（市/省/区/县/州/盟/市辖区/自治区/特别行政区）；
// 3) 在名称索引中做子串匹配（包含关系），命中即返回第一条。
func (s *Store) GetByFuzzy(name string) (Entry, bool) {
    raw := strings.TrimSpace(name)
    if raw == "" {
        return Entry{}, false
    }
    lower := strings.ToLower(raw)

    // 1) 先走现有精确/小写匹配
    if e, ok := s.GetByName(raw); ok {
        return e, true
    }
    if e, ok := s.GetByName(lower); ok {
        return e, true
    }

    // 常见后缀列表
    suffixes := []string{"市", "省", "区", "县", "州", "盟", "市辖区", "自治区", "特别行政区"}

    // 2) 追加后缀再试（如 北京 -> 北京市）
    for _, suf := range suffixes {
        if e, ok := s.GetByName(raw+suf); ok {
            return e, true
        }
        if e, ok := s.GetByName(lower+suf); ok {
            return e, true
        }
    }

    // 3) 去后缀再试（如 广州市 -> 广州）
    trimmed := raw
    for _, suf := range suffixes {
        if strings.HasSuffix(trimmed, suf) {
            trimmed = strings.TrimSuffix(trimmed, suf)
            break
        }
    }
    if trimmed != raw {
        if e, ok := s.GetByName(trimmed); ok {
            return e, true
        }
        if e, ok := s.GetByName(strings.ToLower(trimmed)); ok {
            return e, true
        }
    }

    // 4) 子串包含匹配（遍历 byName 键，命中第一条）
    s.mu.RLock()
    defer s.mu.RUnlock()
    for k, e := range s.byName {
        kk := strings.ToLower(strings.TrimSpace(k))
        if kk == "" {
            continue
        }
        if strings.Contains(kk, lower) || strings.Contains(lower, kk) {
            return e, true
        }
    }
    return Entry{}, false
}
