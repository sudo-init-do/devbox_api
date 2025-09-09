package health

import (
    "encoding/json"
    "net/http"
)

type Response struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
    resp := Response{
        Status:  "ok",
        Message: "Devbox API is healthy and running",
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
