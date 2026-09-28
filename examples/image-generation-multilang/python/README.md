# Python image generation demo

This is a small but production-minded reference client for [AI-ROUTER](https://ai-router.dev). It demonstrates a complete creative loop rather than a one-line API call:

1. `POST /v1/images/generations` turns a written brief into a new image.
2. The returned image is decoded from `b64_json` (or downloaded when the provider returns a URL).
3. `POST /v1/images/edits` receives that saved image as multipart form data and applies a second instruction.
4. Both artifacts are written to an output directory so the result can be inspected, published, or passed to another pipeline.

The project intentionally uses Python's standard library (`urllib`, `json`, and `argparse`) instead of hiding the protocol behind a heavyweight SDK. That makes the example useful with OpenAI, AI-ROUTER, an on-premise gateway, or any other service that implements the OpenAI Images API shape. It also makes the request boundary easy to audit before adapting it to a web worker, queue consumer, or batch job.

## Why this example is useful

* **OpenAI-compatible by design.** The base URL is configurable, while the endpoint paths remain explicit and easy to compare with API documentation.
* **A real two-step workflow.** The generated file becomes the input to the edit call; there is no fake placeholder image and no assumption that an edit is just another text prompt.
* **Safe defaults for a demo.** Credentials come from environment variables, failures include HTTP status without printing the bearer token, and image bytes are written only to the selected output directory.
* **Easy to extend.** Transport (`client.py`), settings (`config.py`), orchestration (`workflow.py`), and presentation (`cli.py`) have separate responsibilities. A service can replace one module without rewriting the others.
* **Marketing-ready output.** The generated and edited artifacts have predictable names, making it straightforward to add a gallery, attach UTM metadata, or publish examples in a product landing page.

## Quick start

Python 3.10 or newer is required. No runtime dependency installation is needed.

```bash
cd marketing/image-generation-multilang/python

# AI-ROUTER users can point the same client at their compatible gateway.
export IMAGE_API_KEY="your-key"
export AI_ROUTER_BASE_URL="https://api.ai-router.dev/v1"
export IMAGE_MODEL="gpt-image-2"

python main.py \
  --prompt "A cinematic editorial illustration of a glass greenhouse in a desert at sunrise" \
  --edit-prompt "Add a small red kite in the sky and keep the greenhouse architecture unchanged" \
  --output-dir ./outputs
```

On Windows PowerShell:

```powershell
$env:IMAGE_API_KEY = "your-key"
$env:AI_ROUTER_BASE_URL = "https://api.ai-router.dev/v1"
python .\main.py --prompt "A paper-cut map of a coastal city" --edit-prompt "Turn the sky into a warm twilight gradient" --output-dir .\outputs
```

The command prints paths similar to:

```text
Generated image: outputs/generated.png
Edited image:    outputs/edited.png
```

You can also install the local console entry point while developing:

```bash
python -m pip install -e .
image-demo --prompt "..." --edit-prompt "..."
```

## Configuration

CLI flags override environment variables. `IMAGE_API_KEY` is preferred; `OPENAI_API_KEY` and `AI_ROUTER_API_KEY` are accepted as convenient fallbacks.

| Setting | Environment variable | Default | Purpose |
| --- | --- | --- | --- |
| API key | `IMAGE_API_KEY` | required | Bearer credential; never commit it |
| API base | `IMAGE_API_BASE_URL` or `AI_ROUTER_BASE_URL` | `https://api.ai-router.dev/v1` | Host or `/v1` root of an OpenAI-compatible gateway |
| Model | `IMAGE_MODEL` / `--model` | `gpt-image-2` | Also works with `gpt-image-2.5`, `gpt-image-2.5-flare`, or `gpt-image-2.5-burst` when enabled upstream |
| Size | `IMAGE_SIZE` / `--size` | `1024x1024` | Provider-supported output size |
| Quality | `IMAGE_QUALITY` / `--quality` | `auto` | Provider-supported quality tier |
| Response | fixed in this demo | `b64_json` | Portable response mode used by both endpoints |
| Timeout | `IMAGE_TIMEOUT_SECONDS` / `--timeout` | `120` | Per-request timeout in seconds |

The URL normalizer accepts either `https://gateway.example` or `https://gateway.example/v1`; it adds `/v1` exactly once. Use an HTTPS URL in production and keep keys in a secret manager or local environment, not in shell history or source control.

## Architecture and request contract

```text
cli.py  ──>  config.py  ──>  ImagesClient  ──>  /v1/images/generations
  │                              │
  └──────>  ImageWorkflow  <─────┘
                                  └────> /v1/images/edits (multipart image)
```

* `config.py` validates and normalizes runtime settings.
* `client.py` owns HTTP headers, JSON/multipart encoding, response decoding, URL downloads, and provider errors. It accepts either `data[0].b64_json` or `data[0].url`.
* `workflow.py` coordinates the mandatory generate → save → edit → save sequence and returns typed paths.
* `cli.py` parses human-facing arguments and turns expected failures into a non-zero exit code.

Generation sends JSON like:

```json
{
  "model": "gpt-image-2",
  "prompt": "...",
  "size": "1024x1024",
  "quality": "auto",
  "response_format": "b64_json"
}
```

The edit request is `multipart/form-data` with `model`, `prompt`, `size`, `quality`, `response_format=b64_json`, and an `image` file field. If your gateway uses a different field name or requires a mask, the only code that needs adapting is `ImagesClient.edit`.

## Testing and extension points

Run the dependency-free unit checks from this directory:

```bash
python -m unittest discover -s tests -v
python -m compileall -q image_demo main.py
```

For an application integration, inject a fake `ImagesClient` into `ImageWorkflow` and test prompts, retries, and storage without making paid API calls. Add retry/backoff at the transport boundary for an asynchronous production service, and record request IDs from your gateway in structured logs. This demo deliberately avoids automatic retries so a marketing experiment cannot silently duplicate a billable generation.

## Responsible use

Follow the selected provider's image policy, obtain rights for source images, and disclose AI-generated or edited assets where your audience expects that context. The project is an educational integration sample, not a promise that every listed model or size is enabled on every account.
