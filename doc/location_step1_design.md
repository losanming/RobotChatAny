# 位置模块 Step 1 设计说明（WGS-84 → GCJ-02 与最小地点信息）

> 目标：以最小代价为 AI 大脑提供“可信、可融合”的位置上下文。第一阶段仅做坐标转换与最小地点信息（不强依赖网络），为后续逆地理/天气/交通与情景推理铺路。

> 时间规范：本项目“所有时间字段一律使用毫秒级时间戳（int64，ms）”。

## 1. 范围（Scope）
- 输入：硬件定位数据（通常为 WGS-84，经纬度 + 精度/速度/航向/时间戳）。
- 处理（策略已确认：高德为主、本地为辅）：
  - 境内优先调用高德坐标转换 API（WGS-84→GCJ-02）；失败或受限时回退到本地算法；境外保持 WGS-84。
- 输出：最小地点信息（经纬度 GCJ-02、是否在中国标记、质量评分），并标准化为 LocationEvent（仅元数据）。
- 不含：逆地理（地址/POI）、天气、交通、地理围栏——均留到 Step 2。

## 2. 目录结构（与 `pkg/screenshot/` 风格对齐）
```
/pkg/location/
  ├── types.go           # 数据结构与错误
  ├── transform/         # 坐标系转换（WGS-84 ↔ GCJ-02，可选 BD-09）
  ├── provider/amap/     # 预留：高德（逆地理/天气/交通）封装（Step 2 使用）
  ├── service.go         # 主服务：采样/去抖/质量评分/缓存/降级/事件输出
  ├── config.go          # 配置项（阈值、TTL、超时、速率限制等）
  ├── metrics.go         # 指标埋点（成功率/延迟/缓存命中/过滤比率）
  └── example/           # 使用示例（风格对齐 screenshot/example）
```

## 3. 数据模型（全部时间使用毫秒时间戳 int64，单位：ms）
- RawLocation（硬件输入）
  - ts(int64, ms), lat_wgs84, lon_wgs84
  - accuracy_m, speed_mps, heading_deg, satellite_count
  - source(gps/wifi/cell/ip), device_id
- TransformedLocation（模块输出最小集）
  - ts(int64, ms)
  - lat_gcj02, lon_gcj02
  - lat_wgs84, lon_wgs84
  - is_in_china(bool)
  - quality(float64 ∈ [0,1])
  - adcode(string，可为空，Step 2 补全)
- LocationEvent（投递黑板/总线）
  - 公共：ts(int64, ms), device_id, source, quality
  - 定位：lat_wgs84, lon_wgs84, lat_gcj02, lon_gcj02, is_in_china
  - 预留：geohash, adcode, reason_codes[]

## 4. 坐标系转换（核心）
- 原则：
  - 中国境内必须 WGS-84 → GCJ-02；境外保持 WGS-84。
  - 是否在中国：粗 bbox + 精修多边形判定；性能优先、准确度可配置。
- 接口：
  - transform.InChina(lat, lon) bool
  - 本地算法：transform.WGS84ToGCJ02(lat, lon) (lat2, lon2)
  - 高德 API（主）：GET https://restapi.amap.com/v3/assistant/coordinate/convert
    - 参数：
      - key（必填）：Web 服务 API KEY
      - locations（必填）："lon,lat|lon,lat|..."（经度在前，最多40对）
      - coordsys（可选）：gps | mapbar | baidu | autonavi（不转换），我们使用 gps
      - output（可选）：json|xml，默认 json
      - sig（可选）：数字签名（按需）
    - 返回：
      - status（"1"成功/"0"失败）、info（"OK"或错误）、infocode、locations（转换后坐标串）
      - 注意：返回坐标对之间可能以“;”分隔，需做健壮解析（兼容“;”或“|”）
- 策略：高德为主，调用失败/限流/配额用尽时回退本地算法（记录 reason_codes）。

