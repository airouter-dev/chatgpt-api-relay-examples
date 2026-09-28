"""Errors that keep provider diagnostics useful without exposing credentials."""


class ImageDemoError(RuntimeError):
    """Base class for expected demo failures."""


class ImageApiError(ImageDemoError):
    """An HTTP or provider response error."""

    def __init__(self, message: str, *, status: int | None = None, body: str = "") -> None:
        self.status = status
        self.body = body
        suffix = f" (HTTP {status})" if status else ""
        super().__init__(f"{message}{suffix}")
