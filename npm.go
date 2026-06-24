// Package npm is a togo dns driver for Nginx Proxy Manager (NPM). It manages
// reverse-proxy hosts over the NPM API (token auth from email+password). DNS
// records and gateway routes return ErrUnsupported.
//
// Install: `togo install togo-framework/dns-npm`, set DNS_DRIVER=npm.
package npm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/togo-framework/dns"
	"github.com/togo-framework/togo"
)

func init() {
	dns.RegisterDriver("npm", func(k *togo.Kernel) (dns.Provider, error) {
		base := os.Getenv("NPM_BASE_URL")
		if base == "" {
			return nil, fmt.Errorf("dns-npm: NPM_BASE_URL not set (e.g. http://localhost:81)")
		}
		p := &provider{
			base:   strings.TrimRight(base, "/"),
			email:  os.Getenv("NPM_EMAIL"),
			secret: os.Getenv("NPM_PASSWORD"),
			token:  os.Getenv("NPM_TOKEN"),
			hc:     &http.Client{Timeout: 20 * time.Second},
		}
		if p.token == "" && (p.email == "" || p.secret == "") {
			return nil, fmt.Errorf("dns-npm: set NPM_TOKEN or NPM_EMAIL+NPM_PASSWORD")
		}
		return p, nil
	})
}

type provider struct {
	base, email, secret string
	token               string
	mu                  sync.Mutex
	hc                  *http.Client
}

func (p *provider) auth(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.token != "" {
		return p.token, nil
	}
	body, _ := json.Marshal(map[string]string{"identity": p.email, "secret": p.secret})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, p.base+"/api/tokens", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var out struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil || out.Token == "" {
		return "", fmt.Errorf("dns-npm: token request failed (status %d)", resp.StatusCode)
	}
	p.token = out.Token
	return p.token, nil
}

func (p *provider) do(ctx context.Context, method, path string, body any, out any) error {
	tok, err := p.auth(ctx)
	if err != nil {
		return err
	}
	var rdr io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rdr = bytes.NewReader(b)
	}
	req, _ := http.NewRequestWithContext(ctx, method, p.base+path, rdr)
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("dns-npm: %s %s -> %d: %s", method, path, resp.StatusCode, string(raw))
	}
	if out != nil && len(raw) > 0 {
		return json.Unmarshal(raw, out)
	}
	return nil
}

// splitUpstream turns http://host:8080 into scheme, host, port.
func splitUpstream(up string) (scheme, host string, port int) {
	scheme = "http"
	port = 80
	if u, err := url.Parse(up); err == nil && u.Host != "" {
		if u.Scheme != "" {
			scheme = u.Scheme
		}
		host = u.Hostname()
		if p := u.Port(); p != "" {
			port, _ = strconv.Atoi(p)
		} else if scheme == "https" {
			port = 443
		}
		return
	}
	host = up
	return
}

func (p *provider) UpsertProxyHost(ctx context.Context, h dns.ProxyHost) (string, error) {
	scheme, host, port := splitUpstream(h.Upstream)
	payload := map[string]any{
		"domain_names":           []string{h.Domain},
		"forward_scheme":         scheme,
		"forward_host":           host,
		"forward_port":           port,
		"ssl_forced":             h.SSL,
		"block_exploits":         true,
		"allow_websocket_upgrade": true,
		"access_list_id":         0,
		"certificate_id":         0,
		"meta":                   map[string]any{"letsencrypt_agree": h.SSL},
		"advanced_config":        "",
		"locations":              []any{},
	}
	var out struct {
		ID int `json:"id"`
	}
	if err := p.do(ctx, http.MethodPost, "/api/nginx/proxy-hosts", payload, &out); err != nil {
		return "", err
	}
	return strconv.Itoa(out.ID), nil
}

func (p *provider) DeleteProxyHost(ctx context.Context, id string) error {
	return p.do(ctx, http.MethodDelete, "/api/nginx/proxy-hosts/"+url.PathEscape(id), nil, nil)
}

func (p *provider) UpsertRecord(context.Context, string, dns.Record) (string, error) {
	return "", dns.ErrUnsupported
}
func (p *provider) DeleteRecord(context.Context, string, string) error { return dns.ErrUnsupported }
func (p *provider) ListRecords(context.Context, string) ([]dns.Record, error) {
	return nil, dns.ErrUnsupported
}
func (p *provider) UpsertRoute(context.Context, dns.Route) (string, error) {
	return "", dns.ErrUnsupported
}
func (p *provider) DeleteRoute(context.Context, string) error { return dns.ErrUnsupported }
