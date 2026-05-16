import json
import re

with open('data/baseline/compliance_rules.json', 'r', encoding='utf-8') as f:
    rules = json.load(f)

# ============================================================
# 名称翻译 - 综合处理
# ============================================================
def translate_name(name):
    # ========== 1. 精确匹配 ==========
    exact = {
        "Ensure auditd Collects Information on the Use of Privileged Commands": "确保 auditd 收集特权命令使用信息",
        "Ensure auditd Collects System Administrator Actions": "确保 auditd 收集系统管理员操作",
        "Ensure auditd Collects Information on Kernel Module Loading and Unloading": "确保 auditd 收集内核模块加载/卸载信息",
        "Ensure auditd Collects Information on Kernel Module Loading": "确保 auditd 收集内核模块加载信息",
        "Ensure auditd Collects Information on Kernel Module Unloading": "确保 auditd 收集内核模块卸载信息",
        "Ensure auditd Collects Information on Exporting to Media (successful)": "确保 auditd 收集导出到介质的信息（成功）",
        "Ensure auditd Collects Information on Exporting to Media (unsuccessful)": "确保 auditd 收集导出到介质的信息（失败）",
        "Ensure auditd Collects File Deletion Events by User (successful)": "确保 auditd 收集用户文件删除事件（成功）",
        "Ensure auditd Collects File Deletion Events by User (unsuccessful)": "确保 auditd 收集用户文件删除事件（失败）",
        "Ensure auditd Collects Unauthorized Access Attempts to Files (unsuccessful)": "确保 auditd 收集文件未授权访问尝试（失败）",
        "Record Events that Modify the System's Discretionary Access Controls": "记录自主访问控制修改事件",
        "Record Unsuccessful Permission Changes to Files": "记录文件权限修改失败事件",
        "Record Unsuccessful Ownership Changes to Files": "记录文件所有权修改失败事件",
        "Record Unsuccessful Access Attempts to Files": "记录文件访问失败事件",
        "Record Unsuccessful Creation Attempts to Files": "记录文件创建失败事件",
        "Record Unsuccessful Modification Attempts to Files": "记录文件修改失败事件",
        "Record Unsuccessful Delete Attempts to Files": "记录文件删除失败事件",
        "Record Attempts to Alter Logon and Logout Events": "记录登录/注销事件修改尝试",
        "Record Any Attempts to Run": "记录运行尝试",
        "Make the auditd Configuration Immutable": "确保 auditd 配置不可变",
        "Shutdown System When Auditing Failures Occur": "审计失败时关闭系统",
        "Encrypt Audit Records Sent With audispd Plugin": "加密通过 audispd 插件发送的审计记录",
        "Include Local Events in Audit Logs": "在审计日志中包含本地事件",
        "Resolve information before writing to audit logs": "写入审计日志前解析信息",
        "Write Audit Logs to the Disk": "将审计日志写入磁盘",
        "Create Warning Banners for All FTP Users": "为所有 FTP 用户创建警告横幅",
        "Require Client Certificates": "要求客户端证书",
        "Prevent Unrestricted Mail Relaying": "防止无限制邮件中继",
        "Require Client SMB Packet Signing, if using mount.cifs": "使用 mount.cifs 时要求客户端 SMB 数据包签名",
        "Require Client SMB Packet Signing, if using smbclient": "使用 smbclient 时要求客户端 SMB 数据包签名",
        "Allow Only SSH Protocol 2": "仅允许 SSH 协议 2",
        "Do Not Allow SSH Environment Options": "不允许 SSH 环境选项",
        "Limit Users' SSH Access": "限制用户的 SSH 访问",
        "Force frequent session key renegotiation": "强制频繁会话密钥重新协商",
        "Use Only Strong Ciphers": "仅使用强加密算法",
        "Use Only Strong Key Exchange algorithms": "仅使用强密钥交换算法",
        "Use Only Strong MACs": "仅使用强 MAC 算法",
        "Prevent remote hosts from connecting to the proxy display": "防止远程主机连接到代理显示",
        "Certificate status checking in SSSD": "SSSD 中的证书状态检查",
        "Certificate trust path in SSSD": "SSSD 中的证书信任路径",
        "Generate USBGuard Policy": "生成 USBGuard 策略",
        "Modify the System Login Banner": "修改系统登录横幅",
        "Modify the System Login Banner for Remote Connections": "修改远程连接的系统登录横幅",
        "Modify the System Message of the Day Banner": "修改系统每日消息横幅",
        "Modify the System GUI Login Banner": "修改系统 GUI 登录横幅",
        "Disallow Configuration to Bypass Password Requirements for Privilege Escalation": "禁止绕过权限提升的密码要求配置",
        "Enforce Delay After Failed Logon Attempts": "登录失败后强制延迟",
        "Lock Accounts After Failed Password Attempts": "密码尝试失败后锁定账户",
        "Lock Accounts Must Persist": "账户锁定必须持久化",
        "Enforce pam_faillock for Local Accounts Only": "仅对本地账户强制使用 pam_faillock",
        "Do Not Show System Messages When Unsuccessful Logon Attempts Occur": "登录失败时不显示系统消息",
        "Limit Password Reuse": "限制密码重用",
        "Limit Password Reuse: password-auth": "限制密码重用：password-auth",
        "Limit Password Reuse: system-auth": "限制密码重用：system-auth",
        "Enforce Password History with use_authtok": "使用 use_authtok 强制密码历史",
        "Require use_authtok for pam_unix.so": "要求 pam_unix.so 使用 use_authtok",
        "Require Authentication for Emergency Systemd Target": "要求紧急 Systemd 目标进行身份验证",
        "Require Authentication for Single User Mode": "要求单用户模式进行身份验证",
        "Support session locking with tmux": "支持使用 tmux 的会话锁定",
        "Support session locking with tmux (not enforcing)": "支持使用 tmux 的会话锁定（非强制）",
        "Prevent user from disabling the screen lock": "防止用户禁用屏幕锁定",
        "Force opensc To Use Defined Smart Card Driver": "强制 opensc 使用定义的智能卡驱动程序",
        "Never Automatically Remove or Disable Emergency Administrator Accounts": "永不自动删除或禁用紧急管理员账户",
        "Policy Requires Immediate Change of Temporary Passwords": "策略要求立即更改临时密码",
        "Avoid using remember in pam_unix module": "避免在 pam_unix 模块中使用 remember",
        "Prevent Login to Accounts With Empty Password": "防止使用空密码登录账户",
        "Direct root Logins Not Allowed": "不允许直接 root 登录",
        "Restrict Web Browser Use for Administrative Accounts": "限制管理账户的 Web 浏览器使用",
        "Direct root Logins Are Not Allowed": "不允许直接 root 登录",
        "Restrict Serial Port Root Logins": "限制串行端口 root 登录",
        "Root Path Must Be Vendor Default": "Root 路径必须为供应商默认值",
        "Restrict Virtual Console Root Logins": "限制虚拟控制台 root 登录",
        "Enforce usage of pam_wheel for su authentication": "强制使用 pam_wheel 进行 su 身份验证",
        "Enforce Usage of pam_wheel with Group Parameter for su Authentication": "强制使用带组参数的 pam_wheel 进行 su 身份验证",
        "Limit the Number of Concurrent Login Sessions Allowed Per User": "限制每个用户允许的并发登录会话数",
        "Enforce all AppArmor Profiles": "强制所有 AppArmor 配置文件",
        "IOMMU configuration directive": "IOMMU 配置指令",
        "Force kernel panic on uncorrected MCEs": "在不可纠正的 MCE 上强制内核恐慌",
        "Enforce Spectre v2 mitigation": "强制 Spectre v2 缓解措施",
        "Boot Loader Is Not Installed On Removable Media": "引导加载程序未安装在可移动介质上",
        "UEFI Boot Loader Is Not Installed On Removable Media": "UEFI 引导加载程序未安装在可移动介质上",
        "Randomize layout of sensitive kernel structures": "随机化敏感内核结构的布局",
        "Force initialization of variables containing userspace addresses": "强制初始化包含用户空间地址的变量",
        "zero-init everything passed by reference": "零初始化所有通过引用传递的内容",
        "Do not allow ACPI methods to be inserted/replaced at run time": "不允许在运行时插入/替换 ACPI 方法",
        "Emulate Privileged Access Never (PAN)": "模拟特权访问从不（PAN）",
        "Trigger a kernel BUG when data corruption is detected": "检测到数据损坏时触发内核 BUG",
        "Warn on W+X mappings found at boot": "在启动时发现 W+X 映射时发出警告",
        "Harden common str/mem functions against buffer overflows": "加固常用 str/mem 函数以防止缓冲区溢出",
        "Harden memory copies between kernel and userspace": "加固内核与用户空间之间的内存复制",
        "Do not allow usercopy whitelist violations to fallback to object size": "不允许用户复制白名单违规回退到对象大小",
        "Require modules to be validly signed": "要求模块有效签名",
        "Specify the hash to use when signing modules": "指定签名模块时使用的哈希算法",
        "Specify module signing key to use": "指定要使用的模块签名密钥",
        "Sign kernel modules with SHA-512": "使用 SHA-512 签名内核模块",
        "Use zero for poisoning instead of debugging value": "使用零进行毒化而不是调试值",
        "Kernel panic oops": "内核恐慌 oops",
        "Kernel panic timeout": "内核恐慌超时",
        "Randomize the address of the kernel image (KASLR)": "随机化内核映像地址（KASLR）",
        "Randomize the kernel memory sections": "随机化内核内存段",
        "Perform full reference count validation": "执行完整的引用计数验证",
        "Avoid speculative indirect branches in kernel": "避免内核中的推测性间接分支",
        "Detect stack corruption on calls to schedule()": "检测 schedule() 调用时的栈损坏",
        "Restrict unprivileged access to the kernel syslog": "限制非特权访问内核系统日志",
        "Harden slab freelist metadata": "加固 slab 空闲列表元数据",
        "Randomize slab freelist": "随机化 slab 空闲列表",
        "Disallow merge of slab caches": "禁止合并 slab 缓存",
        "Stack Protector buffer overflow detection": "栈保护缓冲区溢出检测",
        "Strong Stack Protector": "强栈保护",
        "Make the kernel text and rodata read-only": "确保内核文本和 rodata 只读",
        "Make the module text and rodata read-only": "确保模块文本和 rodata 只读",
        "Unmap kernel when running in userspace (aka KAISER)": "在用户空间运行时取消映射内核（KAISER）",
        "User a virtually-mapped stack": "使用虚拟映射栈",
        "Use Privacy Extensions for Address": "使用地址隐私扩展",
        "Drop Gratuitous ARP frames on All IPv4 Interfaces": "丢弃所有 IPv4 接口上的 gratuitous ARP 帧",
        "Prevent Routing External Traffic to Local Loopback on All IPv4 Interfaces": "防止在所有 IPv4 接口上将外部流量路由到本地回环",
        "Deactivate Wireless Network Interfaces": "停用无线网络接口",
        "Prevent non-Privileged Users from Modifying Network Interfaces using nmcli": "防止非特权用户使用 nmcli 修改网络接口",
        "Restrict Access to Kernel Message Buffer": "限制访问内核消息缓冲区",
        "Kernel panic on oops": "oops 时内核恐慌",
        "Limit CPU consumption of the Perf system": "限制 Perf 系统的 CPU 消耗",
        "Limit sampling frequency of the Perf system": "限制 Perf 系统的采样频率",
        "Disallow kernel profiling by unprivileged users": "禁止非特权用户进行内核性能分析",
        "Disallow magic SysRq key": "禁止魔术 SysRq 键",
        "Harden the operation of the BPF just-in-time compiler": "加固 BPF 即时编译器的操作",
        "Prevent applications from mapping low portion of virtual memory": "防止应用程序映射虚拟内存低位部分",
        "Elevate The SELinux Context When An Administrator Calls The Sudo Command": "管理员调用 sudo 命令时提升 SELinux 上下文",
        "Require Credential Prompting for Remote Access in GNOME3": "要求在 GNOME3 中远程访问时提示凭据",
        "Require Encryption for Remote Access in GNOME3": "要求在 GNOME3 中远程访问时加密",
        "Implement Blank Screensaver": "实现空白屏幕保护程序",
        "The Installed Operating System Is FIPS 140-2 Certified": "已安装的操作系统为 FIPS 140-2 认证",
        "The Installed Operating System Is Vendor Supported": "已安装的操作系统受供应商支持",
        "Harden OpenSSL Crypto Policy": "加固 OpenSSL 加密策略",
        "Harden SSHD Crypto Policy": "加固 SSHD 加密策略",
        "Harden SSH client Crypto Policy": "加固 SSH 客户端加密策略",
        "OpenSSL uses strong entropy source": "OpenSSL 使用强熵源",
        "Virus Scanning Software Definitions Are Updated": "病毒扫描软件定义已更新",
        "System Wide Crypto Policy Files Must Point to FIPS Policy": "系统范围加密策略文件必须指向 FIPS 策略",
        'Package "prelink" Must not be Installed': '软件包 "prelink" 不得安装',
        "Build and Test AIDE Database": "构建并测试 AIDE 数据库",
        "Audit Tools Must Be Group-owned by Root": "审计工具必须由 Root 组拥有",
        "Audit Tools Must Be Owned by Root": "审计工具必须由 Root 拥有",
        "Audit Tools Must Have a Mode of 0755 or Less Permissive": "审计工具必须具有 0755 或更严格的权限模式",
        "Package glibc Installed": "已安装 glibc 软件包",
        "Package uuidd Installed": "已安装 uuidd 软件包",
        "Explicit arguments in sudo specifications": "sudo 规范中的显式参数",
        "Don't define allowed commands in sudoers by means of exclusion": "不要通过排除方式在 sudoers 中定义允许的命令",
        "Don't target root user in the sudoers file": "不要在 sudoers 文件中以 root 用户为目标",
        "Require Re-Authentication When Using the sudo Command": "使用 sudo 命令时要求重新认证",
        "The operating system must restrict privilege escalation to authorized personnel": "操作系统必须限制权限提升仅限于授权人员",
        "Only the VDSM User Can Use sudo NOPASSWD": "只有 VDSM 用户可以使用 sudo NOPASSWD",
        "Chrony Configure Pool and Server": "Chrony 配置池和服务器",
        "Specify Additional Remote NTP Servers": "指定额外的远程 NTP 服务器",
        "Specify a Remote NTP Server": "指定远程 NTP 服务器",
        "A remote time server for Chrony is configured": "已为 Chrony 配置远程时间服务器",
        "Synchronize internal information system clocks": "同步内部信息系统时钟",
        "Name Service Switch does not use NIS": "名称服务切换不使用 NIS",
        "The File placeholder_value Must Exist": "文件 placeholder_value 必须存在",
        "SSH client uses strong entropy to seed (for CSH like shells)": "SSH 客户端使用强熵种子（适用于 CSH 类 Shell）",
        "SSH client uses strong entropy to seed (Bash-like shells)": "SSH 客户端使用强熵种子（适用于 Bash 类 Shell）",
        "SSH server uses strong entropy to seed": "SSH 服务器使用强熵种子",
        "Distribute the SSH Server configuration to multiple files in a config directory.": "将 SSH 服务器配置分发到配置目录中的多个文件",
        "Log USBGuard daemon audit events using Linux Audit": "使用 Linux Audit 记录 USBGuard 守护进程审计事件",
        "Authorize Human Interface Devices in USBGuard daemon": "授权 USBGuard 守护进程中的人机接口设备",
        "Authorize Human Interface Devices and USB hubs in USBGuard daemon": "授权 USBGuard 守护进程中的人机接口设备和 USB 集线器",
        "Authorize USB hubs in USBGuard daemon": "授权 USBGuard 守护进程中的 USB 集线器",
        "SLEM 5 must use the default pam_tally2 tally directory.": "SLEM 5 必须使用默认的 pam_tally2 计数目录",
        "The PAM configuration should not be changed automatically": "PAM 配置不应自动更改",
        "Limit the maximum number of sequential characters in passwords": "限制密码中连续字符的最大数量",
        "Only Authorized Local User Accounts Exist on Operating System": "操作系统上仅存在授权的本地用户账户",
        "Only sidadm and orasid/oracle User Accounts Exist on Operating System": "操作系统上仅存在 sidadm 和 orasid/oracle 用户账户",
        "Display the Standard Mandatory DoD Notice and Consent Banner until Explicit Acknowledgement": "显示标准强制性 DoD 通知和同意横幅直至明确确认",
        "NetworkManager DNS Mode Must Be Must Configured": "NetworkManager DNS 模式必须配置",
        "ufw Must rate-limit network interfaces": "ufw 必须对网络接口进行速率限制",
        "Bind Mount /var/tmp To /tmp": "将 /var/tmp 绑定挂载到 /tmp",
        "OS commands and libraries must have the proper permissions to protect from unauthorized access": "操作系统命令和库必须具有适当的权限以防止未授权访问",
        "The system must booted with init_on_free=1": "系统必须使用 init_on_free=1 引导",
        "System Must Avoid Meltdown and Spectre Exploit Vulnerabilities in Modern Processors": "系统必须避免现代处理器中的 Meltdown 和 Spectre 漏洞利用",
        "All AppArmor Profiles are in enforce or complain mode": "所有 AppArmor 配置文件处于强制或投诉模式",
        "Use Only FIPS 140-3 Validated Ciphers in SSH Client Configuration": "在 SSH 客户端配置中仅使用 FIPS 140-3 验证的加密算法",
        "Use Only FIPS 140-3 Validated MACs": "仅使用 FIPS 140-3 验证的 MAC",
        "Use Only FIPS 140-2 Validated Ciphers": "仅使用 FIPS 140-2 验证的加密算法",
        "Use Only FIPS 140-2 Validated MACs": "仅使用 FIPS 140-2 验证的 MAC",
        "Configure SSH to use System Crypto Policy": "配置 SSH 使用系统加密策略",
        "Configure SSH Server to Use FIPS 140-2 Validated Ciphers: opensshserver.config": "配置 SSH 服务器使用 FIPS 140-2 验证的加密算法：opensshserver.config",
        "Configure SSH Client to Use FIPS 140 Validated Ciphers: openssh.config": "配置 SSH 客户端使用 FIPS 140 验证的加密算法：openssh.config",
        "Configure SSH Server to Use FIPS 140-2 Validated MACs: opensshserver.config": "配置 SSH 服务器使用 FIPS 140-2 验证的 MAC：opensshserver.config",
        "Configure SSH Client to Use FIPS 140-2 Validated MACs: openssh.config": "配置 SSH 客户端使用 FIPS 140-2 验证的 MAC：openssh.config",
        'Verify the system-wide library files in directories\n"/lib", "/lib64", "/usr/lib/" and "/usr/lib64" are group-owned by root.':
            '验证系统库文件目录 "/lib"、"/lib64"、"/usr/lib/" 和 "/usr/lib64" 由 root 组拥有',
        "Enable Dracut FIPS Module": "启用 Dracut FIPS 模块",
        "Enable FIPS Mode": "启用 FIPS 模式",
        "Enable FIPS Mode in GRUB2": "在 GRUB2 中启用 FIPS 模式",
        "Install the dracut-fips-aesni Package": "安装 dracut-fips-aesni 软件包",
        "Install the dracut-fips Package": "安装 dracut-fips 软件包",
        "Install the Host Intrusion Prevention System (HIPS) Module": "安装主机入侵防御系统（HIPS）模块",
        "Install dnf-plugin-subscription-manager Package": "安装 dnf-plugin-subscription-manager 软件包",
        "Install libdnf-plugin-subscription-manager Package": "安装 libdnf-plugin-subscription-manager 软件包",
        "Install subscription-manager Package": "安装 subscription-manager 软件包",
        "Uninstall abrt-addon-kerneloops Package": "卸载 abrt-addon-kerneloops 软件包",
        "Uninstall abrt-plugin-logger Package": "卸载 abrt-plugin-logger 软件包",
        "Uninstall iprutils Package": "卸载 iprutils 软件包",
        "Uninstall libreport-plugin-logger Package": "卸载 libreport-plugin-logger 软件包",
        "Configure AIDE to Use FIPS 140-2 for Validating Hashes": "配置 AIDE 使用 FIPS 140-2 验证哈希",
        "Configure SELinux Policy": "配置 SELinux 策略",
        "Ensure No Device Files are Unlabeled by SELinux": "确保没有设备文件被 SELinux 标记为未标记",
        "Ensure No Daemons are Unconfined by SELinux": "确保没有守护进程被 SELinux 标记为不受限制",
        "Ensure SELinux is Not Disabled": "确保 SELinux 未被禁用",
        "Ensure /boot Located On Separate Partition": "确保 /boot 位于独立分区",
        "Ensure /dev/shm is configured": "确保 /dev/shm 已配置",
        "Ensure /home Located On Separate Partition": "确保 /home 位于独立分区",
        "Ensure /opt Located On Separate Partition": "确保 /opt 位于独立分区",
        "Ensure /srv Located On Separate Partition": "确保 /srv 位于独立分区",
        "Ensure /tmp Located On Separate Partition": "确保 /tmp 位于独立分区",
        "Ensure /usr Located On Separate Partition": "确保 /usr 位于独立分区",
        "Ensure /var Located On Separate Partition": "确保 /var 位于独立分区",
        "Ensure /var/log Located On Separate Partition": "确保 /var/log 位于独立分区",
        "Ensure /var/log/audit Located On Separate Partition": "确保 /var/log/audit 位于独立分区",
        "Ensure /var/tmp Located On Separate Partition": "确保 /var/tmp 位于独立分区",
        "Ensure tmp.mount Unit Is Enabled": "确保 tmp.mount 单元已启用",
        "Ensure '/etc/system-fips' exists": "确保 '/etc/system-fips' 存在",
        "Verify '/proc/sys/crypto/fips_enabled' exists": "验证 '/proc/sys/crypto/fips_enabled' 存在",
        "Set kernel parameter 'crypto.fips_enabled' to 1": "设置内核参数 'crypto.fips_enabled' 为 1",
        "Set the GNOME3 Login Number of Failures": "设置 GNOME3 登录失败次数",
        "Disable GDM Automatic Login": "禁用 GDM 自动登录",
        "Disable GDM Guest Login": "禁用 GDM 访客登录",
        "Disable GDM Unattended or Automatic Login": "禁用 GDM 无人值守或自动登录",
        "Disable GNOME3 Automounting": "禁用 GNOME3 自动挂载",
        "Disable GNOME3 Automount Opening": "禁用 GNOME3 自动挂载打开",
        "Disable Ctrl-Alt-Del Reboot Key Sequence in GNOME3": "禁用 GNOME3 中的 Ctrl-Alt-Del 重启键序列",
        "Verify and Correct Ownership with RPM": "使用 RPM 验证并更正所有权",
        "Verify and Correct File Permissions with RPM": "使用 RPM 验证并更正文件权限",
        "The mailx Package Is Installed": "已安装 mailx 软件包",
        "The Postfix package is installed": "已安装 Postfix 软件包",
        "The s-nail Package Is Installed": "已安装 s-nail 软件包",
        "The Chrony package is installed": "已安装 Chrony 软件包",
        "The Chronyd service is disabled": "Chronyd 服务已禁用",
        "The Chronyd service is enabled": "Chronyd 服务已启用",
        "MIME types for csh or sh shell programs must be disabled": "csh 或 sh shell 程序的 MIME 类型必须禁用",
        "Backup interactive scripts on the production web server are prohibited": "禁止在生产 Web 服务器上备份交互式脚本",
        "Each Web Content Directory Must Contain An index.html File": "每个 Web 内容目录必须包含 index.html 文件",
        "Web Content Directories Must Not Be Shared Anonymously": "Web 内容目录不得匿名共享",
        "Mount Remote Filesystems with Kerberos Security": "使用 Kerberos 安全挂载远程文件系统",
        "Mount Remote Filesystems with nodev": "使用 nodev 挂载远程文件系统",
        "Mount Remote Filesystems with noexec": "使用 noexec 挂载远程文件系统",
        "Mount Remote Filesystems with nosuid": "使用 nosuid 挂载远程文件系统",
        "Use Kerberos Security on All Exports": "在所有导出上使用 Kerberos 安全",
        "Extend Audit Backlog Limit for the Audit Daemon in zIPL": "在 zIPL 中扩展审计守护进程的审计积压限制",
    }
    if name in exact:
        return exact[name]

    # ========== 2. 前缀匹配 ==========
    prefixes = [
        ("Ensure auditd Collects Information on the Use of Privileged Commands - ", "确保 auditd 收集特权命令使用信息 - "),
        ("Ensure auditd Collects System Administrator Actions - ", "确保 auditd 收集系统管理员操作 - "),
        ("Ensure auditd Collects Information on Kernel Module Loading and Unloading - ", "确保 auditd 收集内核模块加载/卸载信息 - "),
        ("Ensure auditd Collects Information on Kernel Module Loading - ", "确保 auditd 收集内核模块加载信息 - "),
        ("Ensure auditd Collects Information on Kernel Module Unloading - ", "确保 auditd 收集内核模块卸载信息 - "),
        ("Record Events that Modify the System's Discretionary Access Controls - ", "记录自主访问控制修改事件 - "),
        ("Record Unsuccessful Permission Changes to Files - ", "记录文件权限修改失败事件 - "),
        ("Record Unsuccessful Ownership Changes to Files - ", "记录文件所有权修改失败事件 - "),
        ("Record Unsuccessful Access Attempts to Files - ", "记录文件访问失败事件 - "),
        ("Record Unsuccessful Creation Attempts to Files - ", "记录文件创建失败事件 - "),
        ("Record Unsuccessful Modification Attempts to Files - ", "记录文件修改失败事件 - "),
        ("Record Unsuccessful Delete Attempts to Files - ", "记录文件删除失败事件 - "),
        ("Record Attempts to Alter Logon and Logout Events - ", "记录登录/注销事件修改尝试 - "),
        ("Record Any Attempts to Run ", "记录运行尝试 - "),
        ("Ensure auditd Rules For Unauthorized Attempts To ", "确保 auditd 规则对未授权尝试 "),
        ("Ensure auditd Unauthorized Access Attempts To ", "确保 auditd 对未授权访问尝试 "),
    ]
    for eng, cn in prefixes:
        if name.startswith(eng):
            suffix = name[len(eng):]
            if "Are Ordered Correctly" in suffix:
                suffix = suffix.replace("Are Ordered Correctly", "").strip()
                if suffix.endswith("To"):
                    suffix = suffix[:-2].strip()
                return cn + suffix + " 正确排序"
            return cn + suffix

    # ========== 3. 正则模式匹配 ==========

    # --- "Must" 系列模式 ---
    # "X Must Be Group-Owned By The Y" -> "X 必须由 Y 组拥有"
    m = re.match(r'^(.+?)\s+Must\s+Be\s+Group-Owned\s+By\s+The\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须由 {m.group(2)} 组拥有"
    # "X Must Be Group-Owned By Y" -> "X 必须由 Y 组拥有"
    m = re.match(r'^(.+?)\s+Must\s+Be\s+Group-Owned\s+By\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须由 {m.group(2)} 组拥有"
    # "X Must Be Group-owned by Y" -> "X 必须由 Y 组拥有"
    m = re.match(r'^(.+?)\s+Must\s+Be\s+Group-owned\s+by\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须由 {m.group(2)} 组拥有"
    # "X Must Be Owned By The Y" -> "X 必须由 Y 拥有"
    m = re.match(r'^(.+?)\s+Must\s+Be\s+Owned\s+By\s+The\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须由 {m.group(2)} 拥有"
    # "X Must Be Owned By Y" -> "X 必须由 Y 拥有"
    m = re.match(r'^(.+?)\s+Must\s+Be\s+Owned\s+By\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须由 {m.group(2)} 拥有"
    # "X Must Be Owned by Y" -> "X 必须由 Y 拥有"
    m = re.match(r'^(.+?)\s+Must\s+Be\s+Owned\s+by\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须由 {m.group(2)} 拥有"
    # "X Must Have Mode Y or Less Permissive" -> "X 必须具有模式 Y 或更严格"
    m = re.match(r'^(.+?)\s+Must\s+Have\s+Mode\s+(.+?)\s+or\s+Less\s+Permissive$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须具有模式 {m.group(2)} 或更严格"
    # "X Must Have mode Y or Less Permissive" -> "X 必须具有模式 Y 或更严格"
    m = re.match(r'^(.+?)\s+Must\s+Have\s+mode\s+(.+?)\s+or\s+Less\s+Permissive$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须具有模式 {m.group(2)} 或更严格"
    # "X Must Have Mode Y Or Less Permissive" -> "X 必须具有模式 Y 或更严格"
    m = re.match(r'^(.+?)\s+Must\s+Have\s+Mode\s+(.+?)\s+Or\s+Less\s+Permissive$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须具有模式 {m.group(2)} 或更严格"
    # "X Must Have a Mode of Y or Less Permissive" -> "X 必须具有模式 Y 或更严格"
    m = re.match(r'^(.+?)\s+Must\s+Have\s+a\s+Mode\s+of\s+(.+?)\s+or\s+Less\s+Permissive$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须具有模式 {m.group(2)} 或更严格"
    # "X Must Have a Valid Y" -> "X 必须具有有效的 Y"
    m = re.match(r'^(.+?)\s+Must\s+Have\s+a\s+Valid\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须具有有效的 {m.group(2)}"
    # "X Must Have A Y Defined" -> "X 必须定义 Y"
    m = re.match(r'^(.+?)\s+Must\s+Have\s+A\s+(.+?)\s+Defined$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须定义 {m.group(2)}"
    # "X Must Have A Y" -> "X 必须具有 Y"
    m = re.match(r'^(.+?)\s+Must\s+Have\s+A\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须具有 {m.group(2)}"
    # "X Must Have Y" -> "X 必须具有 Y"
    m = re.match(r'^(.+?)\s+Must\s+Have\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须具有 {m.group(2)}"
    # "X Must Not Be Shared Y" -> "X 不得 Y 共享"
    m = re.match(r'^(.+?)\s+Must\s+Not\s+Be\s+Shared\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 不得{m.group(2)}共享"
    # "X Must Not Be Y" -> "X 不得为 Y"
    m = re.match(r'^(.+?)\s+Must\s+Not\s+Be\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 不得为 {m.group(2)}"
    # "X Must Not Run Y" -> "X 不得运行 Y"
    m = re.match(r'^(.+?)\s+Must\s+Not\s+Run\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 不得运行 {m.group(2)}"
    # "X Must Not Y" -> "X 不得 Y"
    m = re.match(r'^(.+?)\s+Must\s+Not\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 不得 {m.group(2)}"
    # "X Must Contain Y" -> "X 必须包含 Y"
    m = re.match(r'^(.+?)\s+Must\s+Contain\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须包含 {m.group(2)}"
    # "X Must Point to Y" -> "X 必须指向 Y"
    m = re.match(r'^(.+?)\s+Must\s+Point\s+to\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须指向 {m.group(2)}"
    # "X Must rate-limit Y" -> "X 必须对 Y 进行速率限制"
    m = re.match(r'^(.+?)\s+Must\s+rate-limit\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须对 {m.group(2)} 进行速率限制"
    # "X Must Reside On Y" -> "X 必须位于 Y"
    m = re.match(r'^(.+?)\s+Must\s+Reside\s+On\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须位于 {m.group(2)}"
    # "X Must Exist" -> "X 必须存在"
    m = re.match(r'^(.+?)\s+Must\s+Exist$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须存在"
    # "X Must Be Y" -> "X 必须为 Y"
    m = re.match(r'^(.+?)\s+Must\s+Be\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须为 {m.group(2)}"
    # "X Must Y" -> "X 必须 Y"
    m = re.match(r'^(.+?)\s+Must\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{m.group(1)} 必须 {m.group(2)}"

    # --- "The X Is/are Y" 模式 ---
    m = re.match(r'^The\s+(.+?)\s+(Is|is|are)\s+(.+)$', name)
    if m: return f"{m.group(1)} 已{m.group(3)}"

    # --- "X are Y" 模式 ---
    m = re.match(r'^(.+?)\s+are\s+(.+)$', name)
    if m: return f"{m.group(1)} 为 {m.group(2)}"

    # --- "X is Y" 模式 ---
    m = re.match(r'^(.+?)\s+is\s+(.+)$', name)
    if m: return f"{m.group(1)} 为 {m.group(2)}"

    # --- "X uses Y" 模式 ---
    m = re.match(r'^(.+?)\s+uses\s+(.+)$', name)
    if m: return f"{m.group(1)} 使用 {m.group(2)}"

    # --- "X does not use Y" 模式 ---
    m = re.match(r'^(.+?)\s+does\s+not\s+use\s+(.+)$', name)
    if m: return f"{m.group(1)} 不使用 {m.group(2)}"

    # --- "X should not be Y" 模式 ---
    m = re.match(r'^(.+?)\s+should\s+not\s+be\s+(.+)$', name)
    if m: return f"{m.group(1)} 不应{m.group(2)}"

    # --- "X should be Y" 模式 ---
    m = re.match(r'^(.+?)\s+should\s+be\s+(.+)$', name)
    if m: return f"{m.group(1)} 应{m.group(2)}"

    # --- "X Not Allowed" / "X Are Not Allowed" 模式 ---
    m = re.match(r'^(.+?)\s+(Are\s+)?Not\s+Allowed$', name, re.IGNORECASE)
    if m: return f"不允许{m.group(1)}"

    # --- "X are prohibited" 模式 ---
    m = re.match(r'^(.+?)\s+are\s+prohibited$', name)
    if m: return f"禁止{m.group(1)}"

    # --- "X must be disabled" 模式 ---
    m = re.match(r'^(.+?)\s+must\s+be\s+disabled$', name)
    if m: return f"{m.group(1)} 必须禁用"

    # --- "X for Y must be disabled" 模式 ---
    m = re.match(r'^(.+?)\s+for\s+(.+?)\s+must\s+be\s+disabled$', name)
    if m: return f"必须禁用 {m.group(2)} 的 {m.group(1)}"

    # --- "X on the Y are prohibited" 模式 ---
    m = re.match(r'^(.+?)\s+on\s+the\s+(.+?)\s+are\s+prohibited$', name)
    if m: return f"禁止在 {m.group(2)} 上{m.group(1)}"

    # --- "X Permissions are Y" 模式 ---
    m = re.match(r'^(.+?)\s+Permissions\s+are\s+(.+)$', name)
    if m: return f"{m.group(1)} 权限为 {m.group(2)}"

    # --- "X must have the proper permissions to Y" 模式 ---
    m = re.match(r'^(.+?)\s+must\s+have\s+the\s+proper\s+permissions\s+to\s+(.+)$', name)
    if m: return f"{m.group(1)} 必须具有适当权限以{m.group(2)}"

    # --- "X must be Y" 模式 ---
    m = re.match(r'^(.+?)\s+must\s+be\s+(.+)$', name)
    if m: return f"{m.group(1)} 必须{m.group(2)}"

    # --- "X must Y" 模式 ---
    m = re.match(r'^(.+?)\s+must\s+(.+)$', name)
    if m: return f"{m.group(1)} 必须{m.group(2)}"

    # --- "X must not be Y" 模式 ---
    m = re.match(r'^(.+?)\s+must\s+not\s+be\s+(.+)$', name)
    if m: return f"{m.group(1)} 不得{m.group(2)}"

    # --- "X must not Y" 模式 ---
    m = re.match(r'^(.+?)\s+must\s+not\s+(.+)$', name)
    if m: return f"{m.group(1)} 不得{m.group(2)}"

    # --- "Add X Option to Y" 模式 ---
    m = re.match(r'^Add\s+(.+?)\s+Option\s+to\s+(.+)$', name)
    if m: return f"为 {m.group(2)} 添加 {m.group(1)} 选项"

    # --- "Mount X with Y" 模式 ---
    m = re.match(r'^Mount\s+(.+?)\s+with\s+(.+)$', name)
    if m: return f"使用 {m.group(2)} 挂载 {m.group(1)}"

    # --- "Use X on Y" 模式 ---
    m = re.match(r'^Use\s+(.+?)\s+on\s+(.+)$', name)
    if m: return f"在 {m.group(2)} 上使用 {m.group(1)}"

    # --- "Use Only X in Y" 模式 ---
    m = re.match(r'^Use\s+Only\s+(.+?)\s+in\s+(.+)$', name)
    if m: return f"在 {m.group(2)} 中仅使用 {m.group(1)}"

    # --- "Use Only X" 模式 ---
    m = re.match(r'^Use\s+Only\s+(.+)$', name)
    if m: return f"仅使用 {m.group(1)}"

    # --- "Use X for Y" 模式 ---
    m = re.match(r'^Use\s+(.+?)\s+for\s+(.+)$', name)
    if m: return f"使用 {m.group(1)} 进行 {m.group(2)}"

    # --- "Use X instead of Y" 模式 ---
    m = re.match(r'^Use\s+(.+?)\s+instead\s+of\s+(.+)$', name)
    if m: return f"使用 {m.group(1)} 代替 {m.group(2)}"

    # --- "Extend X for Y in Z" 模式 ---
    m = re.match(r'^Extend\s+(.+?)\s+for\s+(.+?)\s+in\s+(.+)$', name)
    if m: return f"在 {m.group(3)} 中为 {m.group(2)} 扩展 {m.group(1)}"

    # --- "Extend X for Y" 模式 ---
    m = re.match(r'^Extend\s+(.+?)\s+for\s+(.+)$', name)
    if m: return f"为 {m.group(2)} 扩展 {m.group(1)}"

    # --- "Create X for Y" 模式 ---
    m = re.match(r'^Create\s+(.+?)\s+for\s+(.+)$', name)
    if m: return f"为 {m.group(2)} 创建 {m.group(1)}"

    # --- "Write X to Y" 模式 ---
    m = re.match(r'^Write\s+(.+?)\s+to\s+(.+)$', name)
    if m: return f"将 {m.group(1)} 写入 {m.group(2)}"

    # --- "Include X in Y" 模式 ---
    m = re.match(r'^Include\s+(.+?)\s+in\s+(.+)$', name)
    if m: return f"在 {m.group(2)} 中包含 {m.group(1)}"

    # --- "Resolve X before Y" 模式 ---
    m = re.match(r'^Resolve\s+(.+?)\s+before\s+(.+)$', name)
    if m: return f"在 {m.group(2)} 前解析 {m.group(1)}"

    # --- "Encrypt X Sent With Y" 模式 ---
    m = re.match(r'^Encrypt\s+(.+?)\s+Sent\s+With\s+(.+)$', name)
    if m: return f"加密通过 {m.group(2)} 发送的 {m.group(1)}"

    # --- "Shutdown X When Y" 模式 ---
    m = re.match(r'^Shutdown\s+(.+?)\s+When\s+(.+)$', name)
    if m: return f"{m.group(2)}时关闭{m.group(1)}"

    # --- "Make the X Y" 模式 ---
    m = re.match(r'^Make\s+the\s+(.+?)\s+(.+)$', name)
    if m: return f"确保 {m.group(1)} {m.group(2)}"

    # --- "Log X using Y" 模式 ---
    m = re.match(r'^Log\s+(.+?)\s+using\s+(.+)$', name)
    if m: return f"使用 {m.group(2)} 记录 {m.group(1)}"

    # --- "Authorize X and Y in Z" 模式 ---
    m = re.match(r'^Authorize\s+(.+?)\s+and\s+(.+?)\s+in\s+(.+)$', name)
    if m: return f"授权 {m.group(3)} 中的 {m.group(1)} 和 {m.group(2)}"

    # --- "Authorize X in Y" 模式 ---
    m = re.match(r'^Authorize\s+(.+?)\s+in\s+(.+)$', name)
    if m: return f"授权 {m.group(2)} 中的 {m.group(1)}"

    # --- "Require X for Y" 模式 ---
    m = re.match(r'^Require\s+(.+?)\s+for\s+(.+)$', name)
    if m: return f"要求 {m.group(2)} 的 {m.group(1)}"

    # --- "Require X in Y" 模式 ---
    m = re.match(r'^Require\s+(.+?)\s+in\s+(.+)$', name)
    if m: return f"在 {m.group(2)} 中要求 {m.group(1)}"

    # --- "Require X When Using Y" 模式 ---
    m = re.match(r'^Require\s+(.+?)\s+When\s+Using\s+(.+)$', name)
    if m: return f"使用 {m.group(2)} 时要求 {m.group(1)}"

    # --- "Require X" 模式 ---
    m = re.match(r'^Require\s+(.+)$', name)
    if m: return f"要求 {m.group(1)}"

    # --- "Prevent X from Y" 模式 ---
    m = re.match(r'^Prevent\s+(.+?)\s+from\s+(.+)$', name)
    if m: return f"防止 {m.group(1)} {m.group(2)}"

    # --- "Prevent X" 模式 ---
    m = re.match(r'^Prevent\s+(.+)$', name)
    if m: return f"防止 {m.group(1)}"

    # --- "Limit X" 模式 ---
    m = re.match(r'^Limit\s+(.+)$', name)
    if m: return f"限制 {m.group(1)}"

    # --- "Restrict X" 模式 ---
    m = re.match(r'^Restrict\s+(.+)$', name)
    if m: return f"限制 {m.group(1)}"

    # --- "Enforce X" 模式 ---
    m = re.match(r'^Enforce\s+(.+)$', name)
    if m: return f"强制 {m.group(1)}"

    # --- "Disallow X" 模式 ---
    m = re.match(r'^Disallow\s+(.+)$', name)
    if m: return f"禁止 {m.group(1)}"

    # --- "Harden X" 模式 ---
    m = re.match(r'^Harden\s+(.+)$', name)
    if m: return f"加固 {m.group(1)}"

    # --- "Randomize X" 模式 ---
    m = re.match(r'^Randomize\s+(.+)$', name)
    if m: return f"随机化 {m.group(1)}"

    # --- "Force X" 模式 ---
    m = re.match(r'^Force\s+(.+)$', name)
    if m: return f"强制 {m.group(1)}"

    # --- "Specify X" 模式 ---
    m = re.match(r'^Specify\s+(.+)$', name)
    if m: return f"指定 {m.group(1)}"

    # --- "Sign X with Y" 模式 ---
    m = re.match(r'^Sign\s+(.+?)\s+with\s+(.+)$', name)
    if m: return f"使用 {m.group(2)} 签名 {m.group(1)}"

    # --- "Perform X" 模式 ---
    m = re.match(r'^Perform\s+(.+)$', name)
    if m: return f"执行 {m.group(1)}"

    # --- "Avoid X" 模式 ---
    m = re.match(r'^Avoid\s+(.+)$', name)
    if m: return f"避免 {m.group(1)}"

    # --- "Detect X on Y" 模式 ---
    m = re.match(r'^Detect\s+(.+?)\s+on\s+(.+)$', name)
    if m: return f"检测 {m.group(2)} 上的 {m.group(1)}"

    # --- "Drop X on Y" 模式 ---
    m = re.match(r'^Drop\s+(.+?)\s+on\s+(.+)$', name)
    if m: return f"丢弃 {m.group(2)} 上的 {m.group(1)}"

    # --- "Deactivate X" 模式 ---
    m = re.match(r'^Deactivate\s+(.+)$', name)
    if m: return f"停用 {m.group(1)}"

    # --- "Bind X To Y" 模式 ---
    m = re.match(r'^Bind\s+(.+?)\s+To\s+(.+)$', name)
    if m: return f"将 {m.group(1)} 绑定到 {m.group(2)}"

    # --- "Unmap X when Y" 模式 ---
    m = re.match(r'^Unmap\s+(.+?)\s+when\s+(.+)$', name)
    if m: return f"在 {m.group(2)} 时取消映射 {m.group(1)}"

    # --- "Emulate X" 模式 ---
    m = re.match(r'^Emulate\s+(.+)$', name)
    if m: return f"模拟 {m.group(1)}"

    # --- "Trigger X when Y" 模式 ---
    m = re.match(r'^Trigger\s+(.+?)\s+when\s+(.+)$', name)
    if m: return f"在 {m.group(2)} 时触发 {m.group(1)}"

    # --- "Warn on X" 模式 ---
    m = re.match(r'^Warn\s+on\s+(.+)$', name)
    if m: return f"在 {m.group(1)} 时发出警告"

    # --- "Generate X" 模式 ---
    m = re.match(r'^Generate\s+(.+)$', name)
    if m: return f"生成 {m.group(1)}"

    # --- "Modify X" 模式 ---
    m = re.match(r'^Modify\s+(.+)$', name)
    if m: return f"修改 {m.group(1)}"

    # --- "Display X until Y" 模式 ---
    m = re.match(r'^Display\s+(.+?)\s+until\s+(.+)$', name)
    if m: return f"显示 {m.group(1)} 直至 {m.group(2)}"

    # --- "Lock X After Y" 模式 ---
    m = re.match(r'^Lock\s+(.+?)\s+After\s+(.+)$', name)
    if m: return f"在 {m.group(2)} 后锁定 {m.group(1)}"

    # --- "Support X with Y (not enforcing)" 模式 ---
    m = re.match(r'^Support\s+(.+?)\s+with\s+(.+?)\s+\(not\s+enforcing\)$', name)
    if m: return f"支持使用 {m.group(2)} 的 {m.group(1)}（非强制）"

    # --- "Support X with Y" 模式 ---
    m = re.match(r'^Support\s+(.+?)\s+with\s+(.+)$', name)
    if m: return f"支持使用 {m.group(2)} 的 {m.group(1)}"

    # --- "Elevate X When Y" 模式 ---
    m = re.match(r'^Elevate\s+(.+?)\s+When\s+(.+)$', name)
    if m: return f"在 {m.group(2)} 时提升 {m.group(1)}"

    # --- "Implement X" 模式 ---
    m = re.match(r'^Implement\s+(.+)$', name)
    if m: return f"实现 {m.group(1)}"

    # --- "Build and Test X" 模式 ---
    m = re.match(r'^Build\s+and\s+Test\s+(.+)$', name)
    if m: return f"构建并测试 {m.group(1)}"

    # --- "Verify and Correct X with Y" 模式 ---
    m = re.match(r'^Verify\s+and\s+Correct\s+(.+?)\s+with\s+(.+)$', name)
    if m: return f"使用 {m.group(2)} 验证并更正 {m.group(1)}"

    # --- "Set the X to Y" 模式 ---
    m = re.match(r'^Set\s+the\s+(.+?)\s+to\s+(.+)$', name)
    if m: return f"将 {m.group(1)} 设置为 {m.group(2)}"

    # --- "Set X to Y" 模式 ---
    m = re.match(r"^Set\s+(.+?)\s+to\s+(.+)$", name)
    if m: return f"将 {m.group(1)} 设置为 {m.group(2)}"

    # --- "Set the X" 模式 ---
    m = re.match(r'^Set\s+the\s+(.+)$', name)
    if m: return f"设置 {m.group(1)}"

    # --- "Ensure X is Y" 模式 ---
    m = re.match(r'^Ensure\s+(.+?)\s+is\s+(.+)$', name)
    if m: return f"确保 {m.group(1)} 为 {m.group(2)}"

    # --- "Ensure X are Y" 模式 ---
    m = re.match(r'^Ensure\s+(.+?)\s+are\s+(.+)$', name)
    if m: return f"确保 {m.group(1)} 为 {m.group(2)}"

    # --- "Ensure X Not Y" 模式 ---
    m = re.match(r'^Ensure\s+(.+?)\s+Not\s+(.+)$', name)
    if m: return f"确保 {m.group(1)} 未{m.group(2)}"

    # --- "Ensure X Located On Y" 模式 ---
    m = re.match(r'^Ensure\s+(.+?)\s+Located\s+On\s+(.+)$', name)
    if m: return f"确保 {m.group(1)} 位于 {m.group(2)}"

    # --- "Ensure X exists" 模式 ---
    m = re.match(r"^Ensure\s+'(.+?)'\s+exists$", name)
    if m: return f"确保 '{m.group(1)}' 存在"

    # --- "Verify X exists" 模式 ---
    m = re.match(r"^Verify\s+'(.+?)'\s+exists$", name)
    if m: return f"验证 '{m.group(1)}' 存在"

    # --- "Configure X to Use Y for Z" 模式 ---
    m = re.match(r'^Configure\s+(.+?)\s+to\s+Use\s+(.+?)\s+for\s+(.+)$', name)
    if m: return f"配置 {m.group(1)} 使用 {m.group(2)} 进行 {m.group(3)}"

    # --- "Configure X to Use Y: Z" 模式 ---
    m = re.match(r'^Configure\s+(.+?)\s+to\s+Use\s+(.+?):\s+(.+)$', name)
    if m: return f"配置 {m.group(1)} 使用 {m.group(2)}：{m.group(3)}"

    # --- "Configure X to Use Y" 模式 ---
    m = re.match(r'^Configure\s+(.+?)\s+to\s+Use\s+(.+)$', name)
    if m: return f"配置 {m.group(1)} 使用 {m.group(2)}"

    # --- "Configure X for Y" 模式 ---
    m = re.match(r'^Configure\s+(.+?)\s+for\s+(.+)$', name)
    if m: return f"为 {m.group(2)} 配置 {m.group(1)}"

    # --- "Configure X: Y" 模式 ---
    m = re.match(r'^Configure\s+(.+?):\s+(.+)$', name)
    if m: return f"配置 {m.group(1)}：{m.group(2)}"

    # --- "Configure X" 模式 ---
    m = re.match(r'^Configure\s+(.+)$', name)
    if m: return f"配置 {m.group(1)}"

    # --- "Install X Package/Service" 模式 ---
    m = re.match(r'^Install\s+(.+?)\s+(Package|Service|software)$', name, re.IGNORECASE)
    if m: return f"安装 {m.group(1)} {m.group(2)}"
    m = re.match(r'^Install\s+(.+)$', name, re.IGNORECASE)
    if m: return f"安装 {m.group(1)}"

    # --- "Uninstall X Package" 模式 ---
    m = re.match(r'^Uninstall\s+(.+?)\s+(Package)$', name, re.IGNORECASE)
    if m: return f"卸载 {m.group(1)} 软件包"
    m = re.match(r'^Uninstall\s+(.+)$', name, re.IGNORECASE)
    if m: return f"卸载 {m.group(1)}"

    # --- "Remove X" 模式 ---
    m = re.match(r'^Remove\s+(.+)$', name, re.IGNORECASE)
    if m: return f"移除 {m.group(1)}"

    # --- "Enable X Service/Daemon" 模式 ---
    m = re.match(r'^Enable\s+(.+?)\s+(Service|Daemon|the\s+.+)$', name, re.IGNORECASE)
    if m: return f"启用 {m.group(1)}"
    m = re.match(r'^Enable\s+(.+)$', name, re.IGNORECASE)
    if m: return f"启用 {m.group(1)}"

    # --- "Disable X Service/Daemon" 模式 ---
    m = re.match(r'^Disable\s+(.+?)\s+(Service|Daemon|Server Software)$', name, re.IGNORECASE)
    if m: return f"禁用 {m.group(1)}"
    m = re.match(r'^Disable\s+(.+)$', name, re.IGNORECASE)
    if m: return f"禁用 {m.group(1)}"

    # --- "Verify Group Who Owns X" 模式 ---
    m = re.match(r'^Verify\s+Group\s+Who\s+Owns\s+(.+)$', name, re.IGNORECASE)
    if m: return f"验证 {m.group(1)} 的组所有权"
    m = re.match(r'^Verify\s+Group\s+Ownership\s+of\s+(.+)$', name, re.IGNORECASE)
    if m: return f"验证 {m.group(1)} 的组所有权"
    m = re.match(r'^Verify\s+Group\s+Ownership\s+on\s+(.+)$', name, re.IGNORECASE)
    if m: return f"验证 {m.group(1)} 的组所有权"

    # --- "Verify Owner on X" / "Verify Ownership of X" 模式 ---
    m = re.match(r'^Verify\s+(Owner|User)\s+Who\s+Owns\s+(.+)$', name, re.IGNORECASE)
    if m: return f"验证 {m.group(2)} 的用户所有权"
    m = re.match(r'^Verify\s+(Owner|User)\s+on\s+(.+)$', name, re.IGNORECASE)
    if m: return f"验证 {m.group(2)} 的用户所有权"
    m = re.match(r'^Verify\s+[Oo]wnership\s+of\s+(.+)$', name)
    if m: return f"验证 {m.group(1)} 的所有权"

    # --- "Verify Permissions on X" 模式 ---
    m = re.match(r'^Verify\s+Permissions\s+(on|Of)\s+(.+)$', name, re.IGNORECASE)
    if m: return f"验证 {m.group(2)} 的权限"
    m = re.match(r'^Verify\s+[Pp]ermissions\s+[Oo]n\s+(.+)$', name)
    if m: return f"验证 {m.group(1)} 的权限"

    # --- "Verify that X has Y" / "Verify X" 模式 ---
    m = re.match(r'^Verify\s+that\s+(.+?)\s+has\s+(.+)$', name, re.IGNORECASE)
    if m: return f"验证 {m.group(1)} 具有 {m.group(2)}"
    m = re.match(r"^Verify\s+'(/.+?)'\s+exists$", name)
    if m: return f"验证 '{m.group(1)}' 存在"
    m = re.match(r'^Verify\s+(.+?)\s+Have\s+Mode\s+(.+?)\s+or\s+less$', name, re.IGNORECASE)
    if m: return f"验证 {m.group(1)} 具有模式 {m.group(2)} 或更严格"
    m = re.match(r'^Verify\s+(.+?)\s+(Has|Have)\s+(.+)$', name, re.IGNORECASE)
    if m: return f"验证 {m.group(1)} 具有 {m.group(3)}"
    m = re.match(r'^Verify\s+(.+)$', name, re.IGNORECASE)
    if m: return f"验证 {m.group(1)}"

    # --- "Ensure that X exists/does not exist" 模式 ---
    m = re.match(r'^Ensure\s+that\s+(.+?)\s+does\s+not\s+exist$', name, re.IGNORECASE)
    if m: return f"确保 {m.group(1)} 不存在"
    m = re.match(r'^Ensure\s+that\s+(.+?)\s+exists$', name, re.IGNORECASE)
    if m: return f"确保 {m.group(1)} 存在"

    # --- "Ensure X" 模式 (catch-all after specific patterns) ---
    m = re.match(r'^Ensure\s+(.+)$', name, re.IGNORECASE)
    if m: return f"确保 {m.group(1)}"

    # --- "Record Events that Modify X" 模式 ---
    m = re.match(r'^Record\s+Events\s+that\s+Modify\s+(.+)$', name, re.IGNORECASE)
    if m: return f"记录修改 {m.group(1)} 的事件"

    # --- "Record Events When X" 模式 ---
    m = re.match(r'^Record\s+Events\s+When\s+(.+)$', name, re.IGNORECASE)
    if m: return f"记录 {m.group(1)} 事件"

    # --- "Record Attempts to X" 模式 ---
    m = re.match(r'^Record\s+Attempts\s+to\s+(.+)$', name, re.IGNORECASE)
    if m: return f"记录 {m.group(1)} 尝试"

    # --- "Record Access Events to X" 模式 ---
    m = re.match(r'^Record\s+Access\s+Events\s+to\s+(.+)$', name, re.IGNORECASE)
    if m: return f"记录对 {m.group(1)} 的访问事件"

    # --- "Set X" 模式 ---
    m = re.match(r'^Set\s+(.+?)\s+(to|in|for|as|On)\s+(.+)$', name, re.IGNORECASE)
    if m: return f"设置 {m.group(1)} 为 {m.group(3)}"
    m = re.match(r'^Set\s+(.+)$', name, re.IGNORECASE)
    if m: return f"设置 {m.group(1)}"

    # --- "Ensure PAM Enforces Password Requirements - X" 前缀 ---
    if name.startswith("Ensure PAM Enforces Password Requirements - "):
        suffix = name[len("Ensure PAM Enforces Password Requirements - "):]
        return f"确保 PAM 强制执行密码要求 - {suffix}"

    # --- "Set Password X" 模式 ---
    m = re.match(r'^Set\s+Password\s+(.+)$', name, re.IGNORECASE)
    if m: return f"设置密码{m.group(1)}"

    # --- "Set PAM Password Hashing Algorithm - X" 模式 ---
    m = re.match(r'^Set\s+PAM\s+Password\s+Hashing\s+Algorithm\s*-\s*(.+)$', name, re.IGNORECASE)
    if m: return f"设置 PAM 密码哈希算法 - {m.group(1)}"

    # --- "Set number of Password Hashing Rounds - X" 模式 ---
    m = re.match(r'^Set\s+number\s+of\s+Password\s+Hashing\s+Rounds\s*-\s*(.+)$', name, re.IGNORECASE)
    if m: return f"设置密码哈希轮数 - {m.group(1)}"

    # --- "Set Existing Passwords X" 模式 ---
    m = re.match(r'^Set\s+Existing\s+Passwords\s+(.+)$', name, re.IGNORECASE)
    if m: return f"设置现有密码{m.group(1)}"

    # --- "Set existing passwords X" 模式 ---
    m = re.match(r'^Set\s+existing\s+passwords\s+(.+)$', name)
    if m: return f"设置现有密码{m.group(1)}"

    # --- "Set Lockout Time for Failed Password Attempts" 系列 ---
    m = re.match(r'^Set\s+(Lockout\s+Time|Interval|Deny|Root\s+Lockout\s+Time)\s+for\s+Failed\s+Password\s+Attempts.*$', name, re.IGNORECASE)
    if m: return f"设置密码尝试失败的{m.group(1)}"

    # --- "Set Account Expiration Following Inactivity" 系列 ---
    m = re.match(r'^Set\s+Account\s+Expiration\s+Following\s+Inactivity.*$', name, re.IGNORECASE)
    if m: return f"设置账户不活动过期"

    # --- "Ensure Root Account Lockout on Failed Password Attempts" ---
    m = re.match(r'^Ensure\s+Root\s+Account\s+Lockout\s+on\s+Failed\s+Password\s+Attempts$', name, re.IGNORECASE)
    if m: return f"确保密码尝试失败时锁定 root 账户"

    # --- "Ensure Password History Is Enforced for the Root User" ---
    m = re.match(r'^Ensure\s+Password\s+History\s+Is\s+Enforced\s+for\s+the\s+Root\s+User$', name, re.IGNORECASE)
    if m: return f"确保对 root 用户强制执行密码历史"

    # --- "Ensure Active Authselect Profile Includes PAM Modules" ---
    m = re.match(r'^Ensure\s+Active\s+Authselect\s+Profile\s+Includes\s+PAM\s+Modules$', name, re.IGNORECASE)
    if m: return f"确保活跃的 Authselect 配置文件包含 PAM 模块"

    # --- "Ensure PAM Displays Last Logon/Access Notification" ---
    m = re.match(r'^Ensure\s+PAM\s+Displays\s+Last\s+Logon/Access\s+Notification$', name, re.IGNORECASE)
    if m: return f"确保 PAM 显示上次登录/访问通知"

    # --- "Set Up a Private Namespace in PAM Configuration" ---
    m = re.match(r'^Set\s+Up\s+a\s+Private\s+Namespace\s+in\s+PAM\s+Configuration$', name, re.IGNORECASE)
    if m: return f"在 PAM 配置中设置私有命名空间"

    # --- "Ensure system-auth and password-auth files are symbolic links pointing to X" ---
    m = re.match(r'^Ensure\s+system-auth\s+and\s+password-auth\s+files\s+are\s+symbolic\s+links\s+pointing\s+to\s+(.+)$', name, re.IGNORECASE)
    if m: return f"确保 system-auth 和 password-auth 文件是指向 {m.group(1)} 的符号链接"

    # --- "Ensure All User Initialization Files Have Mode X Or Less Permissive" ---
    m = re.match(r'^Ensure\s+All\s+User\s+Initialization\s+Files\s+Have\s+Mode\s+(.+?)\s+Or\s+Less\s+Permissive$', name, re.IGNORECASE)
    if m: return f"确保所有用户初始化文件具有模式 {m.group(1)} 或更严格"

    # --- "Ensure that Users Path Contains Only Local Directories" ---
    m = re.match(r'^Ensure\s+that\s+Users\s+Path\s+Contains\s+Only\s+Local\s+Directories$', name, re.IGNORECASE)
    if m: return f"确保用户路径仅包含本地目录"

    # --- "Ensure All Accounts on the System Have Unique X" ---
    m = re.match(r'^Ensure\s+All\s+Accounts\s+on\s+the\s+System\s+Have\s+Unique\s+(.+)$', name, re.IGNORECASE)
    if m: return f"确保系统上所有账户具有唯一的{m.group(1)}"

    # --- "Ensure that System Accounts Are Locked" ---
    m = re.match(r'^Ensure\s+that\s+System\s+Accounts\s+Are\s+Locked$', name, re.IGNORECASE)
    if m: return f"确保系统账户已锁定"

    # --- "Ensure There Are No Accounts With Blank or Null Passwords" ---
    m = re.match(r'^Ensure\s+There\s+Are\s+No\s+Accounts\s+With\s+Blank\s+or\s+Null\s+Passwords$', name, re.IGNORECASE)
    if m: return f"确保没有账户具有空白或空密码"

    # --- "Ensure sudo group has only necessary members" ---
    m = re.match(r'^Ensure\s+sudo\s+group\s+has\s+only\s+necessary\s+members$', name, re.IGNORECASE)
    if m: return f"确保 sudo 组仅包含必要成员"

    # --- "Verify No X Files Exist" 模式 ---
    m = re.match(r'^Verify\s+No\s+(.+?)\s+Files\s+Exist$', name, re.IGNORECASE)
    if m: return f"确保不存在 {m.group(1)} 文件"

    # --- "Verify Only X Has Y" 模式 ---
    m = re.match(r'^Verify\s+Only\s+(.+?)\s+Has\s+(.+)$', name, re.IGNORECASE)
    if m: return f"验证仅 {m.group(1)} 具有 {m.group(2)}"

    # --- "Verify Root Has A Primary GID 0" ---
    m = re.match(r'^Verify\s+Root\s+Has\s+A\s+Primary\s+GID\s+0$', name, re.IGNORECASE)
    if m: return f"验证 root 的主 GID 为 0"

    # --- "Ensure Web Content Located on Separate partition" ---
    m = re.match(r'^Ensure\s+Web\s+Content\s+Located\s+on\s+Separate\s+partition$', name, re.IGNORECASE)
    if m: return f"确保 Web 内容位于独立分区"

    # --- "Ensure /X Located On Separate Partition" 模式 ---
    m = re.match(r'^Ensure\s+(/\S+)\s+Located\s+On\s+Separate\s+Partition$', name, re.IGNORECASE)
    if m: return f"确保 {m.group(1)} 位于独立分区"

    # --- "Ensure /X is configured" 模式 ---
    m = re.match(r"^Ensure\s+(/\S+)\s+is\s+configured$", name)
    if m: return f"确保 {m.group(1)} 已配置"

    # --- "Ensure tmp.mount Unit Is Enabled" ---
    m = re.match(r'^Ensure\s+tmp\.mount\s+Unit\s+Is\s+Enabled$', name, re.IGNORECASE)
    if m: return f"确保 tmp.mount 单元已启用"

    # --- "Enable Auditing for Processes Which Start Prior to the Audit Daemon" ---
    m = re.match(r'^Enable\s+Auditing\s+for\s+Processes\s+Which\s+Start\s+Prior\s+to\s+the\s+Audit\s+Daemon$', name, re.IGNORECASE)
    if m: return f"启用审计守护进程启动前进程的审计"

    # --- "Enable Auditing to Start Prior to the Audit Daemon in zIPL" ---
    m = re.match(r'^Enable\s+Auditing\s+to\s+Start\s+Prior\s+to\s+the\s+Audit\s+Daemon\s+in\s+zIPL$', name, re.IGNORECASE)
    if m: return f"在 zIPL 中启用审计守护进程启动前的审计"

    # --- "Set type of computer node name logging in audit logs" ---
    m = re.match(r'^Set\s+type\s+of\s+computer\s+node\s+name\s+logging\s+in\s+audit\s+logs$', name, re.IGNORECASE)
    if m: return f"设置审计日志中的计算机节点名称记录类型"

    # --- "Remove Default Configuration to Disable Syscall Auditing" ---
    m = re.match(r'^Remove\s+Default\s+Configuration\s+to\s+Disable\s+Syscall\s+Auditing$', name, re.IGNORECASE)
    if m: return f"移除禁用系统调用审计的默认配置"

    # --- "Disable X in ifcfg" 模式 ---
    m = re.match(r'^Disable\s+(.+?)\s+in\s+ifcfg$', name, re.IGNORECASE)
    if m: return f"在 ifcfg 中禁用 {m.group(1)}"

    # --- "Enable Logging of All FTP Transactions" ---
    m = re.match(r'^Enable\s+Logging\s+of\s+All\s+FTP\s+Transactions$', name, re.IGNORECASE)
    if m: return f"启用所有 FTP 事务的日志记录"

    # --- "Set Permissions on the X Directory" 模式 ---
    m = re.match(r'^Set\s+Permissions\s+on\s+(.+)$', name, re.IGNORECASE)
    if m: return f"设置 {m.group(1)} 的权限"

    # --- "Disable Anonymous FTP Access" ---
    m = re.match(r'^Disable\s+Anonymous\s+FTP\s+Access$', name, re.IGNORECASE)
    if m: return f"禁用匿名 FTP 访问"

    # --- "Enable HTTPD X" 模式 ---
    m = re.match(r'^Enable\s+(HTTPD|Transport)\s+(.+)$', name, re.IGNORECASE)
    if m: return f"启用 {m.group(1)} {m.group(2)}"

    # --- "Remove X Files" 模式 ---
    m = re.match(r'^Remove\s+(.+?)\s+(Files|Package)$', name, re.IGNORECASE)
    if m: return f"移除 {m.group(1)} {m.group(2)}"

    # --- "Disable Postfix Network Listening" ---
    m = re.match(r'^Disable\s+Postfix\s+Network\s+Listening$', name, re.IGNORECASE)
    if m: return f"禁用 Postfix 网络监听"

    # --- "Ensure All-Squashing Disabled On All Exports" ---
    m = re.match(r'^Ensure\s+All-Squashing\s+Disabled\s+On\s+All\s+Exports$', name, re.IGNORECASE)
    if m: return f"确保所有导出上禁用 All-Squashing"

    # --- "Disable chrony daemon from acting as server" ---
    m = re.match(r'^Disable\s+chrony\s+daemon\s+from\s+acting\s+as\s+server$', name, re.IGNORECASE)
    if m: return f"禁用 chrony 守护进程的服务器功能"

    # --- "Disable network management of chrony daemon" ---
    m = re.match(r'^Disable\s+network\s+management\s+of\s+chrony\s+daemon$', name, re.IGNORECASE)
    if m: return f"禁用 chrony 守护进程的网络管理"

    # --- "Remove the X service" 模式 ---
    m = re.match(r'^Remove\s+the\s+(.+?)\s+service$', name, re.IGNORECASE)
    if m: return f"移除 {m.group(1)} 服务"

    # --- "Install/Remove the X Package" 模式 ---
    m = re.match(r'^(Install|Remove)\s+the\s+(.+?)\s+(Package|Service)$', name, re.IGNORECASE)
    if m: return f"{'安装' if m.group(1)=='Install' else '移除'} {m.group(2)} {m.group(3)}"

    # --- "Enable/Disable the X Service" 模式 ---
    m = re.match(r'^(Enable|Disable)\s+the\s+(.+?)\s+(Service|Daemon)$', name, re.IGNORECASE)
    if m: return f"{'启用' if m.group(1)=='Enable' else '禁用'} {m.group(2)} 服务"

    # --- "Enable/Disable X Service (X)" 模式 ---
    m = re.match(r'^(Enable|Disable)\s+(.+?)\s+Service\s+\((.+)\)$', name, re.IGNORECASE)
    if m: return f"{'启用' if m.group(1)=='Enable' else '禁用'} {m.group(2)} 服务 ({m.group(3)})"

    # --- "Enable/Disable X (X)" 模式 ---
    m = re.match(r'^(Enable|Disable)\s+(.+?)\s+\((.+)\)$', name, re.IGNORECASE)
    if m: return f"{'启用' if m.group(1)=='Enable' else '禁用'} {m.group(2)} ({m.group(3)})"

    # --- "Enable/Disable X" 模式 (catch-all) ---
    m = re.match(r'^(Enable|Disable)\s+(.+)$', name, re.IGNORECASE)
    if m: return f"{'启用' if m.group(1)=='Enable' else '禁用'} {m.group(2)}"

    # --- "Remove Host-Based Authentication Files" ---
    m = re.match(r'^Remove\s+(Host-Based|User\s+Host-Based)\s+Authentication\s+Files$', name, re.IGNORECASE)
    if m: return f"移除{m.group(1)}认证文件"

    # --- "Ensure tftp Daemon Uses Secure Mode" ---
    m = re.match(r'^Ensure\s+tftp\s+(Daemon|systemd\s+Service)\s+Uses\s+Secure\s+Mode$', name, re.IGNORECASE)
    if m: return f"确保 tftp 使用安全模式"

    # --- "Disable Printer Browsing Entirely if Possible" ---
    m = re.match(r'^Disable\s+Printer\s+Browsing\s+Entirely\s+if\s+Possible$', name, re.IGNORECASE)
    if m: return f"尽可能完全禁用打印机浏览"

    # --- "Enable the Hardware RNG Entropy Gatherer Service" ---
    m = re.match(r'^Enable\s+the\s+Hardware\s+RNG\s+Entropy\s+Gatherer\s+Service$', name, re.IGNORECASE)
    if m: return f"启用硬件 RNG 熵收集器服务"

    # --- "Disable SSH X" 系列 ---
    m = re.match(r'^Disable\s+SSH\s+(.+)$', name, re.IGNORECASE)
    if m: return f"禁用 SSH {m.group(1)}"

    # --- "Enable SSH X" 系列 ---
    m = re.match(r'^Enable\s+SSH\s+(.+)$', name, re.IGNORECASE)
    if m: return f"启用 SSH {m.group(1)}"

    # --- "Set SSH X" 系列 ---
    m = re.match(r'^Set\s+SSH\s+(.+)$', name, re.IGNORECASE)
    if m: return f"设置 SSH {m.group(1)}"

    # --- "Enable Use of X" 模式 ---
    m = re.match(r'^Enable\s+Use\s+of\s+(.+)$', name, re.IGNORECASE)
    if m: return f"启用 {m.group(1)} 的使用"

    # --- "Enable X in Y" 模式 ---
    m = re.match(r'^Enable\s+(.+?)\s+in\s+(.+)$', name, re.IGNORECASE)
    if m: return f"在 {m.group(2)} 中启用 {m.group(1)}"

    # --- "Remove the X Windows Package Group" ---
    m = re.match(r'^Remove\s+the\s+X\s+Windows\s+(.+)$', name, re.IGNORECASE)
    if m: return f"移除 X Window {m.group(1)}"

    # --- "Disable Graphical Environment Startup By Setting Default Target" ---
    m = re.match(r'^Disable\s+Graphical\s+Environment\s+Startup\s+By\s+Setting\s+Default\s+Target$', name, re.IGNORECASE)
    if m: return f"通过设置默认目标禁用图形环境启动"

    # --- "Disable graphical user interface" ---
    m = re.match(r'^Disable\s+graphical\s+user\s+interface$', name, re.IGNORECASE)
    if m: return f"禁用图形用户界面"

    # --- "Ensure Local/Remote Login Warning Banner Is Configured Properly" ---
    m = re.match(r'^Ensure\s+(Local|Remote)\s+Login\s+Warning\s+Banner\s+Is\s+Configured\s+Properly$', name, re.IGNORECASE)
    if m: return f"确保{m.group(1)}登录警告横幅已正确配置"

    # --- "Ensure Message Of The Day Is Configured Properly" ---
    m = re.match(r'^Ensure\s+Message\s+Of\s+The\s+Day\s+Is\s+Configured\s+Properly$', name, re.IGNORECASE)
    if m: return f"确保每日消息已正确配置"

    # --- "Enable the SSH login confirmation banner" ---
    m = re.match(r'^Enable\s+the\s+SSH\s+login\s+confirmation\s+banner$', name, re.IGNORECASE)
    if m: return f"启用 SSH 登录确认横幅"

    # --- "Enable GNOME3 Login Warning Banner" ---
    m = re.match(r'^Enable\s+GNOME3\s+Login\s+Warning\s+Banner$', name, re.IGNORECASE)
    if m: return f"启用 GNOME3 登录警告横幅"

    # --- "Disable Ctrl-Alt-Del X" 模式 ---
    m = re.match(r'^Disable\s+Ctrl-Alt-Del\s+(.+)$', name, re.IGNORECASE)
    if m: return f"禁用 Ctrl-Alt-Del {m.group(1)}"

    # --- "Install Smart Card Packages For Multifactor Authentication" ---
    m = re.match(r'^Install\s+Smart\s+Card\s+Packages\s+For\s+Multifactor\s+Authentication$', name, re.IGNORECASE)
    if m: return f"安装用于多因素认证的智能卡软件包"

    # --- "Enable Smart Card Logins in PAM" ---
    m = re.match(r'^Enable\s+Smart\s+Card\s+Logins\s+in\s+PAM$', name, re.IGNORECASE)
    if m: return f"在 PAM 中启用智能卡登录"

    # --- "Disable debug-shell SystemD Service" ---
    m = re.match(r'^Disable\s+debug-shell\s+SystemD\s+Service$', name, re.IGNORECASE)
    if m: return f"禁用 debug-shell SystemD 服务"

    # --- "Enable authselect" ---
    m = re.match(r'^Enable\s+authselect$', name, re.IGNORECASE)
    if m: return f"启用 authselect"

    # --- "Install the pam_apparmor Package" ---
    m = re.match(r'^Install\s+the\s+(.+?)\s+Package$', name, re.IGNORECASE)
    if m: return f"安装 {m.group(1)} 软件包"

    # --- "Enable Kernel Page-Table Isolation (KPTI)" ---
    m = re.match(r'^Enable\s+Kernel\s+Page-Table\s+Isolation\s+\(KPTI\)$', name, re.IGNORECASE)
    if m: return f"启用内核页表隔离（KPTI）"

    # --- "Disable vsyscalls" ---
    m = re.match(r'^Disable\s+vsyscalls$', name, re.IGNORECASE)
    if m: return f"禁用 vsyscalls"

    # --- "Disable Recovery Booting" ---
    m = re.match(r'^Disable\s+Recovery\s+Booting$', name, re.IGNORECASE)
    if m: return f"禁用恢复模式启动"

    # --- "Enable randomization of the page allocator" ---
    m = re.match(r'^Enable\s+randomization\s+of\s+the\s+page\s+allocator$', name, re.IGNORECASE)
    if m: return f"启用页分配器随机化"

    # --- "Disable merging of slabs with similar size" ---
    m = re.match(r'^Disable\s+merging\s+of\s+slabs\s+with\s+similar\s+size$', name, re.IGNORECASE)
    if m: return f"禁用相似大小 slab 的合并"

    # --- "Verify X Group/User Ownership/Permissions" 模式 ---
    m = re.match(r'^Verify\s+(the\s+)?(UEFI\s+)?(Boot\s+)?(Loader\s+)?(.+?)\s+(Group|User)\s+Ownership$', name, re.IGNORECASE)
    if m: return f"验证 {m.group(5)} 的{'组' if m.group(6)=='Group' else '用户'}所有权"
    m = re.match(r'^Verify\s+(the\s+)?(UEFI\s+)?(Boot\s+)?(Loader\s+)?(.+?)\s+Permissions$', name, re.IGNORECASE)
    if m: return f"验证 {m.group(5)} 的权限"

    # --- "Set Boot Loader Password in grub2" ---
    m = re.match(r'^Set\s+Boot\s+Loader\s+Password\s+in\s+grub2$', name, re.IGNORECASE)
    if m: return f"在 grub2 中设置引导加载程序密码"

    # --- "Enable X in zIPL" 模式 ---
    m = re.match(r'^Enable\s+(.+?)\s+in\s+zIPL$', name, re.IGNORECASE)
    if m: return f"在 zIPL 中启用 {m.group(1)}"

    # --- "Disable X in zIPL" 模式 ---
    m = re.match(r'^Disable\s+(.+?)\s+in\s+zIPL$', name, re.IGNORECASE)
    if m: return f"在 zIPL 中禁用 {m.group(1)}"

    # --- "Disable kernel support for X" 模式 ---
    m = re.match(r'^Disable\s+kernel\s+support\s+for\s+(.+)$', name, re.IGNORECASE)
    if m: return f"禁用内核支持 {m.group(1)}"

    # --- "Enable support for X" 模式 ---
    m = re.match(r'^Enable\s+support\s+for\s+(.+)$', name, re.IGNORECASE)
    if m: return f"启用对 {m.group(1)} 的支持"

    # --- "Disable compatibility with X" 模式 ---
    m = re.match(r'^Disable\s+compatibility\s+with\s+(.+)$', name, re.IGNORECASE)
    if m: return f"禁用与 {m.group(1)} 的兼容性"

    # --- "Disable the X" 模式 ---
    m = re.match(r'^Disable\s+the\s+(.+)$', name, re.IGNORECASE)
    if m: return f"禁用 {m.group(1)}"

    # --- "Enable checks on X" 模式 ---
    m = re.match(r'^Enable\s+checks\s+on\s+(.+)$', name, re.IGNORECASE)
    if m: return f"启用对 {m.group(1)} 的检查"

    # --- "Disable /dev/X virtual device support" 模式 ---
    m = re.match(r'^Disable\s+(/\S+)\s+virtual\s+device\s+support$', name, re.IGNORECASE)
    if m: return f"禁用 {m.group(1)} 虚拟设备支持"

    # --- "Disable hibernation" ---
    m = re.match(r'^Disable\s+hibernation$', name, re.IGNORECASE)
    if m: return f"禁用休眠"

    # --- "Disable the IPv6 protocol" ---
    m = re.match(r'^Disable\s+the\s+IPv6\s+protocol$', name, re.IGNORECASE)
    if m: return f"禁用 IPv6 协议"

    # --- "Disable kexec system call" ---
    m = re.match(r'^Disable\s+kexec\s+system\s+call$', name, re.IGNORECASE)
    if m: return f"禁用 kexec 系统调用"

    # --- 无匹配时返回原文 ---
    return name


# ============================================================
# 描述翻译 - 基于关键词替换
# ============================================================
def translate_desc(text):
    if not text:
        return text
    # Long phrase replacements first
    long_phrases = {
        "At a minimum, the audit system should collect file permission changes for all users and root.":
            "审计系统应至少收集所有用户和 root 的文件权限变更事件。",
        "If the auditd daemon is configured to use the augenrules program to read audit rules during daemon startup (the default), add the following line to a file with suffix .rules in the directory /etc/audit/rules.d:":
            "如果 auditd 守护进程配置为使用 augenrules 程序在启动时读取审计规则（默认），请将以下行添加到 /etc/audit/rules.d/ 目录下后缀为 .rules 的文件中：",
        "If the system is 64 bit then also add the following line:":
            "如果系统是 64 位，还需添加以下行：",
        "The -S syscall flag selects the system call.":
            "-S syscall 标志用于选择系统调用。",
        "The -F perm=unauthorized_access flag filters for unsuccessful access attempts.":
            "-F perm=unauthorized_access 标志用于过滤未成功的访问尝试。",
        "The -F perm=modifies_permissions flag filters for permission modification attempts.":
            "-F perm=modifies_permissions 标志用于过滤权限修改尝试。",
        "The -F perm=modifies_ownership flag filters for ownership modification attempts.":
            "-F perm=modifies_ownership 标志用于过滤所有权修改尝试。",
        "The -F perm=creates_file flag filters for file creation attempts.":
            "-F perm=creates_file 标志用于过滤文件创建尝试。",
        "The -F perm=deletes_file flag filters for file deletion attempts.":
            "-F perm=deletes_file 标志用于过滤文件删除尝试。",
        "The -F perm=modifies_file flag filters for file modification attempts.":
            "-F perm=modifies_file 标志用于过滤文件修改尝试。",
        "The audit system should collect unauthorized access attempts to files.":
            "审计系统应收集对文件的未授权访问尝试。",
        "The audit system should collect file deletion events.":
            "审计系统应收集文件删除事件。",
        "The audit system should collect information on exporting to media.":
            "审计系统应收集导出到介质的信息。",
        "The audit system should collect system administrator actions.":
            "审计系统应收集系统管理员操作。",
        "The audit system should collect information on the use of privileged commands.":
            "审计系统应收集特权命令的使用信息。",
        "The audit system should collect information on kernel module loading and unloading.":
            "审计系统应收集内核模块加载和卸载的信息。",
        "The audit system should collect logon and logout events.":
            "审计系统应收集登录和注销事件。",
        "AppArmor provide Mandatory Access Controls.":
            "AppArmor 提供强制访问控制。",
        "Install the Asset Configuration Compliance Module (ACCM).":
            "安装资产配置合规模块 (ACCM)。",
        "Install the Policy Auditor (PA) Module.":
            "安装策略审计器 (PA) 模块。",
        "$ sudo service bluetooth stop":
            "蓝牙服务应禁用，使用命令 sudo service bluetooth stop 停止蓝牙服务。",
        'Configure the operating system to audit the execution of the partition management program "fdisk".':
            '配置操作系统以审计分区管理程序 "fdisk" 的执行。',
    }
    for eng, cn in long_phrases.items():
        if eng in text:
            text = text.replace(eng, cn)

    # Short keyword replacements
    keywords = {
        "audit system": "审计系统",
        "auditd daemon": "auditd 守护进程",
        "privileged commands": "特权命令",
        "privileged command": "特权命令",
        "kernel module": "内核模块",
        "unauthorized access": "未授权访问",
        "discretionary access control": "自主访问控制",
        "file permission": "文件权限",
        "file permissions": "文件权限",
        "ownership changes": "所有权变更",
        "permission changes": "权限变更",
        "access attempts": "访问尝试",
        "creation attempts": "创建尝试",
        "modification attempts": "修改尝试",
        "deletion attempts": "删除尝试",
        "logon and logout": "登录和注销",
        "exporting to media": "导出到介质",
        "system administrator": "系统管理员",
        "system administrators": "系统管理员",
        "all users and root": "所有用户和 root",
        "should be configured": "应配置",
        "should be enabled": "应启用",
        "should be disabled": "应禁用",
        "should be installed": "应安装",
        "should be removed": "应移除",
        "should be set": "应设置",
        "should be checked": "应检查",
        "should be verified": "应验证",
        "should be": "应",
        "is required": "是必需的",
        "is recommended": "建议",
        "is necessary": "是必要的",
        "is important": "是重要的",
        "is critical": "是关键的",
        "is enabled": "已启用",
        "is disabled": "已禁用",
        "is installed": "已安装",
        "is configured": "已配置",
        "is set": "已设置",
        "is not set": "未设置",
        "is not configured": "未配置",
        "is not installed": "未安装",
        "is not enabled": "未启用",
        "is not disabled": "未禁用",
        "by default": "默认情况下",
        "by user": "按用户",
        "by users": "按用户",
        "for all users": "对所有用户",
        "for each user": "对每个用户",
        "for the user": "对用户",
        "the following": "以下",
        "the above": "上述",
        "the appropriate": "适当的",
        "the correct": "正确的",
        "the same": "相同的",
        "the different": "不同的",
        "the new": "新的",
        "the old": "旧的",
        "the current": "当前的",
        "the existing": "现有的",
        "the required": "所需的",
        "the specified": "指定的",
        "the configured": "配置的",
        "the installed": "已安装的",
        "the enabled": "已启用的",
        "the disabled": "已禁用的",
        "the file": "文件",
        "the files": "文件",
        "the directory": "目录",
        "the directories": "目录",
        "the configuration": "配置",
        "the configuration file": "配置文件",
        "the configuration files": "配置文件",
        "the parameter": "参数",
        "the parameters": "参数",
        "the option": "选项",
        "the options": "选项",
        "the value": "值",
        "the values": "值",
        "the setting": "设置",
        "the settings": "设置",
        "the service": "服务",
        "the services": "服务",
        "the package": "软件包",
        "the packages": "软件包",
        "the network": "网络",
        "the firewall": "防火墙",
        "the SSH": "SSH",
        "the SSH server": "SSH 服务器",
        "the SSH client": "SSH 客户端",
        "the password": "密码",
        "the passwords": "密码",
        "the account": "账户",
        "the accounts": "账户",
        "the user": "用户",
        "the users": "用户",
        "the group": "组",
        "the groups": "组",
        "the permission": "权限",
        "the permissions": "权限",
        "the ownership": "所有权",
        "the audit": "审计",
        "the audit log": "审计日志",
        "the audit logs": "审计日志",
        "the audit record": "审计记录",
        "the audit records": "审计记录",
        "the audit rule": "审计规则",
        "the audit rules": "审计规则",
        "the kernel": "内核",
        "the system call": "系统调用",
        "the system calls": "系统调用",
        "the module": "模块",
        "the modules": "模块",
        "the process": "进程",
        "the processes": "进程",
        "the port": "端口",
        "the ports": "端口",
        "the protocol": "协议",
        "the protocols": "协议",
        "the address": "地址",
        "the addresses": "地址",
        "the host": "主机",
        "the hosts": "主机",
        "the IP": "IP",
        "the TCP": "TCP",
        "the UDP": "UDP",
        "the HTTP": "HTTP",
        "the HTTPS": "HTTPS",
        "the DNS": "DNS",
        "the NTP": "NTP",
        "the LDAP": "LDAP",
        "the Kerberos": "Kerberos",
        "the SELinux": "SELinux",
        "the AppArmor": "AppArmor",
        "the firewall": "防火墙",
        "the iptables": "iptables",
        "the nftables": "nftables",
        "the firewalld": "firewalld",
        "the cron": "cron",
        "the crontab": "crontab",
        "the syslog": "syslog",
        "the rsyslog": "rsyslog",
        "the syslog-ng": "syslog-ng",
        "the journald": "journald",
        "the logrotate": "logrotate",
        "the tmpfs": "tmpfs",
    }
    # Sort by length descending to match longest first
    for eng, cn in sorted(keywords.items(), key=lambda x: -len(x[0])):
        text = text.replace(eng, cn)

    # ========== 正则模式匹配描述 ==========
    patterns = [
        # "All X must be Y" -> "所有 X 必须 Y"
        (r'^All\s+(.+?)\s+must\s+be\s+(.+)$', r'所有 \1 必须 \2'),
        # "All X must have Y" -> "所有 X 必须具有 Y"
        (r'^All\s+(.+?)\s+must\s+have\s+(.+)$', r'所有 \1 必须具有 \2'),
        # "Configure the X to Y" -> "配置 X 以 Y"
        (r'^Configure\s+the\s+(.+?)\s+to\s+(.+)$', r'配置 \1 以 \2'),
        # "Configure X to Y" -> "配置 X 以 Y"
        (r'^Configure\s+(.+?)\s+to\s+(.+)$', r'配置 \1 以 \2'),
        # "The X service is Y" -> "X 服务是 Y"
        (r'^The\s+(.+?)\s+(service|daemon|system)\s+is\s+(.+)$', r'\1 \2 是 \3'),
        # "The X is Y" -> "X 是 Y"
        (r'^The\s+(.+?)\s+is\s+(.+)$', r'\1 是 \2'),
        # "If the system does not need to have X" -> "如果系统不需要 X"
        (r'^If\s+the\s+system\s+does\s+not\s+need\s+to\s+have\s+(.+)$', r'如果系统不需要 \1'),
        # "If the system does not need X" -> "如果系统不需要 X"
        (r'^If\s+the\s+system\s+does\s+not\s+need\s+(.+)$', r'如果系统不需要 \1'),
        # "To configure X to Y" -> "配置 X 以 Y"
        (r'^To\s+configure\s+(.+?)\s+to\s+(.+)$', r'配置 \1 以 \2'),
        # "To configure X" -> "配置 X"
        (r'^To\s+configure\s+(.+)$', r'配置 \1'),
        # "X provides Y" -> "X 提供 Y"
        (r'^(.+?)\s+provides\s+(.+)$', r'\1 提供 \2'),
        # "X is an essential Y" -> "X 是重要的 Y"
        (r'^(.+?)\s+is\s+an\s+essential\s+(.+)$', r'\1 是重要的 \2'),
        # "X is a Y that can be used to Z" -> "X 是可用于 Z 的 Y"
        (r'^(.+?)\s+is\s+a\s+(.+?)\s+that\s+can\s+be\s+used\s+to\s+(.+)$', r'\1 是可用于 \3 的 \2'),
        # "X is a Y that Z" -> "X 是 Z 的 Y"
        (r'^(.+?)\s+is\s+a\s+(.+?)\s+that\s+(.+)$', r'\1 是 \3 的 \2'),
        # "X attempts to Y" -> "X 尝试 Y"
        (r'^(.+?)\s+attempts\s+to\s+(.+)$', r'\1 尝试 \2'),
        # "X exists to Y" -> "X 用于 Y"
        (r'^(.+?)\s+exists\s+to\s+(.+)$', r'\1 用于 \2'),
        # "X periodically checks for Y" -> "X 定期检查 Y"
        (r'^(.+?)\s+periodically\s+checks\s+for\s+(.+)$', r'\1 定期检查 \2'),
        # "X handles Y on behalf of Z" -> "X 代表 Z 处理 Y"
        (r'^(.+?)\s+handles\s+(.+?)\s+on\s+behalf\s+of\s+(.+)$', r'\1 代表 \3 处理 \2'),
        # "X resets Y in order to Z" -> "X 重置 Y 以 Z"
        (r'^(.+?)\s+resets\s+(.+?)\s+in\s+order\s+to\s+(.+)$', r'\1 重置 \2 以 \3'),
        # "X dispatches Y" -> "X 分发 Y"
        (r'^(.+?)\s+dispatches\s+(.+)$', r'\1 分发 \2'),
        # "X modifies Y" -> "X 修改 Y"
        (r'^(.+?)\s+modifies\s+(.+)$', r'\1 修改 \2'),
        # "X is responsible for Y" -> "X 负责 Y"
        (r'^(.+?)\s+is\s+responsible\s+for\s+(.+)$', r'\1 负责 \2'),
        # "First, X" -> "首先，X"
        (r'^First,\s+(.+)$', r'首先，\1'),
        # "By default, X" -> "默认情况下，X"
        (r'^By\s+default,\s+(.+)$', r'默认情况下，\1'),
        # "Set X to Y in Z" -> "在 Z 中将 X 设置为 Y"
        (r'^Set\s+(.+?)\s+to\s+(.+?)\s+in\s+(.+)$', r'在 \3 中将 \1 设置为 \2'),
        # "This is the default setting" -> "这是默认设置"
        (r'^This\s+is\s+the\s+default\s+setting', r'这是默认设置'),
        # "This is the default" -> "这是默认值"
        (r'^This\s+is\s+the\s+default', r'这是默认值'),
    ]
    for pattern, replacement in patterns:
        try:
            text = re.sub(pattern, replacement, text, flags=re.IGNORECASE)
        except:
            pass

    # ========== 更多正则模式（第二组，处理剩余描述） ==========
    patterns2 = [
        # "Set -c flag so that X" -> "设置 -c 标志以使 X"
        (r'^Set\s+-c\s+flag\s+so\s+that\s+(.+)$', r'设置 -c 标志以使 \1'),
        # "The file X should not exist" -> "文件 X 不应存在"
        (r'^The\s+file\s+(.+?)\s+should\s+not\s+exist', r'文件 \1 不应存在'),
        # "The file X should be Y" -> "文件 X 应 Y"
        (r'^The\s+file\s+(.+?)\s+should\s+be\s+(.+)$', r'文件 \1 应 \2'),
        # "If X exists, it must be Y" -> "如果 X 存在，则必须 Y"
        (r'^If\s+(.+?)\s+exists,\s+it\s+must\s+be\s+(.+)$', r'如果 \1 存在，则必须 \2'),
        # "If X exists, it must have Y" -> "如果 X 存在，则必须具有 Y"
        (r'^If\s+(.+?)\s+exists,\s+it\s+must\s+have\s+(.+)$', r'如果 \1 存在，则必须具有 \2'),
        # "X must be configured to Y" -> "X 必须配置为 Y"
        (r'^(.+?)\s+must\s+be\s+configured\s+to\s+(.+)$', r'\1 必须配置为 \2'),
        # "X should not be Y" -> "X 不应 Y"
        (r'^(.+?)\s+should\s+not\s+be\s+(.+)$', r'\1 不应 \2'),
        # "X should be Y" -> "X 应 Y"
        (r'^(.+?)\s+should\s+be\s+(.+)$', r'\1 应 \2'),
        # "X can be Y" -> "X 可以 Y"
        (r'^(.+?)\s+can\s+be\s+(.+)$', r'\1 可以 \2'),
        # "X prevents Y from Z" -> "X 防止 Y Z"
        (r'^(.+?)\s+prevents\s+(.+?)\s+from\s+(.+)$', r'\1 防止 \2 \3'),
        # "X allows Y to Z" -> "X 允许 Y Z"
        (r'^(.+?)\s+allows\s+(.+?)\s+to\s+(.+)$', r'\1 允许 \2 \3'),
        # "X is used to Y" -> "X 用于 Y"
        (r'^(.+?)\s+is\s+used\s+to\s+(.+)$', r'\1 用于 \2'),
        # "X relies on Y" -> "X 依赖 Y"
        (r'^(.+?)\s+relies\s+on\s+(.+)$', r'\1 依赖 \2'),
        # "X manages Y" -> "X 管理 Y"
        (r'^(.+?)\s+manages\s+(.+)$', r'\1 管理 \2'),
        # "X contains Y" -> "X 包含 Y"
        (r'^(.+?)\s+contains\s+(.+)$', r'\1 包含 \2'),
        # "X can be used to Y" -> "X 可用于 Y"
        (r'^(.+?)\s+can\s+be\s+used\s+to\s+(.+)$', r'\1 可用于 \2'),
        # "X requires Y" -> "X 需要 Y"
        (r'^(.+?)\s+requires\s+(.+)$', r'\1 需要 \2'),
        # "X does not Y" -> "X 不 Y"
        (r'^(.+?)\s+does\s+not\s+(.+)$', r'\1 不 \2'),
        # "X will Y" -> "X 将 Y"
        (r'^(.+?)\s+will\s+(.+)$', r'\1 将 \2'),
        # "X may Y" -> "X 可能 Y"
        (r'^(.+?)\s+may\s+(.+)$', r'\1 可能 \2'),
        # "X must Y" -> "X 必须 Y"
        (r'^(.+?)\s+must\s+(.+)$', r'\1 必须 \2'),
        # "X can Y" -> "X 可以 Y"
        (r'^(.+?)\s+can\s+(.+)$', r'\1 可以 \2'),
        # "X uses Y" -> "X 使用 Y"
        (r'^(.+?)\s+uses\s+(.+)$', r'\1 使用 \2'),
        # "X supports Y" -> "X 支持 Y"
        (r'^(.+?)\s+supports\s+(.+)$', r'\1 支持 \2'),
        # "X includes Y" -> "X 包含 Y"
        (r'^(.+?)\s+includes\s+(.+)$', r'\1 包含 \2'),
        # "X ensures Y" -> "X 确保 Y"
        (r'^(.+?)\s+ensures\s+(.+)$', r'\1 确保 \2'),
        # "X controls Y" -> "X 控制 Y"
        (r'^(.+?)\s+controls\s+(.+)$', r'\1 控制 \2'),
        # "X determines Y" -> "X 确定 Y"
        (r'^(.+?)\s+determines\s+(.+)$', r'\1 确定 \2'),
        # "X specifies Y" -> "X 指定 Y"
        (r'^(.+?)\s+specifies\s+(.+)$', r'\1 指定 \2'),
        # "X defines Y" -> "X 定义 Y"
        (r'^(.+?)\s+defines\s+(.+)$', r'\1 定义 \2'),
        # "X enables Y" -> "X 启用 Y"
        (r'^(.+?)\s+enables\s+(.+)$', r'\1 启用 \2'),
        # "X disables Y" -> "X 禁用 Y"
        (r'^(.+?)\s+disables\s+(.+)$', r'\1 禁用 \2'),
        # "X restricts Y" -> "X 限制 Y"
        (r'^(.+?)\s+restricts\s+(.+)$', r'\1 限制 \2'),
        # "X limits Y" -> "X 限制 Y"
        (r'^(.+?)\s+limits\s+(.+)$', r'\1 限制 \2'),
        # "X prevents Y" -> "X 防止 Y"
        (r'^(.+?)\s+prevents\s+(.+)$', r'\1 防止 \2'),
        # "X protects Y" -> "X 保护 Y"
        (r'^(.+?)\s+protects\s+(.+)$', r'\1 保护 \2'),
        # "X provides Y" -> "X 提供 Y"
        (r'^(.+?)\s+provides\s+(.+)$', r'\1 提供 \2'),
        # "X represents Y" -> "X 表示 Y"
        (r'^(.+?)\s+represents\s+(.+)$', r'\1 表示 \2'),
        # "X indicates Y" -> "X 指示 Y"
        (r'^(.+?)\s+indicates\s+(.+)$', r'\1 指示 \2'),
        # "X checks Y" -> "X 检查 Y"
        (r'^(.+?)\s+checks\s+(.+)$', r'\1 检查 \2'),
        # "X verifies Y" -> "X 验证 Y"
        (r'^(.+?)\s+verifies\s+(.+)$', r'\1 验证 \2'),
        # "X monitors Y" -> "X 监控 Y"
        (r'^(.+?)\s+monitors\s+(.+)$', r'\1 监控 \2'),
        # "X logs Y" -> "X 记录 Y"
        (r'^(.+?)\s+logs\s+(.+)$', r'\1 记录 \2'),
        # "X records Y" -> "X 记录 Y"
        (r'^(.+?)\s+records\s+(.+)$', r'\1 记录 \2'),
        # "X stores Y" -> "X 存储 Y"
        (r'^(.+?)\s+stores\s+(.+)$', r'\1 存储 \2'),
        # "X sends Y" -> "X 发送 Y"
        (r'^(.+?)\s+sends\s+(.+)$', r'\1 发送 \2'),
        # "X receives Y" -> "X 接收 Y"
        (r'^(.+?)\s+receives\s+(.+)$', r'\1 接收 \2'),
        # "X generates Y" -> "X 生成 Y"
        (r'^(.+?)\s+generates\s+(.+)$', r'\1 生成 \2'),
        # "X processes Y" -> "X 处理 Y"
        (r'^(.+?)\s+processes\s+(.+)$', r'\1 处理 \2'),
        # "X implements Y" -> "X 实现 Y"
        (r'^(.+?)\s+implements\s+(.+)$', r'\1 实现 \2'),
        # "X performs Y" -> "X 执行 Y"
        (r'^(.+?)\s+performs\s+(.+)$', r'\1 执行 \2'),
        # "X configures Y" -> "X 配置 Y"
        (r'^(.+?)\s+configures\s+(.+)$', r'\1 配置 \2'),
        # "X sets Y" -> "X 设置 Y"
        (r'^(.+?)\s+sets\s+(.+)$', r'\1 设置 \2'),
        # "X enables Y to Z" -> "X 使 Y 能够 Z"
        (r'^(.+?)\s+enables\s+(.+?)\s+to\s+(.+)$', r'\1 使 \2 能够 \3'),
        # "X allows Y to Z" -> "X 允许 Y Z"
        (r'^(.+?)\s+allows\s+(.+?)\s+to\s+(.+)$', r'\1 允许 \2 \3'),
        # "X requires Y to Z" -> "X 要求 Y Z"
        (r'^(.+?)\s+requires\s+(.+?)\s+to\s+(.+)$', r'\1 要求 \2 \3'),
        # "X is a Y that Z" -> "X 是 Z 的 Y"
        (r'^(.+?)\s+is\s+a\s+(.+?)\s+that\s+(.+)$', r'\1 是 \3 的 \2'),
        # "X is an Y that Z" -> "X 是 Z 的 Y"
        (r'^(.+?)\s+is\s+an\s+(.+?)\s+that\s+(.+)$', r'\1 是 \3 的 \2'),
        # "X is the Y" -> "X 是 Y"
        (r'^(.+?)\s+is\s+the\s+(.+)$', r'\1 是 \2'),
        # "X is not Y" -> "X 不是 Y"
        (r'^(.+?)\s+is\s+not\s+(.+)$', r'\1 不是 \2'),
        # "X are Y" -> "X 是 Y"
        (r'^(.+?)\s+are\s+(.+)$', r'\1 是 \2'),
        # "X has Y" -> "X 具有 Y"
        (r'^(.+?)\s+has\s+(.+)$', r'\1 具有 \2'),
        # "X have Y" -> "X 具有 Y"
        (r'^(.+?)\s+have\s+(.+)$', r'\1 具有 \2'),
        # "X with Y" -> "具有 Y 的 X"
        (r'^(.+?)\s+with\s+(.+)$', r'具有 \2 的 \1'),
        # "X for Y" -> "用于 Y 的 X"
        (r'^(.+?)\s+for\s+(.+)$', r'用于 \2 的 \1'),
        # "X in Y" -> "Y 中的 X"
        (r'^(.+?)\s+in\s+(.+)$', r'\2 中的 \1'),
        # "X on Y" -> "Y 上的 X"
        (r'^(.+?)\s+on\s+(.+)$', r'\2 上的 \1'),
        # "X by Y" -> "由 Y 的 X"
        (r'^(.+?)\s+by\s+(.+)$', r'由 \2 的 \1'),
        # "X from Y" -> "来自 Y 的 X"
        (r'^(.+?)\s+from\s+(.+)$', r'来自 \2 的 \1'),
        # "X to Y" -> "到 Y 的 X"
        (r'^(.+?)\s+to\s+(.+)$', r'到 \2 的 \1'),
        # "X of Y" -> "Y 的 X"
        (r'^(.+?)\s+of\s+(.+)$', r'\2 的 \1'),
        # "X and Y" -> "X 和 Y"
        (r'^(.+?)\s+and\s+(.+)$', r'\1 和 \2'),
        # "X or Y" -> "X 或 Y"
        (r'^(.+?)\s+or\s+(.+)$', r'\1 或 \2'),
    ]
    for pattern, replacement in patterns2:
        try:
            text = re.sub(pattern, replacement, text, flags=re.IGNORECASE)
        except:
            pass

    # 清理 placeholder_value 占位符
    text = re.sub(r'placeholder_value\s*', '', text)
    # 清理多余空格
    text = re.sub(r'\s+', ' ', text).strip()
    # 清理末尾的标点多余
    text = re.sub(r'[\.\s]+$', '。', text)
    if not text.endswith('。') and not text.endswith('。') and text:
        text += '。'

    return text


# ============================================================
# 风险描述翻译
# ============================================================
def translate_risk(risk_desc, name, risk_level):
    """Translate risk description to Chinese."""
    risk_names = {1: "高危", 2: "中危", 3: "低危"}
    rn = risk_names.get(risk_level, "未知")
    return f"{rn}风险: {name}。建议在适当时候修复此安全问题。"


# ============================================================
# 主处理逻辑
# ============================================================
translated = []
for r in rules:
    rule = dict(r)
    # Check if already Chinese
    name = rule.get('name', '')
    if not any('\u4e00' <= c <= '\u9fff' for c in name):
        rule['name'] = translate_name(name)
        rule['riskDescription'] = translate_risk(
            rule.get('riskDescription', ''),
            rule['name'],
            rule.get('risk', 2)
        )
    # 始终翻译描述（无论名称是否已翻译）
    desc = rule.get('description', '')
    if desc and not any('\u4e00' <= c <= '\u9fff' for c in desc):
        rule['description'] = translate_desc(desc)
    translated.append(rule)

# Save
output_path = 'data/baseline/compliance_rules_zh.json'
with open(output_path, 'w', encoding='utf-8') as f:
    json.dump(translated, f, ensure_ascii=False, indent=2)

# Stats
total = len(translated)
cn_names = sum(1 for r in translated if any('\u4e00' <= c <= '\u9fff' for c in r.get('name', '')))
print(f"Total rules: {total}")
print(f"Chinese names: {cn_names}")
print(f"Saved to: {output_path}")

# Sample
print("\n=== Sample translated rules ===")
for r in translated[:3]:
    print(f"  name: {r['name']}")
    print(f"  desc: {r['description'][:100]}...")
    print(f"  risk: {r['riskDescription']}")
    print()