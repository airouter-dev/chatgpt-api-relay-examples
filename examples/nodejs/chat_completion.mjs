import OpenAI from "openai";

const client = new OpenAI({
  apiKey: process.env.AI_ROUTER_API_KEY,
  baseURL: "https://api.ai-router.dev/v1",
});

const response = await client.chat.completions.create({
  model: "gpt-5.4-mini",
  messages: [
    {
      role: "user",
      content: "Explain what an OpenAI-compatible API relay is in one sentence.",
    },
  ],
});

console.log(response.choices[0].message.content);

