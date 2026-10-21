package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

// cityEntry 表示从表格抽取的一行
type cityEntry struct {
	Name   string `json:"name"`
	Adcode string `json:"adcode"`
}

func main() {
	in := flag.String("in", "asset/AMap_adcode_citycode.xlsx", "input Excel path")
	sheet := flag.String("sheet", "Sheet1", "sheet name to read")
	out := flag.String("out", "asset/cities_adcode.json", "output JSON file path")
	flag.Parse()

	if err := run(*in, *sheet, *out); err != nil {
		fmt.Println("adcode_gen error:", err)
		os.Exit(1)
	}
	fmt.Println("generated:", *out)
}

func run(inPath, sheet, outPath string) error {
	// 打开 Excel
	f, err := excelize.OpenFile(inPath)
	if err != nil {
		return err
	}
	defer f.Close()

	rows, err := f.GetRows(sheet)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return errors.New("sheet has no rows")
	}

	// 找表头列索引
	iName, iAd := -1, -1
	head := rows[0]
	for i, v := range head {
		vv := strings.TrimSpace(v)
		if vv == "中文名" {
			iName = i
		}
		if strings.EqualFold(vv, "adcode") {
			iAd = i
		}
	}
	if iName < 0 || iAd < 0 {
		return fmt.Errorf("header not found: 中文名 at %d, adcode at %d", iName, iAd)
	}

	var out []cityEntry
	for r := 1; r < len(rows); r++ { // 从第2行开始
		row := rows[r]
		// 防止列不足
		name := pick(row, iName)
		ad := pick(row, iAd)
		name = strings.TrimSpace(name)
		ad = strings.TrimSpace(ad)
		if name == "" || ad == "" {
			continue
		}
		// 过滤非数字或长度异常的 adcode（一般为6位）
		if len(ad) < 6 {
			continue
		}
		out = append(out, cityEntry{Name: name, Adcode: ad})
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	fp, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer fp.Close()
	enc := json.NewEncoder(fp)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func pick(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return row[idx]
}
