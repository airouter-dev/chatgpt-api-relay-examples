# ChatGPT API for Scripts and Automation

AI ROUTER can be used as an OpenAI-compatible ChatGPT API relay for scripts and automation workflows.

Scripts often need a simple endpoint, an API key, and a way to check usage after a run. AI ROUTER focuses on that path: create a key, call the endpoint, and review usage in the dashboard.

## Good fit

- Local scripts.
- Batch tests.
- Prompt evaluation.
- Internal automation.
- Documentation helpers.
- Coding-agent support tasks.

## Basic script pattern

```bash
export AI_ROUTER_API_KEY="replace_with_your_api_key"
export AI_ROUTER_MODEL="model_id_from_your_dashboard_or_models_response"

curl https://api.ai-router.dev/v1/chat/completions \
  -H "Authorization: Bearer $AI_ROUTER_API_KEY" \
  -H "Content-Type: application/json" \
  --data-binary @- <<JSON
  {
    "model": "$AI_ROUTER_MODEL",
    "messages": [
      { "role": "user", "content": "Summarize this task in one sentence." }
    ]
  }
JSON
```

## Operational advice

- Use a separate API key per script.
- Add retry limits.
- Log request IDs and errors where available.
- Check dashboard usage after test runs.
- Start with a smaller plan before running large batches.

Related:

- [API key usage tracking](../api-key-usage-tracking.md)
- [Daily and weekly plans](../daily-weekly-chatgpt-api-plans.md)
