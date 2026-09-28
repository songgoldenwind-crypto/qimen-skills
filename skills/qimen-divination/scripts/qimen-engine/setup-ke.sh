#!/usr/bin/env bash
set -euo pipefail

engine_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source_dir="$engine_dir/.ke-src"
revision="e6680ac4ca0b0da5ce3fe637e05f9fc32066ec5a"

if ! command -v git >/dev/null 2>&1; then
  echo 'error: ke engine setup needs git' >&2
  exit 2
fi

if [[ ! -x "$engine_dir/.venv/bin/python" ]]; then
  bash "$engine_dir/setup.sh"
fi

if [[ -d "$source_dir" ]]; then
  actual="$(git -C "$source_dir" rev-parse HEAD 2>/dev/null || true)"
  if [[ "$actual" != "$revision" ]]; then
    echo 'error: existing .ke-src is not the pinned KinQiMen revision; inspect it before removing or moving it' >&2
    exit 2
  fi
else
  git clone --quiet https://github.com/kentang2017/kinqimen.git "$source_dir"
  git -C "$source_dir" checkout --quiet "$revision"
fi

"$engine_dir/.venv/bin/python" "$engine_dir/paipan.py" --family 刻家 --datetime 2026-09-21T14:09:59 \
  | "$engine_dir/.venv/bin/python" -c 'import json,sys; p=json.load(sys.stdin); assert p["chart_granularity_minutes"] == 10; assert p["raw"]["干支"].endswith("庚午分")'
"$engine_dir/.venv/bin/python" "$engine_dir/paipan.py" --family 刻家 --datetime 2026-09-21T14:10:00 \
  | "$engine_dir/.venv/bin/python" -c 'import json,sys; p=json.load(sys.stdin); assert p["raw"]["干支"].endswith("辛未分")'
echo 'ke qimen engine ready: pinned 10-minute boundary samples passed'
