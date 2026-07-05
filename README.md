# ChatGPT API Relay Examples

OpenAI-compatible ChatGPT API relay examples for AI ROUTER.

AI ROUTER lets developers call ChatGPT-style models through an OpenAI-compatible endpoint, manage API keys, and check usage, quota, balance, and subscription status from the dashboard.

Website:

- English: https://ai-router.dev
- 中文: https://ai-router.dev/cn
- Русский: https://ai-router.dev/ru
- فارسی: https://ai-router.dev/fa

API endpoint:

```text
https://api.ai-router.dev/v1
```

## Examples

| Runtime | File |
| --- | --- |
| curl | [examples/curl/chat-completions.sh](examples/curl/chat-completions.sh) |
| Python | [examples/python/chat_completion.py](examples/python/chat_completion.py) |
| Node.js | [examples/nodejs/chat_completion.mjs](examples/nodejs/chat_completion.mjs) |

## Quick start

Set your API key:

```bash
export AI_ROUTER_API_KEY="replace_with_your_api_key"
```

Call the OpenAI-compatible chat completions endpoint:

```bash
curl https://api.ai-router.dev/v1/chat/completions \
  -H "Authorization: Bearer $AI_ROUTER_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-5.4-mini",
    "messages": [
      { "role": "user", "content": "Write one sentence about OpenAI-compatible APIs." }
    ]
  }'
```

Use the model name shown in your AI ROUTER dashboard. The example uses `gpt-5.4-mini` as a low-cost default.

## Documentation

- English: [docs/en](docs/en/README.md)
- 中文: [docs/cn](docs/cn/README.md)
- Русский: [docs/ru](docs/ru/README.md)
- فارسی: [docs/fa](docs/fa/README.md)
- All supported language landing pages: [docs/languages.md](docs/languages.md)
- OpenAI-compatible API notes: [docs/openai-compatible-api.md](docs/openai-compatible-api.md)
- International SEO link map: [docs/international-seo.md](docs/international-seo.md)

## What this repository contains

This repository contains public examples and developer documentation only. It does not contain the AI ROUTER backend, routing logic, billing logic, account scheduling, or security implementation.

AI ROUTER is not an official OpenAI service. It provides an independent relay service with OpenAI-compatible request patterns. Availability, model names, plans, and pricing are shown in the product dashboard.

## Keywords

ChatGPT API relay, OpenAI-compatible API, AI API relay, LLM API, API key management, usage tracking, daily ChatGPT API plan, weekly ChatGPT API plan, developer AI API.

