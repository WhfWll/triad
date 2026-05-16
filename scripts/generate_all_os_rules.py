"""
为所有4种操作系统类型生成规则:
  osType=1: Linux/Unix (已有)
  osType=2: Windows (新建)
  osType=3: 国产操作系统 (从Linux复制)
  osType=4: 嵌入式OS (从Linux子集)
"""
import json
import os

def main():
    src_file = r'd:\goproject\triad\data\baseline\compliance_rules.json'
    output_dir = r'd:\goproject\triad\data\baseline'
    
    with open(src_file, 'r', encoding='utf-8') as f:
        linux_rules = json.load(f)
    
    print(f"Linux规则数: {len(linux_rules)}")
    
    all_rules = []
    
    # 1. Linux/Unix 规则 (osType=1)
    for r in linux_rules:
        rule = dict(r)
        rule['osType'] = 1
        all_rules.append(rule)
    
    # 2. 国产操作系统 规则 (osType=3) - 从Linux复制，国产OS基于Linux
    for r in linux_rules:
        rule = dict(r)
        rule['osType'] = 3
        rule['id'] = abs(hash(f"{r['id']}_domestic")) % 1000000
        rule['name'] = r['name']
        rule['description'] = r.get('description', '')
        rule['commands'] = list(r['commands'])
        rule['expectedValue'] = r['expectedValue']
        rule['matchType'] = r['matchType']
        rule['fixSuggestion'] = r.get('fixSuggestion', '请参考系统安全最佳实践进行修复')
        rule['riskDescription'] = r.get('riskDescription', '')
        all_rules.append(rule)
    
    # 3. 嵌入式OS 规则 (osType=4) - 从Linux复制相关子集
    # 嵌入式OS通常关注: 内核安全、文件权限、审计日志、网络服务、密码策略
    embedded_categories = {1, 4, 5, 6, 7}
    embedded_count = 0
    for r in linux_rules:
        if r.get('category', 99) in embedded_categories:
            rule = dict(r)
            rule['osType'] = 4
            rule['id'] = abs(hash(f"{r['id']}_embedded")) % 1000000
            rule['name'] = r['name']
            rule['description'] = r.get('description', '')
            rule['commands'] = list(r['commands'])
            rule['expectedValue'] = r['expectedValue']
            rule['matchType'] = r['matchType']
            rule['fixSuggestion'] = r.get('fixSuggestion', '请参考系统安全最佳实践进行修复')
            rule['riskDescription'] = r.get('riskDescription', '')
            all_rules.append(rule)
            embedded_count += 1
    
    print(f"嵌入式OS规则数: {embedded_count}")
    
    # 4. Windows 规则 (osType=2)
    windows_rules = generate_windows_rules()
    all_rules.extend(windows_rules)
    print(f"Windows规则数: {len(windows_rules)}")
    
    # 保存所有规则
    output_file = os.path.join(output_dir, 'compliance_rules_all.json')
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(all_rules, f, ensure_ascii=False, indent=2)
    
    print(f"\n总规则数: {len(all_rules)}")
    
    # 按OS类型统计
    os_counts = {}
    for r in all_rules:
        ot = r.get('osType', 0)
        os_counts[ot] = os_counts.get(ot, 0) + 1
    
    print("\n各OS类型规则数:")
    os_names = {1: 'Linux/Unix', 2: 'Windows', 3: '国产操作系统', 4: '嵌入式OS'}
    for ot in sorted(os_counts.keys()):
        name = os_names.get(ot, f'未知({ot})')
        print(f"  {name}: {os_counts[ot]}")
    
    # 按分类统计
    from collections import Counter
    cat_counter = Counter(r.get('category', 99) for r in all_rules)
    print("\n分类分布:")
    for cat, count in sorted(cat_counter.items()):
        print(f"  分类{cat}: {count}")


