# API contract used by every demo

The examples intentionally target the OpenAI-compatible Images API exposed by [AI-ROUTER](https://ai-router.dev). This makes the same workflow portable across languages and lets a developer keep an existing OpenAI-style integration while changing only the base URL and API key.

## Text-to-image

```http
POST {base_url}/images/generations
Authorization: Bearer $AI_ROUTER_API_KEY
Content-Type: application/json
```

```json
{
  "model": "gpt-image-2",
  "prompt": "A small glass greenhouse on a rainy rooftop, editorial photography",
  "size": "1024x1024",
  "n": 1,
  "response_format": "b64_json"
}
```

The successful response is expected to contain `data`, with each item containing `b64_json` or (depending on account/model configuration) a downloadable `url`. The demos prefer `b64_json` so the result can be saved without an additional unauthenticated download step.

## Image-to-image

```http
POST {base_url}/images/edits
Authorization: Bearer $AI_ROUTER_API_KEY
Content-Type: multipart/form-data
```

Multipart fields:

| Field | Required | Description |
| --- | --- | --- |
| `model` | yes | Image model exposed to the account, for example `gpt-image-2` or `gpt-image-2.5-flare`. |
| `prompt` | yes | The transformation instruction. |
| `size` | no | A supported output size such as `1024x1024`, `1536x1024`, or `2048x2048`; availability is account-dependent. |
| `response_format` | no | `b64_json` is the portable choice used by these demos. |
| `image` | yes | The image written by the text-to-image step. |

The edit response uses the same `data` shape as generations. A demo saves the first result as `<language>-edited.<extension>`.

## Compatibility and model selection

The model is deliberately configurable. Use a model that appears in the authenticated account's `/models` response and that is enabled for image generation. Some accounts may expose `gpt-image-2`, `gpt-image-2.5`, `gpt-image-2.5-flare`, or `gpt-image-2.5-burst`; the examples do not hard-code access to any one catalog. `AI_ROUTER_MODEL` and the `--model` option are passed through unchanged.

AI-ROUTER is an independent OpenAI-compatible API service, not an OpenAI product. Availability, quotas, pricing, response formats, and optional parameters depend on the account and current service configuration. Review the [AI-ROUTER documentation](https://ai-router.dev) before using the examples in production.
