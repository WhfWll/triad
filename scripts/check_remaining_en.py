import json

with open('data/baseline/compliance_rules_zh.json', 'r', encoding='utf-8') as f:
    rules = json.load(f)

# Find rules still with English names
en_rules = [r for r in rules if not any('\u4e00' <= c <= '\u9fff' for c in r.get('name', ''))]
print(f"Rules still with English names: {len(en_rules)}")
print("\n=== Sample English-named rules ===")
for r in en_rules[:20]:
    print(f"  [{r['osType']}] {r['name']}")

# Check by osType
from collections import Counter
os_counts = Counter()
for r in en_rules:
    os_counts[r['osType']] += 1
print(f"\n=== By osType ===")
for k, v in sorted(os_counts.items()):
    print(f"  osType={k}: {v}")

# Check if they have Chinese in description
cn_desc = sum(1 for r in en_rules if any('\u4e00' <= c <= '\u9fff' for c in r.get('description', '')))
print(f"\nEnglish-name rules with Chinese description: {cn_desc}")

# Check total by osType
total_os = Counter()
for r in rules:
    total_os[r['osType']] += 1
print(f"\n=== Total by osType ===")
for k, v in sorted(total_os.items()):
    print(f"  osType={k}: {v}")