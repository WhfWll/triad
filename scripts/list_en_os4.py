import json

with open('data/baseline/compliance_rules_zh.json', 'r', encoding='utf-8') as f:
    rules = json.load(f)

en_rules = [r for r in rules if not any('\u4e00' <= c <= '\u9fff' for c in r.get('name', ''))]

# Print osType=4 rules
print("osType=4 rules:")
for r in en_rules:
    if r['osType'] == 4:
        print(f"  {r['name']}")