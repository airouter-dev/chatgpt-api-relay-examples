#!/usr/bin/env bash
set -euo pipefail

: "${AI_ROUTER_API_KEY:?Set AI_ROUTER_API_KEY before running this example.}"

curl https://api.ai-router.dev/v1/chat/completions \
  -H "Authorization: Bearer ${AI_ROUTER_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-5.4-mini",
    "messages": [
      {
        "role": "user",
        "content": "Explain what an OpenAI-compatible API relay is in one sentence."
      }
    ]
  }'

