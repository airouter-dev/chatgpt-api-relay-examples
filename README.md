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
| Image generation and editing | [Python, Node.js, C++, Go, Rust, and PHP](examples/image-generation-multilang/README.md) |

## Quick start

Set your API key and choose a model available to that key:

```bash
export AI_ROUTER_API_KEY="replace_with_your_api_key"
export AI_ROUTER_MODEL="model_id_from_your_dashboard_or_models_response"
```

With an authenticated API key, you can list the model IDs available to it:

```bash
curl -sS https://api.ai-router.dev/v1/models \
  -H "Authorization: Bearer $AI_ROUTER_API_KEY"
```

Copy an `id` from the response into `AI_ROUTER_MODEL`. The available catalog can
vary by API key, account configuration, and time; the dashboard is an equivalent
source of the current model ID.

Call the OpenAI-compatible chat completions endpoint:

```bash
curl https://api.ai-router.dev/v1/chat/completions \
  -H "Authorization: Bearer $AI_ROUTER_API_KEY" \
  -H "Content-Type: application/json" \
  --data-binary @- <<JSON
  {
    "model": "$AI_ROUTER_MODEL",
    "messages": [
      { "role": "user", "content": "Write one sentence about OpenAI-compatible APIs." }
    ]
  }
JSON
```

Use the model name returned by `/v1/models` or shown in your AI ROUTER dashboard.

## Documentation

- English: [docs/en](docs/en/README.md)
- 中文: [docs/cn](docs/cn/README.md)
- Русский: [docs/ru](docs/ru/README.md)
- فارسی: [docs/fa](docs/fa/README.md)
- SEO topic index: [docs/seo-index.md](docs/seo-index.md)
- ChatGPT API relay: [docs/chatgpt-api-relay.md](docs/chatgpt-api-relay.md)
- OpenAI-compatible endpoint: [docs/openai-compatible-endpoint.md](docs/openai-compatible-endpoint.md)
- API key usage tracking: [docs/api-key-usage-tracking.md](docs/api-key-usage-tracking.md)
- Daily and weekly plans: [docs/daily-weekly-chatgpt-api-plans.md](docs/daily-weekly-chatgpt-api-plans.md)
- Coding-agent workflows: [docs/coding-agent-chatgpt-api.md](docs/coding-agent-chatgpt-api.md)
- FAQ: [docs/faq.md](docs/faq.md)
- Glossary: [docs/glossary.md](docs/glossary.md)
- All supported language landing pages: [docs/languages.md](docs/languages.md)
- Localized GitHub docs for all supported languages: [docs/localized-docs.md](docs/localized-docs.md)
- OpenAI-compatible API notes: [docs/openai-compatible-api.md](docs/openai-compatible-api.md)
- International SEO link map: [docs/international-seo.md](docs/international-seo.md)
- Multilingual image generation API examples: [examples/image-generation-multilang](examples/image-generation-multilang/README.md)

## Use-case pages

- [ChatGPT API for AI app development](docs/use-cases/ai-app-development.md)
- [ChatGPT API for scripts and automation](docs/use-cases/scripts-and-automation.md)

## Integration pages

- [OpenAI Python SDK with AI ROUTER](docs/integrations/openai-python-sdk.md)
- [OpenAI Node.js SDK with AI ROUTER](docs/integrations/openai-nodejs-sdk.md)
- [Cursor compatibility status](docs/integrations/cursor.md)
- [Continue](docs/integrations/continue.md)
- [LiteLLM](docs/integrations/litellm.md)
- [Open WebUI](docs/integrations/open-webui.md)
- [BYOK smoke-test checklist](docs/integrations/byok-smoke-test-checklist.md)

## What this repository contains

This repository contains public examples and developer documentation only. It does not contain the AI ROUTER backend, routing logic, billing logic, account scheduling, or security implementation.

AI ROUTER is not an official OpenAI service. It provides an independent relay service with OpenAI-compatible request patterns. Availability, model names, plans, and pricing are shown in the product dashboard.

## Keywords

ChatGPT API relay, OpenAI-compatible API, AI API relay, LLM API, API key management, usage tracking, daily ChatGPT API plan, weekly ChatGPT API plan, developer AI API.
