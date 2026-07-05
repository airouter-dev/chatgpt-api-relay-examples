# OpenAI-Compatible Endpoint

AI ROUTER provides an OpenAI-compatible endpoint for ChatGPT API relay usage:

```text
https://api.ai-router.dev/v1
```

For many projects, the migration path is to keep the OpenAI SDK and change the SDK base URL to AI ROUTER.

## Python base URL

```python
from openai import OpenAI

client = OpenAI(
    api_key="replace_with_your_ai_router_api_key",
    base_url="https://api.ai-router.dev/v1",
)
```

## Node.js base URL

```js
import OpenAI from "openai";

const client = new OpenAI({
  apiKey: process.env.AI_ROUTER_API_KEY,
  baseURL: "https://api.ai-router.dev/v1",
});
```

## What compatibility means

OpenAI-compatible means AI ROUTER follows familiar OpenAI-style API patterns where supported. It does not mean AI ROUTER is an official OpenAI service, and it does not mean every model, parameter, or behavior is identical. Check the dashboard for available models and plan rules.

## Recommended testing flow

1. Create a small test API key.
2. Set the SDK base URL.
3. Send one chat completion request.
4. Check usage in the dashboard.
5. Increase plan size only after the workflow is stable.

Related:

- [ChatGPT API relay](chatgpt-api-relay.md)
- [Python SDK integration](integrations/openai-python-sdk.md)
- [Node.js SDK integration](integrations/openai-nodejs-sdk.md)

