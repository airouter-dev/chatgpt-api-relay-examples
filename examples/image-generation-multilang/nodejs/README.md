# Node.js demo

This directory contains a dependency-free Node.js implementation of the shared generate-then-edit workflow. It uses the built-in `fetch`, `FormData`, `Blob`, and `node:test` APIs available in Node.js 18.17+, so the transport remains visible and easy to move into an HTTP server, serverless function, or queue worker.

```bash
cd nodejs
npm test
AI_ROUTER_API_KEY=... npm run demo -- \
  --prompt "A tiny glass greenhouse on a rainy rooftop" \
  --edit-prompt "Change the lighting to warm golden hour" \
  --model gpt-image-2
```

The command writes `outputs/generated.png` and `outputs/edited.png`. The first file is deliberately persisted before the edit request, so a user can review or reuse the draft even when the second call fails.

## Configuration

| Setting | Environment variable or flag | Default |
| --- | --- | --- |
| API key | `AI_ROUTER_API_KEY`, `OPENAI_API_KEY`, `IMAGE_API_KEY`, `--api-key` | required |
| API base | `AI_ROUTER_BASE_URL`, `IMAGE_API_BASE_URL`, `--base-url` | `https://api.ai-router.dev/v1` |
| Model | `AI_ROUTER_MODEL`, `IMAGE_MODEL`, `--model` | `gpt-image-2` |
| Size | `AI_ROUTER_SIZE`, `IMAGE_SIZE`, `--size` | `1024x1024` |
| Timeout | `IMAGE_TIMEOUT_MS`, `--timeout` | `120000` |

The client sends JSON to `/images/generations` and multipart form data to `/images/edits`, always asking for `response_format=b64_json`. It also accepts a provider `url` response and downloads it without forwarding the bearer token to the asset host. HTTP status handling, bounded diagnostics, base64 decoding, and URL materialization stay in `src/client.mjs`; `src/workflow.mjs` only orchestrates the two calls.

Run the two focused tests with `npm test`. They validate URL normalization and credential fallback without making a paid request. See the shared [API contract](../docs/API_CONTRACT.md), [architecture notes](../docs/ARCHITECTURE.md), and [AI-ROUTER site](https://ai-router.dev).
