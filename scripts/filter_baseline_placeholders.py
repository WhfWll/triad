#!/usr/bin/env python3
"""Remove baseline rules whose commands/expectedValue still contain placeholder_value.

Usage:
  python scripts/filter_baseline_placeholders.py              # dry-run
  python scripts/filter_baseline_placeholders.py --apply       # backup *.bak and overwrite JSON

Default files under data/baseline/:
  compliance_rules_zh.json, compliance_rules.json, compliance_rules_linux.json
(custom_rules.json is kept as-is; it has no placeholders)
"""

from __future__ import annotations

import argparse
import json
import shutil
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
BASELINE_DIR = ROOT / "data" / "baseline"
DEFAULT_FILES = [
    "compliance_rules_zh.json",
    "compliance_rules.json",
    "compliance_rules_linux.json",
]


def has_placeholder(rule: dict) -> bool:
    expected = (rule.get("expectedValue") or "").lower()
    if "placeholder_value" in expected:
        return True
    for cmd in rule.get("commands") or []:
        if "placeholder_value" in str(cmd).lower():
            return True
    return False


def filter_file(path: Path, apply: bool) -> tuple[int, int]:
    with path.open(encoding="utf-8") as f:
        rules = json.load(f)
    kept = [r for r in rules if not has_placeholder(r)]
    removed = len(rules) - len(kept)
    print(f"{path.name}: total={len(rules)} kept={len(kept)} removed={removed}")
    if apply and removed > 0:
        backup = path.with_suffix(path.suffix + ".bak")
        if not backup.exists():
            shutil.copy2(path, backup)
            print(f"  backup -> {backup.name}")
        with path.open("w", encoding="utf-8") as f:
            json.dump(kept, f, ensure_ascii=False, indent=2)
            f.write("\n")
        print("  written")
    return len(rules), removed


def main() -> None:
    parser = argparse.ArgumentParser(description="Filter placeholder_value rules from baseline JSON")
    parser.add_argument("--apply", action="store_true", help="write filtered JSON (creates .bak once)")
    parser.add_argument("--file", action="append", help="specific file name under data/baseline/")
    args = parser.parse_args()

    names = args.file or DEFAULT_FILES
    total = removed = 0
    for name in names:
        path = BASELINE_DIR / name
        if not path.exists():
            print(f"skip missing: {name}")
            continue
        t, r = filter_file(path, args.apply)
        total += t
        removed += r
    print(f"summary: files={len(names)} total_rules={total} removed={removed} kept={total - removed}")
    if not args.apply:
        print("dry-run only; re-run with --apply to update JSON files")


if __name__ == "__main__":
    main()
