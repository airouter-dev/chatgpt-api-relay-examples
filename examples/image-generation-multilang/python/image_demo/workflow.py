"""Business workflow: text-to-image first, then image-to-image editing."""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path

from .client import ImageAsset, ImagesClient, save_asset


@dataclass(frozen=True)
class WorkflowResult:
    generated_path: Path
    edited_path: Path
    generated: ImageAsset
    edited: ImageAsset


class ImageWorkflow:
    """Orchestrate a two-step image pipeline without coupling it to the CLI."""

    def __init__(self, client: ImagesClient) -> None:
        self.client = client

    def run(self, prompt: str, edit_prompt: str, output_dir: Path) -> WorkflowResult:
        output_dir.mkdir(parents=True, exist_ok=True)
        generated = self.client.generate(prompt)
        generated_path = save_asset(generated, output_dir / f"generated{generated.extension}")
        edited = self.client.edit(generated_path, edit_prompt)
        edited_path = save_asset(edited, output_dir / f"edited{edited.extension}")
        return WorkflowResult(generated_path, edited_path, generated, edited)
