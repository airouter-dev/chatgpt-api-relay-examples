import os

from openai import OpenAI


client = OpenAI(
    api_key=os.environ["AI_ROUTER_API_KEY"],
    base_url="https://api.ai-router.dev/v1",
)

response = client.chat.completions.create(
    model="gpt-5.4-mini",
    messages=[
        {
            "role": "user",
            "content": "Explain what an OpenAI-compatible API relay is in one sentence.",
        }
    ],
)

print(response.choices[0].message.content)

