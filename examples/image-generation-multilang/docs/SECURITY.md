# Security notes

- Put `AI_ROUTER_API_KEY` in the environment or a local ignored `.env` file. Never commit it, print it, or paste it into an issue.
- The examples use HTTPS by default and send the key only in the `Authorization` header.
- API error bodies are truncated by each client so an upstream diagnostic cannot flood logs. Keys are not included in error strings.
- Generated images may contain sensitive prompts or personal data. The demos write them to a local ignored directory; review and delete artifacts according to your retention policy.
- Treat downloaded image URLs, if your account returns them, as untrusted input. Production applications should enforce an outbound allow-list and size/time limits before downloading.
- These are educational integrations. Add authentication, request validation, retries with jitter, observability, virus scanning, and durable storage before exposing a public upload endpoint.
