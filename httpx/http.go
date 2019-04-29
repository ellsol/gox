package httpx

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"mime"
	"net/http"
	"strings"
	"time"
)

var LogRequest = true

const (
	ContentTypeJson = "application/json"
	ContentTypeCsb  = "text/csv"
)

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

func CorsHandler(next http.Handler) http.Handler {
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

func StartServer(router *mux.Router, port string) {
	log.Println("Starting server on port: ", port)

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(port, handler))
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func HasContentType(r *http.Request, mimetype string) bool {
	contentType := r.Header.Get("Content-type")
	if contentType == "" {
		return mimetype == "application/octet-stream"
	}

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}
	return false
}

func ShowAllRoutes(router *mux.Router) {
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		fmt.Println(t)
		return nil
	})
}

type HttpStatusAndError struct {
	StatusCode int
	Error      error
}

func HttpResponseOk() *HttpStatusAndError {
	return &HttpStatusAndError{
		StatusCode: http.StatusOK,
	}
}

func HttpInternalError(err error) *HttpStatusAndError {
	return &HttpStatusAndError{
		Error:      err,
		StatusCode: http.StatusInternalServerError,
	}
}

func HttpResourceNotFound(label string, resource string) *HttpStatusAndError {
	return &HttpStatusAndError{
		Error:      fmt.Errorf("%v %v not found", label, resource),
		StatusCode: http.StatusNotFound,
	}
}

func HttpConflict(msg string) *HttpStatusAndError {
	return &HttpStatusAndError{
		Error:      fmt.Errorf("conflict: %v", msg),
		StatusCode: http.StatusConflict,
	}
}

func (it HttpStatusAndError) IsError() bool {
	return it.Error != nil
}

func HttpBadRequest(msg string) *HttpStatusAndError {
	return &HttpStatusAndError{
		Error:      fmt.Errorf("bad request: %v", msg),
		StatusCode: http.StatusBadRequest,
	}
}
