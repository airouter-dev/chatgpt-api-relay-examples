# Continue with AI ROUTER

Continue configuration formats and setting names vary by release. In the current
provider schema or settings UI, configure an OpenAI-compatible provider with
these values:

| Setting | Value |
| --- | --- |
| Provider | `openai` or the release's OpenAI-compatible provider |
| API base URL | `https://api.ai-router.dev/v1` |
| API key | an AI ROUTER API key supplied through a local secret or environment variable |
| Model | a current `id` returned by authenticated `GET /v1/models` |

## `config.yaml` example

For Continue releases using the current `config.yaml` schema, the equivalent
local configuration has this shape. Replace both placeholders locally; do not
commit a real key.

```yaml
name: AI ROUTER
version: 0.0.1
schema: v1
models:
  - name: AI ROUTER selected model
    provider: openai
    model: MODEL_ID_FROM_AI_ROUTER
    apiBase: https://api.ai-router.dev/v1
    apiKey: YOUR_AI_ROUTER_KEY
```

`MODEL_ID_FROM_AI_ROUTER` must be an actual current ID returned by `/v1/models`
or shown in the dashboard. When your Continue version supports a secret
reference, use it instead of placing the key directly in the local file.

Do not commit a plaintext key to a project-level configuration file. If your
Continue release supports a secret reference, use that reference; otherwise keep
the local configuration outside version control.

## Verify the configuration

1. Use a dedicated test key and a single current model ID.
2. Run a short, non-streaming prompt first.
3. Verify that the key's usage is visible in the AI ROUTER dashboard.
4. Test streaming, tool calls, and autocomplete independently if your workflow
   relies on them.

If automatic model discovery is incomplete, use the manual model-ID field to
enter a current ID discovered through `/v1/models` or the dashboard. The
configured base URL must end with `/v1`.

Related: [client integration overview](README.md) and [BYOK smoke-test
checklist](byok-smoke-test-checklist.md).