## 5. 采样/去抖与质量控制
- 采样节流（默认值，可配置）：
  - min_interval_sec: 60
  - min_distance_m: 200
  - min_heading_deg: 35（速度>阈值时生效）
- 质量评分 quality：
  - 基于 accuracy_m、satellite_count、speed_mps、与上一点的合理性（反漂移）
  - 输出 [0,1]，低于 0.3 直接丢弃或延迟处理
- 去漂移：
  - 瞬移/极值（如 >2km 且短时内）触发最近可信位置兜底

## 6. 缓存与最近可信位置
- recent_ttl_sec: 600（默认）
- 弱网场景：返回最近可信 TransformedLocation，并在 reason_codes 中标注（e.g. "fallback_recent_cache"）
 - 转换结果缓存：相同原始坐标点的转换结果短 TTL（例如 5 分钟）以降低外呼；批量请求优先命中缓存。

## 7. 配置（`config.yaml` 建议字段，单位已标注）
```yaml
location:
  enabled: true
  mode: "amap_primary"     # amap_primary | local_primary
  sampling:
    min_interval_sec: 60   # 单位：秒
    min_distance_m: 200    # 单位：米
    min_heading_deg: 35    # 单位：度（仅在速度>阈值时生效）
  quality:
    min_score: 0.3         # [0,1]
  cache:
    recent_ttl_sec: 600    # 单位：秒
amap:
  key: ""              # Step 2 使用
  timeout_ms: 800       # 单位：毫秒（Step 2 使用）
  rate_limit_rps: 5     # 每秒请求数（Step 2 使用）
  retry_ms: [200, 500, 1000]
  convert:
    batch_max: 40           # 单次最多40对坐标
    request_sep: "|"        # 请求locations分隔符
    response_sep_fallbacks: [";", "|"]  # 响应locations可能的分隔符
    cache_ttl_sec: 300      # 转换结果缓存TTL（秒）
```

## 8. 服务编排与接口
- Service 对外：
  - OnRawLocation(raw RawLocation) (TransformedLocation, error)
- 内部流程：
  1) 采样与去抖筛选 → 2) 质量评分 → 3) InChina 判定 → 4) 若在境内则根据 mode：
     - amap_primary：优先批量加入高德转换队列（≤40/批）；若调用失败/超时/配额受限，则回退到本地算法
     - local_primary：本地算法优先；异常时按策略调用高德校验
     5) 形成 TransformedLocation（统一毫秒时间戳）
  6) 更新最近可信位置缓存与转换结果缓存 → 7) 可选：输出 LocationEvent 至 NSQ → 8) 返回结果

## 9. 事件与总线（与现有 NSQ 设计对齐）
- 消息体只含元数据（不含图片等大对象），字段见 LocationEvent。
- 主题建议：`events.location.transformed`
- 消费方：Step 2 的逆地理/天气/交通处理器；AI 大脑上下文黑板聚合器。

## 10. 指标与日志（`metrics.go`）
- location_transform_success_total / failed_total
- location_transform_latency_ms
- location_recent_cache_hit_total
- location_filtered_total（按原因细分：interval/distance/heading/quality）
- location_out_of_china_total
 - amap_convert_success_total / amap_convert_failed_total
 - amap_convert_latency_ms（summary/histogram）
 - amap_convert_batch_size
 - amap_convert_cache_hit_total
 - amap_rate_limited_total / amap_retry_total

## 11. 边界与降级
- 出境：is_in_china=false，lat_gcj02/lon_gcj02 与输入一致；reason_codes=["out_of_china"]
- 弱网/无网：不影响本地转换；无新数据时返回最近可信位置
- 异常值：精度>1000m、速度极值、瞬移，直接丢弃或降级
 - 高德失败/配额：回退本地算法；记录 reason_codes=["amap_failed" 或 "amap_rate_limited"]

## 12. 隐私与合规（默认最小化）
- 本地优先：本阶段不上传任何与人相关的敏感信息；仅坐标元数据。
- 模糊化预留：后续上报可采用 geohash-6/7；门牌级地址需二次授权（Step 2）。
- 加密：传输 TLS，服务侧最小权限访问与访问审计（与项目基线一致）。

