# -*- coding: utf-8 -*-
"""Fix UTF-8 mojibake in security Vue files (U+FFFD + ?lt; style breaks)."""
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1] / "src" / "pages" / "security"

FIXES_APPSPECIFIC = [
    ("\ufffd?lt;/div>", "测</div>", 1),
    ("\ufffd?lt;/el-button>", "务</el-button>", 1),
    ('label="发现漏洞\ufffd">', 'label="发现漏洞数">', 1),
    ('label="状\ufffd">', 'label="状态">', 1),
    ('title="新建专项应用检测任\ufffd ', 'title="新建专项应用检测任务" ', 1),
    ('placeholder="请输入任务名\ufffd>', 'placeholder="请输入任务名称">', 1),
    ('label="云时\ufffd ', 'label="云时空" ', 1),
    ('label="亿赛\ufffd ', 'label="亿赛通" ', 1),
    ('label="检测类\ufffd prop=', 'label="检测类型" prop=', 1),
    ('placeholder="请选择检测类\ufffd>', 'placeholder="请选择检测类型">', 1),
    ('label="未授权访\ufffd ', 'label="未授权访问" ', 1),
    ('label="弱口\ufffd ', 'label="弱口令" ', 1),
    ('title="检测结果详\ufffd ', 'title="检测结果详情" ', 1),
    ("message: '请输入任务名\ufffd, trigger", "message: '请输入任务名称', trigger", 1),
    ("4: '云时\ufffd, 5: '亿赛\ufffd,", "4: '云时空', 5: '亿赛通',", 1),
    ("1: '等待检\ufffd, 2: '检测中', 3: '已完\ufffd }", "1: '等待检测', 2: '检测中', 3: '已完成' }", 1),
    ("2: '未授权访\ufffd, 3: 'SQL注入', 4: '文件上传', 5: '弱口\ufffd,", "2: '未授权访问', 3: 'SQL注入', 4: '文件上传', 5: '弱口令',", 1),
]

FIXES_DYNAMIC = [
    ("\ufffd?lt;/div>", "描</div>", 1),
    ('label="爬取页面\ufffd">', 'label="爬取页面数">', 1),
    ('label="发现漏洞\ufffd">', 'label="发现漏洞数">', 1),
    ('label="状\ufffd">', 'label="状态">', 1),
    ('title="新建动态扫描任\ufffd ', 'title="新建动态扫描任务" ', 1),
    ('placeholder="请输入任务名\ufffd>', 'placeholder="请输入任务名称">', 1),
    ("\ufffd?lt;/el-radio>", "合</el-radio>", 1),
    ('label="最大深\ufffd prop=', 'label="最大深度" prop=', 1),
    ('label="并发\ufffd prop=', 'label="并发数" prop=', 1),
    ('placeholder="用户\ufffd size', 'placeholder="用户名" size', 1),
    ('label="检测类\ufffd prop=', 'label="检测类型" prop=', 1),
    ('placeholder="请选择检测类\ufffd>', 'placeholder="请选择检测类型">', 1),
    ('爬取页面\ufffd?lt;/span>', '爬取页面数</span>', 1),
    ("message: '请输入任务名\ufffd, trigger", "message: '请输入任务名称', trigger", 1),
    ("2: '扫描\ufffd, 3: '已完\ufffd }", "2: '扫描中', 3: '已完成' }", 1),
]


def apply_fixes(path: Path, fixes):
    text = path.read_text(encoding="utf-8")
    orig = text
    for old, new, count in fixes:
        if count == -1:
            text = text.replace(old, new)
        else:
            n = text.count(old)
            if n != count and n > 0:
                text = text.replace(old, new)
            elif n == count:
                text = text.replace(old, new)
            else:
                # try without explicit U+FFFD (file may use ? only)
                alt_old = old.replace("\ufffd", "")
                if alt_old in text:
                    text = text.replace(alt_old, new.replace("\ufffd", ""))
                elif old.replace("\ufffd", "?") in text:
                    text = text.replace(old.replace("\ufffd", "?"), new)
    if text != orig:
        path.write_text(text, encoding="utf-8")
        print(f"Updated {path.name}")
    else:
        print(f"No changes {path.name} (patterns may differ)")


def main():
    apply_fixes(ROOT / "AppSpecificScan.vue", FIXES_APPSPECIFIC)
    apply_fixes(ROOT / "DynamicScan.vue", FIXES_DYNAMIC)


if __name__ == "__main__":
    main()