def generate_windows_rules():
    """生成Windows安全配置核查规则"""
    rules = []
    
    windows_rules_data = [
        # 密码策略 (category=10)
        {
            "name": "检查密码复杂度是否启用",
            "description": "检查Windows系统是否启用了密码复杂度要求",
            "category": 10, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System' -Name 'PasswordComplexity' | Select-Object -ExpandProperty PasswordComplexity\""],
            "expectedValue": "1", "matchType": "contains",
            "fixSuggestion": "请启用密码复杂度策略: gpedit.msc → 计算机配置 → Windows设置 → 安全设置 → 账户策略 → 密码策略 → 密码必须符合复杂性要求",
            "riskDescription": "高危风险: 未启用密码复杂度要求，系统容易受到暴力破解攻击"
        },
        {
            "name": "检查密码最小长度",
            "description": "检查Windows系统密码最小长度设置",
            "category": 10, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System' -Name 'MinPasswordLength' | Select-Object -ExpandProperty MinPasswordLength\""],
            "expectedValue": "8", "matchType": "contains",
            "fixSuggestion": "请设置密码最小长度为8位以上",
            "riskDescription": "高危风险: 密码最小长度不足，建议设置为8位以上"
        },
        {
            "name": "检查密码最长使用期限",
            "description": "检查Windows系统密码最长使用期限",
            "category": 10, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System' -Name 'MaxPasswordAge' | Select-Object -ExpandProperty MaxPasswordAge\""],
            "expectedValue": "90", "matchType": "contains",
            "fixSuggestion": "请设置密码最长使用期限为90天以内",
            "riskDescription": "中危风险: 密码使用期限过长，建议设置为90天以内"
        },
        {
            "name": "检查密码历史记录",
            "description": "检查Windows系统是否启用了密码历史记录",
            "category": 10, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System' -Name 'PasswordHistorySize' | Select-Object -ExpandProperty PasswordHistorySize\""],
            "expectedValue": "5", "matchType": "contains",
            "fixSuggestion": "请设置密码历史记录为5个以上",
            "riskDescription": "中危风险: 密码历史记录不足，建议设置为记住5个以上密码"
        },
        {
            "name": "检查账户锁定阈值",
            "description": "检查Windows系统账户锁定阈值设置",
            "category": 10, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System' -Name 'LockoutBadCount' | Select-Object -ExpandProperty LockoutBadCount\""],
            "expectedValue": "5", "matchType": "contains",
            "fixSuggestion": "请设置账户锁定阈值为5次以内",
            "riskDescription": "高危风险: 未设置账户锁定阈值，系统容易受到暴力破解攻击"
        },
        {
            "name": "检查账户锁定时间",
            "description": "检查Windows系统账户锁定时间设置",
            "category": 10, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System' -Name 'LockoutDuration' | Select-Object -ExpandProperty LockoutDuration\""],
            "expectedValue": "30", "matchType": "contains",
            "fixSuggestion": "请设置账户锁定时间为30分钟以上",
            "riskDescription": "中危风险: 账户锁定时间过短"
        },
        {
            "name": "检查最小密码长度(域策略)",
            "description": "检查Windows域策略中的最小密码长度",
            "category": 10, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-ADDefaultDomainPasswordPolicy | Select-Object -ExpandProperty MinPasswordLength\""],
            "expectedValue": "8", "matchType": "contains",
            "fixSuggestion": "请设置域密码策略最小长度为8位以上",
            "riskDescription": "高危风险: 域密码最小长度不足"
        },
        
        # 用户权限 (category=11)
        {
            "name": "检查Guest账户状态",
            "description": "检查Windows系统Guest账户是否已禁用",
            "category": 11, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-LocalUser -Name 'Guest' | Select-Object -ExpandProperty Enabled\""],
            "expectedValue": "False", "matchType": "contains",
            "fixSuggestion": "请禁用Guest账户: net user Guest /active:no",
            "riskDescription": "高危风险: Guest账户未禁用，存在安全风险"
        },
        {
            "name": "检查Administrator账户重命名",
            "description": "检查Windows系统Administrator账户是否已重命名",
            "category": 11, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-LocalUser -Name 'Administrator' | Select-Object -ExpandProperty Name\""],
            "expectedValue": "NOT_FOUND", "matchType": "not_contains",
            "fixSuggestion": "请重命名Administrator账户，避免使用默认管理员账户名",
            "riskDescription": "中危风险: 未重命名Administrator账户，建议修改为其他名称"
        },
        {
            "name": "检查空密码账户",
            "description": "检查Windows系统是否存在空密码账户",
            "category": 11, "risk": 0, "osType": 2,
            "commands": ["powershell -Command \"Get-LocalUser | Where-Object { $_.PasswordRequired -eq $false -and $_.Enabled -eq $true } | Select-Object -ExpandProperty Name\""],
            "expectedValue": "ALL_OK", "matchType": "contains",
            "fixSuggestion": "请为所有账户设置密码，不允许空密码账户存在",
            "riskDescription": "严重风险: 存在空密码账户，系统极易被入侵"
        },
        {
            "name": "检查管理员组成员",
            "description": "检查Windows系统Administrators组成员",
            "category": 11, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-LocalGroupMember -Group 'Administrators' | Select-Object -ExpandProperty Name\""],
            "expectedValue": "Administrator", "matchType": "contains",
            "fixSuggestion": "请检查Administrators组成员，确保只有必要用户拥有管理员权限",
            "riskDescription": "中危风险: 请检查管理员组成员是否合理"
        },
        {
            "name": "检查远程桌面用户组",
            "description": "检查Windows系统Remote Desktop Users组成员",
            "category": 11, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-LocalGroupMember -Group 'Remote Desktop Users' 2>$null | Select-Object -ExpandProperty Name\""],
            "expectedValue": "ALL_OK", "matchType": "contains",
            "fixSuggestion": "请检查远程桌面用户组成员，确保只有必要用户拥有远程访问权限",
            "riskDescription": "中危风险: 远程桌面用户组存在成员，请确认是否必要"
        },
        
        # 防火墙规则 (category=12)
        {
            "name": "检查Windows防火墙状态",
            "description": "检查Windows防火墙是否已启用",
            "category": 12, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-NetFirewallProfile -Profile Domain,Public,Private | Select-Object Name,Enabled\""],
            "expectedValue": "True", "matchType": "contains",
            "fixSuggestion": "请启用Windows防火墙: Set-NetFirewallProfile -All -Enabled True",
            "riskDescription": "高危风险: Windows防火墙未启用，系统暴露在网络攻击风险中"
        },
        {
            "name": "检查远程桌面是否启用",
            "description": "检查Windows系统远程桌面是否已启用",
            "category": 12, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SYSTEM\\CurrentControlSet\\Control\\Terminal Server' -Name 'fDenyTSConnections' | Select-Object -ExpandProperty fDenyTSConnections\""],
            "expectedValue": "1", "matchType": "contains",
            "fixSuggestion": "如非必要，请禁用远程桌面",
            "riskDescription": "中危风险: 远程桌面已启用，如非必要请禁用"
        },
        {
            "name": "检查远程桌面端口",
            "description": "检查Windows远程桌面端口是否已更改",
            "category": 12, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SYSTEM\\CurrentControlSet\\Control\\Terminal Server\\WinStations\\RDP-Tcp' -Name 'PortNumber' | Select-Object -ExpandProperty PortNumber\""],
            "expectedValue": "3389", "matchType": "not_contains",
            "fixSuggestion": "建议更改远程桌面默认端口(3389)为其他端口",
            "riskDescription": "中危风险: 远程桌面使用默认端口3389，建议修改"
        },
        {
            "name": "检查ICMP重定向",
            "description": "检查Windows系统是否禁用了ICMP重定向",
            "category": 12, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters' -Name 'EnableICMPRedirect' | Select-Object -ExpandProperty EnableICMPRedirect\""],
            "expectedValue": "0", "matchType": "contains",
            "fixSuggestion": "请禁用ICMP重定向",
            "riskDescription": "中危风险: ICMP重定向未禁用"
        },
        {
            "name": "检查SMBv1协议",
            "description": "检查Windows系统是否禁用了SMBv1协议",
            "category": 12, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SYSTEM\\CurrentControlSet\\Services\\LanmanServer\\Parameters' -Name 'SMB1' | Select-Object -ExpandProperty SMB1\""],
            "expectedValue": "0", "matchType": "contains",
            "fixSuggestion": "请禁用SMBv1协议，使用SMBv2或更高版本",
            "riskDescription": "高危风险: SMBv1协议未禁用，存在严重安全漏洞"
        },
        {
            "name": "检查NetBIOS over TCP/IP",
            "description": "检查Windows系统是否禁用了NetBIOS over TCP/IP",
            "category": 12, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SYSTEM\\CurrentControlSet\\Services\\NetBT\\Parameters\\Interfaces\\*' -Name 'NetbiosOptions' | Select-Object -ExpandProperty NetbiosOptions\""],
            "expectedValue": "2", "matchType": "contains",
            "fixSuggestion": "请禁用NetBIOS over TCP/IP",
            "riskDescription": "中危风险: NetBIOS over TCP/IP未禁用"
        },
        
        # 服务配置 (category=13)
        {
            "name": "检查不必要的服务-Telnet",
            "description": "检查Windows Telnet服务是否已禁用",
            "category": 13, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-Service -Name 'Telnet' -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Status\""],
            "expectedValue": "NOT_FOUND", "matchType": "not_contains",
            "fixSuggestion": "请卸载或禁用Telnet服务，使用SSH替代",
            "riskDescription": "高危风险: Telnet服务未禁用，Telnet使用明文传输"
        },
        {
            "name": "检查不必要的服务-FTP",
            "description": "检查Windows FTP服务是否已禁用",
            "category": 13, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-Service -Name 'FTPSVC' -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Status\""],
            "expectedValue": "NOT_FOUND", "matchType": "not_contains",
            "fixSuggestion": "如非必要，请禁用FTP服务，使用SFTP替代",
            "riskDescription": "中危风险: FTP服务未禁用，FTP使用明文传输"
        },
        {
            "name": "检查不必要的服务-IIS",
            "description": "检查Windows IIS服务是否已禁用(如非Web服务器)",
            "category": 13, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-Service -Name 'W3SVC' -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Status\""],
            "expectedValue": "NOT_FOUND", "matchType": "not_contains",
            "fixSuggestion": "如非必要，请禁用IIS服务",
            "riskDescription": "中危风险: IIS服务未禁用"
        },
        {
            "name": "检查Windows自动更新",
            "description": "检查Windows自动更新是否已启用",
            "category": 13, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\WindowsUpdate\\Auto Update' -Name 'AUOptions' | Select-Object -ExpandProperty AUOptions\""],
            "expectedValue": "3", "matchType": "contains",
            "fixSuggestion": "请启用Windows自动更新，设置为自动下载并安装",
            "riskDescription": "中危风险: Windows自动更新未启用"
        },
        {
            "name": "检查Windows Defender状态",
            "description": "检查Windows Defender防病毒是否已启用",
            "category": 13, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-MpComputerStatus | Select-Object -ExpandProperty RealTimeProtectionEnabled\""],
            "expectedValue": "True", "matchType": "contains",
            "fixSuggestion": "请启用Windows Defender实时保护",
            "riskDescription": "高危风险: Windows Defender实时保护未启用"
        },
        {
            "name": "检查Windows Defender病毒库",
            "description": "检查Windows Defender病毒库是否最新",
            "category": 13, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-MpComputerStatus | Select-Object AntivirusSignatureLastUpdated\""],
            "expectedValue": "AntivirusSignatureLastUpdated", "matchType": "contains",
            "fixSuggestion": "请更新Windows Defender病毒库",
            "riskDescription": "中危风险: Windows Defender病毒库可能不是最新"
        },
        {
            "name": "检查远程注册表服务",
            "description": "检查Windows Remote Registry服务是否已禁用",
            "category": 13, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-Service -Name 'RemoteRegistry' | Select-Object -ExpandProperty Status\""],
            "expectedValue": "Stopped", "matchType": "contains",
            "fixSuggestion": "请禁用Remote Registry服务",
            "riskDescription": "高危风险: Remote Registry服务未禁用，允许远程修改注册表"
        },
        {
            "name": "检查Windows时间同步",
            "description": "检查Windows时间同步是否已配置",
            "category": 13, "risk": 3, "osType": 2,
            "commands": ["powershell -Command \"w32tm /query /configuration | Select-String 'NtpServer'\""],
            "expectedValue": "NtpServer", "matchType": "contains",
            "fixSuggestion": "请配置NTP时间同步服务器",
            "riskDescription": "低危风险: NTP时间同步未配置"
        },
        
        # 审计策略 (category=14)
        {
            "name": "检查审计策略-账户登录",
            "description": "检查Windows是否启用了账户登录审计",
            "category": 14, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-AuditPolicy -SubCategory 'Account Logon' | Select-Object -ExpandProperty IncludeSuccess\""],
            "expectedValue": "True", "matchType": "contains",
            "fixSuggestion": "请启用账户登录审计策略",
            "riskDescription": "中危风险: 账户登录审计未启用"
        },
        {
            "name": "检查审计策略-账户管理",
            "description": "检查Windows是否启用了账户管理审计",
            "category": 14, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-AuditPolicy -SubCategory 'Account Management' | Select-Object -ExpandProperty IncludeSuccess\""],
            "expectedValue": "True", "matchType": "contains",
            "fixSuggestion": "请启用账户管理审计策略",
            "riskDescription": "中危风险: 账户管理审计未启用"
        },
        {
            "name": "检查审计策略-登录事件",
            "description": "检查Windows是否启用了登录事件审计",
            "category": 14, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-AuditPolicy -SubCategory 'Logon' | Select-Object -ExpandProperty IncludeSuccess\""],
            "expectedValue": "True", "matchType": "contains",
            "fixSuggestion": "请启用登录事件审计策略",
            "riskDescription": "中危风险: 登录事件审计未启用"
        },
        {
            "name": "检查审计策略-对象访问",
            "description": "检查Windows是否启用了对象访问审计",
            "category": 14, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-AuditPolicy -SubCategory 'Object Access' | Select-Object -ExpandProperty IncludeSuccess\""],
            "expectedValue": "True", "matchType": "contains",
            "fixSuggestion": "请启用对象访问审计策略",
            "riskDescription": "中危风险: 对象访问审计未启用"
        },
        {
            "name": "检查审计策略-策略更改",
            "description": "检查Windows是否启用了策略更改审计",
            "category": 14, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-AuditPolicy -SubCategory 'Policy Change' | Select-Object -ExpandProperty IncludeSuccess\""],
            "expectedValue": "True", "matchType": "contains",
            "fixSuggestion": "请启用策略更改审计策略",
            "riskDescription": "中危风险: 策略更改审计未启用"
        },
        {
            "name": "检查审计策略-特权使用",
            "description": "检查Windows是否启用了特权使用审计",
            "category": 14, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-AuditPolicy -SubCategory 'Privilege Use' | Select-Object -ExpandProperty IncludeSuccess\""],
            "expectedValue": "True", "matchType": "contains",
            "fixSuggestion": "请启用特权使用审计策略",
            "riskDescription": "中危风险: 特权使用审计未启用"
        },
        {
            "name": "检查审计策略-进程跟踪",
            "description": "检查Windows是否启用了进程跟踪审计",
            "category": 14, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-AuditPolicy -SubCategory 'Detailed Tracking' | Select-Object -ExpandProperty IncludeSuccess\""],
            "expectedValue": "True", "matchType": "contains",
            "fixSuggestion": "请启用进程跟踪审计策略",
            "riskDescription": "中危风险: 进程跟踪审计未启用"
        },
        {
            "name": "检查审计日志大小",
            "description": "检查Windows安全日志最大大小",
            "category": 14, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SYSTEM\\CurrentControlSet\\Services\\EventLog\\Security' -Name 'MaxSize' | Select-Object -ExpandProperty MaxSize\""],
            "expectedValue": "209715200", "matchType": "contains",
            "fixSuggestion": "请设置安全日志最大大小为200MB以上",
            "riskDescription": "中危风险: 安全日志大小不足，可能导致日志丢失"
        },
        {
            "name": "检查审计日志保留策略",
            "description": "检查Windows安全日志保留策略",
            "category": 14, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SYSTEM\\CurrentControlSet\\Services\\EventLog\\Security' -Name 'Retention' | Select-Object -ExpandProperty Retention\""],
            "expectedValue": "0", "matchType": "contains",
            "fixSuggestion": "请设置安全日志按需覆盖，避免日志丢失",
            "riskDescription": "中危风险: 安全日志保留策略配置不当"
        },
        
        # 系统安全配置
        {
            "name": "检查UAC是否启用",
            "description": "检查Windows用户账户控制(UAC)是否已启用",
            "category": 11, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System' -Name 'EnableLUA' | Select-Object -ExpandProperty EnableLUA\""],
            "expectedValue": "1", "matchType": "contains",
            "fixSuggestion": "请启用UAC(用户账户控制)",
            "riskDescription": "高危风险: UAC未启用，系统容易受到权限提升攻击"
        },
        {
            "name": "检查屏幕保护程序超时",
            "description": "检查Windows屏幕保护程序超时设置",
            "category": 11, "risk": 3, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKCU:\\Control Panel\\Desktop' -Name 'ScreenSaveTimeOut' | Select-Object -ExpandProperty ScreenSaveTimeOut\""],
            "expectedValue": "900", "matchType": "contains",
            "fixSuggestion": "请设置屏幕保护程序超时为15分钟以内",
            "riskDescription": "低危风险: 屏幕保护程序超时过长"
        },
        {
            "name": "检查屏幕保护程序密码",
            "description": "检查Windows屏幕保护程序是否要求密码",
            "category": 11, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKCU:\\Control Panel\\Desktop' -Name 'ScreenSaverIsSecure' | Select-Object -ExpandProperty ScreenSaverIsSecure\""],
            "expectedValue": "1", "matchType": "contains",
            "fixSuggestion": "请设置屏幕保护程序恢复时需要密码",
            "riskDescription": "中危风险: 屏幕保护程序未设置密码保护"
        },
        {
            "name": "检查自动登录",
            "description": "检查Windows是否配置了自动登录",
            "category": 11, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Winlogon' -Name 'AutoAdminLogon' | Select-Object -ExpandProperty AutoAdminLogon\""],
            "expectedValue": "0", "matchType": "contains",
            "fixSuggestion": "请禁用自动登录功能",
            "riskDescription": "高危风险: 系统配置了自动登录，存在安全风险"
        },
        {
            "name": "检查关机事件跟踪",
            "description": "检查Windows是否启用了关机事件跟踪",
            "category": 11, "risk": 3, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System' -Name 'ShutdownReasonUI' | Select-Object -ExpandProperty ShutdownReasonUI\""],
            "expectedValue": "1", "matchType": "contains",
            "fixSuggestion": "请启用关机事件跟踪",
            "riskDescription": "低危风险: 关机事件跟踪未启用"
        },
        {
            "name": "检查LSASS保护",
            "description": "检查Windows LSASS进程是否启用了保护模式",
            "category": 11, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SYSTEM\\CurrentControlSet\\Control\\Lsa' -Name 'RunAsPPL' | Select-Object -ExpandProperty RunAsPPL\""],
            "expectedValue": "1", "matchType": "contains",
            "fixSuggestion": "请启用LSASS保护模式(RunAsPPL)",
            "riskDescription": "高危风险: LSASS未启用保护模式，容易受到凭据窃取攻击"
        },
        {
            "name": "检查匿名访问",
            "description": "检查Windows是否允许匿名枚举SAM账户",
            "category": 11, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SYSTEM\\CurrentControlSet\\Control\\Lsa' -Name 'RestrictAnonymousSAM' | Select-Object -ExpandProperty RestrictAnonymousSAM\""],
            "expectedValue": "1", "matchType": "contains",
            "fixSuggestion": "请禁止匿名枚举SAM账户",
            "riskDescription": "高危风险: 允许匿名枚举SAM账户"
        },
        {
            "name": "检查匿名枚举",
            "description": "检查Windows是否允许匿名枚举",
            "category": 11, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SYSTEM\\CurrentControlSet\\Control\\Lsa' -Name 'RestrictAnonymous' | Select-Object -ExpandProperty RestrictAnonymous\""],
            "expectedValue": "1", "matchType": "contains",
            "fixSuggestion": "请禁止匿名枚举",
            "riskDescription": "高危风险: 允许匿名枚举"
        },
        {
            "name": "检查C盘共享",
            "description": "检查Windows默认管理共享是否已禁用",
            "category": 12, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-SmbShare -Name 'C$' -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Name\""],
            "expectedValue": "C$", "matchType": "not_contains",
            "fixSuggestion": "如非必要，请禁用默认管理共享",
            "riskDescription": "中危风险: 默认管理共享(C$)存在"
        },
        {
            "name": "检查PowerShell执行策略",
            "description": "检查Windows PowerShell执行策略",
            "category": 11, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-ExecutionPolicy\""],
            "expectedValue": "Restricted", "matchType": "contains",
            "fixSuggestion": "请设置PowerShell执行策略为Restricted或RemoteSigned",
            "riskDescription": "中危风险: PowerShell执行策略过于宽松"
        },
        {
            "name": "检查WinRM服务状态",
            "description": "检查Windows WinRM服务是否已启用",
            "category": 13, "risk": 3, "osType": 2,
            "commands": ["powershell -Command \"Get-Service -Name 'WinRM' | Select-Object -ExpandProperty Status\""],
            "expectedValue": "Running", "matchType": "contains",
            "fixSuggestion": "WinRM服务用于远程管理，请确保已正确配置",
            "riskDescription": "信息: WinRM服务状态检查"
        },
        {
            "name": "检查WinRM认证方式",
            "description": "检查Windows WinRM是否启用了基本认证",
            "category": 12, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\WSMAN\\Client' -Name 'AllowBasic' | Select-Object -ExpandProperty AllowBasic\""],
            "expectedValue": "0", "matchType": "contains",
            "fixSuggestion": "请禁用WinRM基本认证，使用Kerberos认证",
            "riskDescription": "中危风险: WinRM启用了基本认证"
        },
        {
            "name": "检查BitLocker状态",
            "description": "检查Windows BitLocker驱动器加密状态",
            "category": 5, "risk": 2, "osType": 2,
            "commands": ["powershell -Command \"Get-BitLockerVolume -MountPoint 'C:' -ErrorAction SilentlyContinue | Select-Object -ExpandProperty ProtectionStatus\""],
            "expectedValue": "1", "matchType": "contains",
            "fixSuggestion": "请启用BitLocker驱动器加密",
            "riskDescription": "中危风险: 系统盘未启用BitLocker加密"
        },
        {
            "name": "检查Windows版本",
            "description": "检查Windows操作系统版本信息",
            "category": 8, "risk": 3, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion' | Select-Object ProductName,ReleaseId,CurrentBuild\""],
            "expectedValue": "ProductName", "matchType": "contains",
            "fixSuggestion": "请确保使用受支持的Windows版本",
            "riskDescription": "信息: Windows版本信息"
        },
        {
            "name": "检查已安装的更新",
            "description": "检查Windows系统最近安装的安全更新",
            "category": 8, "risk": 1, "osType": 2,
            "commands": ["powershell -Command \"Get-HotFix | Sort-Object InstalledOn -Descending | Select-Object -First 10 HotFixID,InstalledOn\""],
            "expectedValue": "HotFixID", "matchType": "contains",
            "fixSuggestion": "请确保系统已安装最新的安全更新",
            "riskDescription": "高危风险: 请检查系统是否缺少安全更新"
        },
        {
            "name": "检查Windows激活状态",
            "description": "检查Windows系统是否已激活",
            "category": 8, "risk": 3, "osType": 2,
            "commands": ["powershell -Command \"Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion' -Name 'ProductId' | Select-Object -ExpandProperty ProductId\""],
            "expectedValue": "ProductId", "matchType": "contains",
            "fixSuggestion": "请确保Windows系统已激活",
            "riskDescription": "信息: Windows激活状态"
        },
    ]
    
    for i, r in enumerate(windows_rules_data):
        rule = {
            "id": 200000 + i,
            "name": r["name"],
            "description": r["description"],
            "category": r["category"],
            "risk": r["risk"],
            "osType": 2,
            "commands": r["commands"],
            "expectedValue": r["expectedValue"],
            "matchType": r["matchType"],
            "fixSuggestion": r["fixSuggestion"],
            "riskDescription": r["riskDescription"],
        }
        rules.append(rule)
    
    return rules


if __name__ == '__main__':
    main()