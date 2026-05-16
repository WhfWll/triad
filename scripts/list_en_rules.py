import json

with open('data/baseline/compliance_rules_zh.json', 'r', encoding='utf-8') as f:
    rules = json.load(f)

en_rules = [r for r in rules if not any('\u4e00' <= c <= '\u9fff' for c in r.get('name', ''))]

with open('scripts/en_rules_remaining.txt', 'w', encoding='utf-8') as f:
    f.write(f"Total English-named rules remaining: {len(en_rules)}\n\n")
    for ot in sorted(set(r['osType'] for r in en_rules)):
        subset = [r for r in en_rules if r['osType'] == ot]
        f.write(f"\n{'='*60}\n")
        f.write(f"osType={ot} ({len(subset)} rules)\n")
        f.write(f"{'='*60}\n")
        for r in subset:
            f.write(f"  {r['name']}\n")

print(f"Saved {len(en_rules)} remaining English rules to scripts/en_rules_remaining.txt")