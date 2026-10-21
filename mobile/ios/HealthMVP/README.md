# HealthMVP (iOS HealthKit Minimal Demo)

目的：以最少依赖实现 HealthKit 授权与核心指标采集（体重、体脂、步数、心率、血压、睡眠），在 App 内本地展示列表，不进行上传，便于快速验证。

## 使用步骤
1) 在 Xcode 新建 iOS App（SwiftUI + Swift），Bundle ID 自定，iOS 16+。
2) 在项目 Signing & Capabilities 添加 HealthKit 能力；在 Info 添加 `NSHealthShareUsageDescription`（必需）与 `NSHealthUpdateUsageDescription`（可选）文案。
3) 将本目录下的 Swift 文件添加到你的 Xcode 工程：
   - HealthMVPApp.swift
   - ContentView.swift
   - HealthManager.swift
4) 首次运行：点击页面按钮触发授权与一次性拉取最近 7 天数据并在页面列表展示。

## 权限提示文案建议（Info.plist）
- NSHealthShareUsageDescription: 需要读取体征与活动数据，以便为你的主动式AI助手提供健康上下文与贴合建议。
- NSHealthUpdateUsageDescription: 仅在需要向健康数据写入标记或记录时使用（当前MVP不写入）。

## 指标范围（读取）
- 体重（HKQuantityTypeIdentifier.bodyMass）
- 体脂（HKQuantityTypeIdentifier.bodyFatPercentage）
- 步数（HKQuantityTypeIdentifier.stepCount）
- 心率（HKQuantityTypeIdentifier.heartRate）
- 血压（systolic/diastolic）
- 睡眠分析（HKCategoryTypeIdentifier.sleepAnalysis）

## 数据展示
列表项包含：
- type（weight/bp/activity/sleep）
- ts（时间戳秒）
- payload 摘要（如 weight_kg、body_fat_pct、sys/dia、hr、steps、duration_s 等）

## 兼容性
- 设备：iPhone，iOS 16+。
- 需真机测试；模拟器无 HealthKit 数据。

## 故障排查
- 未弹出授权：确认已添加 HealthKit 能力与 Info 权限。
- 无数据：确认 Apple 健康中存在对应数据；或将拉取窗口改为更长（修改 `collectRecentDays(days:)`）。
- 权限仅弹一次：到 设置 -> 健康 -> 应用 中手动开启相应读取权限。
