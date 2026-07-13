# BYOK Smoke-Test Checklist

Use this checklist when validating an OpenAI-compatible client against AI ROUTER
or another OpenAI-compatible endpoint. It is designed to catch setup problems
before a client is used with production traffic.

## Minimum contract

1. The client can configure an API base URL.
2. The client can configure a bearer/API key without committing it to source.
3. The client permits a manual model ID.
4. Model discovery does not prevent a user from selecting a current permitted
   model ID manually.

## Endpoint checks

Set a test key, then list the models available to that key:

```bash
export AI_ROUTER_API_KEY="replace_with_your_api_key"

curl -sS https://api.ai-router.dev/v1/models \
  -H "Authorization: Bearer $AI_ROUTER_API_KEY"
```

For a model ID returned by that response, verify one short completion:

```bash
export AI_ROUTER_MODEL="model_id_from_your_models_response"

curl -sS https://api.ai-router.dev/v1/chat/completions \
  -H "Authorization: Bearer $AI_ROUTER_API_KEY" \
  -H "Content-Type: application/json" \
  --data-binary @- <<JSON
{
  "model": "$AI_ROUTER_MODEL",
  "messages": [{"role": "user", "content": "Reply with the word: ready"}]
}
JSON
```

Confirm the following:

- The base URL is `https://api.ai-router.dev/v1`.
- A failed model lookup or completion exposes a useful error rather than silently
  switching to another provider.
- Authentication, unavailable-model, quota/billing, and network failures can be
  distinguished in the client logs.
- The requested model ID is visible in redacted diagnostics where your client
  supports it.
- Usage is visible for the expected test key in the AI ROUTER dashboard.

## Feature checks

Test each feature separately rather than inferring support from a successful
basic chat completion:

- Non-streaming chat completion
- Streaming response
- Tool/function calling
- File or image inputs
- Token and usage fields used by your application

Do not describe account credit, plan terms, pricing, or model availability from
an old integration note. Confirm current terms in the product dashboard before
communicating them to users.
