package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"os"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type Response struct {
	Headers http.Header    `json:"headers"`
	Body    map[string]any `json:"body"`
	Path    string         `json:"path"`
	Params  url.Values     `json:"params"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET")
}

func handler(w http.ResponseWriter, r *http.Request) {
	logger.Info("request received", "method", r.Method, "path", r.URL.Path, "params", r.URL.Query())

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var body map[string]any
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		body = map[string]any{}
	}
	
	enableCors(&w)

	res := Response{
		Headers: r.Header,
		Body:    body,
		Path:    r.URL.Path,
		Params:  r.URL.Query(),
	}

	json.NewEncoder(w).Encode(res)
	logger.Info("response sent", "status", http.StatusOK, "response", res)
}

func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}
	logger.Info("server started", "port", port)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)
}
