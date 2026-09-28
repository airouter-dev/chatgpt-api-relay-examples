"""Minimal OpenAI-compatible Images API client using only Python's stdlib.

The client deliberately contains transport and response decoding only. Prompt
orchestration belongs in :mod:`image_demo.workflow`, making it easy to replace
this client with an SDK or a mocked transport in a production application.
"""

from __future__ import annotations

import base64
import binascii
import json
import mimetypes
import uuid
from dataclasses import dataclass
from pathlib import Path
from typing import Any
from urllib.error import HTTPError, URLError
from urllib.parse import urlparse
from urllib.request import Request, urlopen

from .config import ImageConfig
from .errors import ImageApiError, ImageDemoError


@dataclass(frozen=True)
class ImageAsset:
    """Decoded image bytes plus metadata returned by an Images endpoint."""

    content: bytes
    media_type: str = "image/png"
    source_url: str | None = None
    revised_prompt: str | None = None

    @property
    def extension(self) -> str:
        subtype = self.media_type.split("/", 1)[-1].lower()
        return {"jpeg": ".jpg", "jpg": ".jpg", "webp": ".webp"}.get(subtype, ".png")


def _json_error_body(raw: bytes) -> str:
    text = raw.decode("utf-8", errors="replace")
    try:
        parsed = json.loads(text)
    except json.JSONDecodeError:
        return text[:1000]
    error = parsed.get("error") if isinstance(parsed, dict) else None
    if isinstance(error, dict) and error.get("message"):
        return str(error["message"])
    return text[:1000]


def decode_image_response(payload: dict[str, Any], *, fetch_url=urlopen) -> ImageAsset:
    """Decode the first image from ``b64_json`` or ``url`` response formats."""

    data = payload.get("data")
    if not isinstance(data, list) or not data or not isinstance(data[0], dict):
        raise ImageApiError("Images API returned no image data")
    item = data[0]
    revised_prompt = item.get("revised_prompt") if isinstance(item.get("revised_prompt"), str) else None
    encoded = item.get("b64_json")
    if isinstance(encoded, str) and encoded:
        try:
            content = base64.b64decode(encoded, validate=True)
        except (binascii.Error, ValueError) as exc:
            raise ImageApiError("Images API returned invalid base64 image data") from exc
        return ImageAsset(content=content, revised_prompt=revised_prompt)

    image_url = item.get("url")
    if isinstance(image_url, str) and image_url:
        if image_url.startswith("data:"):
            header, encoded_data = image_url.split(",", 1)
            media_type = header[5:].split(";", 1)[0] or "image/png"
            try:
                content = base64.b64decode(encoded_data, validate=True)
            except (binascii.Error, ValueError) as exc:
                raise ImageApiError("Images API returned invalid data URL") from exc
            return ImageAsset(content=content, media_type=media_type, source_url=image_url, revised_prompt=revised_prompt)
        request = Request(image_url, headers={"User-Agent": "ai-router-image-demo/0.1"})
        try:
            with fetch_url(request, timeout=120) as response:
                content = response.read()
                media_type = response.headers.get_content_type() or "image/png"
        except (HTTPError, URLError, TimeoutError) as exc:
            raise ImageApiError(f"Could not download image URL: {exc}") from exc
        return ImageAsset(content=content, media_type=media_type, source_url=image_url, revised_prompt=revised_prompt)
    raise ImageApiError("Images API item contained neither b64_json nor url")


class ImagesClient:
    """Transport adapter for ``/v1/images/generations`` and ``/v1/images/edits``."""

    def __init__(self, config: ImageConfig) -> None:
        self.config = config

    def _request(self, request: Request) -> dict[str, Any]:
        try:
            with urlopen(request, timeout=self.config.timeout_seconds) as response:
                raw = response.read()
                status = getattr(response, "status", 200)
        except HTTPError as exc:
            body = _json_error_body(exc.read())
            raise ImageApiError("Images API request failed", status=exc.code, body=body) from exc
        except (URLError, TimeoutError) as exc:
            raise ImageApiError(f"Could not reach image API: {exc}") from exc
        try:
            payload = json.loads(raw.decode("utf-8"))
        except (UnicodeDecodeError, json.JSONDecodeError) as exc:
            raise ImageApiError(f"Images API returned invalid JSON (HTTP {status})") from exc
        if not isinstance(payload, dict):
            raise ImageApiError("Images API returned a non-object JSON response")
        return payload

    def generate(self, prompt: str, *, model: str | None = None, size: str | None = None, quality: str | None = None) -> ImageAsset:
        """Create an image from text via ``POST /images/generations``."""

        if not prompt.strip():
            raise ValueError("prompt cannot be empty")
        payload = {
            "model": model or self.config.model,
            "prompt": prompt,
            "size": size or self.config.size,
            "quality": quality or self.config.quality,
            "response_format": "b64_json",
        }
        body = json.dumps(payload).encode("utf-8")
        request = Request(
            f"{self.config.base_url}/images/generations",
            data=body,
            method="POST",
            headers={
                "Authorization": f"Bearer {self.config.api_key}",
                "Content-Type": "application/json",
                "Accept": "application/json",
                "User-Agent": "ai-router-image-demo/0.1",
            },
        )
        return decode_image_response(self._request(request))

    def edit(self, image_path: Path, prompt: str, *, model: str | None = None, size: str | None = None, quality: str | None = None) -> ImageAsset:
        """Edit a generated image via multipart ``POST /images/edits``."""

        if not prompt.strip():
            raise ValueError("edit prompt cannot be empty")
        path = image_path.expanduser()
        if not path.is_file():
            raise FileNotFoundError(f"source image does not exist: {path}")
        mime_type = mimetypes.guess_type(path.name)[0] or "application/octet-stream"
        boundary = f"----ai-router-image-demo-{uuid.uuid4().hex}"
        parts: list[bytes] = []

        def field(name: str, value: str) -> None:
            parts.extend([
                f"--{boundary}\r\n".encode(),
                f'Content-Disposition: form-data; name="{name}"\r\n\r\n'.encode(),
                value.encode("utf-8"),
                b"\r\n",
            ])

        def file_field(name: str, filename: str, content_type: str, content: bytes) -> None:
            parts.extend([
                f"--{boundary}\r\n".encode(),
                f'Content-Disposition: form-data; name="{name}"; filename="{filename}"\r\n'.encode(),
                f"Content-Type: {content_type}\r\n\r\n".encode(),
                content,
                b"\r\n",
            ])

        field("model", model or self.config.model)
        field("prompt", prompt)
        field("size", size or self.config.size)
        field("quality", quality or self.config.quality)
        field("response_format", "b64_json")
        file_field("image", path.name, mime_type, path.read_bytes())
        parts.append(f"--{boundary}--\r\n".encode())
        request = Request(
            f"{self.config.base_url}/images/edits",
            data=b"".join(parts),
            method="POST",
            headers={
                "Authorization": f"Bearer {self.config.api_key}",
                "Content-Type": f"multipart/form-data; boundary={boundary}",
                "Accept": "application/json",
                "User-Agent": "ai-router-image-demo/0.1",
            },
        )
        return decode_image_response(self._request(request))


def save_asset(asset: ImageAsset, destination: Path) -> Path:
    """Persist bytes atomically enough for a small demo and return the path."""

    target = destination.with_suffix(destination.suffix or asset.extension)
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_bytes(asset.content)
    return target
