# -*- coding: utf-8 -*-
from pathlib import Path
p = Path(r"D:/goproject/triad/web/src/pages/security/AppSpecificScan.vue")
for i, line in enumerate(p.read_text(encoding="utf-8").splitlines()):
    if i + 1 in (34, 41, 68, 71, 78, 79, 94, 95, 97, 100, 116, 199, 286, 302, 310):
        print(i + 1, [hex(ord(c)) for c in line[-40:]])
