# OpenAI Node.js SDK with AI ROUTER

This example shows how to use the OpenAI Node.js SDK with AI ROUTER's OpenAI-compatible ChatGPT API relay endpoint.

Endpoint:

```text
https://api.ai-router.dev/v1
```

## Install

```bash
npm install openai
```

## Set environment variable

```bash
export AI_ROUTER_API_KEY="replace_with_your_api_key"
export AI_ROUTER_MODEL="model_id_from_your_dashboard_or_models_response"
```

## Node.js example

```js
import OpenAI from "openai";

if (!process.env.AI_ROUTER_API_KEY || !process.env.AI_ROUTER_MODEL) {
  throw new Error("Set AI_ROUTER_API_KEY and AI_ROUTER_MODEL before running this example.");
}

const client = new OpenAI({
  apiKey: process.env.AI_ROUTER_API_KEY,
  baseURL: "https://api.ai-router.dev/v1",
});

const response = await client.chat.completions.create({
  model: process.env.AI_ROUTER_MODEL,
  messages: [
    { role: "user", content: "Say hello from AI ROUTER." },
  ],
});

console.log(response.choices[0].message.content);
```

Use a model ID returned by authenticated `GET /v1/models` or shown in your AI
ROUTER dashboard.

Related:

- [OpenAI-compatible endpoint](../openai-compatible-endpoint.md)
- [ChatGPT API relay](../chatgpt-api-relay.md)
