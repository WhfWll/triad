#!/usr/bin/env python3
"""Fix baseline JSON rules: strip sudo, repair expectedValue copied from command.

Usage:
  python scripts/fix_baseline_rules.py           # dry-run stats
  python scripts/fix_baseline_rules.py --apply   # write (creates .pre-fix.bak once)
"""

from __future__ import annotations

import argparse
import json
import re
import shutil
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
BASELINE_DIR = ROOT / "data" / "baseline"
DEFAULT_FILES = [
    "compliance_rules_zh.json",
    "compliance_rules.json",
    "compliance_rules_linux.json",
]

NOT_KEYWORDS = re.compile(
    r"(not used|not be|not enabled|disabled|disable|remove|no default|"
    r"未使用|未启用|禁用|不应|未配置|remove any|ensure no|without)",
    re.I,
)


def normalize_cmd_prefix(cmd: str) -> str:
    cmd = cmd.strip()
    while cmd.startswith("$ "):
        cmd = cmd[2:].strip()
    while cmd.startswith("$"):
        cmd = cmd[1:].strip()
    return cmd


def strip_sudo(cmd: str) -> str:
    cmd = normalize_cmd_prefix(cmd)
    changed = True
    while changed:
        changed = False
        lower = cmd.lower()
        for prefix in ("sudo -n ", "sudo "):
            if lower.startswith(prefix):
                cmd = cmd[len(prefix) :].strip()
                changed = True
                break
    return cmd


def extract_grep_patterns(cmd: str) -> list[str]:
    patterns: list[str] = []
    for m in re.finditer(
        r"grep\s+(?:-(?:E|P|e|i|w|v|F|h|q)\s+)*['\"]([^'\"]+)['\"]", cmd, re.I
    ):
        patterns.append(m.group(1))
    for m in re.finditer(
        r"grep\s+(?:-(?:E|P|e|i|w|v|F|h|q)\s+)*(-?/\S+|/\S+|\S+)", cmd, re.I
    ):
        token = m.group(1)
        if token.startswith("-") and not token.startswith("-/"):
            continue
        if token in {"|", "2>/dev/null", "head", "tail", "wc", "awk", "sed"}:
            continue
        patterns.append(token.strip("'\" "))
    return patterns


def infer_expected(rule: dict, cmd: str) -> tuple[str, str]:
    """Return (expectedValue, matchType) when current expected duplicates command."""
    name = (rule.get("name") or "") + " " + (rule.get("description") or "")
    mt = rule.get("matchType") or "contains"
    lower_cmd = cmd.lower()

    if re.search(r"public\|private|'public\|private'", cmd, re.I):
        return "public", "not_contains"

    if "grep" in lower_cmd:
        patterns = extract_grep_patterns(cmd)
        if patterns:
            pick = patterns[-1]
            # path-like grep -w '/sbin/foo'
            if pick.startswith("/"):
                token = pick.rstrip("/").split("/")[-1] or pick
            else:
                token = pick.split("|")[0].strip()
            if NOT_KEYWORDS.search(name) or re.search(
                r"grep -v|ensure.*not|default.*not", name + " " + cmd, re.I
            ):
                return token, "not_contains"
            if "|" in pick and NOT_KEYWORDS.search(name):
                return pick.split("|")[0].strip(), "not_contains"
            return pick if not pick.startswith("/") else token, "contains"

    if "auditctl" in lower_cmd and "grep" in lower_cmd:
        patterns = extract_grep_patterns(cmd)
        if patterns:
            p = patterns[-1]
            token = p.split("/")[-1] if "/" in p else p
            return token, "contains"
        m = re.search(r"grep\s+(\S+)", cmd)
        if m:
            return m.group(1).strip("'\""), "contains"

    if "systemctl" in lower_cmd:
        m = re.search(r"systemctl\s+(is-enabled|is-active|show)\s+(\S+)", cmd, re.I)
        if m:
            return m.group(2), "contains"

    m = re.search(r"/etc/[\w./-]+", cmd)
    if m:
        seg = m.group(0).split("/")[-1]
        if seg:
            return seg, "contains"

    return rule.get("expectedValue") or "", mt


def fix_rule(rule: dict) -> tuple[bool, list[str]]:
    changes: list[str] = []
    cmds = list(rule.get("commands") or [])
    if not cmds:
        return False, changes

    new_cmds = []
    for i, c in enumerate(cmds):
        nc = strip_sudo(c)
        if nc != c:
            changes.append(f"cmd[{i}] strip sudo")
        new_cmds.append(nc)
    if new_cmds != cmds:
        rule["commands"] = new_cmds
        cmds = new_cmds

    exp = (rule.get("expectedValue") or "").strip()
    cmd0 = cmds[0].strip()
    exp_norm = strip_sudo(exp)

    if exp == cmds[0] or exp_norm == cmd0 or exp.startswith("$ sudo"):
        new_exp, new_mt = infer_expected(rule, cmd0)
        if new_exp and (new_exp != exp or new_mt != (rule.get("matchType") or "contains")):
            rule["expectedValue"] = new_exp
            rule["matchType"] = new_mt
            changes.append(f"expected -> {new_mt}:{new_exp[:48]}")
        elif exp.startswith("$ sudo") and exp_norm != exp:
            rule["expectedValue"] = exp_norm
            changes.append("expected strip sudo")

    return bool(changes), changes


def process_file(path: Path, apply: bool) -> dict:
    with path.open(encoding="utf-8") as f:
        rules = json.load(f)
    sudo_n = exp_n = 0
    samples: list[str] = []
    for r in rules:
        changed, notes = fix_rule(r)
        if not changed:
            continue
        if any("strip sudo" in n for n in notes):
            sudo_n += 1
        if any(n.startswith("expected") for n in notes):
            exp_n += 1
        if len(samples) < 5:
            samples.append(f"  - {r.get('name', '')[:50]} | {', '.join(notes)}")

    stats = {
        "file": path.name,
        "total": len(rules),
        "rules_changed": sudo_n,
        "expected_fixed": exp_n,
        "samples": samples,
    }
    if apply:
        backup = path.with_suffix(path.suffix + ".pre-fix.bak")
        if not backup.exists():
            shutil.copy2(path, backup)
        with path.open("w", encoding="utf-8") as f:
            json.dump(rules, f, ensure_ascii=False, indent=2)
            f.write("\n")
        stats["written"] = True
    return stats


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--apply", action="store_true")
    parser.add_argument("--file", action="append")
    args = parser.parse_args()
    names = args.file or DEFAULT_FILES
    for name in names:
        path = BASELINE_DIR / name
        if not path.exists():
            print(f"skip missing {name}")
            continue
        st = process_file(path, args.apply)
        print(
            f"{st['file']}: changed_rules~={st['rules_changed']} "
            f"expected_fixed={st['expected_fixed']}/{st['total']}"
        )
        for s in st["samples"]:
            print(s)
        if args.apply:
            print("  -> written (backup .pre-fix.bak)")
    if not args.apply:
        print("dry-run; use --apply to write")


if __name__ == "__main__":
    main()
