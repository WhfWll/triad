package services

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestShouldSkipNucleiTemplateSkipsWorkflowTemplate(t *testing.T) {
	doc := nucleiTemplateDoc{
		ID:        "xiuno-workflow",
		Workflows: []map[string]any{{"template": "http/exposed-panels/xiuno-panel.yaml"}},
	}
	doc.Info.Name = "Xiuno Workflow"
	doc.Info.Severity = "medium"

	reason, skip := shouldSkipNucleiTemplate(doc, "workflows/xiuno-workflow.yaml", doc.Info.Name, nil, nil)
	if !skip {
		t.Fatalf("expected workflow template to be skipped")
	}
	if reason == "" {
		t.Fatalf("expected skip reason to be populated")
	}
}

func TestShouldSkipNucleiTemplateKeepsRealVulnerabilityTemplate(t *testing.T) {
	doc := nucleiTemplateDoc{
		ID: "github-workflows-disclosure",
	}
	doc.Info.Name = "GitHub Workflows Disclosure"
	doc.Info.Severity = "medium"

	reason, skip := shouldSkipNucleiTemplate(doc, "http/exposures/configs/github-workflows-disclosure.yaml", doc.Info.Name, []string{"exposure", "config"}, nil)
	if skip {
		t.Fatalf("expected vulnerability template to be kept, got skip reason: %s", reason)
	}
}

func TestShouldSkipNucleiTemplateKeepsDetectionTemplate(t *testing.T) {
	doc := nucleiTemplateDoc{
		ID: "springboot-actuator-detection",
	}
	doc.Info.Name = "Spring Boot Actuator Detection"
	doc.Info.Severity = "medium"

	reason, skip := shouldSkipNucleiTemplate(doc, "http/technologies/springboot-actuator-detection.yaml", doc.Info.Name, []string{"springboot", "actuator", "detection"}, nil)
	if skip {
		t.Fatalf("expected medium severity detection template to be kept, got skip reason: %s", reason)
	}
}

func TestShouldSkipNucleiTemplateKeepsOtherTypeWithoutCVE(t *testing.T) {
	doc := nucleiTemplateDoc{
		ID: "custom-panel-exposure",
	}
	doc.Info.Name = "Custom Panel Exposure"
	doc.Info.Severity = "low"

	reason, skip := shouldSkipNucleiTemplate(doc, "http/exposures/panels/custom-panel-exposure.yaml", doc.Info.Name, []string{"panel", "exposure"}, []string{"https://example.com/advisory"})
	if skip {
		t.Fatalf("expected low severity non-CVE vulnerability template to be kept, got skip reason: %s", reason)
	}
}

func TestShouldTreatZipAsNuclei(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	fw, err := zw.Create("workflows/wordpress-workflow.yaml")
	if err != nil {
		t.Fatalf("create workflow entry: %v", err)
	}
	_, _ = fw.Write([]byte("id: wordpress-workflow\ninfo:\n  name: Wordpress Workflow\n  severity: medium\nworkflows:\n  - template: http/exposed-panels/wordpress-panel.yaml\n"))

	tf, err := zw.Create("http/exposures/configs/github-workflows-disclosure.yaml")
	if err != nil {
		t.Fatalf("create vuln entry: %v", err)
	}
	_, _ = tf.Write([]byte("id: github-workflows-disclosure\ninfo:\n  name: GitHub Workflows Disclosure\n  severity: medium\n"))

	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}

	if !shouldTreatZipAsNuclei(buf.Bytes()) {
		t.Fatalf("expected nuclei source zip to be routed to nuclei importer")
	}
}

func TestBuildNucleiVulIDKeepsCVE(t *testing.T) {
	got := buildNucleiVulID("http/cves/2024/example.yaml", "CVE-2024-12345")
	if got != "CVE-2024-12345" {
		t.Fatalf("expected cve-based vul id, got %s", got)
	}
}

func TestBuildNucleiVulIDCompressesLongNonCVETemplate(t *testing.T) {
	pocname := "nuclei-templates-main/cloud/alibaba/ram/password-policy-expiration-unconfigured.yaml"
	got := buildNucleiVulID(pocname, "")
	if len(got) > 50 {
		t.Fatalf("expected vul id <= 50 chars, got len=%d id=%s", len(got), got)
	}
	if !strings.HasPrefix(got, "NUCLEI-") {
		t.Fatalf("expected nuclei prefix, got %s", got)
	}
	if !strings.Contains(got, "-") {
		t.Fatalf("expected compact id to include hash separator, got %s", got)
	}
}

func TestLocalizeNucleiTemplateFieldsUsesZhMap(t *testing.T) {
	doc := nucleiTemplateDoc{ID: "iis-directory-browsing"}
	doc.Info.Name = "IIS Directory Browsing Detection"
	doc.Info.Description = "Ensures IIS directory browsing is disabled to prevent exposure of file structures."
	doc.Info.Remediation = ""

	name, desc, fix := localizeNucleiTemplateFields("file/audit/iis/iis-directory-browsing.yaml", doc, []string{"iis", "windows", "file", "hardening"}, nil)
	if name != "IIS 目录浏览检测" {
		t.Fatalf("expected mapped chinese name, got %s", name)
	}
	if !strings.Contains(desc, "目录结构") {
		t.Fatalf("expected mapped chinese description, got %s", desc)
	}
	if !strings.Contains(fix, "关闭 IIS 目录浏览功能") {
		t.Fatalf("expected mapped chinese fix, got %s", fix)
	}
}

