import json

with open('data/baseline/compliance_rules.json', 'r') as f:
    rules = json.load(f)

print(f'总规则数: {len(rules)}')
print()

# 检查一些示例规则
for r in rules[:5]:
    print(f'  [{r["id"]}] {r["name"][:60]}')
    print(f'    命令: {r["commands"][0][:80]}...')
    print(f'    期望值: {r["expectedValue"][:40]}')
    print(f'    匹配类型: {r["matchType"]}')
    print()

# 检查命令质量
bad_cmds = 0
for r in rules:
    for cmd in r['commands']:
        if any(w in cmd.lower() for w in ['to the', 'for the', 'of the', 'run the following', 'if the']):
            bad_cmds += 1
            break
print(f'包含无效文本的命令数: {bad_cmds}')

# 检查空描述
empty_desc = sum(1 for r in rules if not r.get('description'))
print(f'空描述数: {empty_desc}')

# 检查空修复建议
empty_fix = sum(1 for r in rules if not r.get('fixSuggestion'))
print(f'空修复建议数: {empty_fix}')

# 分类统计
from collections import Counter
cat_counter = Counter(r['category'] for r in rules)
print(f'\n分类分布:')
cat_names = {1: '密码策略', 2: '用户权限', 3: '防火墙规则', 4: '内核安全',
             5: '文件权限', 6: '审计日志', 7: '网络服务', 8: '系统更新',
             9: 'SSH配置', 99: '其他'}
for cat, count in sorted(cat_counter.items()):
    name = cat_names.get(cat, f'未知({cat})')
    print(f'  {name}: {count}')

# 风险等级统计
risk_counter = Counter(r['risk'] for r in rules)
risk_names = {0: '严重', 1: '高危', 2: '中危', 3: '低危', 4: '信息'}
print(f'\n风险等级分布:')
for risk, count in sorted(risk_counter.items()):
    name = risk_names.get(risk, f'未知({risk})')
    print(f'  {name}: {count}')