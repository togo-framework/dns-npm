<!-- togo-header -->
<div align="center">
  <picture><source media="(prefers-color-scheme: dark)" srcset=".github/assets/togo-mark-dark.svg" /><img src=".github/assets/togo-mark.svg" alt="ToGO" height="64" /></picture>
  <h1>togo-framework/dns-npm</h1>
  <p>Nginx Proxy Manager driver for togo dns.</p>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1F8A99" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/dns-npm"><img src="https://pkg.go.dev/badge/github.com/togo-framework/dns-npm.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/dns-npm
```
<!-- /togo-header -->

[Nginx Proxy Manager](https://nginxproxymanager.com/) driver for togo's
[`dns`](https://github.com/togo-framework/dns) subsystem. Manages reverse-proxy
hosts over the NPM API.

## Config

| Env | Meaning |
|-----|---------|
| `DNS_DRIVER` | set to `npm` |
| `NPM_BASE_URL` | NPM admin URL, e.g. `http://localhost:81` (required) |
| `NPM_EMAIL` / `NPM_PASSWORD` | admin credentials (mint a token), **or** |
| `NPM_TOKEN` | a pre-minted API token |

```go
svc, _ := dns.FromKernel(k)
svc.UpsertProxyHost(ctx, dns.ProxyHost{Domain: "app.example.com", Upstream: "http://127.0.0.1:8080", SSL: true})
```

`SSL: true` requests a forced-TLS host with Let's Encrypt agreement. DNS records
and gateway routes return `dns.ErrUnsupported`.

<!-- togo-sponsors -->
---
<div align="center">
  <h3>Premium sponsors</h3>
  <p><a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp; <a href="https://one-studio.co"><strong>One Studio</strong></a></p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