## 13. 与 Step 2/3 的衔接
- Step 2：`provider/amap` 增加 ReverseGeocode/Weather/Traffic 封装 + TTL 缓存 + 限流/重试。
- Step 3：AI 大脑读取黑板上的 LocationEvent + Weather + Traffic + 历史行为，生成候选意图与话术。

## 14. 验收标准（Step 1）
- 准确度：境内转换误差符合行业通用标准（数十米量级）。
- 时延：单次转换 < 5ms；对外调用（含采样逻辑）< 20ms。
- 稳定性：弱网/无网/出境可用；有清晰日志与指标。
- 一致性：所有时间字段使用毫秒时间戳（int64, ms）。

## 15. 开发清单（无代码实现指引）
- 创建目录 `pkg/location/` 与子模块骨架
- 实现 `transform.InChina` 与 `transform.WGS84ToGCJ02`
- 实现 `service.OnRawLocation` 的采样/质量/转换/缓存流程
- 定义 `types.go` 与 `LocationEvent` Schema（与 NSQ 对齐）
- 补充 `config.go` 默认值与读取
- 写 `metrics.go` 指标与日志位点
- example：为开发者提供最小调用示例（与 screenshot/example 风格一致）

## 16. 参考与对齐
- 与现有 `pkg/screenshot/` 的服务编排、质量过滤、去重、埋点风格保持一致
- 与 NSQ “指针+元数据”的设计理念相容（位置事件本身为纯元数据）
- 与 Eino/ReAct/Agent 编排思路兼容：位置作为上下文信号输入小模型候选与大模型解释

## 17. 附录（字段/阈值/埋点/日志）

### 17.1 字段字典（节选）
- RawLocation
  - ts(int64, ms)：硬件定位时间
  - lat_wgs84/lon_wgs84(float64)：原始经纬度（WGS-84）
  - accuracy_m(float64)：精度（米）
  - speed_mps(float64)：速度（米/秒）
  - heading_deg(float64)：航向（度，0–360）
  - satellite_count(int)：卫星数量（可选）
  - source(string)：gps/wifi/cell/ip
  - device_id(string)：设备标识
- TransformedLocation
  - ts(int64, ms)、lat_gcj02/lon_gcj02、lat_wgs84/lon_wgs84
  - is_in_china(bool)、quality([0,1])、adcode(string, 可空)
- LocationEvent（NSQ）
  - ts(int64, ms), device_id, source, quality
  - lat_wgs84, lon_wgs84, lat_gcj02, lon_gcj02, is_in_china
  - geohash(optional), adcode(optional), reason_codes[](optional)

### 17.2 默认阈值（可配置）
- 采样节流：min_interval_sec=60、min_distance_m=200、min_heading_deg=35
- 质量阈值：quality.min_score=0.3
- 最近可信缓存：recent_ttl_sec=600

### 17.3 指标（Metrics）命名
- 计数：
  - location_transform_success_total
  - location_transform_failed_total
  - location_recent_cache_hit_total
  - location_filtered_total{reason="interval|distance|heading|quality"}
  - location_out_of_china_total
- 时延：
  - location_transform_latency_ms（summary/histogram）

### 17.4 日志格式（建议 JSON 行）
```json
{"ts": 1727170000123, "level":"info", "module":"location", "event":"transform", "device_id":"dev-001", "in_china":true, "lat_wgs84":31.2304, "lon_wgs84":121.4737, "lat_gcj02":31.2245, "lon_gcj02":121.4800, "quality":0.86}
```
```json
{"ts": 1727170060456, "level":"warn", "module":"location", "event":"filtered", "reason":"quality", "score":0.21, "threshold":0.3}
```
```json
{"ts": 1727170120789, "level":"info", "module":"location", "event":"fallback_recent_cache", "age_ms":42000}
```
