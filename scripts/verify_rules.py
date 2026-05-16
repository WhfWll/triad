import json

with open('data/baseline/compliance_rules.json', 'r', encoding='utf-8') as f:
    rules = json.load(f)

print(f'Total rules: {len(rules)}')

os_counts = {}
for r in rules:
    ot = r.get('osType', 0)
    os_counts[ot] = os_counts.get(ot, 0) + 1

os_names = {1: 'Linux/Unix', 2: 'Windows', 3: 'Domestic OS', 4: 'Embedded OS'}
print('\nRules by OS type:')
for ot in sorted(os_counts.keys()):
    name = os_names.get(ot, f'Unknown({ot})')
    print(f'  {name}: {os_counts[ot]}')

# Check for invalid commands
invalid = 0
empty_desc = 0
for r in rules:
    cmds = r.get('commands', [])
    if not cmds or all(not c.strip() for c in cmds):
        invalid += 1
    if not r.get('description', '').strip():
        empty_desc += 1

print(f'\nInvalid command rules: {invalid}')
print(f'Empty description rules: {empty_desc}')

# Check Windows rules
win_rules = [r for r in rules if r.get('osType') == 2]
print(f'\nWindows rules sample (first 3):')
for r in win_rules[:3]:
    print(f'  ID={r["id"]}, Name={r["name"]}, Category={r["category"]}, Risk={r["risk"]}')
    print(f'  Command: {r["commands"][0][:100]}...')

# Check domestic OS rules
dom_rules = [r for r in rules if r.get('osType') == 3]
print(f'\nDomestic OS rules sample (first 2):')
for r in dom_rules[:2]:
    print(f'  ID={r["id"]}, Name={r["name"]}, Category={r["category"]}')

# Check embedded OS rules
emb_rules = [r for r in rules if r.get('osType') == 4]
print(f'\nEmbedded OS rules sample (first 2):')
for r in emb_rules[:2]:
    print(f'  ID={r["id"]}, Name={r["name"]}, Category={r["category"]}')