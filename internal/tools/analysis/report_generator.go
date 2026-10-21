package analysis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ReportGenerator 报告生成工具
type ReportGenerator struct {
	enableLearning bool
	category       string
	toolName       string
}

// NewReportGenerator 创建报告生成工具
func NewReportGenerator() *ReportGenerator {
	return &ReportGenerator{
		enableLearning: true,
		category:       "analysis",
		toolName:       "report_generator",
	}
}

// Info 返回工具信息
func (r *ReportGenerator) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "report_generator",
		Desc: "报告生成工具，支持模板化报告生成、多格式导出、数据可视化",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {
				Desc:     "操作类型",
				Type:     schema.String,
				Required: true,
				Enum:     []string{"generate", "list_templates", "create_template"},
			},
			"template": {
				Desc:     "模板类型",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"weekly_report", "monthly_report", "project_summary", "data_analysis", "custom"},
			},
			"data": {
				Desc:     "报告数据",
				Type:     schema.Object,
				Required: false,
			},
			"output_format": {
				Desc:     "输出格式",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"markdown", "html", "json", "text"},
			},
			"title": {
				Desc:     "报告标题",
				Type:     schema.String,
				Required: false,
			},
			"author": {
				Desc:     "作者",
				Type:     schema.String,
				Required: false,
			},
			"include_charts": {
				Desc:     "是否包含图表",
				Type:     schema.Boolean,
				Required: false,
			},
		}),
	}, nil
}

// InvokableRun 执行工具调用
func (r *ReportGenerator) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	action := params["action"].(string)
	switch action {
	case "generate":
		return r.generateReport(params)
	case "list_templates":
		return r.listTemplates()
	case "create_template":
		return r.createTemplate(params)
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", action)
	}
}

// generateReport 生成报告
func (r *ReportGenerator) generateReport(params map[string]interface{}) (string, error) {
	template := getStringParam(params, "template", "weekly_report")
	outputFormat := getStringParam(params, "output_format", "markdown")
	title := getStringParam(params, "title", "报告")
	author := getStringParam(params, "author", "AI助手")
	includeCharts := getBoolParam(params, "include_charts", false)

	data, ok := params["data"].(map[string]interface{})
	if !ok {
		data = make(map[string]interface{})
	}

	// 根据模板生成报告
	var content string
	switch template {
	case "weekly_report":
		content = r.generateWeeklyReport(data, title, author, includeCharts)
	case "monthly_report":
		content = r.generateMonthlyReport(data, title, author, includeCharts)
	case "project_summary":
		content = r.generateProjectSummary(data, title, author, includeCharts)
	case "data_analysis":
		content = r.generateDataAnalysisReport(data, title, author, includeCharts)
	case "custom":
		content = r.generateCustomReport(data, title, author, includeCharts)
	default:
		return "", fmt.Errorf("不支持的模板类型: %s", template)
	}

	// 格式转换
	var finalContent string
	switch outputFormat {
	case "markdown":
		finalContent = content
	case "html":
		finalContent = r.markdownToHTML(content)
	case "json":
		finalContent = r.reportToJSON(content, title, author, template)
	case "text":
		finalContent = r.markdownToText(content)
	default:
		finalContent = content
	}

	result := map[string]interface{}{
		"success":       true,
		"template":      template,
		"output_format": outputFormat,
		"title":         title,
		"author":        author,
		"content":       finalContent,
		"generated_at":  time.Now().Format(time.RFC3339),
		"content_length": len(finalContent),
	}

	return toJSON(result), nil
}

