# Cursor with AI ROUTER

Use this guide only when your Cursor version exposes an OpenAI-compatible or
custom provider configuration. The names and availability of Cursor settings can
change between releases.

## Configuration values

Enter these values in the applicable custom-provider fields:

- API base URL: `https://api.ai-router.dev/v1`
- API key: an AI ROUTER API key stored in Cursor's secret/configuration store
- Model: an `id` returned by authenticated `GET /v1/models`, or the same current
  model ID shown in the AI ROUTER dashboard

The base URL must include `/v1`. Do not use the website URL as an API base URL.

## Verify the connection

1. Create a dedicated test API key rather than reusing a shared production key.
2. Select or enter one current model ID.
3. Send one short, non-sensitive prompt.
4. Confirm that usage appears for the expected key in the AI ROUTER dashboard.

If Cursor does not list a model automatically, use its manual model-entry field
to enter a current ID that you discovered through `/v1/models` or the dashboard.
If the installed version does not permit a custom OpenAI-compatible endpoint,
use the provider configuration supported by that version instead of attempting
to modify its network behavior.

Related: [client integration overview](README.md) and [BYOK smoke-test
checklist](byok-smoke-test-checklist.md).
