package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mozillazg/go-pinyin"
	"log/slog"
	"new-robot-self/cmd/testFile/utils"
	"time"
	"xorm.io/xorm"
)

const (
	FoodListUrl = "https://one.boohee.com/store/pages/food_list_json"   //大类列表
	FoodInfoUrl = "https://food.boohee.com/fb/v1/food_rankings/%d.json" //食物信息-拼接id
	DetailUrl   = "https://www.miaofoods.com/detail/%s.html"
	AllInfoUrl  = "https://www.miaofoods.com/miaozsPc/foodHome/getFoodHomeInfo"
)

type AllData struct {
	Status struct {
		StatusCode   int    `json:"statusCode"`
		StatusReason string `json:"statusReason"`
	} `json:"status"`
	Result ResultData `json:"result"`
}

type ResultData struct {
	List           []ListData `json:"list"`
	Total          int        `json:"total"`
	FoodCateList   []CateList `json:"foodCateList"`   //大类列表
	HotFoodList    []ListData `json:"hotFoodList"`    //热门食物列表
	FriendLinkList []CateList `json:"friendLinkList"` //友情链接
}

type ListData struct {
	FoodId               int     `json:"foodId" xorm:"pk autoincr"`              // 食物ID
	CateName             string  `json:"cateName" xorm:"comment('分类名称')"`        // 分类名称
	FoodToken            string  `json:"foodToken" xorm:"comment('食物标识') index"` // 食物标识
	FoodName             string  `json:"foodName"`                               // 食物名称
	Foodcalory           string  `json:"foodcalory"`                             // 热量描述
	FoodIcon             string  `json:"foodIcon"`                               // 食物图标
	FoodNum              int     `json:"foodNum"`                                // 食物数量
	FoodWeightVal        int     `json:"foodWeightVal"`                          // 食物重量值
	FoodWeightUnit       string  `json:"foodWeightUnit"`                         // 食物重量单位
	FoodCaloriesVal      float64 `json:"foodCaloriesVal"`                        // 热量值
	FoodCaloriesUnit     string  `json:"foodCaloriesUnit"`                       // 热量单位
	FoodProteinVal       float64 `json:"foodProteinVal"`                         // 蛋白质值
	FoodProteinUnit      string  `json:"foodProteinUnit"`                        // 蛋白质单位
	FoodFatVal           float64 `json:"foodFatVal"`                             // 脂肪值
	FoodFatUnit          string  `json:"foodFatUnit"`                            // 脂肪单位
	FoodCarbohydrateVal  float64 `json:"foodCarbohydrateVal"`                    // 碳水化合物值
	FoodCarbohydrateUnit string  `json:"foodCarbohydrateUnit"`                   // 碳水化合物单位
}

type CateList struct {
	CateName  string `json:"cateName"`
	CateToken string `json:"cateToken"`
	ImageUrl  string `json:"imageUrl"`
}

