#!/bin/bash

# Batch translation example: translate multiple texts from a file
# Prerequisites: export DEEPL_AUTH_KEY="your-api-key"

set -e

# Create a temporary input file with texts to translate (one per line)
INPUT_FILE=$(mktemp)
cat > "$INPUT_FILE" << 'EOF'
Hello, world!
How are you?
Thank you for your help.
EOF

echo "=== Batch Translation Example ==="
echo ""
echo "Input file ($INPUT_FILE):"
cat "$INPUT_FILE"
echo ""

# Process each line
echo "Translating to German..."
while IFS= read -r text; do
    result=$(doopl translate "$text" DE --output json | jq -r '.text')
    echo "  '$text' => '$result'"
done < "$INPUT_FILE"

echo ""
echo "Cleanup:"
rm "$INPUT_FILE"
echo "Done!"
