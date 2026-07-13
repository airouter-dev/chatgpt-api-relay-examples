# Open WebUI with AI ROUTER

In the Open WebUI administration area, add an OpenAI-compatible connection when
your installed version provides one. The exact menu labels vary by release.

Use these connection values:

- API base URL: `https://api.ai-router.dev/v1`
- API key: an AI ROUTER API key
- Model: a current `id` returned by authenticated `GET /v1/models`, or the same
  model ID shown in the AI ROUTER dashboard

Use the exact API base URL above, without a trailing slash, rather than
`https://ai-router.dev`. Open WebUI appends endpoint paths such as `/models`, so
the configured URL must already include the single `/v1` suffix.

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

Related: [Open WebUI's official OpenAI-compatible provider documentation](https://docs.openwebui.com/getting-started/quick-start/connect-a-provider/starting-with-openai-compatible),
[client integration overview](README.md), and [BYOK smoke-test
checklist](byok-smoke-test-checklist.md).
