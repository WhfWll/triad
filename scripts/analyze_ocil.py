import os, yaml, re
from collections import Counter

guide_dir = 'linux_os/guide'
ocil_patterns = Counter()
ocil_examples = []

for root, dirs, files in os.walk(guide_dir):
    if 'rule.yml' not in files:
        continue
    
    try:
        with open(os.path.join(root, 'rule.yml'), encoding='utf-8') as f:
            data = yaml.safe_load(f)
    except:
        continue
    
    if not data or not isinstance(data, dict):
        continue
    if not data.get('title'):
        continue
    
    ocil = data.get('ocil', '')
    if not ocil:
        continue
    
    # Check if our current regex would match
    current_re = r'(?:^|[.;])\s*(sudo\s+)?(grep|cat|ls|find|stat|sysctl|systemctl|rpm|dpkg|auditctl|ss|netstat|getenforce|sestatus|mount|ps|df|du|chkconfig|readlink|file|sha256sum|md5sum|awk|sed|chkconfig)\s+[-/\w].*?(?=[.;]|$)'
    if re.search(current_re, ocil, re.IGNORECASE):
        continue  # Already matched
    
    # Check for other command patterns
    cmd_patterns = [
        r'\b(check|verify|ensure|see|determine|review|confirm|inspect|examine|look for|search for)\b',
        r'\b(run|execute|use|type|enter|issue)\s+the\s+following\s+command',
        r'`([^`]+)`',
        r'\"([^\"]+)\"',
        r'#\s*(grep|cat|ls|find|stat|sysctl|systemctl|rpm|dpkg|auditctl|ss|netstat|mount|ps|df|du|readlink|file|sha256sum|md5sum|awk|sed|echo|printf|test|if|for|while|which|whereis|type|command)\b',
    ]
    
    matched = False
    for p in cmd_patterns:
        if re.search(p, ocil, re.IGNORECASE):
            matched = True
            break
    
    if not matched:
        # Check first 100 chars
        first_line = ocil.strip()[:100].replace('\n', ' ')
        ocil_patterns[first_line] += 1

print("=== 最常见的OCIL开头文本(未匹配) ===")
for text, count in ocil_patterns.most_common(20):
    print(f"  [{count}] {text}...")