// listTemplates 列出可用模板
func (r *ReportGenerator) listTemplates() (string, error) {
	templates := []map[string]interface{}{
		{
			"name":        "weekly_report",
			"description": "周报模板，包含本周工作总结、下周计划、问题与风险",
			"fields":      []string{"week_start", "week_end", "tasks_completed", "tasks_in_progress", "issues", "next_week_plan"},
		},
		{
			"name":        "monthly_report",
			"description": "月报模板，包含月度工作总结、关键成果、数据分析",
			"fields":      []string{"month", "key_achievements", "metrics", "challenges", "next_month_plan"},
		},
		{
			"name":        "project_summary",
			"description": "项目总结模板，包含项目概况、进度、成果、问题",
			"fields":      []string{"project_name", "status", "progress", "milestones", "risks", "next_steps"},
		},
		{
			"name":        "data_analysis",
			"description": "数据分析报告模板，包含数据概览、统计结果、趋势分析",
			"fields":      []string{"data_source", "time_range", "statistics", "trends", "insights"},
		},
		{
			"name":        "custom",
			"description": "自定义模板，根据提供的数据生成报告",
			"fields":      []string{"custom_data"},
		},
	}

	result := map[string]interface{}{
		"success":   true,
		"templates": templates,
		"total":     len(templates),
	}

	return toJSON(result), nil
}

// createTemplate 创建自定义模板
func (r *ReportGenerator) createTemplate(params map[string]interface{}) (string, error) {
	templateName := getStringParam(params, "template", "custom")
	description := getStringParam(params, "description", "自定义模板")

	// 简化实现：返回模板创建确认
	result := map[string]interface{}{
		"success":     true,
		"template":    templateName,
		"description": description,
		"message":     "模板创建成功（演示功能）",
		"created_at":  time.Now().Format(time.RFC3339),
	}

	return toJSON(result), nil
}

// 报告生成函数
func (r *ReportGenerator) generateWeeklyReport(data map[string]interface{}, title, author string, includeCharts bool) string {
	weekStart := getStringFromData(data, "week_start", time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	weekEnd := getStringFromData(data, "week_end", time.Now().Format("2006-01-02"))
	tasksCompleted := getStringFromData(data, "tasks_completed", "5个主要任务")
	tasksInProgress := getStringFromData(data, "tasks_in_progress", "3个进行中任务")
	issues := getStringFromData(data, "issues", "无重大问题")
	nextWeekPlan := getStringFromData(data, "next_week_plan", "继续推进重点项目")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))
	sb.WriteString(fmt.Sprintf("**作者**: %s  \n", author))
	sb.WriteString(fmt.Sprintf("**报告周期**: %s 至 %s  \n", weekStart, weekEnd))
	sb.WriteString(fmt.Sprintf("**生成时间**: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	sb.WriteString("## 📋 本周工作总结\n\n")
	sb.WriteString(fmt.Sprintf("- **完成任务**: %s\n", tasksCompleted))
	sb.WriteString(fmt.Sprintf("- **进行中任务**: %s\n", tasksInProgress))
	sb.WriteString(fmt.Sprintf("- **问题与风险**: %s\n\n", issues))

	sb.WriteString("## 📊 关键数据\n\n")
	if includeCharts {
		sb.WriteString("```\n")
		sb.WriteString("任务完成情况图表（演示）\n")
		sb.WriteString("```\n\n")
	}
	sb.WriteString("- 任务完成率: 85%\n")
	sb.WriteString("- 工作时长: 40小时\n")
	sb.WriteString("- 效率评分: 4.5/5.0\n\n")

	sb.WriteString("## 🎯 下周计划\n\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", nextWeekPlan))

	sb.WriteString("## 💡 改进建议\n\n")
	sb.WriteString("- 优先处理高价值任务\n")
	sb.WriteString("- 加强团队协作\n")
	sb.WriteString("- 持续学习新技术\n\n")

	return sb.String()
}

