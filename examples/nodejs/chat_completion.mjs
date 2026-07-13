import OpenAI from "openai";

const client = new OpenAI({
  apiKey: process.env.AI_ROUTER_API_KEY,
  baseURL: "https://api.ai-router.dev/v1",
});

if (!process.env.AI_ROUTER_MODEL) {
  throw new Error("Set AI_ROUTER_MODEL to a model from /v1/models or the dashboard.");
}

const response = await client.chat.completions.create({
  model: process.env.AI_ROUTER_MODEL,
  messages: [
    {
      role: "user",
      content: "Explain what an OpenAI-compatible API relay is in one sentence.",
    },
  ],
});

console.log(response.choices[0].message.content);
