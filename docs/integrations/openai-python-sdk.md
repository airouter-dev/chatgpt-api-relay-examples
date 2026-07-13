# OpenAI Python SDK with AI ROUTER

This example shows how to use the OpenAI Python SDK with AI ROUTER's OpenAI-compatible ChatGPT API relay endpoint.

Endpoint:

```text
https://api.ai-router.dev/v1
```

## Install

```bash
pip install openai
```

## Set environment variable

```bash
export AI_ROUTER_API_KEY="replace_with_your_api_key"
export AI_ROUTER_MODEL="model_id_from_your_dashboard_or_models_response"
```

## Python example

```python
import os
from openai import OpenAI

client = OpenAI(
    api_key=os.environ["AI_ROUTER_API_KEY"],
    base_url="https://api.ai-router.dev/v1",
)

response = client.chat.completions.create(
    model=os.environ["AI_ROUTER_MODEL"],
    messages=[
        {"role": "user", "content": "Say hello from AI ROUTER."}
    ],
)

print(response.choices[0].message.content)
```

Use a model ID returned by authenticated `GET /v1/models` or shown in your AI
ROUTER dashboard.

Related:

- [OpenAI-compatible endpoint](../openai-compatible-endpoint.md)
- [ChatGPT API relay](../chatgpt-api-relay.md)
