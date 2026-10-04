"""Чтение базы клинических рекомендаций из knowledge/.

- guidelines/manifest.json — версия базы рекомендаций;
- guidelines/*.md         — фрагменты. Каждый начинается со строки `## [ID] Ссылка для врача`,
                            дальше до следующего заголовка — текст;
- instructions.md         — системная инструкция для LLM.
"""

import json
import re
from pathlib import Path

from app.domain.guidelines import GuidelineFragment, Guidelines

_FRAGMENT_HEADER = re.compile(r"^## \[(?P<id>[A-Z0-9-]+)\]\s*(?P<ref>.+?)\s*$", re.MULTILINE)


def load_guidelines(root: Path) -> Guidelines:
    guidelines_dir = root / "guidelines"
    manifest = json.loads((guidelines_dir / "manifest.json").read_text(encoding="utf-8"))
    if not manifest.get("version"):
        raise ValueError("guidelines manifest must declare a non-empty version")

    fragments: dict[str, GuidelineFragment] = {}
    for path in sorted(guidelines_dir.glob("*.md")):
        for fragment in _parse_fragments(path.read_text(encoding="utf-8")):
            if fragment.id in fragments:
                raise ValueError(f"duplicate guideline fragment id {fragment.id!r} in {path.name}")
            fragments[fragment.id] = fragment
    if not fragments:
        raise ValueError("guidelines must not be empty")

    return Guidelines(
        version=manifest["version"],
        fragments=fragments,
        instructions=(root / "instructions.md").read_text(encoding="utf-8").strip(),
    )


def _parse_fragments(text: str) -> list[GuidelineFragment]:
    headers = list(_FRAGMENT_HEADER.finditer(text))
    fragments = []
    for i, header in enumerate(headers):
        end = headers[i + 1].start() if i + 1 < len(headers) else len(text)
        body = text[header.end() : end].strip()
        fragments.append(GuidelineFragment(id=header["id"], ref=header["ref"], text=body))
    return fragments
