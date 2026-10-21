import Foundation
import HealthKit

struct HealthMetric: Codable {
    let type: String
    let ts: Int64
    let payload: [String: Double]
    let external_id: String?
    let raw: [String: String]?
}

@MainActor
final class HealthManager: ObservableObject {
    private let healthStore = HKHealthStore()
    
    func requestAuthorization() async throws {
        guard HKHealthStore.isHealthDataAvailable() else { throw HKError(.errorHealthDataUnavailable) }
        let readTypes: Set<HKObjectType> = [
            HKObjectType.quantityType(forIdentifier: .bodyMass)!,
            HKObjectType.quantityType(forIdentifier: .bodyFatPercentage)!,
            HKObjectType.quantityType(forIdentifier: .stepCount)!,
            HKObjectType.quantityType(forIdentifier: .heartRate)!,
            HKObjectType.quantityType(forIdentifier: .bloodPressureSystolic)!,
            HKObjectType.quantityType(forIdentifier: .bloodPressureDiastolic)!,
            HKObjectType.categoryType(forIdentifier: .sleepAnalysis)!
        ]
        try await healthStore.requestAuthorization(toShare: nil, read: readTypes)
    }
    
    func collectRecentDays(days: Int) async throws -> [HealthMetric] {
        let from = Calendar.current.date(byAdding: .day, value: -days, to: Date()) ?? Date(timeIntervalSinceNow: -Double(days)*86400)
        var results: [HealthMetric] = []
        
        // Quantity samples
        let quantityIdentifiers: [(HKQuantityTypeIdentifier, String, HKUnit, (Double)->Double)] = [
            (.bodyMass, "weight", .gramUnit(with: .kilo)) { $0 },
            (.bodyFatPercentage, "body_fat_pct", .percent()) { $0 * 100.0 },
            (.stepCount, "steps", .count()) { $0 },
            (.heartRate, "hr", HKUnit(from: "count/min")) { $0 },
            (.bloodPressureSystolic, "sys", .millimeterOfMercury()) { $0 },
            (.bloodPressureDiastolic, "dia", .millimeterOfMercury()) { $0 }
        ]
        
        for (id, key, unit, transform) in quantityIdentifiers {
            if let qt = HKObjectType.quantityType(forIdentifier: id) {
                let qs = try await queryQuantitySamples(type: qt, from: from)
                for s in qs {
                    let v = s.quantity.doubleValue(for: unit)
                    let ts = Int64(s.endDate.timeIntervalSince1970)
                    results.append(HealthMetric(type: mapTypeKey(for: id), ts: ts, payload: [key: transform(v)], external_id: s.uuid.uuidString, raw: nil))
                }
            }
        }
        
        // Sleep (category)
        if let ct = HKObjectType.categoryType(forIdentifier: .sleepAnalysis) {
            let cs = try await queryCategorySamples(type: ct, from: from)
            for s in cs {
                let dur = s.endDate.timeIntervalSince(s.startDate)
                let ts = Int64(s.endDate.timeIntervalSince1970)
                let stage = mapSleep(value: s.value)
                results.append(HealthMetric(type: "sleep", ts: ts, payload: ["duration_s": dur], external_id: s.uuid.uuidString, raw: ["stage": stage]))
            }
        }
        
        return results.sorted(by: { $0.ts < $1.ts })
    }
    
    private func queryQuantitySamples(type: HKQuantityType, from: Date) async throws -> [HKQuantitySample] {
        return try await withCheckedThrowingContinuation { cont in
            let predicate = HKQuery.predicateForSamples(withStart: from, end: Date(), options: [])
            let q = HKSampleQuery(sampleType: type, predicate: predicate, limit: HKObjectQueryNoLimit, sortDescriptors: nil) { _, samples, error in
                if let error = error { cont.resume(throwing: error); return }
                cont.resume(returning: samples as? [HKQuantitySample] ?? [])
            }
            healthStore.execute(q)
        }
    }
    
    private func queryCategorySamples(type: HKCategoryType, from: Date) async throws -> [HKCategorySample] {
        return try await withCheckedThrowingContinuation { cont in
            let predicate = HKQuery.predicateForSamples(withStart: from, end: Date(), options: [])
            let q = HKSampleQuery(sampleType: type, predicate: predicate, limit: HKObjectQueryNoLimit, sortDescriptors: nil) { _, samples, error in
                if let error = error { cont.resume(throwing: error); return }
                cont.resume(returning: samples as? [HKCategorySample] ?? [])
            }
            healthStore.execute(q)
        }
    }
    
    private func mapTypeKey(for id: HKQuantityTypeIdentifier) -> String {
        switch id {
        case .bodyMass: return "weight"
        case .bodyFatPercentage: return "weight" // 体脂作为 weight 类型的 payload 字段之一
        case .stepCount: return "activity"
        case .heartRate: return "activity"
        case .bloodPressureSystolic, .bloodPressureDiastolic: return "bp"
        default: return "unknown"
        }
    }
    
    private func mapSleep(value: Int) -> String {
        switch value {
        case HKCategoryValueSleepAnalysis.inBed.rawValue: return "in_bed"
        case HKCategoryValueSleepAnalysis.asleepUnspecified.rawValue: return "asleep"
        case HKCategoryValueSleepAnalysis.awake.rawValue: return "awake"
        default: return "unknown"
        }
    }
}
