#!/usr/bin/env python3
"""Run a pinned upstream KinQiMen ke chart in an isolated interpreter."""

import json
import subprocess
import sys
from datetime import datetime
from pathlib import Path


REVISION = "e6680ac4ca0b0da5ce3fe637e05f9fc32066ec5a"
SOURCE = Path(__file__).resolve().parent / ".ke-src"


def main():
    if len(sys.argv) != 3:
        raise ValueError("expected ISO local datetime and upstream method number")
    commit = subprocess.run(
        ["git", "-C", str(SOURCE), "rev-parse", "HEAD"],
        capture_output=True, text=True, check=False,
    )
    if commit.returncode or commit.stdout.strip() != REVISION:
        raise RuntimeError("刻家源码提交不匹配；运行 setup-ke.sh 重新安装固定版本")
    local = datetime.fromisoformat(sys.argv[1])
    option = int(sys.argv[2])
    if option != 2:
        raise ValueError("此入口仅使用上游 pan_minute(2)")
    sys.path.insert(0, str(SOURCE))
    from kinqimen import Qimen  # noqa: E402 — upstream source, not the PyPI wheel

    raw = Qimen(local.year, local.month, local.day, local.hour, local.minute).pan_minute(option)
    print(json.dumps(raw, ensure_ascii=False))


if __name__ == "__main__":
    try:
        main()
    except (ValueError, RuntimeError, ImportError, OSError) as exc:
        sys.exit(f"error: {exc}")
