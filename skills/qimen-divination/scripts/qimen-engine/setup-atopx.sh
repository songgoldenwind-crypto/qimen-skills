#!/usr/bin/env bash
set -euo pipefail

engine_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source_dir="$engine_dir/atopx-src"
bin_dir="$engine_dir/.bin"

if ! command -v go >/dev/null 2>&1; then
  echo 'error: atopx engine needs Go; see scripts/qimen-engine/README.md' >&2
  exit 2
fi

mkdir -p "$bin_dir"
go -C "$source_dir" test ./...
go -C "$source_dir" build -o "$bin_dir/atopx-qimen" ./cmd/qimen-json

python3 "$engine_dir/paipan.py" --engine atopx --family 时家 --datetime 2020-04-18T14:00 --timezone Asia/Shanghai \
  | python3 -c 'import json,sys; p=json.load(sys.stdin); assert p["raw"]["局数"] == 7; assert p["palaces"]["7"]["stem_pair"] == "乙+戊"'
echo 'atopx qimen engine ready: pinned source compiled and reference sample passed'
