#!/bin/bash
# Quick testing script for doopl with your own DeepL API token
# Usage: ./QUICK_TEST.sh "your-deepl-api-key"

set -e

API_KEY="${1:-$DEEPL_AUTH_KEY}"

if [ -z "$API_KEY" ]; then
    echo "Usage: ./QUICK_TEST.sh <your-deepl-api-key>"
    echo "  or: DEEPL_AUTH_KEY=... ./QUICK_TEST.sh"
    echo ""
    echo "Get your API key from: https://www.deepl.com/pro-api"
    exit 1
fi

echo "🚀 doopl Quick Test Suite"
echo "========================="
echo ""

# Export for all commands
export DEEPL_AUTH_KEY="$API_KEY"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Phase 1: Unit Tests (no API calls needed)${NC}"
echo "Run: task test"
task test
echo -e "${GREEN}✓ Unit tests passed${NC}"
echo ""

echo -e "${BLUE}Phase 2: Build CLI and MCP${NC}"
echo "Run: task cli:build && task mcp:build"
task cli:build
task mcp:build
echo -e "${GREEN}✓ Both binaries built successfully${NC}"
echo ""

echo -e "${BLUE}Phase 3: Doctor Check${NC}"
echo "Run: ./bin/doopl doctor"
./bin/doopl doctor
echo -e "${GREEN}✓ Authentication verified${NC}"
echo ""

echo -e "${BLUE}Phase 4: Quick Translation Test${NC}"
echo "Run: ./bin/doopl translate 'Hello, world!' --target-lang ES"
RESULT=$(./bin/doopl translate 'Hello, world!' --target-lang ES)
echo "Result: $RESULT"
echo -e "${GREEN}✓ Translation works${NC}"
echo ""

echo -e "${BLUE}Phase 5: Languages${NC}"
echo "Run: ./bin/doopl languages --resource translate"
./bin/doopl languages --resource translate | head -5
echo "..."
echo -e "${GREEN}✓ Languages endpoint works${NC}"
echo ""

echo -e "${BLUE}Phase 6: Usage Check${NC}"
echo "Run: ./bin/doopl usage"
./bin/doopl usage
echo -e "${GREEN}✓ Usage endpoint works${NC}"
echo ""

echo -e "${BLUE}Phase 7: JSON Output${NC}"
echo "Run: ./bin/doopl translate 'Hello' --target-lang FR --output json"
./bin/doopl translate 'Hello' --target-lang FR --output json | head -3
echo "..."
echo -e "${GREEN}✓ JSON output works${NC}"
echo ""

echo "╔════════════════════════════════════════════════════════════╗"
echo "║  ✅ ALL TESTS PASSED                                       ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo "Next steps:"
echo "  - See TESTING_PLAYBOOK.md for full testing guide"
echo "  - Run './bin/doopl --help' to see all commands"
echo "  - Try: ./bin/doopl translate 'Your text' --target-lang <CODE>"
echo ""
