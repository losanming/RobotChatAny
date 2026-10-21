package resource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	neturl "net/url"

	"github.com/PuerkitoBio/goquery"
)

// Fetcher 资源获取器
type Fetcher struct {
	sources    []ResourceSource
	cache      map[ResourceType][]Resource
	cacheMutex sync.RWMutex
	lastUpdate int64
	updateLock sync.Mutex
	httpClient *http.Client
}

// NewFetcher 创建资源获取器
func NewFetcher(configPath string) (*Fetcher, error) {
	// 读取配置文件
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config struct {
		Sources []ResourceSource `json:"sources"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &Fetcher{
		sources:    config.Sources,
		cache:      make(map[ResourceType][]Resource),
		lastUpdate: 0, // 确保第一次调用时会更新缓存
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// GetRandomResources 获取指定类型的随机资源
func (f *Fetcher) GetRandomResources(types []ResourceType, count int) ([]Resource, error) {
	// 检查缓存是否需要更新
	if err := f.checkAndUpdateCache(); err != nil {
		return nil, err
	}

	f.cacheMutex.RLock()
	defer f.cacheMutex.RUnlock()

	// 创建一个map来记录已添加的资源
	addedResources := make(map[string]bool)
	var result []Resource

	// 遍历所有类型
	for _, t := range types {
		if resources, ok := f.cache[t]; ok {
			// 遍历该类型的所有资源
			for _, res := range resources {
				// 使用URL作为唯一标识
				if !addedResources[res.URL] {
					result = append(result, res)
					addedResources[res.URL] = true
				}
			}
		}
	}
	if count != 0 && len(result) > count {
		result = result[:count]
	}
	return result, nil
}

// checkAndUpdateCache 检查并更新缓存
// 默认缓存更新间隔（4小时）
const defaultCacheUpdateInterval = 4 * 3600

// 最大并发更新数
const maxConcurrentUpdates = 3

func (f *Fetcher) checkAndUpdateCache() error {
	now := time.Now().Unix()
	// 如果距离上次更新不足4小时，则不更新
	if now-f.lastUpdate < defaultCacheUpdateInterval {
		return nil
	}

	// 获取更新锁
	f.updateLock.Lock()
	defer f.updateLock.Unlock()

	// 再次检查，避免并发更新
	if now-f.lastUpdate < defaultCacheUpdateInterval {
		return nil
	}

	// 使用信号量控制并发数
	sem := make(chan struct{}, maxConcurrentUpdates)
	var wg sync.WaitGroup
	errorChan := make(chan error, len(f.sources))
	newCache := make(map[ResourceType][]Resource)

	for _, source := range f.sources {
		wg.Add(1)
		go func(s ResourceSource) {
			// 获取信号量，限制并发
			sem <- struct{}{}
			defer func() { <-sem }()
			defer wg.Done()
			resources, err := f.fetchFromSource(s)
			if err != nil {
				errorChan <- fmt.Errorf("从%s获取资源失败: %v", s.Name, err)
				return
			}
			f.cacheMutex.Lock()
			newCache[s.Type] = append(newCache[s.Type], resources...)
			f.cacheMutex.Unlock()
		}(source)
	}

	wg.Wait()
	close(errorChan)

	// 检查是否有错误
	var errors []error
	for err := range errorChan {
		// 只记录错误，不中断执行
		fmt.Printf("警告: %v\n", err)
		errors = append(errors, err)
	}

	// 更新缓存
	f.cacheMutex.Lock()
	f.cache = newCache
	f.lastUpdate = now
	f.cacheMutex.Unlock()

	return nil
}

// fetchFromSource 从指定源获取资源
func (f *Fetcher) fetchFromSource(source ResourceSource) ([]Resource, error) {
	// 创建HTTP请求
	req, err := http.NewRequest("GET", source.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request error: %v", err)
	}

	// 添加请求头
	for key, value := range source.Headers {
		req.Header.Set(key, value)
	}

	// 发送HTTP请求
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch url error: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 检查内容类型
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(strings.ToLower(contentType), "application/json") {
		// 处理JSON响应
		var data interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("decode json error: %v", err)
		}

		// 处理统一API响应
		if jsonData, ok := data.(map[string]interface{}); ok {
			// 检查是否在维护中
			if code, ok := jsonData["code"].(string); ok {
				if code == "-1" {
					// 如果在维护中，返回空列表
					return []Resource{}, nil
				} else if code == "200" {
					if items, ok := jsonData["data"].([]interface{}); ok {
						var resources []Resource
						for _, item := range items {
							if hot, ok := item.(map[string]interface{}); ok {
								title := fmt.Sprintf("%v", hot["title"])
								content := fmt.Sprintf("%v", hot["desc"])
								// 使用分类器对内容进行智能分类
								types := classifyByKeywords(title, content)
								// 确保原始类型也包含在内
								typeExists := false
								for _, t := range types {
									if t == source.Type {
										typeExists = true
										break
									}
								}
								if !typeExists {
									types = append(types, source.Type)
								}
								// 获取URL
								url := fmt.Sprintf("%v", hot["url"])

								// 如果是掘金或知乎，尝试获取详细内容
								if strings.Contains(url, "juejin.cn") || strings.Contains(url, "zhihu.com") {
									// 发起请求获取详细内容
									req, err := http.NewRequest("GET", url, nil)
									if err == nil {
										// 添加请求头
										for key, value := range source.Headers {
											req.Header.Set(key, value)
										}

										// 发送请求
										resp, err := f.httpClient.Do(req)
										if err == nil {
											defer resp.Body.Close()

											// 解析HTML
											doc, err := goquery.NewDocumentFromReader(resp.Body)
											if err == nil {
												// 根据不同网站选择内容选择器
												var articleContent string
												if strings.Contains(url, "juejin.cn") {
													// 掘金文章内容
													articleContent = doc.Find(".article-content").Text()
												} else if strings.Contains(url, "zhihu.com") {
													// 知乎回答内容
													articleContent = doc.Find(".RichContent-inner").Text()
												} else if strings.Contains(url, "ithome.com") {
													// IT之家文章内容
													articleContent = doc.Find(".news_info").Text()
												}

												// 如果获取到内容，使用详细内容
												if articleContent != "" {
													// 截取前500个字符
													if len(articleContent) > 500 {
														articleContent = articleContent[:500] + "..."
													}
													content = articleContent
												}
											}
										}
									}
								}

								// 如果没有获取到详细内容，使用标题
								if content == "" || content == "<nil>" {
									content = title
								}

								// 转换热度值
								var hotValue int64
								switch v := hot["hot"].(type) {
								case float64:
									hotValue = int64(v)
								case int:
									hotValue = int64(v)
								case int64:
									hotValue = v
								case string:
									if val, err := strconv.ParseInt(v, 10, 64); err == nil {
										hotValue = val
									}
								}

								// 添加热度信息
								if hotValue > 0 {
									content = fmt.Sprintf("%s\n\n热度：%d", content, hotValue)
								}

								// 添加来源信息
								content = fmt.Sprintf("%s\n\n来源：%s", content, source.Name)

								resources = append(resources, Resource{
									Types:     types,
									Title:     title,
									Content:   content,
									URL:       fmt.Sprintf("%v", hot["url"]),
									HotValue:  hotValue,
									Timestamp: time.Now().Unix(),
								})
							}
						}
						return resources, nil
					}
				}
			}
		}

		return nil, fmt.Errorf("unsupported json format")
	}

	// 使用goquery解析HTML页面
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse html error: %v", err)
	}

	var resources []Resource
	doc.Find(source.Selector).Each(func(i int, s *goquery.Selection) {
		// 提取标题
		var title string
		if source.TitleSelector != "" {
			title = strings.TrimSpace(s.Find(source.TitleSelector).First().Text())
		} else {
			title = strings.TrimSpace(s.Text())
		}

		// 提取链接
		url, exists := s.Attr("href")
		if !exists && source.URLSelector != "" {
			url, _ = s.Find(source.URLSelector).First().Attr("href")
		}

		// 规范化URL
		if url != "" {
			if !isValidURL(url) {
				// 处理相对URL
				baseURL, err := neturl.Parse(source.URL)
				if err != nil {
					return
				}

				if strings.HasPrefix(url, "//") {
					url = baseURL.Scheme + ":" + url
				} else if strings.HasPrefix(url, "/") {
					url = baseURL.Scheme + "://" + baseURL.Host + url
				} else {
					// 处理相对路径
					relative, err := neturl.Parse(url)
					if err != nil {
						return
					}
					url = baseURL.ResolveReference(relative).String()
				}
			}
		}

		// 提取内容
		var content string
		if source.ContentSelector != "" {
			content = strings.TrimSpace(s.Find(source.ContentSelector).First().Text())
		}
		if content == "" {
			content = title
		}

		// 过滤无效数据
		if title == "" || url == "" {
			return
		}

		// 使用分类器对内容进行智能分类
		types := classifyByKeywords(title, content)
		// 确保原始类型也包含在内
		typeExists := false
		for _, t := range types {
			if t == source.Type {
				typeExists = true
				break
			}
		}
		if !typeExists {
			types = append(types, source.Type)
		}

		// 处理内容
		if content == "" {
			content = title
		}

		// 提取热度值
		var hotValue int64
		if source.HotSelector != "" {
			hotText := strings.TrimSpace(s.Find(source.HotSelector).First().Text())
			// 尝试不同的热度格式
			// 先尝试直接解析数字
			if val, err := strconv.ParseInt(hotText, 10, 64); err == nil {
				hotValue = val
			} else {
				// 尝试提取数字部分
				re := regexp.MustCompile(`\d+`)
				if matches := re.FindStringSubmatch(hotText); len(matches) > 0 {
					if val, err := strconv.ParseInt(matches[0], 10, 64); err == nil {
						hotValue = val
					}
				}
			}
		}

		// 添加热度信息
		if hotValue > 0 {
			content = fmt.Sprintf("%s\n\n热度：%d", content, hotValue)
		}

		// 添加来源信息
		content = fmt.Sprintf("%s\n\n来源：%s", content, source.Name)

		// 添加资源
		resources = append(resources, Resource{
			Types:     types,
			Title:     title,
			URL:       url,
			Content:   content,
			HotValue:  hotValue,
			Timestamp: time.Now().Unix(),
		})
	})

	// 检查是否获取到资源
	if len(resources) == 0 {
		return nil, fmt.Errorf("no resources found for selector: %s", source.Selector)
	}

	return resources, nil
}

// isValidURL 检查URL是否合法
func isValidURL(url string) bool {
	_, err := neturl.ParseRequestURI(url)
	return err == nil
}

// getInt64Value 安全地获取int64值
func getInt64Value(v interface{}) int64 {
	if v == nil {
		return 0
	}

	switch value := v.(type) {
	case float64:
		return int64(value)
	case float32:
		return int64(value)
	case int64:
		return value
	case int32:
		return int64(value)
	case int:
		return int64(value)
	case string:
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
	}

	return 0
}
