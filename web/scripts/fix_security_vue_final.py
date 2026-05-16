# -*- coding: utf-8 -*-
"""Fix remaining mojibake in AppSpecificScan.vue and DynamicScan.vue."""
import re
from pathlib import Path

DIR = Path(__file__).resolve().parents[1] / "src" / "pages" / "security"


def fix_app_specific(text: str) -> str:
    text = re.sub(
        r'<el-table-column prop="vulnCount" label="发现漏洞[^>]*>',
        '<el-table-column prop="vulnCount" label="发现漏洞数">',
        text,
    )
    text = re.sub(
        r'<el-table-column prop="status" label="状[^>]*>',
        '<el-table-column prop="status" label="状态">',
        text,
    )
    text = re.sub(
        r'<el-dialog title="新建专项应用检测任[^"]*"\s+:visible\.sync="dialogVisible"',
        '<el-dialog title="新建专项应用检测任务" :visible.sync="dialogVisible"',
        text,
    )
    text = re.sub(
        r'placeholder="请输入任务名[^"]*">',
        'placeholder="请输入任务名称">',
        text,
    )
    text = re.sub(
        r'<el-option label="云时[^"]*"\s+:value="4"',
        '<el-option label="云时空" :value="4"',
        text,
    )
    text = re.sub(
        r'<el-option label="亿赛[^"]*"\s+:value="5"',
        '<el-option label="亿赛通" :value="5"',
        text,
    )
    text = re.sub(
        r'<el-form-item label="检测类[^"]*"\s+prop="checkTypes"',
        '<el-form-item label="检测类型" prop="checkTypes"',
        text,
    )
    text = re.sub(
        r'placeholder="请选择检测类[^"]*">',
        'placeholder="请选择检测类型">',
        text,
    )
    text = re.sub(
        r'<el-option label="未授权访[^"]*"\s+:value="2"',
        '<el-option label="未授权访问" :value="2"',
        text,
    )
    text = re.sub(
        r'<el-option label="弱口[^"]*"\s+:value="5"',
        '<el-option label="弱口令" :value="5"',
        text,
    )
    text = re.sub(
        r'title="检测结果详[^"]*"\s+:visible\.sync="detailVisible"',
        'title="检测结果详情" :visible.sync="detailVisible"',
        text,
    )
    text = re.sub(
        r"message: '请输入任务名[^']*',\s*trigger:\s*'blur'",
        "message: '请输入任务名称', trigger: 'blur'",
        text,
    )
    text = re.sub(
        r"4: '云时[^']*',\s*5: '亿赛[^']*',",
        "4: '云时空', 5: '亿赛通',",
        text,
    )
    text = re.sub(
        r"1: '等待检[^']*',\s*2: '检测中',\s*3: '已完[^']*'\s*}",
        "1: '等待检测', 2: '检测中', 3: '已完成' }",
        text,
    )
    text = re.sub(
        r"2: '未授权访[^']*',\s*3: 'SQL注入',\s*4: '文件上传',\s*5: '弱口[^']*',",
        "2: '未授权访问', 3: 'SQL注入', 4: '文件上传', 5: '弱口令',",
        text,
    )
    return text


def fix_dynamic(text: str) -> str:
    text = re.sub(
        r'<el-table-column prop="pageCount" label="爬取页面[^>]*>',
        '<el-table-column prop="pageCount" label="爬取页面数">',
        text,
    )
    text = re.sub(
        r'<el-table-column prop="vulnCount" label="发现漏洞[^>]*>',
        '<el-table-column prop="vulnCount" label="发现漏洞数">',
        text,
    )
    text = re.sub(
        r'<el-table-column prop="status" label="状[^>]*>',
        '<el-table-column prop="status" label="状态">',
        text,
    )
    text = re.sub(
        r'title="新建动态扫描任[^"]*"\s+:visible\.sync="dialogVisible"',
        'title="新建动态扫描任务" :visible.sync="dialogVisible"',
        text,
    )
    text = re.sub(
        r'placeholder="请输入任务名[^"]*">',
        'placeholder="请输入任务名称">',
        text,
    )
    text = re.sub(
        r'<el-form-item label="最大深[^"]*"\s+prop="maxDepth"',
        '<el-form-item label="最大深度" prop="maxDepth"',
        text,
    )
    text = re.sub(
        r'<el-form-item label="并发[^"]*"\s+prop="concurrency"',
        '<el-form-item label="并发数" prop="concurrency"',
        text,
    )
    text = re.sub(
        r'placeholder="用户[^"]*"\s+size="small"',
        'placeholder="用户名" size="small"',
        text,
    )
    text = re.sub(
        r'<el-form-item label="检测类[^"]*"\s+prop="checkTypes"',
        '<el-form-item label="检测类型" prop="checkTypes"',
        text,
    )
    text = re.sub(
        r'placeholder="请选择检测类[^"]*">',
        'placeholder="请选择检测类型">',
        text,
    )
    text = re.sub(
        r"message: '请输入任务名[^']*',\s*trigger:\s*'blur'",
        "message: '请输入任务名称', trigger: 'blur'",
        text,
    )
    text = re.sub(
        r"2: '扫描[^']*',\s*3: '已完[^']*'\s*}",
        "2: '扫描中', 3: '已完成' }",
        text,
    )
    return text


def main():
    for name, fn in (
        ("AppSpecificScan.vue", fix_app_specific),
        ("DynamicScan.vue", fix_dynamic),
    ):
        p = DIR / name
        raw = p.read_text(encoding="utf-8")
        fixed = fn(raw)
        if fixed != raw:
            p.write_text(fixed, encoding="utf-8")
            print("OK", name)
        else:
            print("skip", name)


if __name__ == "__main__":
    main()
