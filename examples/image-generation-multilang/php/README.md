# AI-ROUTER image generation — PHP

This PHP sample is a reusable mini-client, not a hard-coded HTTP snippet:

- `Config` reads the endpoint, API key, and timeout from the environment.
- `ImageClient` owns the OpenAI-compatible JSON and multipart contracts.
- `ImageData` normalises `b64_json` and URL responses.
- `ArtifactStore` owns filesystem output only.
- `Workflow` composes generation and `/images/edits` without mixing CLI parsing
  into transport code.

That separation keeps each module cohesive and makes it straightforward to add
streaming, moderation, retries, or another language implementation later.

## Requirements

- PHP 8.1+
- `curl` and `json` extensions
- An AI-ROUTER API key

```bash
export AI_ROUTER_API_KEY="your-key"
# Optional; defaults to https://api.ai-router.dev/v1
export AI_ROUTER_BASE_URL="https://api.ai-router.dev/v1"
```

No runtime package is required: `src/autoload.php` is a tiny PSR-4 loader for
the demo. For an application, run `composer dump-autoload` and use the PSR-4
mapping in `composer.json`.

## Run it

Text-to-image:

```bash
php bin/image-demo.php generate \
  --model gpt-image-2 \
  --size 1024x1024 \
  --prompt "A cinematic glass greenhouse on Mars at sunrise" \
  --output out/generated.png
```

Image-to-image through the edit endpoint:

```bash
php bin/image-demo.php edit \
  --image out/generated.png \
  --prompt "Add warm hanging lights and a small orange cat" \
  --output out/edited.png
```

Or run the full generate → edit workflow:

```bash
php bin/image-demo.php workflow \
  --prompt "A clean product photo of a red backpack" \
  --edit-prompt "Change the backpack to cobalt blue and keep the composition" \
  --generated out/generated.png \
  --edited out/edited.png
```

The client accepts either base64 or URL image responses, reports provider
errors with their HTTP status, and never prints the API key. Try
`gpt-image-2.5`, `gpt-image-2.5-flare`, or `gpt-image-2.5-burst` when your
account has access; model availability and pricing are account-dependent.

## Learn more

Use the [AI-ROUTER platform](https://ai-router.dev/) for API access, model
availability, and production integration guidance.
