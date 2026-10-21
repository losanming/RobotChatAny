import SwiftUI

struct ContentView: View {
    @StateObject private var hm = HealthManager()
    @State private var log: String = ""
    @State private var metrics: [HealthMetric] = []
    
    var body: some View {
        NavigationView {
            VStack(spacing: 16) {
                Button("请求授权并采集（最近7天，本地展示）") {
                    Task {
                        do {
                            try await hm.requestAuthorization()
                            let result = try await hm.collectRecentDays(days: 7)
                            metrics = result
                            append("Collected metrics: \(result.count)")
                        } catch {
                            append("Error: \(error.localizedDescription)")
                        }
                    }
                }
                .buttonStyle(.borderedProminent)
                .padding(.top, 24)
                
                List(metrics, id: \.external_id) { m in
                    VStack(alignment: .leading, spacing: 4) {
                        Text("type: \(m.type)  ts: \(m.ts)")
                            .font(.subheadline)
                        Text(summary(of: m))
                            .font(.caption.monospaced())
                            .foregroundColor(.secondary)
                    }
                }
                .frame(maxHeight: 260)

                ScrollView {
                    Text(log)
                        .font(.system(.footnote, design: .monospaced))
                        .frame(maxWidth: .infinity, alignment: .leading)
                }
            }
            .padding()
            .navigationTitle("HealthMVP")
        }
    }
    
    private func append(_ line: String) {
        log.append("\n" + line)
        print(line)
    }
    
    private func summary(of m: HealthMetric) -> String {
        let kv = m.payload.map { "\($0.key)=\(String(format: "%.2f", $0.value))" }.sorted().joined(separator: ", ")
        return kv.isEmpty ? "<no-payload>" : kv
    }
}

#Preview {
    ContentView()
}
