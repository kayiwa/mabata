# mabata

Mabata is a starter kit for a self-hosted Go service that uses Microsoft Entra ID for OIDC login and DuckDB as an embedded analytics engine.

## Features

- OIDC sign-in with Microsoft Entra ID
- Signed cookie session
- Embedded DuckDB
- Restricted query registry
- Simple HTML UI
- Caddy and systemd deployment examples

## Requirements

- Go 1.24+
- C toolchain for CGO
- DuckDB Go driver dependencies
- Microsoft Entra ID app registration

## Entra app registration

Create a web app registration and configure:

- Redirect URI: `http://127.0.0.1:8080/auth/callback`
- ID tokens: enabled
- Client secret: create one and copy it to `.env`

Optional:

- Add group claims if you want group-based authorization

## Quick start

1. Copy `.env.example` to `.env`
2. Fill in Entra settings
3. Run:

   ```bash
   just tidy
   just run
   ```

4. Open [http://127.0.0.1:8080](http://127.0.0.1:8080)

### Security model

By design no arbitrary SQL execution is exposed. A fixed query registry can be found at `internal/duck/queries.go`
