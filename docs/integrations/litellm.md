# LiteLLM with AI ROUTER

LiteLLM Proxy can route an OpenAI-compatible model entry to AI ROUTER. First
export an API key and select a current model ID returned by authenticated
`GET /v1/models`.

```bash
export AI_ROUTER_API_KEY="replace_with_your_api_key"
export AI_ROUTER_MODEL="model_id_from_your_models_response"
```

Use the model ID in the LiteLLM configuration. Replace the placeholder below
with the value of `AI_ROUTER_MODEL`; do not include the literal placeholder in
the final YAML.

```yaml
model_list:
  - model_name: ai-router-selected-model
    litellm_params:
      model: openai/MODEL_ID_FROM_AI_ROUTER
      api_base: https://api.ai-router.dev/v1
      api_key: os.environ/AI_ROUTER_API_KEY
```

`api_base` must point to the API endpoint, including `/v1`, rather than the
website root. Keep the API key in the environment or another secret manager.
This is a LiteLLM Proxy configuration: client applications call the LiteLLM
endpoint, and LiteLLM forwards the selected model to AI ROUTER.

## Verify the route

1. Start LiteLLM with the configuration.
2. Send one short request to `ai-router-selected-model`.
3. Confirm the requested model ID and usage in your LiteLLM logs and the AI
   ROUTER dashboard.
4. Test streaming and tool calls separately if you enable them.

Related: [LiteLLM's official OpenAI-compatible endpoint documentation](https://docs.litellm.ai/docs/providers/openai_compatible),
[client integration overview](README.md), and [BYOK smoke-test
checklist](byok-smoke-test-checklist.md).