func TestLocalizeNucleiTemplateFieldsGeneratesChineseFallback(t *testing.T) {
	doc := nucleiTemplateDoc{ID: "generic-rce"}
	doc.Info.Name = "Tomcat Remote Code Execution"
	doc.Info.Description = "This template detects remote code execution in Apache Tomcat."

	name, desc, fix := localizeNucleiTemplateFields("http/cves/tomcat/generic-rce.yaml", doc, []string{"tomcat", "rce"}, nil)
	if !strings.Contains(name, "远程代码执行") {
		t.Fatalf("expected chinese translated name, got %s", name)
	}
	if !strings.Contains(desc, "检测") {
		t.Fatalf("expected chinese translated description, got %s", desc)
	}
	if !strings.Contains(fix, "升级到安全版本") && !strings.Contains(fix, "官方安全公告") {
		t.Fatalf("expected chinese fallback remediation, got %s", fix)
	}
}

func TestAutoTranslateNucleiNamePatterns(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Unrestricted PostgreSQL Access in EC2", "EC2 上的 PostgreSQL 未受限访问"},
		{"Restrict EC2 RDP Access", "限制 EC2 RDP 访问"},
		{"IAM Access Analyzer is not Used", "IAM Access Analyzer 未使用"},
		{"IAM Password Policy Not Configured", "IAM 密码策略 未配置"},
		{"MFA not enabled on AWS Root Account", "AWS Root Account 未启用 MFA"},
		{"AceNet AceReporter Report - Arbitrary File Download", "AceNet AceReporter Report - 任意文件下载"},
	}
	for _, tc := range cases {
		got := autoTranslateNucleiName(tc.in)
		if got != tc.want {
			t.Fatalf("translate %q => %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestAutoTranslateNucleiNameVendors(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Landray-OA - Local File Inclusion", "蓝凌-OA - 本地文件包含"},
		{"UFIDA NC BeanShell Remote Command Execution", "用友 NC BeanShell 远程命令执行"},
		{"Pan Micro E-office File Upload", "泛微 E-Office 文件上传"},
		{"ZenTao CMS - SQL Injection", "禅道 CMS - SQL 注入"},
		{"Wuzhicms 4.1.0 - SQL Injection", "五指CMS 4.1.0 - SQL 注入"},
		{"Hongjing Human Resource Management System - SQL Injection", "宏景 人力资源管理系统 - SQL 注入"},
	}
	for _, tc := range cases {
		got := autoTranslateNucleiName(tc.in)
		if got != tc.want {
			t.Fatalf("vendor translate %q => %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestAutoTranslateNucleiNameCloudAndBaselinePatterns(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Public Access to ACK Cluster's API Server - Enabled", "ACK 集群 的 API Server 已开启公网访问"},
		{"Cluster Support for Network Policies - Missing", "集群 缺少 网络策略 支持"},
		{"Global Service (Multi-Region) Logging - Disabled", "全局服务 (跨区域) 日志 已禁用"},
		{"OS Patches - Outdated", "OS 补丁 过时"},
		{"Encryption for Unattached Disks - Disabled", "未挂载磁盘 的 加密 已禁用"},
		{"RAM Password Policy Expiration - Unconfigured", "RAM 密码策略过期时间 未配置"},
		{"RAM Password Policy requires Minimum Length 14 or Greater", "RAM 密码策略 要求最小长度为 14 或更大"},
		{"RAM Password Policy requires atleast One Lowercase - Unconfigured", "RAM 密码策略 要求至少包含一个 小写字母"},
		{"PowerShell Script Block Logging - Disabled", "PowerShell Script Block 日志 已禁用"},
		{"Remote Desktop Connections Allowed Without Password", "远程桌面连接 允许无需密码"},
		{"Remote Desktop Users Can Redirect Drives", "远程桌面用户 可重定向 驱动器"},
		{"Network Level Authentication for RDP Disabled", "RDP 的 网络级别身份验证 已禁用"},
		{"Remote Desktop Enabled on Non-Server OS", "在 非服务器操作系统 上启用了 远程桌面"},
		{"Store Passwords Using Reversible Encryption Enabled", "使用可逆加密存储密码 已启用"},
		{"Safe DLL Search Mode Disabled", "Safe DLL Search Mode 已禁用"},
		{"System Allows Shutdown Without Logging On", "系统 允许在未登录时关机"},
		{"Unencrypted Passwords to SMB Servers Allowed", "允许向 SMB 服务器 发送 未加密密码"},
	}
	for _, tc := range cases {
		got := autoTranslateNucleiName(tc.in)
		if got != tc.want {
			t.Fatalf("pattern translate %q => %q, want %q", tc.in, got, tc.want)
		}
	}
}
