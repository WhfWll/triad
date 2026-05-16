import json

with open('data/baseline/compliance_rules_zh.json', 'r', encoding='utf-8') as f:
    rules = json.load(f)

print(f'Total rules: {len(rules)}')
print()

# Stats by osType
for ot in [1, 2, 3, 4]:
    subset = [r for r in rules if r['osType'] == ot]
    print(f'osType={ot}: {len(subset)} rules')

print()

# Show samples from different osTypes
for ot in [1, 2, 3, 4]:
    subset = [r for r in rules if r['osType'] == ot]
    print(f'=== osType={ot} samples ===')
    for r in subset[:2]:
        name = r['name']
        print(f'  name: {name[:80]}')
        print(f'  risk: {r["riskDescription"][:60]}')
        print()

# Check for any remaining English names
en_count = 0
for r in rules:
    name = r.get('name', '')
    if not any('\u4e00' <= c <= '\u9fff' for c in name):
        en_count += 1
        if en_count <= 5:
            print(f'  ENGLISH: [{r["osType"]}] {name[:100]}')

print(f'\nRemaining English names: {en_count}')
print('All rules translated successfully!' if en_count == 0 else f'{en_count} rules still need translation')