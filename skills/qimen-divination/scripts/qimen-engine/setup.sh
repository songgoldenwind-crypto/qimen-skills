#!/usr/bin/env bash
set -euo pipefail

engine_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
env_dir="$engine_dir/.venv"
if ! command -v uv >/dev/null 2>&1; then
  echo 'error: setup needs uv (https://docs.astral.sh/uv/)' >&2
  exit 2
fi
if ! command -v git >/dev/null 2>&1; then
  echo 'error: setup needs git' >&2
  exit 2
fi

source_dir="$(mktemp -d "${TMPDIR:-/tmp}/sxtwl-qimen.XXXXXXXX")"
trap 'rm -rf "$source_dir"' EXIT

uv venv --python 3.12 "$env_dir"
git clone --quiet https://github.com/yuangu/sxtwl_cpp.git "$source_dir/upstream"
git -C "$source_dir/upstream" checkout --quiet 7598b0601a76cfdaa9266257b1b5690720c1e2ce
uv pip install --python "$env_dir/bin/python" "$source_dir/upstream/python" bidict==0.23.1 ephem==4.1.6
uv pip install --python "$env_dir/bin/python" --no-deps kinqimen==0.0.6.6

"$env_dir/bin/python" "$engine_dir/paipan.py" --purpose chart-only --engine kinqimen --family 时家 --datetime 2020-10-07T14:23 --timezone Asia/Shanghai \
  | "$env_dir/bin/python" -c 'import json,sys; p=json.load(sys.stdin); assert p["raw"]["排局"] == "陰遁六局上元"; assert p["palaces"]["3"]["stem_pair"] == "壬+辛"'
echo 'qimen engine ready: pinned version and reference sample passed'
