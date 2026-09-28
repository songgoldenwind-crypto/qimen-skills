#!/usr/bin/env python3
"""Validate a standalone, method-only Skill and its optional release archive."""

from __future__ import annotations

import re
import sys
import zipfile
from pathlib import Path
from urllib.parse import unquote

PROJECT = Path(__file__).resolve().parents[1]
SKILL = PROJECT / "skills" / "qimen-divination"
EXCLUDED = {".git", ".venv", ".ke-src", ".bin", "__pycache__", ".pytest_cache", "dist"}
SOURCE_MARKERS = re.compile(
    r"《[^》\n]+》|\bQ[0-9]{2}\b|\bT0[12]\b(?!:)|书图|物理页|ISBN|出版社|"
    r"(?:作者|编者|著者)\s*[:：]|(?:PDF|OCR)[^\n]{0,24}(?:页|报告)|"
    r"第\s*[0-9一二三四五六七八九十百]+(?:\s*[–－~～-]\s*[0-9一二三四五六七八九十百]+)?\s*(?:页|卷|章|堂)|"
    r"qimendengshu|/Users/|04-source-map|10-source-reading|search_sources"
)
LINK = re.compile(r"(?<!!)\[[^\]]+\]\(([^)]+)\)")
REQUIRED = (
    "SKILL.md", "LICENSE", "THIRD_PARTY_NOTICES.md", "agents/openai.yaml",
    "scripts/phone_palace.py", "scripts/qimen-engine/paipan.py",
    "scripts/qimen-engine/ke_paipan.py", "scripts/qimen-engine/setup-atopx.sh",
    "scripts/qimen-engine/atopx-src/LICENSE", "scripts/qimen-engine/atopx-src/go.mod",
    "references/04-method-compatibility.md", "references/10-validation-boundaries.md",
)


def release_files():
    return sorted(p for p in SKILL.rglob("*") if p.is_file()
                  and not any(part in EXCLUDED for part in p.relative_to(SKILL).parts)
                  and p.name != ".DS_Store" and p.suffix not in {".pyc", ".pyo"})


def inspect_text(name: str, data: bytes) -> list[str]:
    if Path(name).name == "LICENSE":
        return []  # Required software copyright and license text is retained.
    try:
        content = data.decode("utf-8")
    except UnicodeDecodeError:
        return [f"交付文件不应为二进制: {name}"]
    errors = []
    if SOURCE_MARKERS.search(name) or SOURCE_MARKERS.search(content):
        errors.append(f"发现资料出处标识: {name}")
    if Path(name).suffix.lower() in {".pdf", ".epub", ".djvu", ".doc", ".docx"}:
        errors.append(f"不得分发原始文献: {name}")
    return errors


def check_tree() -> list[str]:
    errors = [f"缺少必要文件: {name}" for name in REQUIRED if not (SKILL / name).is_file()]
    for path in release_files():
        rel = path.relative_to(SKILL).as_posix()
        errors.extend(inspect_text(rel, path.read_bytes()))
        if path.suffix != ".md":
            continue
        for target in LINK.findall(path.read_text(encoding="utf-8")):
            target = unquote(target.split("#", 1)[0]).strip("<>")
            if not target or "://" in target or target.startswith("mailto:"):
                continue
            resolved = (path.parent / target).resolve()
            if not resolved.is_relative_to(SKILL.resolve()) or not resolved.exists():
                errors.append(f"引用未在独立Skill内落实: {rel} -> {target}")
    if (PROJECT / "README.md").is_file():
        errors.extend(inspect_text("README.md", (PROJECT / "README.md").read_bytes()))
    return errors


def check_archive(path: Path) -> list[str]:
    expected = {"qimen-divination/" + p.relative_to(SKILL).as_posix(): p for p in release_files()}
    errors = []
    with zipfile.ZipFile(path) as archive:
        names = archive.namelist()
        if set(names) != set(expected) or len(names) != len(expected):
            errors.append("ZIP文件集合与独立Skill不一致")
        for item in archive.infolist():
            if item.filename not in expected:
                errors.append(f"ZIP出现额外文件: {item.filename}")
                continue
            data = archive.read(item)
            if data != expected[item.filename].read_bytes():
                errors.append(f"ZIP内容与Skill不同: {item.filename}")
            errors.extend(inspect_text(item.filename, data))
    return errors


def main() -> int:
    errors = check_tree()
    if len(sys.argv) > 1:
        errors.extend(check_archive(Path(sys.argv[1])))
    if errors:
        print("\n".join(errors), file=sys.stderr)
        return 1
    print(f"结构、方法隐私、引用与许可证检查通过：{len(release_files())}个交付文件")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
