# AI-ROUTER image generation — Rust

This example is intentionally shaped like a small production client rather
than a one-file snippet. `config` owns environment parsing, `client` owns the
OpenAI-compatible HTTP contract, `workflow` owns the generate-to-edit
composition, `types` owns response models, and `artifacts` owns filesystem
output. The CLI is only an adapter. The same client exposes both text-to-image
and image-to-image operations, so adding another language or a new image
option does not couple transport code to the CLI.

## Requirements

- Rust 1.80+ (stable)
- An AI-ROUTER API key

```bash
export AI_ROUTER_API_KEY="your-key"
# Optional: defaults to https://api.ai-router.dev/v1
export AI_ROUTER_BASE_URL="https://api.ai-router.dev/v1"
```

## Run it

Generate an image from text:

```bash
cargo run --release -- generate \
  --model gpt-image-2 \
  --size 1024x1024 \
  --prompt "A cinematic glass greenhouse on Mars at sunrise" \
  --output out/generated.png
```

Then send that file to the edit endpoint:

```bash
cargo run --release -- edit \
  --image out/generated.png \
  --prompt "Add warm hanging lights and a small orange cat" \
  --output out/edited.png
```

For a single end-to-end example, `workflow` calls generation first and then
`/images/edits` with the generated artifact:

```bash
cargo run --release -- workflow \
  --prompt "A clean product photo of a red backpack" \
  --edit-prompt "Change the backpack to cobalt blue and keep the composition" \
  --generated out/generated.png \
  --edited out/edited.png
```

The client accepts `b64_json` responses and also downloads URL responses. API
errors include the HTTP status and provider message, while the API key is read
only from the environment and is never printed.

## Model choices

The default is `gpt-image-2`; try `gpt-image-2.5`,
`gpt-image-2.5-flare`, or `gpt-image-2.5-burst` when those models are enabled
for your account. Model availability, size limits, and pricing are controlled
by the account and provider, not by this demo.

## Learn more

AI-ROUTER provides one OpenAI-compatible integration surface for production
image workflows. See the [AI-ROUTER platform](https://ai-router.dev/) for API
access, model availability, and account documentation.
