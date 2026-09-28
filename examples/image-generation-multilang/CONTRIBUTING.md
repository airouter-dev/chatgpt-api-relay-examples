# Contributing a language implementation

Keep every implementation aligned with [`docs/API_CONTRACT.md`](docs/API_CONTRACT.md): one JSON generation request, one persisted artifact, then one multipart edit request using that artifact. A new language should expose the same three operations where the runtime permits them:

- `generate(prompt)` returns the first image data item;
- `edit(image_path, prompt)` sends the generated file to `/images/edits`;
- `workflow(prompt, edit_prompt)` composes the two operations and preserves the first artifact if editing fails.

Use four small layers—configuration, API client, workflow, and CLI/application adapter. Keep API keys out of source and logs, accept both `b64_json` and `url` responses, bound response/download sizes, and avoid forwarding the bearer token to a returned asset URL on another host.

Before opening a pull request, update the root language table and add a language README containing requirements, an AI-ROUTER link, installation/run commands, model selection notes, and the relevant test command. Tests must use a local HTTP fixture or mocks; they must not make paid requests or require a real API key.
