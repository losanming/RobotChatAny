# 位置模块 Step 2 设计说明（AMap 逆地理 / 天气 / 实时路况）

> 目标：在 Step1 的基础上，基于 GCJ-02 坐标接入高德服务，最小代价地获取地点（逆地理）、天气与路况上下文，统一投递到上下文黑板，为 AI 大脑做情境分析与推理提供可信环境信号。

> 时间规范：所有时间字段一律使用毫秒级时间戳（int64, ms）。

## 1. 范围（Scope）
- 输入：`TransformedLocation`（来自 Step1，含 `lat_gcj02/lon_gcj02`）。
- 输出：`ReverseGeocodeResult`（adcode/addr/poi_topk）、`WeatherNow/Forecast`、`TrafficSnapshot`，统一封装为事件，写入黑板/总线。
- 限制：面向中国境内（GCJ-02）。境外则仅提供城市级兜底（可选）。

## 2. 目录结构（与 pkg/screenshot 风格对齐）
```
/pkg/location/provider/amap/
  ├── client.go          # 通用 HTTP 客户端（超时/重试/鉴权/签名）
  ├── reverse.go         # 逆地理封装（GCJ-02 → 地址/POI/adcode）
  ├── weather.go         # 实时天气与预报（Now/Forecast）
  ├── traffic.go         # 实时路况（圆形/矩形/道路）
  ├── cache.go           # 本地缓存（TTL 可配置）
  ├── rate_limit.go      # 速率限制（RPS）
  ├── types.go           # 数据结构与错误定义
  └── metrics.go         # 指标埋点
/pkg/location/service.go  # 串联 Step2：装配 AMap Provider，形成组合输出
```

## 3. AMap 接口（待你调研结果后补齐具体 URL/签名与字段）
- 逆地理（Reverse Geocoding）
  - 输入：`lat_gcj02/lon_gcj02`
  - 输出：`adcode`、`formatted_address`、`poi_topk[{id,name,type,distance}]`
  - 缓存：geohash-7/8 维度 3–10 分钟
- 天气（Weather）
  - 实时天气 Now：查询 `adcode` 或经纬度
  - 预报 Forecast：查询 `adcode`
  - 缓存：Now 5–10 分钟；Forecast 1–3 小时
- 实时路况（Traffic）
  - 查询范围：圆形（lat/lon/radius）或矩形（bounds）
  - 输出：拥堵指数、平均速度、拥堵路段 Top-K（可截断）
  - 缓存：1–3 分钟（速度>0 时更频繁刷新）

> 备注：具体 API 名称、版本、URL、鉴权方式（key/签名）与字段命名，待你完成高德文档确认后在本文件相应小节中补齐。

## 4. 对外数据结构（统一 Schema，时间戳 ms）
- ReverseGeocodeResult
  - `ts, lat_gcj02, lon_gcj02, adcode, addr, poi_topk[], provider="amap"`
- WeatherNow / WeatherForecast
  - `ts, adcode, temp_c, feels_like_c, humidity, wind, precip_prob, condition, provider="amap"`
  - Forecast 追加：`daily[] {date_ts, high_c, low_c, condition}`
- TrafficSnapshot
  - `ts, center_lat, center_lon, radius_m, jam_index, avg_speed_kph, segments_topk[], provider="amap"`

## 5. 服务编排（Step2 总流程）
1) 从 Step1 的 `TransformedLocation` 进入（含 ts、lat/lon GCJ-02）。
2) 命中缓存则直接返回（逆地理/天气/路况各自独立缓存）。
3) 触发远程调用（尊重速率限制与超时/重试）。
4) 产出 `ReverseGeocodeResult / WeatherNow/Forecast / TrafficSnapshot`。
5) 写入上下文黑板与/或投递总线（NSQ 主题建议见下）。

## 6. 触发与节流策略（默认值，可配置）
- 逆地理：
  - 时间阈值：≥60s；距离阈值：≥150–300m；停留触发；围栏进出
- 天气：
  - Now：5–10 分钟；Forecast：1–3 小时；天气突变（降雨/温差/风力阈值）触发加急刷新
- 路况：
  - 速度>0 时每 1–3 分钟；静止时拉长至 ≥5 分钟；拥堵指数突变触发加急

## 7. 缓存/限流/重试
- 缓存层级：
  - 本地内存（快速 TTL）；可选 Redis（跨进程共享）
- 限流：
  - RPS 按模块独立配置（逆地理/天气/路况）
- 重试：
  - 指数退避：200ms → 500ms → 1s；幂等请求仅在网络/5xx 场景重试

## 8. NSQ 事件（主题建议）
- `events.location.reverse_geocode`
- `events.location.weather_now`
- `events.location.weather_forecast`
- `events.location.traffic`

消息体均包含：`ts(int64, ms), device_id, provider, payload{...}`。

## 9. 指标与日志（Metrics & Logging）
- 计数：
  - amap_reverse_success_total / failed_total
  - amap_weather_now_success_total / failed_total
  - amap_weather_forecast_success_total / failed_total
  - amap_traffic_success_total / failed_total
- 时延：
  - amap_reverse_latency_ms、amap_weather_now_latency_ms、amap_weather_forecast_latency_ms、amap_traffic_latency_ms
- 缓存：
  - amap_cache_hit_total{type="reverse|weather_now|weather_forecast|traffic"}
- 限流/重试：
  - amap_rate_limited_total、amap_retry_total
- 日志：建议 JSON 行，包含 trace_id、url、status_code、cost_ms、err

## 10. 错误与降级
- 无网/超时/5xx：返回最近可信缓存；必要时降级到城市级天气；路况可返回“未知/上次值+时间差”。
- 配额/限流：尊重限流，延迟触发或跳过非关键查询。

## 11. 隐私与合规
- 均以最小化返回并落地黑板，仅供 AI 大脑决策；不存用户敏感信息（门牌级地址需授权）。
- 传输 TLS；密钥管理与最小权限；访问审计。

## 12. 验收标准（Step2）
- 准确性：逆地理地址与 POI 的匹配率满足预期（以抽样评估）。
- 时延：单接口 P95 延迟达预设（逆地理 < 500ms、天气 < 400ms、路况 < 700ms，视配额情况可调整）。
- 可用性：缓存命中率 ≥ 60%（城市/通勤场景），失败可降级且对话链路可继续。
- 一致性：时间一律毫秒时间戳（int64, ms）。

## 13. 待补充（由你调研的 AMap 文档确认后填写）
- 具体 API 名称/URL/版本/鉴权方式（Key/签名）
- 字段映射表（AMap → 我们的 Schema）
- 错误码与重试白名单
- 费用与配额（免费/付费额度、QPS、日限额）
- 示例请求与响应（脱敏）
