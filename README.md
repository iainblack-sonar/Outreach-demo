# Outreach Platform Services — Security Demo

This repo contains a Go backend API and a TypeScript AI agent integration layer that together
represent the kind of polyglot, agent-assisted codebase modern engineering teams are building.

It's intentionally realistic: the code is structured, commented, and looks like it was written
by a capable engineer (or a capable AI). That's the point. The vulnerabilities here aren't
obvious toy examples — they're the kind of issues that ship in production when SAST is absent
or noisy enough that developers tune it out.

---

## The App

| Module | Language | Purpose |
|--------|----------|---------|
| `backend/` | Go | REST API — user management, auth, file handling, proxy |
| `agent/` | TypeScript | AI agent integration layer — analysis pipeline, tool calls |

---

## What Sonar Finds

### Critical — fix before merge

| File | Issue | Why it matters |
|------|-------|----------------|
| `backend/handlers/users.go:14` | SQL injection (string concat) | HTTP param flows directly into `db.Query()` — full DB read/write |
| `backend/handlers/users.go:26` | SQL injection (`fmt.Sprintf`) | Same impact, different pattern — both must be caught |
| `backend/handlers/users.go:46` | Second-order SQL injection | Value read from DB used in new query — harder to spot manually |
| `backend/handlers/files.go:14` | Path traversal | User-supplied filename → `os.ReadFile("/var/data/" + filename)` — reads arbitrary files |
| `backend/handlers/files.go:37` | SSRF | User-controlled URL passed to `http.Get()` — can reach internal services |
| `backend/handlers/auth.go:56` | Command injection | User input passed unsanitized to `exec.Command("grep", userInput, ...)` |
| `agent/src/analyzer.ts:10` | Code injection (`eval`) | Agent-supplied expression evaluated directly — full code execution |
| `agent/src/analyzer.ts:30` | Code injection (`eval`) | Same pattern in JSON response handler |
| `agent/src/auth.ts:33` | JWT alg:none accepted | `algorithms: ['none', 'HS256']` — unsigned tokens bypass all auth |

### High

| File | Issue | Why it matters |
|------|-------|----------------|
| `backend/handlers/auth.go:13-16` | Hardcoded credentials (4) | DB password, JWT secret, admin key, internal token in source |
| `agent/src/auth.ts:6-9` | Hardcoded API keys (2) | Sonar + OpenAI keys committed — exposed on clone |
| `agent/src/analyzer.ts:22` | XSS (innerHTML) | User-supplied HTML rendered without sanitization |
| `backend/utils/crypto.go:16` | Insecure random (OTP) | `math/rand` for 2FA codes — predictable, enumerable |
| `backend/utils/crypto.go:27` | MD5 password hashing | MD5 is broken for passwords — trivially crackable with GPU |
| `agent/src/api-client.ts:6` | Prototype pollution | `for...in` merge without `hasOwnProperty` guard — overrides `Object.prototype` |

---

## The Snyk Comparison

Snyk is a dependency scanner with limited SAST. Here's what each tool sees in this repo:

| Finding category | Snyk | Sonar |
|-----------------|------|-------|
| Dependency CVEs (lodash, axios, qs, jwt-go) | ✅ | ✅ |
| SQL injection (taint flow, HTTP → DB) | ❌ | ✅ |
| Command injection | ❌ | ✅ |
| Path traversal | ❌ | ✅ |
| SSRF | ❌ | ✅ |
| `eval()` / code injection | ❌ | ✅ |
| JWT alg:none | Limited | ✅ |
| Hardcoded secrets in source | Limited | ✅ |
| Prototype pollution in your code (not deps) | ❌ | ✅ |
| Insecure random for crypto operations | ❌ | ✅ |

**The dependency CVEs are table stakes.** You already have Wiz Code and GHAS covering that
ground. The SAST findings — the ones in *your code* — are what neither of those tools catches.

---

## Vulnerable Dependencies (SCA)

| Package | Version | CVE | Impact |
|---------|---------|-----|--------|
| `lodash` | 4.17.15 | CVE-2021-23337 | Prototype pollution via `merge()` |
| `axios` | 0.21.1 | CVE-2021-3749 | SSRF via open redirect |
| `qs` | 6.5.2 | CVE-2022-24999 | Prototype pollution |
| `minimist` | 1.2.5 | CVE-2021-44906 | Prototype pollution |
| `github.com/dgrijalva/jwt-go` | v3.2.0 | CVE-2020-26160 | JWT audience claim validation bypass |

---

## Quality Gate

The Sonar Quality Gate is configured to **block merge** if any of the following are present
on new code:

- Any issue rated **Blocker** or **Critical**
- Any hardcoded credential or secret
- Security hotspot review coverage < 80%

This PR would currently **fail** the gate. The CI badge below shows live status.

[![Quality Gate Status](https://sonarqube.io/api/project_badges/measure?project=iainblack-sonar_Outreach-demo&metric=alert_status&token=YOUR_BADGE_TOKEN)](https://sonarqube.io/dashboard?id=iainblack-sonar_Outreach-demo)

---

## MCP — Connect Your Agent to Sonar

See [`mcp/README.md`](./mcp/README.md) for how to wire Cursor, Claude Code, or any MCP-compatible
agent to query Sonar findings live from your editor.

**Example agent prompt once connected:**
> *"What are the critical security issues in `backend/handlers/users.go`? Show me the exact line and explain why it's exploitable."*
