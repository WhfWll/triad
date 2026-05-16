# -*- coding: utf-8 -*-
"""Second pass: fix remaining U+FFFD corruption in security Vue files."""
from pathlib import Path

U = "\ufffd"

def fix_text(text: str) -> str:
    pairs = [
        (f"发现漏洞{U}?>", '发现漏洞数">'),
        (f"状{U}?>", '状态">'),
        (f'爬取页面{U}?>', '爬取页面数">'),
        (f'title="新建专项应用检测任{U}? ', 'title="新建专项应用检测任务" '),
        (f'title="新建动态扫描任{U}? ', 'title="新建动态扫描任务" '),
        (f'title="检测结果详{U}? ', 'title="检测结果详情" '),
        (f'placeholder="请输入任务名{U}?>', 'placeholder="请输入任务名称">'),
        (f'label="云时{U}? :value', 'label="云时空" :value'),
        (f'label="亿赛{U}? :value', 'label="亿赛通" :value'),
        (f'label="检测类{U}? prop=', 'label="检测类型" prop='),
        (f'placeholder="请选择检测类{U}?>', 'placeholder="请选择检测类型">'),
        (f'label="未授权访{U}? :value', 'label="未授权访问" :value'),
        (f'label="弱口{U}? :value', 'label="弱口令" :value'),
        (f'label="最大深{U}? prop=', 'label="最大深度" prop='),
        (f'label="并发{U}? prop=', 'label="并发数" prop='),
        (f'placeholder="用户{U}? size', 'placeholder="用户名" size'),
        (f"message: '请输入任务名{U}?, trigger", "message: '请输入任务名称', trigger"),
        (f"4: '云时{U}?, 5: '亿赛{U}?,", "4: '云时空', 5: '亿赛通',"),
        (f"1: '等待检{U}?, 2: '检测中', 3: '已完{U}? }}", "1: '等待检测', 2: '检测中', 3: '已完成' }"),
        (f"2: '未授权访{U}?, 3: 'SQL注入', 4: '文件上传', 5: '弱口{U}?,", "2: '未授权访问', 3: 'SQL注入', 4: '文件上传', 5: '弱口令',"),
        (f"2: '扫描{U}?, 3: '已完{U}? }}", "2: '扫描中', 3: '已完成' }"),
        (f'爬取页面{U}?lt;/span>', '爬取页面数</span>'),
    ]
    for old, new in pairs:
        text = text.replace(old, new)
    return text


def main():
    for name in ("AppSpecificScan.vue", "DynamicScan.vue"):
        p = Path(__file__).resolve().parents[1] / "src" / "pages" / "security" / name
        raw = p.read_text(encoding="utf-8")
        fixed = fix_text(raw)
        if fixed != raw:
            p.write_text(fixed, encoding="utf-8")
            print("fixed", name)
        else:
            print("no change", name)


if __name__ == "__main__":
    main()
