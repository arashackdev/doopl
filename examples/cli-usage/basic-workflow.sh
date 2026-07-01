#!/bin/bash

# Basic doopl CLI workflow
# Prerequisites: export DEEPL_AUTH_KEY="your-api-key"

set -e

echo "=== doopl CLI Basic Workflow ==="
echo ""

# 1. Check available languages
echo "1. List supported languages for translation:"
doopl languages translate --output text | head -10
echo "   ... (and more)"
echo ""

# 2. Check API usage
echo "2. Check your API quota:"
doopl usage --output text
echo ""

# 3. Simple translation
echo "3. Translate a single text:"
echo "   Input: 'Hello, world!'"
doopl translate "Hello, world!" DE --output text
echo ""

# 4. Translate with options
echo "4. Translate with formality option:"
doopl translate "Can you help me?" DE --formality more --output text
echo ""

# 5. JSON output (useful for scripting)
echo "5. Get JSON output for parsing:"
doopl translate "Hello" FR --output json
echo ""

# 6. Rich terminal output
echo "6. Rich TUI output:"
doopl translate "Hello" ES --output tui
