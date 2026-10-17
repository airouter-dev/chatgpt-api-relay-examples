# x402 paid image calls

AI ROUTER's image routes return a machine-readable HTTP 402 challenge before
any money moves. The x402 SDK verifies the challenge, signs an exact EVM
authorization, and retries the same request. You do not put a merchant secret
or PayAI secret in the client.

## Verified production contract

Use Base mainnet, USDC, and one of the four configured models. `n` is required
and must be between 1 and 4. The amount is `4500 * n` atomic USDC units.

Generations use JSON. Edits use multipart form data with an `image` file. The
Go, Python, and Node clients all read `.env` or TOML configuration and keep
paid mode opt-in.

## Private-key safety

Use a dedicated payer wallet with a small balance. Store
`AI_ROUTER_X402_EVM_PRIVATE_KEY` only in a local ignored `.env` file or secret
manager. Never paste it into source, shell history, an issue, or a log. Start
with probe mode, check the challenge, then explicitly choose pay mode.

See the [English landing page](https://ai-router.dev) for account and API-key
setup. See the [Chinese landing page](https://ai-router.dev/cn) and the
[Russian landing page](https://ai-router.dev/ru) for localized product context.
