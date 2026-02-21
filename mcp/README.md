# @solafon/mcp-server

MCP (Model Context Protocol) server for the Solafon platform. Gives AI assistants (Claude, Cursor, VS Code, etc.) direct access to Solafon developer tools.

## What is this?

When connected to your AI assistant, it gets tools to:

- **Send and receive messages** via Bot API
- **Manage webhooks** and bot commands
- **Read platform documentation** without leaving the chat
- **Generate bot code** from templates (Node.js, Python, Go)
- **Make any API request** to Solafon endpoints
- **Debug issues** with health checks and webhook inspection

## Quick Setup

### 1. Build

```bash
cd mcp
npm install
npm run build
```

### 2. Configure Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "solafon": {
      "command": "node",
      "args": ["/absolute/path/to/solanafon-back/mcp/dist/index.js"],
      "env": {
        "SOLAFON_BOT_TOKEN": "your-api-token",
        "SOLAFON_API_URL": "https://api.solafon.com/api/v1"
      }
    }
  }
}
```

### 3. Configure Cursor

Add to `.cursor/mcp.json` in your project:

```json
{
  "mcpServers": {
    "solafon": {
      "command": "node",
      "args": ["./mcp/dist/index.js"],
      "env": {
        "SOLAFON_BOT_TOKEN": "your-api-token"
      }
    }
  }
}
```

### 4. Configure VS Code (Claude Code)

Add to `.vscode/mcp.json`:

```json
{
  "servers": {
    "solafon": {
      "command": "node",
      "args": ["./mcp/dist/index.js"],
      "env": {
        "SOLAFON_BOT_TOKEN": "your-api-token"
      }
    }
  }
}
```

## Available Tools

| Tool | Description |
|------|-------------|
| `solafon_read_docs` | Read documentation by topic |
| `solafon_search_docs` | Search across all documentation |
| `solafon_list_docs` | List available doc topics |
| `solafon_get_bot_info` | Get your bot info |
| `solafon_send_message` | Send message to a user |
| `solafon_get_updates` | Get pending messages (polling) |
| `solafon_set_webhook` | Configure webhook URL |
| `solafon_delete_webhook` | Remove webhook |
| `solafon_get_webhook_info` | Check webhook status |
| `solafon_set_commands` | Define bot commands |
| `solafon_get_commands` | List bot commands |
| `solafon_api_request` | Make any API request |
| `solafon_health_check` | Check API status |
| `solafon_scaffold_bot` | Generate bot code template |
| `solafon_list_templates` | List code templates |

## Available Prompts

| Prompt | Description |
|--------|-------------|
| `create-solafon-bot` | Guided bot creation walkthrough |
| `debug-solafon-bot` | Diagnose and fix bot issues |
| `solafon-api-explorer` | Interactive API exploration |

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SOLAFON_BOT_TOKEN` | For Bot API | — | Your bot API token |
| `SOLAFON_API_URL` | No | `https://api.solafon.com/api/v1` | API base URL |

## Examples

### Create a bot with AI assistance

Just tell your AI assistant:
> "Create a new Solafon bot that responds to /start and /help commands"

The AI will use `solafon_scaffold_bot` to generate code and `solafon_set_commands` to configure commands.

### Debug webhook issues

> "My Solafon bot isn't receiving messages. Help me debug."

The AI will check health, verify token, inspect webhook config, and diagnose the issue.

### Explore the API

> "Show me how to send button messages on Solafon"

The AI will search docs and show you examples with metadata format.

## License

MIT
