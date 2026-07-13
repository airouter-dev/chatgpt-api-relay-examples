# Client Integrations

AI ROUTER exposes OpenAI-compatible request patterns at:

```text
https://api.ai-router.dev/v1
```

Before configuring any client, create an API key and choose a model available to
that key. The model catalog can change and can differ by account configuration,
so do not copy a model ID from an old example.

```bash
export AI_ROUTER_API_KEY="replace_with_your_api_key"

curl -sS https://api.ai-router.dev/v1/models \
  -H "Authorization: Bearer $AI_ROUTER_API_KEY"
```

Copy an `id` from the response into the client configuration, or select the
same current model ID in the AI ROUTER dashboard. Keep keys in a local secret
store or environment variable; never commit a production key to a repository.

## Guides

- [OpenAI Python SDK](openai-python-sdk.md)
- [OpenAI Node.js SDK](openai-nodejs-sdk.md)
- [Cursor compatibility status](cursor.md)
- [Continue](continue.md)
- [LiteLLM](litellm.md)
- [Open WebUI](open-webui.md)
- [BYOK smoke-test checklist](byok-smoke-test-checklist.md)

Each client has its own provider schema and feature support. Test a non-streaming
chat completion first, then separately validate streaming, tool calls, and any
other feature your application requires.

AI ROUTER is an independent service, not an official OpenAI service. OpenAI-
compatible request patterns do not guarantee that every model, parameter, or
client feature is available for every account.
