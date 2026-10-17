from __future__ import annotations

import os
from dataclasses import dataclass

from dotenv import load_dotenv


@dataclass(frozen=True)
class Config:
    base_url: str = "https://api.ai-router.dev"
    mode: str = "probe"
    endpoint: str = "generations"
    model: str = "gpt-image-2"
    prompt: str = "A small red fox under a moonlit sky"
    n: int = 1
    image_path: str = ""
    private_key: str = ""
    rpc_url: str = "https://mainnet.base.org"
    timeout: float = 60.0

    @classmethod
    def from_env(cls) -> Config:
        load_dotenv()
        value = cls(
            base_url=os.getenv("AI_ROUTER_X402_BASE_URL", cls.base_url).rstrip("/"),
            mode=os.getenv("AI_ROUTER_X402_MODE", cls.mode).lower(),
            endpoint=os.getenv("AI_ROUTER_X402_ENDPOINT", cls.endpoint).lower(),
            model=os.getenv("AI_ROUTER_X402_MODEL", cls.model),
            prompt=os.getenv("AI_ROUTER_X402_PROMPT", cls.prompt),
            n=int(os.getenv("AI_ROUTER_X402_N", str(cls.n))),
            image_path=os.getenv("AI_ROUTER_X402_IMAGE_PATH", ""),
            private_key=os.getenv("AI_ROUTER_X402_EVM_PRIVATE_KEY", ""),
            rpc_url=os.getenv("AI_ROUTER_X402_EVM_RPC_URL", cls.rpc_url),
            timeout=float(os.getenv("AI_ROUTER_X402_TIMEOUT_SECONDS", str(cls.timeout))),
        )
        value.validate()
        return value

    def validate(self) -> None:
        if self.mode not in {"probe", "pay"}: raise ValueError("mode must be probe or pay")
        if self.endpoint not in {"generations", "edits"}: raise ValueError("endpoint must be generations or edits")
        if not self.model or not self.prompt or not 1 <= self.n <= 4: raise ValueError("model, prompt and n (1-4) are required")
        if self.timeout <= 0: raise ValueError("timeout must be positive")
        if self.mode == "pay" and not self.private_key.startswith("0x"): raise ValueError("pay mode requires AI_ROUTER_X402_EVM_PRIVATE_KEY")
        if self.endpoint == "edits" and not self.image_path: raise ValueError("edits requires AI_ROUTER_X402_IMAGE_PATH")
