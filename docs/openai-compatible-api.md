# OpenAI-Compatible API Notes

AI ROUTER exposes an OpenAI-compatible API endpoint:

```text
https://api.ai-router.dev/v1
```

For many SDKs, the integration change is the `base_url` or `baseURL` field plus an AI ROUTER API key.

## Python

```python
from openai import OpenAI

client = OpenAI(
    api_key="replace_with_your_api_key",
    base_url="https://api.ai-router.dev/v1",
)
```

## Node.js

```js
import OpenAI from "openai";

const client = new OpenAI({
  apiKey: process.env.AI_ROUTER_API_KEY,
  baseURL: "https://api.ai-router.dev/v1",
});
```

## Operational notes

- Use the model names available in your dashboard.
- Store API keys in environment variables or secret managers.
- Track quota and usage from the dashboard before scaling traffic.
- Start with a small plan when testing a new workflow.

AI ROUTER is not an official OpenAI service. It provides an independent relay service with OpenAI-compatible request patterns.

