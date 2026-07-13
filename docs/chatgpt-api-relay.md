# ChatGPT API Relay

AI ROUTER is a ChatGPT API relay with an OpenAI-compatible endpoint. It is built for developers who want to create an API key, call ChatGPT-style models from existing code, and track usage, quota, balance, and subscription status in one dashboard.

Website: https://ai-router.dev

## What a ChatGPT API relay does

A ChatGPT API relay sits between your application and model providers. Your app sends an OpenAI-style request to the relay endpoint. The relay handles the API key, quota, usage accounting, and service-side routing details shown in the product dashboard.

For developers, the practical goal is simple: keep the request format familiar and reduce the amount of setup needed to start testing.

## Endpoint

```text
https://api.ai-router.dev/v1
```

## Why developers use AI ROUTER

- OpenAI-compatible request format.
- API key management from the dashboard.
- Usage, balance, quota, and subscription visibility.
- Daily and weekly ChatGPT API plans for testing and short projects.
- Examples for curl, Python, and Node.js.

## Basic curl request

```bash
export AI_ROUTER_MODEL="model_id_from_your_dashboard_or_models_response"

curl https://api.ai-router.dev/v1/chat/completions \
  -H "Authorization: Bearer $AI_ROUTER_API_KEY" \
  -H "Content-Type: application/json" \
  --data-binary @- <<JSON
  {
    "model": "$AI_ROUTER_MODEL",
    "messages": [
      { "role": "user", "content": "Say hello from a ChatGPT API relay." }
    ]
  }
JSON
```

Use a model ID returned by `GET /v1/models` for your API key or shown in your
AI ROUTER dashboard.

## Related pages

- [OpenAI-compatible endpoint](openai-compatible-endpoint.md)
- [API key usage tracking](api-key-usage-tracking.md)
- [Daily and weekly plans](daily-weekly-chatgpt-api-plans.md)
- [OpenAI Python SDK setup](integrations/openai-python-sdk.md)
- [OpenAI Node.js SDK setup](integrations/openai-nodejs-sdk.md)
- [Client integration guides](integrations/README.md)

AI ROUTER is not an official OpenAI service. It provides an independent relay service with OpenAI-compatible request patterns.
