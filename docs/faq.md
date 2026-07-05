# FAQ

## Is AI ROUTER an official OpenAI service?

No. AI ROUTER is an independent ChatGPT API relay service with OpenAI-compatible request patterns. Model availability, plan rules, and pricing are shown in the dashboard.

## What endpoint should I use?

Use:

```text
https://api.ai-router.dev/v1
```

## Can I use existing OpenAI SDK code?

In many cases, yes. Keep the SDK and set the base URL to AI ROUTER's OpenAI-compatible endpoint. See:

- [OpenAI Python SDK integration](integrations/openai-python-sdk.md)
- [OpenAI Node.js SDK integration](integrations/openai-nodejs-sdk.md)

## Which model should I use?

Use the model names available in your AI ROUTER dashboard. Examples in this repository use `gpt-5.4-mini` as a low-cost default example.

## Can I track API key usage?

Yes. AI ROUTER is designed around API key management and usage visibility, including usage, balance, quota, and subscription status where available.

## Are plans unlimited?

No. AI ROUTER does not market unlimited usage. Plans have quota and period rules shown in the product dashboard.

## Which languages are supported?

The website has localized paths for 48 language options. See [supported language landing pages](languages.md).

Priority SEO languages:

- English: https://ai-router.dev
- Chinese: https://ai-router.dev/cn
- Russian: https://ai-router.dev/ru
- Persian: https://ai-router.dev/fa

