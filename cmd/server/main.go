package main

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "strings"
    "time"
)

// This server is intentionally minimal and does not import internal packages.
// It performs a simple direct request to an IP detection service and returns
// the parsed result as JSON. This avoids module resolution issues during
// early development.

type ipResponse struct {
    IP string `json:"ip"`
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/check/public-ip", publicIPHandler)

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 60 * time.Second,
    }

    log.Println("ProxyDoctor server starting on :8080")
    if err := srv.ListenAndServe(); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}

func publicIPHandler(w http.ResponseWriter, r *http.Request) {
    // Try ipify first, then icanhazip
    services := []string{
        "https://api.ipify.org?format=json",
        "https://icanhazip.com/",
        "https://ifconfig.me/ip",
    }

    var ip string
    client := &http.Client{Timeout: 10 * time.Second}
    for _, url := range services {
        resp, err := client.Get(url)
        if err != nil {
            continue
        }
        body, _ := io.ReadAll(resp.Body)
        resp.Body.Close()

        // Try parse JSON {"ip":"x.x.x.x"}
        if strings.HasPrefix(strings.TrimSpace(string(body)), "{") {
            var r ipResponse
            if err := json.Unmarshal(body, &r); err == nil && r.IP != "" {
                ip = r.IP
                break
            }
        }

        // Fallback: plain text
        txt := strings.TrimSpace(string(body))
        if txt != "" {
            ip = txt
            break
        }
    }

    w.Header().Set("Content-Type", "application/json")
    if ip == "" {
        w.WriteHeader(http.StatusInternalServerError)
        _ = json.NewEncoder(w).Encode(map[string]string{"error": "unable to detect public IP"})
        return
    }

    _ = json.NewEncoder(w).Encode(map[string]string{"ip": ip})
}
