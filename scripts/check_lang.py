import json

with open('data/baseline/compliance_rules.json', 'r', encoding='utf-8') as f:
    rules = json.load(f)

# Show some English Linux rules
linux_rules = [r for r in rules if r['osType'] == 1]
print(f"Linux rules: {len(linux_rules)}")
print("\n=== First 5 Linux rules (name + description) ===")
for r in linux_rules[:5]:
    print(f"  name: {r['name']}")
    print(f"  desc: {r['description']}")
    print(f"  fix:  {r['fixSuggestion']}")
    print(f"  risk: {r['riskDescription']}")
    print()

# Show some Windows rules
win_rules = [r for r in rules if r['osType'] == 2]
print(f"\nWindows rules: {len(win_rules)}")
if win_rules:
    print("=== First 3 Windows rules ===")
    for r in win_rules[:3]:
        print(f"  name: {r['name']}")
        print(f"  desc: {r['description']}")
        print()

# Check what categories exist
cats = set()
for r in rules:
    cats.add((r['category'], r.get('category', '')))
print(f"\nCategories: {sorted(cats)}")

# Check what's English vs Chinese
en_count = 0
cn_count = 0
for r in linux_rules:
    name = r.get('name', '')
    if any(ord(c) > 127 for c in name):
        cn_count += 1
    else:
        en_count += 1
print(f"\nLinux rules - English names: {en_count}, Chinese names: {cn_count}")