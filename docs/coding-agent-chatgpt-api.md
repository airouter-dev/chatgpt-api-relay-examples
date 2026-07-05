# ChatGPT API for Coding Agents

Coding agents and developer automation tools can generate frequent ChatGPT API calls. AI ROUTER is useful when you want an OpenAI-compatible endpoint, API key visibility, and quota tracking for those workflows.

Common examples include local coding assistants, repository analysis scripts, documentation generators, test-data helpers, and internal developer tools.

## Why coding-agent workflows need usage visibility

Coding-agent traffic can be bursty. One task may make many small requests. A separate API key and visible quota help you understand:

- Which workflow used the key.
- Whether quota remains.
- When a daily or weekly plan should be adjusted.
- Whether a script is retrying too aggressively.

## Setup pattern

1. Create one API key per coding-agent tool.
2. Set the OpenAI-compatible base URL:

```text
https://api.ai-router.dev/v1
```

3. Use the model name shown in the dashboard.
4. Run a small task first.
5. Review usage before running larger automation.

## Related examples

- [OpenAI Python SDK integration](integrations/openai-python-sdk.md)
- [OpenAI Node.js SDK integration](integrations/openai-nodejs-sdk.md)
- [ChatGPT API for scripts and automation](use-cases/scripts-and-automation.md)

AI ROUTER supports legitimate development and automation use cases. It is not intended for spam, abuse, or policy-violating activity.

