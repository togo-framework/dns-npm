# dns-npm — docs

**Nginx Proxy Manager.** Reverse-proxy hosts with Let's Encrypt via the NPM admin API.

## Install

```bash
togo install togo-framework/dns-npm
```

Registers on the [`dns`](https://github.com/togo-framework/dns) base; select it with **dns.provider in togo.yaml (or DNS_DRIVER)**, then use **`togo proxy`**.

## Interface

`Provider` — `UpsertRecord`/`DeleteRecord`/`ListRecords`, `UpsertProxyHost`/`DeleteProxyHost`, `UpsertRoute`/`DeleteRoute`.

## Configuration

| Env var | Description |
|---|---|
| `NPM_BASE_URL` | Nginx Proxy Manager base URL, e.g. `http://localhost:81` (required). |
| `NPM_EMAIL` | NPM admin email for token auth (required unless NPM_TOKEN set). |
| `NPM_PASSWORD` | NPM admin password for token auth (required unless NPM_TOKEN set). |
| `NPM_TOKEN` | Pre-minted NPM API token (alternative to email+password). Optional. |

## Usage & notes

Mints a token from email+password (or uses `NPM_TOKEN`), then creates proxy hosts (domain→upstream) with optional TLS. DNS-record ops return `ErrUnsupported`.

## Example

```bash
togo proxy:host:add app.example.com http://localhost:3000 --provider npm --dry-run
```

## Links

- [Nginx Proxy Manager](https://nginxproxymanager.com/)
- [Marketplace](https://to-go.dev/marketplace)
- [Source](https://github.com/togo-framework/dns-npm)
