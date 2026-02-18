#!/usr/bin/env bash
set -euo pipefail

CONFIG="$HOME/.picoclaw/config.json"
API_KEY=""
API_BASE=""

if [ -f "$CONFIG" ]; then
  if command -v jq >/dev/null 2>&1; then
    API_KEY=$(jq -r '.providers.openrouter.api_key // empty' "$CONFIG")
    API_BASE=$(jq -r '.providers.openrouter.api_base // empty' "$CONFIG")
  else
    echo "jq is required to read $CONFIG. Install jq or set OPENROUTER_API_KEY and OPENROUTER_API_BASE." >&2
  fi
fi

API_KEY="${API_KEY:-${OPENROUTER_API_KEY:-}}"
API_BASE="${API_BASE:-${OPENROUTER_API_BASE:-}}"

if [ -z "$API_KEY" ]; then
  echo "No OpenRouter API key found. Set OPENROUTER_API_KEY or add it to $CONFIG" >&2
  exit 2
fi

MODEL="arcee-ai/trinity-large-preview:free"
PROMPT="${1:-What is 2+2?}"

if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required to format JSON output. Install jq to get pretty output." >&2
fi

PAYLOAD=$(jq -nc --arg model "$MODEL" --arg content "$PROMPT" '{model:$model, messages:[{role:"user", content:$content}], max_tokens:100}')

# Determine endpoint order:
# 1) OPENROUTER_API_ENDPOINT env
# 2) API_BASE from config if it appears to be a full model URL
# 3) API_BASE + /chat/completions or /api/v1/chat/completions
# 4) fallback to model-specific URL

ENDPOINT=""
if [ -n "${OPENROUTER_API_ENDPOINT:-}" ]; then
  ENDPOINT="$OPENROUTER_API_ENDPOINT"
elif [ -n "$API_BASE" ]; then
  # If API_BASE already looks like a model path (contains "arcee-ai/"), use it directly
  if echo "$API_BASE" | grep -q "arcee-ai/"; then
    ENDPOINT="$API_BASE"
  elif echo "$API_BASE" | grep -q "/api"; then
    ENDPOINT="${API_BASE%/}/chat/completions"
  else
    ENDPOINT="${API_BASE%/}/api/v1/chat/completions"
  fi
else
  ENDPOINT="https://openrouter.ai/arcee-ai/trinity-large-preview:free"
fi

echo "Using endpoint: $ENDPOINT"
curl -sS -X POST "$ENDPOINT" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d "$PAYLOAD" | jq .
