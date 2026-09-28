"""Configuration with one source of truth for CLI flags and environment variables."""

from __future__ import annotations

import os
from dataclasses import dataclass, replace
from urllib.parse import urlparse


def normalize_base_url(value: str) -> str:
    """Return a base URL that ends in ``/v1`` without a trailing slash."""

    candidate = value.strip().rstrip("/")
    if not candidate:
        raise ValueError("IMAGE_API_BASE_URL cannot be empty")
    parsed = urlparse(candidate)
    if parsed.scheme not in {"http", "https"} or not parsed.netloc:
        raise ValueError("IMAGE_API_BASE_URL must be an http(s) URL")
    return candidate if parsed.path.rstrip("/").endswith("/v1") else f"{candidate}/v1"


@dataclass(frozen=True)
class ImageConfig:
    """Runtime settings shared by generation and edit requests."""

    api_key: str
    base_url: str = "https://api.ai-router.dev/v1"
    model: str = "gpt-image-2"
    size: str = "1024x1024"
    quality: str = "auto"
    timeout_seconds: float = 120.0

    @classmethod
    def from_env(cls) -> "ImageConfig":
        api_key = os.getenv("IMAGE_API_KEY") or os.getenv("OPENAI_API_KEY") or os.getenv("AI_ROUTER_API_KEY")
        if not api_key:
            raise ValueError("Set AI_ROUTER_API_KEY (or IMAGE_API_KEY/OPENAI_API_KEY) before running the demo")
        try:
            timeout = float(os.getenv("IMAGE_TIMEOUT_SECONDS", "120"))
        except ValueError as exc:
            raise ValueError("IMAGE_TIMEOUT_SECONDS must be a number") from exc
        if timeout <= 0:
            raise ValueError("IMAGE_TIMEOUT_SECONDS must be greater than zero")
        return cls(
            api_key=api_key,
            base_url=normalize_base_url(
                os.getenv("IMAGE_API_BASE_URL")
                or os.getenv("AI_ROUTER_BASE_URL")
                or "https://api.ai-router.dev/v1"
            ),
            model=os.getenv("IMAGE_MODEL") or os.getenv("AI_ROUTER_MODEL", "gpt-image-2"),
            size=os.getenv("IMAGE_SIZE") or os.getenv("AI_ROUTER_SIZE", "1024x1024"),
            quality=os.getenv("IMAGE_QUALITY", "auto"),
            timeout_seconds=timeout,
        )

    def with_overrides(self, **values: object) -> "ImageConfig":
        """Apply non-None CLI values while preserving validated defaults."""

        updates = {key: value for key, value in values.items() if value is not None}
        if "base_url" in updates:
            updates["base_url"] = normalize_base_url(str(updates["base_url"]))
        if "timeout_seconds" in updates and float(updates["timeout_seconds"]) <= 0:
            raise ValueError("timeout must be greater than zero")
        return replace(self, **updates)
