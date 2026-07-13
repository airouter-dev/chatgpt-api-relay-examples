#!/usr/bin/env bash
set -euo pipefail

: "${AI_ROUTER_API_KEY:?Set AI_ROUTER_API_KEY before running this example.}"
: "${AI_ROUTER_MODEL:?Set AI_ROUTER_MODEL to a model from /v1/models or the dashboard.}"

curl https://api.ai-router.dev/v1/chat/completions \
  -H "Authorization: Bearer ${AI_ROUTER_API_KEY}" \
  -H "Content-Type: application/json" \
  --data-binary @- <<JSON
  {
    "model": "${AI_ROUTER_MODEL}",
    "messages": [
      {
        "role": "user",
        "content": "Explain what an OpenAI-compatible API relay is in one sentence."
      }
    ]
  }
JSON
