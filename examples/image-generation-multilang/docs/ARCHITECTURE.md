# Architecture notes

Each language follows the same four-layer shape:

1. **Configuration** reads environment variables and command-line overrides. It owns defaults and validates required secrets.
2. **API client** owns HTTP transport, authentication, JSON/multipart encoding, status handling, and response decoding. It knows nothing about filenames or CLI output.
3. **Workflow** composes two domain operations: generate from a prompt, then edit the generated file. It owns output naming and the hand-off between operations.
4. **CLI entrypoint** parses user input, invokes the workflow, and renders concise progress messages. It contains no HTTP details.

This separation makes it straightforward to embed the client in a web server, queue worker, notebook, or desktop app later. It also gives each language a natural place for retries, tracing, rate-limit handling, and alternate storage without duplicating protocol code.

## Why the demos save locally

The first response is persisted before the edit call. That makes the two-stage workflow inspectable, supports a human review step, and gives users a concrete input asset to reuse. In an application, replace the local `ArtifactStore`/output directory with object storage while keeping the client and workflow unchanged.

## Failure boundaries

- Configuration failures happen before a request is made.
- HTTP failures include status code and a bounded response snippet, but never print the API key.
- Decode failures identify the missing response field or unsupported URL form.
- A failed edit never deletes a successfully generated source image.

## Extending the project

To add another language, implement the contract in `docs/API_CONTRACT.md`, keep the same `generate → edit` workflow, add a language-specific README, and register the command in the root README's comparison table. Do not copy transport code into the CLI; keep it behind a small client interface.
