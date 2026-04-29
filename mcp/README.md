# Sonar MCP Server — Setup Guide

The SonarQube MCP server lets your AI coding agent (Cursor, Claude Code, Copilot) query
Sonar findings directly from your editor. Instead of checking the SonarQube Cloud dashboard
separately, your agent can ask: *"What are the critical security issues in this file?"* and
get findings in context, inline, while it writes or reviews code.

## Quick Setup

### 1. Get your Sonar token
Go to [SonarQube Cloud](https://sonarqube.io) → My Account → Security → Generate Token.
Save the token — you'll need it below.

### 2. Configure your editor

**Cursor** — edit `.cursor/mcp.json` in this repo (already included):
```json
{
  "mcpServers": {
    "sonarqube": {
      "command": "npx",
      "args": ["-y", "@sonar/sonarqube-mcp-server@latest"],
      "env": {
        "SONARQUBE_URL": "https://sonarqube.io",
        "SONAR_TOKEN": "your-token-here"
      }
    }
  }
}
```

**Claude Code** — edit `.claude/mcp.json` (already included), same format.

### 3. Restart your editor and confirm the MCP server appears in the tools list.

## What You Can Ask Your Agent

Once connected, try these prompts:

```
What are the critical security issues in this project?

Show me all SQL injection vulnerabilities in the backend.

Which issues in auth.go would a CISO consider blockers?

What would fail the quality gate right now?

Show me the hotspots in the agent/src directory.
```

The agent pulls live findings from SonarQube Cloud — the same results the CI scan produces —
and can reason about them in the context of the code you're currently editing.

## How It Works in an Agentic Workflow

```
Developer (or Copilot/Rovo/Kiro) writes code
        ↓
Agent calls Sonar MCP: "any issues with this function?"
        ↓
Sonar returns: SQL injection in users.go:14 — taint from HTTP param to db.Query()
        ↓
Agent fixes the issue before the code ever reaches CI
        ↓
Quality Gate passes on first push
```

This is the shift-left model at the agentic layer — Sonar becomes part of the agent's
feedback loop, not a post-hoc gate.

## Docs
- [SonarQube MCP Server](https://docs.sonarsource.com/sonarqube-cloud/improving/sonarqube-mcp-server/)
- [SonarQube Cloud API](https://sonarqube.io/web_api)