func (r *ReportGenerator) generateMonthlyReport(data map[string]interface{}, title, author string, includeCharts bool) string {
	month := getStringFromData(data, "month", time.Now().Format("2006年01月"))
	keyAchievements := getStringFromData(data, "key_achievements", "完成了3个主要项目")
	metrics := getStringFromData(data, "metrics", "性能提升20%")
	challenges := getStringFromData(data, "challenges", "资源分配挑战")
	nextMonthPlan := getStringFromData(data, "next_month_plan", "优化工作流程")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))
	sb.WriteString(fmt.Sprintf("**作者**: %s  \n", author))
	sb.WriteString(fmt.Sprintf("**报告月份**: %s  \n", month))
	sb.WriteString(fmt.Sprintf("**生成时间**: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	sb.WriteString("## 🎯 月度总结\n\n")
	sb.WriteString(fmt.Sprintf("**关键成果**: %s\n\n", keyAchievements))

	sb.WriteString("## 📊 数据指标\n\n")
	sb.WriteString(fmt.Sprintf("- **核心指标**: %s\n", metrics))
	sb.WriteString("- 项目完成数: 8个\n")
	sb.WriteString("- 团队协作评分: 4.7/5.0\n\n")

	if includeCharts {
		sb.WriteString("## 📈 月度趋势\n\n")
		sb.WriteString("```\n")
		sb.WriteString("月度数据趋势图（演示）\n")
		sb.WriteString("```\n\n")
	}

	sb.WriteString("## ⚠️ 挑战与问题\n\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", challenges))

	sb.WriteString("## 🚀 下月计划\n\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", nextMonthPlan))

	return sb.String()
}

func (r *ReportGenerator) generateProjectSummary(data map[string]interface{}, title, author string, includeCharts bool) string {
	projectName := getStringFromData(data, "project_name", "智能助手项目")
	status := getStringFromData(data, "status", "进行中")
	progress := getStringFromData(data, "progress", "75%")
	milestones := getStringFromData(data, "milestones", "完成核心功能开发")
	risks := getStringFromData(data, "risks", "技术风险可控")
	nextSteps := getStringFromData(data, "next_steps", "进入测试阶段")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))
	sb.WriteString(fmt.Sprintf("**作者**: %s  \n", author))
	sb.WriteString(fmt.Sprintf("**项目名称**: %s  \n", projectName))
	sb.WriteString(fmt.Sprintf("**当前状态**: %s  \n", status))
	sb.WriteString(fmt.Sprintf("**项目进度**: %s  \n", progress))
	sb.WriteString(fmt.Sprintf("**更新时间**: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	sb.WriteString("## 📋 项目概况\n\n")
	sb.WriteString(fmt.Sprintf("项目目前处于**%s**状态，整体进度为**%s**。\n\n", status, progress))

	sb.WriteString("## 🎯 关键里程碑\n\n")
	sb.WriteString(fmt.Sprintf("- %s\n\n", milestones))

	sb.WriteString("## 📊 项目数据\n\n")
	if includeCharts {
		sb.WriteString("```\n")
		sb.WriteString("项目进度甘特图（演示）\n")
		sb.WriteString("```\n\n")
	}
	sb.WriteString("- 任务完成: 75%\n")
	sb.WriteString("- 代码质量: 优秀\n")
	sb.WriteString("- 测试覆盖: 80%\n\n")

	sb.WriteString("## ⚠️ 风险与问题\n\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", risks))

	sb.WriteString("## 🚀 下一步行动\n\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", nextSteps))

	return sb.String()
}

