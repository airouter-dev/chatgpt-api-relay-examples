"""Command-line entry point; all provider-specific work stays in client/workflow."""

from __future__ import annotations

import argparse
import sys
from pathlib import Path

from .client import ImagesClient
from .config import ImageConfig
from .errors import ImageDemoError
from .workflow import ImageWorkflow


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        description="Generate an image, then edit that image through OpenAI-compatible endpoints."
    )
    parser.add_argument("--prompt", required=True, help="Text-to-image prompt")
    parser.add_argument("--edit-prompt", required=True, help="Instruction for the image-to-image edit")
    parser.add_argument("--output-dir", type=Path, default=Path("outputs"))
    parser.add_argument("--model", help="Model id (default: IMAGE_MODEL or gpt-image-2)")
    parser.add_argument("--size", help="Image size, for example 1024x1024")
    parser.add_argument("--quality", help="Image quality, for example auto or high")
    parser.add_argument("--base-url", help="OpenAI-compatible API base URL")
    parser.add_argument("--api-key", help="API key; prefer AI_ROUTER_API_KEY in the environment")
    parser.add_argument("--timeout", type=float, help="HTTP timeout in seconds")
    return parser


def main(argv: list[str] | None = None) -> int:
    args = build_parser().parse_args(argv)
    try:
        config = ImageConfig.from_env()
        config = config.with_overrides(
            api_key=args.api_key,
            base_url=args.base_url,
            model=args.model,
            size=args.size,
            quality=args.quality,
            timeout_seconds=args.timeout,
        )
        result = ImageWorkflow(ImagesClient(config)).run(args.prompt, args.edit_prompt, args.output_dir)
    except (ImageDemoError, OSError, ValueError) as exc:
        print(f"image demo failed: {exc}", file=sys.stderr)
        return 1
    print(f"Generated image: {result.generated_path}")
    print(f"Edited image:    {result.edited_path}")
    if result.generated.revised_prompt:
        print(f"Provider revised prompt: {result.generated.revised_prompt}")
    return 0
