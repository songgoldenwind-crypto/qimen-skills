#!/usr/bin/env python3
"""Build a portable Skill ZIP after checking its actual contents."""

from __future__ import annotations

import hashlib
import json
import zipfile

from check_release import PROJECT, SKILL, check_archive, check_tree, release_files


def main() -> None:
    errors = check_tree()
    if errors:
        raise SystemExit("\n".join(errors))
    output = PROJECT / "dist"
    output.mkdir(exist_ok=True)
    archive = output / "qimen-divination.zip"
    with zipfile.ZipFile(archive, "w", compression=zipfile.ZIP_DEFLATED) as bundle:
        for path in release_files():
            bundle.write(path, "qimen-divination/" + path.relative_to(SKILL).as_posix())
    errors = check_archive(archive)
    if errors:
        raise SystemExit("\n".join(errors))
    manifest = {
        "skill": "qimen-divination",
        "archive": archive.name,
        "sha256": hashlib.sha256(archive.read_bytes()).hexdigest(),
        "files": len(release_files()),
    }
    (output / "manifest.json").write_text(json.dumps(manifest, ensure_ascii=False, indent=2) + "\n")
    print(json.dumps(manifest, ensure_ascii=False))


if __name__ == "__main__":
    main()
