# MCP Integration: Using doopl with Claude and AI Clients

doopl provides an MCP (Model Context Protocol) server that exposes translation capabilities to Claude, Claude Desktop, and other MCP-compatible AI clients.

## Quick Setup (Claude Code)

### 1. Build the MCP Server

```bash
cd /path/to/doopl
task mcp:build
# Binary created at: ./bin/doopl-mcp
```

### 2. Configure Claude Code

Add to your Claude Code settings (`.claude/settings.json`):

```json
{
  "mcpServers": {
    "doopl": {
      "command": "/absolute/path/to/bin/doopl-mcp",
      "args": ["serve"],
      "env": {
        "DEEPL_AUTH_KEY": "your-deepl-api-key"
      }
    }
  }
}
```

Or use your saved config file:
```json
{
  "mcpServers": {
    "doopl": {
      "command": "/absolute/path/to/bin/doopl-mcp",
      "env": {
        "DEEPL_AUTH_KEY": "your-deepl-api-key"
      }
    }
  }
}
```

### 3. Restart Claude Code

Claude will detect the new MCP server and make these tools available to AI agents.

## Available Tools

### translate
Translate text to a target language.

**Input:**
```json
{
  "text": "Text to translate",
  "target_lang": "DE",
  "source_lang": "EN",      // optional, auto-detected if omitted
  "formality": "more",       // optional: default|more|less|prefer_more|prefer_less
  "glossary_id": "abc123"    // optional: glossary ID
}
```

**Output:**
```json
{
  "text": "Translated text",
  "detected_lang": "EN"
}
```

**Example:**
```
Human: "Translate this to German: Hello, world!"
Claude: [calls translate with text="Hello, world!", target_lang="DE"]
Result: { "text": "Hallo, Welt!", "detected_lang": "EN" }
Claude: "The translation is: Hallo, Welt!"
```

### languages
List supported languages for a resource type.

**Input:**
```json
{
  "resource": "translate"  // optional: translate|document|glossary|write (default: translate)
}
```

**Output:**
```json
{
  "languages": [
    { "code": "EN", "name": "English" },
    { "code": "DE", "name": "German" },
    { "code": "FR", "name": "French" },
    ...
  ]
}
```

### usage
Check API quota and character usage for your DeepL account.

**Input:**
```json
{}
```

**Output:**
```json
{
  "character_count": 12345,
  "character_limit": 500000,
  "document_count": 3,
  "document_limit": 50
}
```

## Use Cases

1. **Document Localization** — Translate README.md or documentation to multiple languages
2. **Multilingual Support** — Reply to users or customers in their language
3. **Content Translation** — Bulk translate blog posts, articles, or content
4. **Language Availability** — Verify which languages DeepL supports for a task
5. **Quota-Aware Workflows** — Check usage before performing batch operations

## Example Conversation

```
Human: "I need to translate my README.md to German, French, and Spanish. 
How much of my DeepL quota will I use?"

Claude: [calls usage]
Result: character_count=5000, character_limit=100000

Claude: [reads README.md, estimates 2500 characters]
Claude: "Your README is about 2500 characters. Translating to 3 languages 
would use ~7500 characters (2500 × 3), leaving you with ~87,500 remaining 
from your monthly quota."
```
