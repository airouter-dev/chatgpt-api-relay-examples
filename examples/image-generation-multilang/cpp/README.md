# C++ image generation demo

This sample is a complete text-to-image and image-to-image workflow using a
small C++17 library. It calls the OpenAI-compatible endpoints directly:

```text
POST /v1/images/generations  JSON: model, prompt, size, response_format
POST /v1/images/edits        multipart: image, model, prompt, size, response_format
```

The command first writes `generated.png`, then uploads that file and writes
`edited.png`. The response resolver accepts both `data[0].b64_json` (including
data URLs) and `data[0].url`. HTTP errors include the provider's JSON message,
requests have connect and total timeouts, and response bodies are capped at 100
MiB. Bearer credentials are never attached to a returned image URL download.

## Try it

Install a C++17 compiler, CMake 3.20+, and libcurl development files. On a
vcpkg-based Windows setup, for example:

```powershell
vcpkg install curl:x64-windows
```

Then configure (add the vcpkg toolchain flag when needed):

```bash
export AI_ROUTER_API_KEY="your-key"
cmake -S . -B build -DCMAKE_BUILD_TYPE=Release \
  -DCMAKE_TOOLCHAIN_FILE="$VCPKG_ROOT/scripts/buildsystems/vcpkg.cmake"
cmake --build build --config Release
./build/imagegen \
  --base-url https://api.ai-router.dev/v1 \
  --model gpt-image-2 \
  --prompt "A cinematic product photograph of a glass spaceship on a mossy desk" \
  --edit-prompt "Add a soft sunrise rim light and keep the product silhouette unchanged"
```

Get an API key from [AI-ROUTER](https://ai-router.dev). The same binary can
target OpenAI or another compatible gateway by changing `--base-url`.

## Configuration

CLI flags override environment defaults:

| Flag | Environment | Default |
| --- | --- | --- |
| `--api-key` | `OPENAI_API_KEY` or `AI_ROUTER_API_KEY` | required |
| `--base-url` | `OPENAI_BASE_URL` or `AI_ROUTER_BASE_URL` | `https://api.ai-router.dev/v1` |
| `--model` | `IMAGE_MODEL` or `AI_ROUTER_MODEL` | `gpt-image-2` |
| `--size` | `IMAGE_SIZE` or `AI_ROUTER_SIZE` | `1024x1024` |
| `--response-format` | `IMAGE_RESPONSE_FORMAT` | `b64_json` |
| `--timeout` | `IMAGE_API_TIMEOUT` or `AI_ROUTER_IMAGE_TIMEOUT` | 120 seconds |
| `--output` | — | `generated.png` |
| `--edited-output` | — | `edited.png` |

Supported model names depend on the gateway. This example is ready for
`gpt-image-2`, `gpt-image-2.5`, `gpt-image-2.5-flare`, and
`gpt-image-2.5-burst`; use the size values documented by the selected model.

## Architecture

```text
include/image_client.hpp  public, reusable Config and ImageClient API
src/image_client.cpp      libcurl transport, JSON string extraction, base64,
                          multipart upload, URL download, error/timeout policy
src/main.cpp              CLI parsing and the generation → edit orchestration
tests/image_client_tests.cpp  offline tests for encoding helpers
```

The public client is intentionally independent of the CLI. A web service can
construct `Config`, call `generate`, persist the bytes, and call `edit` without
adopting command-line code. If the API grows (masks, variations, async jobs),
add operations beside `generate`/`edit` and reuse the same response resolver and
transport policy.

## Verify locally

```bash
cmake -S . -B build -DBUILD_TESTING=ON
cmake --build build --config Release
ctest --test-dir build --output-on-failure
```

The test target is offline and does not require a key. It checks standard and
data-URL base64 decoding plus JSON escaping. For an integration test, point the
client at a local HTTP fixture that implements the two documented endpoints.

## Security notes

Keep keys in the environment or a secret manager; do not pass them in source
or commit generated images that may contain sensitive material. Returned URLs
can be hosted on a different object-storage domain, so the implementation
intentionally omits the API Authorization header when downloading them.
