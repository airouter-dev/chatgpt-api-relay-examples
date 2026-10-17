from __future__ import annotations

import base64
import json
from pathlib import Path
from typing import Any

import requests
from eth_account import Account
from x402 import x402ClientSync
from x402.http.clients import wrapRequestsWithPayment
from x402.mechanisms.evm import EthAccountSignerWithRPC
from x402.mechanisms.evm.exact import register_exact_evm_client

from .config import Config

NETWORK = "eip155:8453"
ASSET = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913".lower()
PAY_TO = "0xEe027fE58800Bcef81469EDF0C2E7f2ab54965De".lower()


def _field(value: Any, name: str, fallback: str) -> Any:
    if isinstance(value, dict):
        return value.get(name, value.get(fallback))
    return getattr(value, name, getattr(value, fallback, None))


def select_requirements(_version: int, requirements: list[Any]) -> Any:
    candidates = [r for r in requirements if _field(r, "scheme", "scheme") == "exact" and _field(r, "network", "network") == NETWORK and str(_field(r, "asset", "asset")).lower() == ASSET and str(_field(r, "pay_to", "payTo")).lower() == PAY_TO]
    if len(candidates) != 1: raise ValueError("402 challenge did not contain the expected exact Base USDC recipient")
    amount = int(_field(candidates[0], "amount", "amount"))
    if not 1 <= amount <= 18000: raise ValueError(f"402 amount {amount} is outside the configured 1-4 image range")
    return candidates[0]


def request_parts(cfg: Config) -> tuple[str, dict[str, Any]]:
    url = f"{cfg.base_url}/x402/v1/images/{cfg.endpoint}"
    if cfg.endpoint == "generations":
        return url, {"json": {"model": cfg.model, "prompt": cfg.prompt, "n": cfg.n}}
    return url, {"data": {"model": cfg.model, "prompt": cfg.prompt, "n": str(cfg.n)}, "files": {"image": (Path(cfg.image_path).name, open(cfg.image_path, "rb"), "image/png")}}


def run(cfg: Config) -> dict[str, Any]:
    url, kwargs = request_parts(cfg)
    try:
        if cfg.mode == "probe":
            response = requests.post(url, timeout=cfg.timeout, **kwargs)
        else:
            account = Account.from_key(cfg.private_key)
            signer = EthAccountSignerWithRPC(account, rpc_url=cfg.rpc_url)
            client = x402ClientSync(select_requirements)
            register_exact_evm_client(client, signer, networks=NETWORK)
            response = wrapRequestsWithPayment(requests.Session(), client).post(url, timeout=cfg.timeout, **kwargs)
    finally:
        files = kwargs.get("files")
        if files:
            files["image"][1].close()
    if response.status_code == 402:
        value = response.headers.get("PAYMENT-REQUIRED")
        if not value: raise ValueError("402 response has no PAYMENT-REQUIRED header")
        payload = json.loads(base64.urlsafe_b64decode(value + "=" * (-len(value) % 4)))
        chosen = select_requirements(payload["x402Version"], payload["accepts"])
        result = {"mode": cfg.mode, "status": 402, "endpoint": cfg.endpoint, "model": cfg.model, "amount": _field(chosen, "amount", "amount"), "network": _field(chosen, "network", "network"), "asset": _field(chosen, "asset", "asset"), "payTo": _field(chosen, "pay_to", "payTo")}
        print(json.dumps(result)); return result
    if not response.ok: raise RuntimeError(f"HTTP {response.status_code}: {response.text[:500]}")
    result = response.json(); print(json.dumps({"mode": cfg.mode, "status": response.status_code, "data": result})); return result
