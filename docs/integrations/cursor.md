# Cursor Compatibility Status for AI ROUTER

Cursor's current official BYOK documentation lists API-key configuration for
OpenAI, Anthropic, Google, Azure OpenAI, and AWS Bedrock. It does not document an
arbitrary OpenAI-compatible base URL or custom provider. AI ROUTER therefore
does not claim a supported direct Cursor integration at this time.

Do not paste an AI ROUTER key into Cursor's OpenAI key field or attempt to alter
Cursor's network behavior. Use a client that officially permits a custom
OpenAI-compatible endpoint, such as Continue, LiteLLM, or Open WebUI, or call the
AI ROUTER API directly with an OpenAI SDK.

If Cursor publishes a custom-endpoint workflow in the future, revalidate its
base-URL, model-discovery, feature, and key-handling behavior before documenting
AI ROUTER setup steps here.

Source: [Cursor's official BYOK documentation](https://cursor.com/help/models-and-usage/api-keys).

Related: [client integration overview](README.md) and [BYOK smoke-test
checklist](byok-smoke-test-checklist.md).