func main() {
	cate_data := InitData()
	db, err := initDB()
	defer db.Close()
	if err != nil {
		slog.Error("init db err is ", err)
		return
	}
	// 根据分类获取食物 以及详情
	var task []utils.Task
	//var all []*AllData
	pool := utils.NewGoPool()
	for _, v := range cate_data {
		task = append(task, func() error {
			// 获取第一个分类食物数据`
			var cate = v.CateToken
			var map_headers = map[string]string{
				"Accept":          "application/json, text/plain, */*",
				"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36 Edg/136.0.0.0",
				"Content-Type":    "application/json;charset=UTF-8",
				"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
				"Accept-Encoding": "gzip, deflate, br",
				"Host":            "www.miaofoods.com",
				"Connection":      "keep-alive",
			}
			list_client := resty.New()
			list_client.SetHeaders(map_headers)
			list_client.SetProxy("http://127.0.0.1:7890")
			list_client.SetTimeout(5 * time.Second)
			//先走一次请求
			params := fmt.Sprintf("{\"cateName\":\"%s\",\"range\":{\"page\":%d,\"pageSize\":100}}", cate, 1)
			list_client.SetQueryParam("param", params)
			get, err := list_client.R().SetResult(&AllData{}).Get("https://www.miaofoods.com/miaozsPc/foodHome/getFoodHomeInfo")
			if err != nil {
				slog.Error("err is ", err)
				return err
			}
			var resp AllData
			err = json.Unmarshal(get.Body(), &resp)
			if err != nil {
				slog.Error("解析JSON失败:", err)
				return err
			}
			count := resp.Result.Total
			// 已经获取了第1页的数据，从第2页开始循环
			for i := 2; i <= (count+99)/100; i++ { // 使用(count+99)/100来向上取整
				// 设置查询参数，请求第i页数据
				list_client.SetQueryParam("param", fmt.Sprintf("{\"cateName\":\"%s\",\"range\":{\"page\":%d,\"pageSize\":100}}", cate, i))

				// 发送请求获取数据
				pageResp, err := list_client.R().SetResult(&AllData{}).Get("https://www.miaofoods.com/miaozsPc/foodHome/getFoodHomeInfo")
				if err != nil {
					slog.Error("获取第", i, "页数据失败:", err)
					continue
				}

				// 解析响应数据
				var pageData AllData
				err = json.Unmarshal(pageResp.Body(), &pageData)
				if err != nil {
					slog.Error("解析第", i, "页JSON失败:", err)
					continue
				}

				// 将新获取的数据添加到结果中
				resp.Result.List = append(resp.Result.List, pageData.Result.List...)
				if pageData.Result.Total == 0 {
					fmt.Printf("没有数据了，catename: %s catetoken: %s, index: %d, pageSize: %d \n", v.CateName, v.CateToken, i, 100)
				}
			}
			slog.Info("获取数据成功", "count", count, "resp len ", len(resp.Result.List), "task ", v.CateName)
			// 分批插入数据，每批最多100条
			if len(resp.Result.List) > 0 {
				// 解除注释，启用数据库插入
				batchSize := 100
				totalItems := len(resp.Result.List)

				for i := 0; i < totalItems; i += batchSize {
					end := i + batchSize
					if end > totalItems {
						end = totalItems
					}

					batch := resp.Result.List[i:end]
					slog.Info("正在插入数据批次", "分类", v.CateName, "批次", i/batchSize+1, "数据量", len(batch))

					// 插入数据
					_, err1 := db.Insert(&batch)
					if err1 != nil {
						slog.Error("插入数据批次失败:", "分类", v.CateName, "批次", i/batchSize+1, "错误", err1)
						// 继续处理下一批，不中断整个任务
						continue
					}
				}
			}
			return nil
		})
	}
	submitAllErr := pool.SubmitAll(task)
	if len(submitAllErr) > 0 {
		for _, err := range submitAllErr {
			slog.Error("任务执行失败:", err)
		}
	}

}

func InitData() []CateList {
	var map_headers = map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36 Edg/136.0.0.0",
		"Content-Type":    "application/json;charset=UTF-8",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"Accept-Encoding": "gzip, deflate, br",
		"Host":            "www.miaofoods.com",
		"Connection":      "keep-alive",
	}
	list_client := resty.New()
	list_client.SetHeaders(map_headers)
	list_client.SetProxy("http://127.0.0.1:7890")
	list_client.SetTimeout(5 * time.Second)
	list_client.SetQueryParam("param", "{\"cateName\":\"UCRVG0PAR4\",\"range\":{\"page\":1,\"pageSize\":100}}")
	get, err := list_client.R().SetResult(&AllData{}).Get("https://www.miaofoods.com/miaozsPc/foodHome/getFoodHomeInfo")
	if err != nil {
		slog.Error("err is ", err)
		return []CateList{}
	}
	var resp AllData
	err = json.Unmarshal(get.Body(), &resp)
	if err != nil {
		slog.Error("解析JSON失败:", err)
		return []CateList{}
	}
	return resp.Result.FoodCateList
}

// 将中文转换为拼音URL
func generatePinyinURL(chineseName string) string {
	// 设置拼音转换选项
	args := pinyin.NewArgs()
	args.Style = pinyin.Normal // 普通风格，不带声调
	args.Separator = ""        // 不使用分隔符

	// 转换为拼音
	pinyinSlice := pinyin.Pinyin(chineseName, args)

	// 将二维切片转换为一维字符串
	var pinyinStr string
	for _, p := range pinyinSlice {
		if len(p) > 0 {
			pinyinStr += p[0]
		}
	}

	// 构建详情页URL
	return "https://www.miaofoods.com/detail/" + pinyinStr + ".html"
}

func initDB() (*xorm.Engine, error) {
	// 替换成您的数据库连接信息
	dsn := "./food.db"
	engine, err := xorm.NewEngine("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}

	// 设置连接池
	engine.SetMaxIdleConns(5)
	engine.SetMaxOpenConns(10)
	err = engine.CreateTables(&ListData{})
	if err != nil {
		return nil, fmt.Errorf("创建表失败: %v", err)
	}
	return engine, nil
}
