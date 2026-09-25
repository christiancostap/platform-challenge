package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_ReturnsJSONResponse(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/teste?nome=joao", strings.NewReader(`{"hello":"world"}`))
	rec := httptest.NewRecorder()

	handler(rec, req)

	var res Response
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("error decoding response: %v", err)
	}

	if res.Path != "/teste" {
		t.Fatalf("path = %q, want %q", res.Path, "/teste")
	}

	if res.Params.Get("nome") != "joao" {
		t.Fatalf("param nome = %q, want %q", res.Params.Get("nome"), "joao")
	}

	if res.Body["hello"] != "world" {
		t.Fatalf("body.hello = %v, want %q", res.Body["hello"], "world")
	}
}
