# Open WebUI with AI ROUTER

In the Open WebUI administration area, add an OpenAI-compatible connection when
your installed version provides one. The exact menu labels vary by release.

Use these connection values:

- API base URL: `https://api.ai-router.dev/v1`
- API key: an AI ROUTER API key
- Model: a current `id` returned by authenticated `GET /v1/models`, or the same
  model ID shown in the AI ROUTER dashboard

Use the API base URL, not `https://ai-router.dev`. The `/v1` suffix is required
for OpenAI-compatible requests. If an Open WebUI release automatically appends
`/v1`, enter the value its connection test expects rather than creating a
duplicated `/v1/v1` path.

## Verification sequence

1. Add a dedicated, low-scope test key through the connection's secret field.
2. Refresh the model list, or add a current ID discovered through `/v1/models`
   manually if the UI supports manual entry.
3. Send one short, non-streaming message.
4. Confirm the API key's usage in the AI ROUTER dashboard before adding more
   users or enabling additional features.

Never place a shared production key in a public workspace export or a committed
configuration file. Test streaming, file features, and tool calls independently
because OpenAI-compatible endpoints can support different feature subsets.

Related: [client integration overview](README.md) and [BYOK smoke-test
checklist](byok-smoke-test-checklist.md).
