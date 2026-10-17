"""No-payment production smoke matrix for the configured x402 image routes."""
from __future__ import annotations

import base64
import json
import os
import sys
import urllib.request

BASE = os.getenv("AI_ROUTER_X402_BASE_URL", "https://api.ai-router.dev").rstrip("/")
MODELS = ["gpt-image-2", "gpt-image-2.5", "gpt-image-2.5-sunburst", "gpt-image-2.5-flare"]
ROUTES = ["generations", "edits"]


def decode(value: str) -> dict:
    return json.loads(base64.urlsafe_b64decode(value + "=" * (-len(value) % 4)))


def probe(route: str, model: str, n: int) -> dict:
    url = f"{BASE}/x402/v1/images/{route}"
    body = json.dumps({"model": model, "prompt": "x402 probe; do not settle", "n": n}).encode()
    request = urllib.request.Request(url, data=body, headers={"content-type": "application/json"})
    try:
        urllib.request.urlopen(request, timeout=30)
    except urllib.error.HTTPError as error:
        if error.code != 402:
            raise RuntimeError(f"{route}/{model}/n={n}: HTTP {error.code}") from error
        payload = decode(error.headers.get("PAYMENT-REQUIRED", ""))
        accepts = payload.get("accepts", [])
        if len(accepts) != 1:
            raise RuntimeError(f"{route}/{model}/n={n}: expected one payment option")
        amount = int(accepts[0]["amount"])
        expected = 4500 * n
        if amount != expected:
            raise RuntimeError(f"{route}/{model}/n={n}: amount {amount}, expected {expected}")
        return {"route": route, "model": model, "n": n, "amount": amount, "network": accepts[0]["network"]}
    raise RuntimeError(f"{route}/{model}/n={n}: expected HTTP 402")


def main() -> int:
    results = [probe(route, model, n) for route in ROUTES for model in MODELS for n in (1, 2)]
    print(json.dumps({"base_url": BASE, "probes": results}, indent=2))
    return 0


if __name__ == "__main__":
    sys.exit(main())
