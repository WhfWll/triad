import json
from collections import Counter

with open('data/baseline/compliance_rules.json', 'r', encoding='utf-8') as f:
    rules = json.load(f)

# Analyze Linux rules (English)
linux_rules = [r for r in rules if r['osType'] == 1]

# Extract common name prefixes
name_prefixes = Counter()
for r in linux_rules:
    name = r['name']
    # Get first few words
    words = name.split()
    prefix = ' '.join(words[:4])
    name_prefixes[prefix] += 1

print("=== Most common name prefixes ===")
for p, c in name_prefixes.most_common(30):
    print(f"  [{c:4d}] {p}")

# Extract common description patterns
desc_starts = Counter()
for r in linux_rules:
    desc = r.get('description', '')
    if desc:
        first_60 = desc[:60]
        desc_starts[first_60] += 1

print("\n=== Most common description starts ===")
for d, c in desc_starts.most_common(20):
    print(f"  [{c:4d}] {d}...")

# Check categories
cat_names = {1: '系统审计', 2: '文件权限', 3: '账户策略', 4: '网络服务', 5: '内核安全',
             6: '日志审计', 7: '密码策略', 8: 'SSH配置', 10: '系统更新', 11: '服务管理',
             12: '防火墙', 13: 'SELinux', 14: '身份认证', 99: '其他'}
cat_counts = Counter()
for r in linux_rules:
    cat_counts[r['category']] += 1
print("\n=== Rules by category ===")
for c, cnt in sorted(cat_counts.items()):
    print(f"  [{c:4d}] {cat_names.get(c, '未知')}: {cnt}")

# Check unique risk levels
risks = Counter()
for r in linux_rules:
    risks[r['risk']] += 1
print(f"\n=== Risk levels: {dict(risks)} ===")

# Check unique match types
match_types = Counter()
for r in linux_rules:
    match_types[r.get('matchType', 'N/A')] += 1
print(f"=== Match types: {dict(match_types)} ===")

# Sample all unique names
print("\n=== All unique rule names (first 100) ===")
for r in linux_rules[:100]:
    print(f"  {r['name']}")