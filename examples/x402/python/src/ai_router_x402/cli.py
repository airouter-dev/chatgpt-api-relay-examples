from __future__ import annotations

import argparse

from .client import run
from .config import Config


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--mode", choices=["probe", "pay"])
    parser.add_argument("--endpoint", choices=["generations", "edits"])
    parser.add_argument("--image")
    args = parser.parse_args()
    cfg = Config.from_env()
    values = {"mode": args.mode or cfg.mode, "endpoint": args.endpoint or cfg.endpoint, "image_path": args.image or cfg.image_path}
    cfg = Config(**{**cfg.__dict__, **values})
    cfg.validate()
    run(cfg)


if __name__ == "__main__":
    main()
