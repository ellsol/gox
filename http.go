package gox

import (
	"encoding/json"
	"net/http"
	"time"
	"log"
)

var LogRequest = true

type healthcheckResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(healthcheckResponse{Status: "OK", Code: 200})
}

func LoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		if (LogRequest) {
			log.Printf("[%s] %q  %v", r.Method, r.URL.String(), t2.Sub(t1))
		}
	}

	return http.HandlerFunc(fn)
}
