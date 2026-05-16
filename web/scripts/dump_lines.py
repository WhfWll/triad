# -*- coding: utf-8 -*-
import pathlib
p = pathlib.Path(r"D:/goproject/triad/web/src/pages/security/AppSpecificScan.vue")
lines = p.read_text(encoding="utf-8").splitlines()
out = pathlib.Path(r"D:/goproject/triad/web/tmp_repr.txt")
with out.open("w", encoding="utf-8") as f:
    for i in list(range(93, 104)) + [33, 40, 67, 70, 77, 78, 93, 94, 115, 198, 285, 301, 309]:
        f.write(f"{i+1}: {repr(lines[i])}\n")
