# Go image generation demo

This example is a production-minded, two-step image workflow in Go:

1. `POST /v1/images/generations` turns a text prompt into an image.
2. The generated file is sent as multipart form data to `POST /v1/images/edits`.
3. The first `data[0].b64_json` or `data[0].url` item is resolved and written to disk.

It is intentionally built on Go's standard library. There is no vendor lock-in or
SDK magic: the small `internal/client` package owns HTTP, authentication,
multipart encoding, response decoding, URL downloads, size limits, and errors;
`internal/config` owns environment defaults; `cmd/imagegen` owns CLI concerns.
That separation makes it straightforward to embed the client in a worker, API
handler, or queue consumer without copying the demo's command-line code.

## Try it in two minutes

Get an API key from [AI-ROUTER](https://ai-router.dev) (the API is
OpenAI-compatible), then run:

```bash
export AI_ROUTER_API_KEY="your-key"
go run ./cmd/imagegen \
  --base-url https://api.ai-router.dev/v1 \
  --model gpt-image-2 \
  --prompt "A cinematic product photograph of a glass spaceship on a mossy desk" \
  --edit-prompt "Add a soft sunrise rim light and keep the product silhouette unchanged"
```

The command creates `generated.png`, then `edited.png`. Never commit an API key,
and use a throwaway output directory when running untrusted prompts.

## Configuration

Every flag overrides its environment default:

| Flag | Environment | Default | Meaning |
| --- | --- | --- | --- |
| `--api-key` | `OPENAI_API_KEY` or `AI_ROUTER_API_KEY` | — | Bearer token (required) |
| `--base-url` | `OPENAI_BASE_URL` or `AI_ROUTER_BASE_URL` | `https://api.ai-router.dev/v1` | OpenAI-compatible API root |
| `--model` | `IMAGE_MODEL` or `AI_ROUTER_MODEL` | `gpt-image-2` | `gpt-image-2`, `gpt-image-2.5`, `gpt-image-2.5-flare`, or `gpt-image-2.5-burst` |
| `--size` | `IMAGE_SIZE` or `AI_ROUTER_SIZE` | `1024x1024` | Size accepted by the selected provider |
| `--response-format` | `IMAGE_RESPONSE_FORMAT` | `b64_json` | `b64_json` or `url` |
| `--timeout` | `IMAGE_API_TIMEOUT` | `120` seconds | End-to-end HTTP timeout |
| `--output` | — | `generated.png` | First image path |
| `--edited-output` | — | `edited.png` | Edited image path |

`response_format=url` is useful for large images. The client downloads the
returned URL without forwarding your API bearer token to the object-storage
host. Responses are capped at 100 MiB to protect long-lived workers from an
accidental unbounded allocation.

## Project layout

```text
go/
├── cmd/imagegen/main.go       # flags and generation → edit orchestration
├── internal/config/config.go  # environment-backed defaults
├── internal/client/client.go   # Images API transport and decoding
├── internal/client/client_test.go
└── go.mod
```

## Verify locally

```bash
go test ./...
go vet ./...
go build ./cmd/imagegen
```

The tests use `httptest.Server` to assert JSON fields, multipart fields and
uploaded bytes, and to cover both base64 and URL responses. They do not call a
real provider or require an API key.

## Extending the example

Add provider-specific parameters in `internal/client` rather than in `main`.
Keep the response resolver (`b64_json`/`url`) and transport error handling in
one place so any future operation—variations, masks, or asynchronous jobs—can
reuse the same authentication and timeout policy.
