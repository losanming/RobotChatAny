package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type HealthMetric struct {
	Type       string             `json:"type"`
	Ts         int64              `json:"ts"`
	Payload    map[string]float64 `json:"payload"`
	ExternalID *string            `json:"external_id"`
	Raw        map[string]string  `json:"raw"`
}

type UploadEnvelope struct {
	UserID   string         `json:"user_id"`
	Source   string         `json:"source"`
	DeviceID string         `json:"device_id"`
	Metrics  []HealthMetric `json:"metrics"`
}

type Resp struct {
	OK    bool `json:"ok"`
	Dedup int  `json:"dedup"`
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var env UploadEnvelope
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&env); err != nil {
		log.Printf("decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("/api/health/report user=%s source=%s device=%s metrics=%d ts=%d", env.UserID, env.Source, env.DeviceID, len(env.Metrics), time.Now().Unix())
	for i, m := range env.Metrics {
		if i >= 5 { // limit log
			break
		}
		b, _ := json.Marshal(m)
		log.Printf("  metric[%d]=%s", i, string(b))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Resp{OK: true, Dedup: 0})
}

func main() {
	addr := ":8088"
	if v := os.Getenv("ADDR"); v != "" {
		addr = v
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health/report", reportHandler)
	log.Printf("healthmock listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
