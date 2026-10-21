import Foundation

struct UploadEnvelope: Codable {
    let user_id: String
    let source: String
    let device_id: String
    let metrics: [HealthMetric]
}

final class Uploader {
    static let shared = Uploader()
    private init() {}
    
    private var baseURL: URL!
    private var tokenProvider: (() -> String)?
    private var debugLog = true
    
    func configure(baseURL: String, tokenProvider: @escaping () -> String) async throws {
        guard let url = URL(string: baseURL) else { throw URLError(.badURL) }
        self.baseURL = url
        self.tokenProvider = tokenProvider
    }
    
    func upload(userId: String, deviceId: String, source: String, metrics: [HealthMetric]) async throws {
        let env = UploadEnvelope(user_id: userId, source: source, device_id: deviceId, metrics: metrics)
        let data = try JSONEncoder().encode(env)
        var req = URLRequest(url: baseURL.appendingPathComponent("/api/health/report"))
        req.httpMethod = "POST"
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        if let token = tokenProvider?() { req.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization") }
        req.setValue(String(Int(Date().timeIntervalSince1970)), forHTTPHeaderField: "X-Client-Ts")
        req.httpBody = data
        
        let (respData, resp) = try await URLSession.shared.data(for: req)
        if debugLog {
            print(String(data: respData, encoding: .utf8) ?? "<no-body>")
            print(resp)
        }
    }
}
