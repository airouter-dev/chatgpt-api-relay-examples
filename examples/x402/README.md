# AI ROUTER x402 client examples

This directory is a modular reference client for AI ROUTER's x402 image API.
It demonstrates the same payment flow in Go, Python, Node.js, and curl:

1. send an unpaid request;
2. receive `402 Payment Required` and inspect `PAYMENT-REQUIRED`;
3. let the official x402 SDK sign an exact Base USDC payment;
4. retry the original request with `PAYMENT-SIGNATURE`;
5. inspect `PAYMENT-RESPONSE` and save the image response.

The examples default to **probe mode**. Probe mode never signs a transaction and
is safe for CI or documentation screenshots. Paid mode is opt-in and requires a
funded EVM wallet on Base mainnet.

## Production routes currently configured

The production catalog was verified with unpaid probes on 2026-10-17/18. Only
the following combinations are documented as working:

| Endpoint | Models | Billing |
| --- | --- | --- |
| `POST /x402/v1/images/generations` | `gpt-image-2`, `gpt-image-2.5`, `gpt-image-2.5-sunburst`, `gpt-image-2.5-flare` | `0.0045 USDC × n` |
| `POST /x402/v1/images/edits` | the same four models | `0.0045 USDC × n` |

The settlement requirement is Base mainnet (`eip155:8453`) USDC with six
decimals. `n` must be sent explicitly and is currently limited to 1–4.

`/x402/v1/images/variations` may appear in generic discovery metadata, but it
is not a working production route and is intentionally not included here.

## Configuration

Copy `.env.example` to `.env` or copy `config.example.toml` to `config.toml`.
Both files are ignored by Git. Never commit a real private key, API key, or
wallet seed.

```dotenv
AI_ROUTER_X402_BASE_URL=https://api.ai-router.dev
AI_ROUTER_X402_MODE=probe
AI_ROUTER_X402_MODEL=gpt-image-2
AI_ROUTER_X402_N=1
AI_ROUTER_X402_PROMPT=A small red fox under a moonlit sky
# Paid mode only:
# AI_ROUTER_X402_EVM_PRIVATE_KEY=0x...
# AI_ROUTER_X402_EVM_RPC_URL=https://mainnet.base.org
```

Use `probe` first. Set `pay` only after verifying the decoded amount, network,
asset, and recipient. The private key is the **payer's** wallet key; it is not
the merchant API key and is never sent to AI ROUTER.

## Run each client

### Go

```bash
cd examples/x402/go
go mod download
go run ./cmd/ai-router-x402 --endpoint generations --mode probe
# Paid mode:
go run ./cmd/ai-router-x402 --endpoint generations --mode pay
```

### Python

```bash
cd examples/x402/python
python -m venv .venv
. .venv/bin/activate
pip install -e '.[dev]'
python -m ai_router_x402 --endpoint generations --mode probe
python -m ai_router_x402 --endpoint edits --mode probe --image ./input.png
```

For paid mode set `AI_ROUTER_X402_EVM_PRIVATE_KEY` in `.env` and run with
`--mode pay`.

### Node.js

```bash
cd examples/x402/nodejs
npm install
node src/cli.mjs --endpoint generations --mode probe
node src/cli.mjs --endpoint edits --mode probe --image ./input.png
```

The Node client uses the official `@x402/fetch` and `@x402/evm` packages for
paid mode. Go and Python likewise use official x402 Foundation SDKs; no payment
headers are hand-crafted in the examples.

### curl probe (never pays)

```bash
curl -i https://api.ai-router.dev/x402/v1/images/generations \
  -H 'content-type: application/json' \
  --data '{"model":"gpt-image-2","prompt":"a red fox","n":1}'
```

The response should be `402` and contain `PAYMENT-REQUIRED`. A paid curl call
is intentionally not provided because signing exact EVM authorizations in a
shell script is unsafe.

## Tests

Unit tests validate configuration, endpoint construction, 402 decoding, and
secret redaction without spending funds. The live probe script validates the
four configured models and both working routes but never signs a payment:

```bash
python scripts/probe_matrix.py
```

The live probe is not run automatically in GitHub Actions. A successful probe
is not proof that paid settlement will succeed: the payer still needs Base
USDC, gas, a valid private key, and permission to sign the advertised payment.

## Language guides

- [English](docs/en.md) · [中文](docs/zh-CN.md) · [Русский](docs/ru.md)
- [العربية](docs/ar.md) · [Français](docs/fr.md) · [Español](docs/es.md)

Landing pages:

- [English](https://ai-router.dev) · [中文](https://ai-router.dev/cn)
- [Русский](https://ai-router.dev/ru) · [العربية](https://ai-router.dev/ar)
- [Français](https://ai-router.dev/fr) · [Español](https://ai-router.dev/es)

## Security model

- `.env`, `config.toml`, wallet files, generated images, and local logs are
  ignored by Git.
- Paid mode requires an explicit `--mode pay`.
- The client validates the challenge before signing: exact scheme, Base mainnet,
  Base USDC, six decimals, expected recipient, and a positive amount.
- The client never logs private keys, API keys, or full payment headers.
- Do not use a production treasury key as a payer key. Use a dedicated wallet
  with a narrowly bounded balance.