func (r *ReportGenerator) generateDataAnalysisReport(data map[string]interface{}, title, author string, includeCharts bool) string {
	dataSource := getStringFromData(data, "data_source", "用户行为数据")
	timeRange := getStringFromData(data, "time_range", "最近30天")
	statistics := getStringFromData(data, "statistics", "数据量: 10万条")
	trends := getStringFromData(data, "trends", "呈上升趋势")
	insights := getStringFromData(data, "insights", "用户活跃度提升")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))
	sb.WriteString(fmt.Sprintf("**作者**: %s  \n", author))
	sb.WriteString(fmt.Sprintf("**数据来源**: %s  \n", dataSource))
	sb.WriteString(fmt.Sprintf("**时间范围**: %s  \n", timeRange))
	sb.WriteString(fmt.Sprintf("**分析时间**: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	sb.WriteString("## 📊 数据概览\n\n")
	sb.WriteString(fmt.Sprintf("- **统计信息**: %s\n", statistics))
	sb.WriteString("- 数据完整性: 98%\n")
	sb.WriteString("- 数据质量: 优秀\n\n")

	sb.WriteString("## 📈 统计结果\n\n")
	if includeCharts {
		sb.WriteString("```\n")
		sb.WriteString("数据分布直方图（演示）\n")
		sb.WriteString("```\n\n")
	}
	sb.WriteString("- 均值: 45.6\n")
	sb.WriteString("- 中位数: 42.3\n")
	sb.WriteString("- 标准差: 12.8\n\n")

	sb.WriteString("## 🔄 趋势分析\n\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", trends))

	sb.WriteString("## 💡 关键洞察\n\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", insights))

	sb.WriteString("## 📝 结论与建议\n\n")
	sb.WriteString("- 数据质量良好，可用于进一步分析\n")
	sb.WriteString("- 建议继续监控关键指标\n")
	sb.WriteString("- 可考虑增加数据收集维度\n\n")

	return sb.String()
}

func (r *ReportGenerator) generateCustomReport(data map[string]interface{}, title, author string, includeCharts bool) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))
	sb.WriteString(fmt.Sprintf("**作者**: %s  \n", author))
	sb.WriteString(fmt.Sprintf("**生成时间**: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	sb.WriteString("## 📋 自定义报告\n\n")
	sb.WriteString("以下是根据提供的数据生成的报告内容：\n\n")

	// 遍历数据生成内容
	for key, value := range data {
		sb.WriteString(fmt.Sprintf("### %s\n\n", key))
		sb.WriteString(fmt.Sprintf("%v\n\n", value))
	}

	return sb.String()
}

// 格式转换函数
func (r *ReportGenerator) markdownToHTML(markdown string) string {
	// 简化的Markdown到HTML转换
	html := markdown
	
	// 标题转换
	html = strings.ReplaceAll(html, "# ", "<h1>")
	html = strings.ReplaceAll(html, "## ", "<h2>")
	html = strings.ReplaceAll(html, "### ", "<h3>")
	
	// 换行转换
	html = strings.ReplaceAll(html, "\n\n", "</p><p>")
	html = strings.ReplaceAll(html, "\n", "<br>")
	
	// 粗体转换
	html = strings.ReplaceAll(html, "**", "<strong>")
	
	// 代码块转换
	html = strings.ReplaceAll(html, "```", "<pre>")
	
	return fmt.Sprintf("<html><body><p>%s</p></body></html>", html)
}

func (r *ReportGenerator) markdownToText(markdown string) string {
	// 移除Markdown格式
	text := markdown
	
	// 移除标题符号
	text = strings.ReplaceAll(text, "# ", "")
	text = strings.ReplaceAll(text, "## ", "")
	text = strings.ReplaceAll(text, "### ", "")
	
	// 移除粗体符号
	text = strings.ReplaceAll(text, "**", "")
	
	// 移除代码块符号
	text = strings.ReplaceAll(text, "```", "")
	
	return text
}

func (r *ReportGenerator) reportToJSON(content, title, author, template string) string {
	report := map[string]interface{}{
		"title":       title,
		"author":      author,
		"template":    template,
		"content":     content,
		"generated_at": time.Now().Format(time.RFC3339),
		"metadata": map[string]interface{}{
			"content_length": len(content),
			"word_count":     len(strings.Fields(content)),
		},
	}
	
	bytes, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Sprintf("{\"error\": \"JSON序列化失败: %v\"}", err)
	}
	return string(bytes)
}

// AfterRun 工具调用后的后处理
func (r *ReportGenerator) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，暂时不做学习集成
}

// 辅助函数
func getStringFromData(data map[string]interface{}, key, defaultValue string) string {
	if val, ok := data[key]; ok {
		return val.(string)
	}
	return defaultValue
}