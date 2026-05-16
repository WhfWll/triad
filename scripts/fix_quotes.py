import re

with open('scripts/generate_all_os_rules.py', 'r', encoding='utf-8') as f:
    content = f.read()

# Fix the specific pattern: \"SOMETHING\"""] -> \"SOMETHING\""],
content = content.replace('\\"""],', '\\""],')

# Also fix any remaining unescaped patterns
# Pattern: "powershell - Command "SOMETHING""],
content = re.sub(
    r'"powershell -Command "([^"]+)"\],',
    r'"powershell -Command \\"\1\\"],',
    content
)

with open('scripts/generate_all_os_rules.py', 'w', encoding='utf-8') as f:
    f.write(content)

print('Fixed quotes')