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
```

## Node.js example

```js
import OpenAI from "openai";

const client = new OpenAI({
  apiKey: process.env.AI_ROUTER_API_KEY,
  baseURL: "https://api.ai-router.dev/v1",
});

const response = await client.chat.completions.create({
  model: "gpt-5.4-mini",
  messages: [
    { role: "user", content: "Say hello from AI ROUTER." },
  ],
});

console.log(response.choices[0].message.content);
```

Use the model names available in your AI ROUTER dashboard.

Related:

- [OpenAI-compatible endpoint](../openai-compatible-endpoint.md)
- [ChatGPT API relay](../chatgpt-api-relay.md)

