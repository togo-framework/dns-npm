package npm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/togo-framework/dns"
)

func TestUpsertProxyHost(t *testing.T) {
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/tokens":
			json.NewEncoder(w).Encode(map[string]string{"token": "jwt"})
		case r.URL.Path == "/api/nginx/proxy-hosts" && r.Method == http.MethodPost:
			if r.Header.Get("Authorization") != "Bearer jwt" {
				t.Errorf("auth=%q", r.Header.Get("Authorization"))
			}
			json.NewDecoder(r.Body).Decode(&body)
			json.NewEncoder(w).Encode(map[string]any{"id": 7})
		default:
			http.Error(w, "no", 404)
		}
	}))
	defer srv.Close()

	p := &provider{base: srv.URL, email: "a@b.c", secret: "pw", hc: srv.Client()}
	id, err := p.UpsertProxyHost(context.Background(), dns.ProxyHost{Domain: "app.example.com", Upstream: "http://127.0.0.1:8080", SSL: true})
	if err != nil {
		t.Fatalf("upsert: %v", err)
	}
	if id != "7" {
		t.Fatalf("id=%q", id)
	}
	if body["forward_host"] != "127.0.0.1" || body["forward_port"].(float64) != 8080 {
		t.Fatalf("body=%v", body)
	}
}

func TestSplitUpstream(t *testing.T) {
	s, h, port := splitUpstream("https://api.internal:9443")
	if s != "https" || h != "api.internal" || port != 9443 {
		t.Fatalf("got %s %s %d", s, h, port)
	}
}
