"""
ComplianceAsCode -> Triad BaselineRule

Usage:
    python scripts/convert_compliance_rules.py <content_dir> [output_file]

Example:
    python scripts/convert_compliance_rules.py D:\shellprojects\content-master data/baseline/compliance_rules.json
"""

import json
import os
import re
import sys
import yaml
import xml.etree.ElementTree as ET

SEVERITY_MAP = {
    'low': 3, 'medium': 2, 'high': 1, 'critical': 0, 'unknown': 4,
}

CATEGORY_MAP = {
    'accounts': 1, 'accounts-pam': 1, 'accounts-physical': 2,
    'accounts-restrictions': 2, 'accounts-session': 2, 'accounts-banners': 2,
    'selinux': 4, 'bootloader-grub2': 4, 'bootloader-zipl': 4,
    'kernel_build_config': 4, 'kernel': 4,
    'permissions': 5, 'files': 5, 'partitions': 5, 'mounting': 5,
    'restrictions': 5, 'coredumps': 5,
    'logging': 6, 'journald': 6, 'audit': 6,
    'network': 7, 'network-kernel': 7, 'network-ipv6': 7,
    'network-firewalld': 3, 'network-iptables': 3, 'network-nftables': 3,
    'network-ufw': 3, 'network-wireless': 7, 'networkmanager': 7,
    'ssh': 9, 'ssh_server': 9, 'ssh_client': 9,
    'updates': 8, 'apparmor': 4, 'entropy': 4, 'secureboot': 4,
}

DEFAULT_CATEGORY = 99

JINJA_DEFAULTS = {
    'pam_lastlog_path': '/etc/pam.d/postlogin',
    'pam_lastlog': 'pam_lastlog.so',
    'control': '[default=1]',
    'full_name': 'the system',
    'rule_id': 'rule_001',
    'rule_title': 'Security Rule',
    'login_defs_path': '/etc/login.defs',
    'grub2_boot_path': '/boot/grub2',
    'grub2_uefi_boot_path': '/boot/efi/EFI/redhat',
    'sshd_sysconfig_file': '/etc/sysconfig/sshd',
    'openssh_client_crypto_policy_config_file': '/etc/ssh/ssh_config.d/openssh.config',
    'openssh_server_crypto_policy_config_file': '/etc/ssh/sshd_config.d/openssh.config',
    'pam_lastlog2': 'pam_lastlog2.so',
    'pam_lastlog2_path': '/etc/pam.d/postlogin-session',
    'after_match': '^\\s*session.*pam_succeed_if\\.so.*',
    'pam_files': 'password-auth system-auth',
    'pam_sections': 'auth account password session',
    'audit_rules_file': '/etc/audit/audit.rules',
    'audit_rules_dir': '/etc/audit/rules.d',
    'sshd_config': '/etc/ssh/sshd_config',
    'ssh_config': '/etc/ssh/ssh_config',
    'sysctl_default': '/etc/sysctl.conf',
    'sysctl_dir': '/etc/sysctl.d',
    'pwquality_conf': '/etc/security/pwquality.conf',
    'limits_conf': '/etc/security/limits.conf',
    'sudoers_file': '/etc/sudoers',
    'sudoers_dir': '/etc/sudoers.d',
    'shadow_file': '/etc/shadow',
    'passwd_file': '/etc/passwd',
    'group_file': '/etc/group',
    'fstab_file': '/etc/fstab',
    'hosts_allow': '/etc/hosts.allow',
    'hosts_deny': '/etc/hosts.deny',
    'nsswitch_conf': '/etc/nsswitch.conf',
    'resolv_conf': '/etc/resolv.conf',
    'cron_dir': '/etc/cron.d',
    'crontab_file': '/etc/crontab',
    'logrotate_conf': '/etc/logrotate.conf',
    'logrotate_dir': '/etc/logrotate.d',
    'rsyslog_conf': '/etc/rsyslog.conf',
    'rsyslog_dir': '/etc/rsyslog.d',
    'profile_dir': '/etc/profile.d',
    'bashrc_file': '/etc/bashrc',
    'bash_profile': '/etc/profile',
    'environment_file': '/etc/environment',
    'securetty_file': '/etc/securetty',
    'issue_file': '/etc/issue',
    'issue_net_file': '/etc/issue.net',
    'motd_file': '/etc/motd',
    'gdm_conf': '/etc/gdm/custom.conf',
    'dconf_db': '/etc/dconf/db',
    'dconf_profile': '/etc/dconf/profile',
    'crypto_policy': '/etc/crypto-policies/config',
    'crypto_policy_dir': '/etc/crypto-policies/policies',
    'sub_policy': '/etc/crypto-policies/policies/modules',
    'ntp_conf': '/etc/ntp.conf',
    'chrony_conf': '/etc/chrony.conf',
    'timesyncd_conf': '/etc/systemd/timesyncd.conf',
    'journald_conf': '/etc/systemd/journald.conf',
    'system_conf': '/etc/systemd/system.conf',
    'user_conf': '/etc/systemd/user.conf',
    'logind_conf': '/etc/systemd/logind.conf',
    'resolved_conf': '/etc/systemd/resolved.conf',
    'coredump_conf': '/etc/systemd/coredump.conf',
    'sleep_conf': '/etc/systemd/sleep.conf',
    'tmpfiles_dir': '/etc/tmpfiles.d',
    'sysusers_dir': '/etc/sysusers.d',
    'modules_load_dir': '/etc/modules-load.d',
    'modprobe_dir': '/etc/modprobe.d',
    'depmod_dir': '/etc/depmod.d',
    'udev_dir': '/etc/udev/rules.d',
    'pam_d_dir': '/etc/pam.d',
    'security_dir': '/etc/security',
    'skel_dir': '/etc/skel',
    'X11_dir': '/etc/X11',
    'fonts_dir': '/etc/fonts',
    'ldap_conf': '/etc/ldap/ldap.conf',
    'nslcd_conf': '/etc/nslcd.conf',
    'nscd_conf': '/etc/nscd.conf',
    'autofs_conf': '/etc/autofs.conf',
    'auto_master': '/etc/auto.master',
    'smb_conf': '/etc/samba/smb.conf',
    'vsftpd_conf': '/etc/vsftpd/vsftpd.conf',
    'httpd_conf_dir': '/etc/httpd/conf.d',
    'httpd_conf': '/etc/httpd/conf/httpd.conf',
    'nginx_conf': '/etc/nginx/nginx.conf',
    'nginx_conf_dir': '/etc/nginx/conf.d',
    'mail_alias': '/etc/aliases',
    'postfix_main_cf': '/etc/postfix/main.cf',
    'postfix_master_cf': '/etc/postfix/master.cf',
    'dhcpd_conf': '/etc/dhcp/dhcpd.conf',
    'named_conf': '/etc/named.conf',
    'snmpd_conf': '/etc/snmp/snmpd.conf',
    'syslog_conf': '/etc/syslog.conf',
    'syslog_ng_conf': '/etc/syslog-ng/syslog-ng.conf',
    'yum_conf': '/etc/yum.conf',
    'yum_repos_dir': '/etc/yum.repos.d',
    'dnf_conf': '/etc/dnf/dnf.conf',
    'apt_conf_dir': '/etc/apt/apt.conf.d',
    'sources_list': '/etc/apt/sources.list',
    'sources_list_d': '/etc/apt/sources.list.d',
    'zypp_repos_dir': '/etc/zypp/repos.d',
    'zypp_conf': '/etc/zypp/zypp.conf',
    'pam_pwquality_conf': '/etc/security/pwquality.conf',
    'pam_faillock_conf': '/etc/security/faillock.conf',
    'pam_access_conf': '/etc/security/access.conf',
    'pam_time_conf': '/etc/security/time.conf',
    'pam_limits_conf': '/etc/security/limits.conf',
    'pam_namespace_conf': '/etc/security/namespace.conf',
    'pam_namespace_dir': '/etc/security/namespace.d',
    'sealert_conf': '/etc/sealert.conf',
    'sestatus_conf': '/etc/sestatus.conf',
    'selinux_config': '/etc/selinux/config',
    'fips_mode_file': '/etc/system-fips',
    'prelink_conf': '/etc/prelink.conf',
    'prelink_conf_dir': '/etc/prelink.conf.d',
    'kdump_conf': '/etc/kdump.conf',
    'kdump_sysconfig': '/etc/sysconfig/kdump',
    'network_scripts': '/etc/sysconfig/network-scripts',
    'network_config': '/etc/sysconfig/network',
    'ifcfg_dir': '/etc/sysconfig/network-scripts',
    'sysconfig_dir': '/etc/sysconfig',
    'default_dir': '/etc/default',
    'grub_default': '/etc/default/grub',
    'grub_dir': '/boot/grub2',
    'grub_efi_dir': '/boot/efi/EFI/redhat',
    'kernel_opts': '/etc/kernel/cmdline',
    'proc_cmdline': '/proc/cmdline',
    'sysctl_kernel_file': '/etc/sysctl.d/99-kernel.conf',
    'sysctl_net_file': '/etc/sysctl.d/99-network.conf',
    'auditd_conf': '/etc/audit/auditd.conf',
    'audisp_dir': '/etc/audisp',
    'audisp_plugins_dir': '/etc/audisp/plugins.d',
    'sudo_lecture_file': '/etc/sudoers.d/lecture',
    'sudo_log_file': '/etc/sudoers.d/log',
    'sudo_timestamp_file': '/etc/sudoers.d/timestamp',
    'sudo_umask_file': '/etc/sudoers.d/umask',
    'sudo_requiretty_file': '/etc/sudoers.d/requiretty',
    'sudo_passwd_timeout_file': '/etc/sudoers.d/passwd_timeout',
    'sudo_ticket_timeout_file': '/etc/sudoers.d/ticket_timeout',
    'sudo_max_attempts_file': '/etc/sudoers.d/max_attempts',
    'sudo_ignore_local_shell_file': '/etc/sudoers.d/ignore_local_shell',
    'sudo_use_pty_file': '/etc/sudoers.d/use_pty',
    'sudo_env_reset_file': '/etc/sudoers.d/env_reset',
    'sudo_secure_path_file': '/etc/sudoers.d/secure_path',
    'sudo_secure_path': '/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin',
    'sudo_mail_badpass_file': '/etc/sudoers.d/mail_badpass',
    'sudo_mail_always_file': '/etc/sudoers.d/mail_always',
    'sudo_noexec_file': '/etc/sudoers.d/noexec',
    'auid': '1000',
    'auid_threshold': '1000',
    'ansible_facts': '',
    'product': 'rhel8',
    'families': 'rhel',
}


def preprocess_yaml(content):
    """Preprocess YAML content to remove Jinja2 syntax before parsing."""
    lines = content.split('\n')
    new_lines = []
    for line in lines:
        stripped = line.strip()
        if line.startswith('{{{') and stripped.endswith('}}}'):
            continue
        if line.startswith('{{%') and stripped.endswith('%}}'):
            continue
        if line.startswith('{{') and stripped.endswith('}}') and not line.startswith('{{ '):
            continue
        new_lines.append(line)
    content = '\n'.join(new_lines)
    content = re.sub(r'\{\{\{[^}]*\}\}\}', 'placeholder_value', content)
    content = re.sub(r'\{\%[^%]*\%\}', '', content)
    content = re.sub(r'\{\{[^}]*\}\}', 'placeholder_value', content)
    return content


def replace_jinja_vars(text):
    if not text:
        return text

    def replace_var(match):
        var_name = match.group(1).strip()
        return JINJA_DEFAULTS.get(var_name, var_name)

    text = re.sub(r'\{\{\{\s*(\w+)\s*\}\}\}', replace_var, text)
    text = re.sub(r'\{\%-\s*if\s+.*?\%-\s*\}.*?\{\%-\s*endif\s*\%-\s*\}', '', text, flags=re.DOTALL)
    text = re.sub(r'\{\%-\s*for\s+.*?\%-\s*\}.*?\{\%-\s*endfor\s*\%-\s*\}', '', text, flags=re.DOTALL)
    text = re.sub(r'\{\%-\s*macro\s+.*?\%-\s*\}.*?\{\%-\s*endmacro\s*\%-\s*\}', '', text, flags=re.DOTALL)
    text = re.sub(r'\{\%[+-]?\s*\w+.*?[+-]?\%\}', '', text)
    text = re.sub(r'\{\{\{\s*.*?\s*\}\}\}', '', text)
    text = re.sub(r'\{\{\s*.*?\s*\}\}', '', text)
    text = re.sub(r'\n\s*\n', '\n', text)
    text = re.sub(r' +', ' ', text)
    return text.strip()


def clean_html(text):
    if not text:
        return text
    text = re.sub(r'<[^>]+>', '', text)
    text = re.sub(r'&lt;', '<', text)
    text = re.sub(r'&gt;', '>', text)
    text = re.sub(r'&amp;', '&', text)
    text = re.sub(r'&quot;', '"', text)
    text = re.sub(r'&#39;', "'", text)
    return text


def clean_text(text):
    if not text:
        return text
    text = clean_html(text)
    text = replace_jinja_vars(text)
    text = re.sub(r'\s+', ' ', text)
    return text.strip()


OVAL_NS = {
    'oval': 'http://oval.mitre.org/XMLSchema/oval-definitions-5',
    'ind': 'http://oval.mitre.org/XMLSchema/oval-definitions-5#independent',
    'unix': 'http://oval.mitre.org/XMLSchema/oval-definitions-5#unix',
    'linux': 'http://oval.mitre.org/XMLSchema/oval-definitions-5#linux',
}


def oval_to_shell_commands(oval_xml):
    commands = []
    expected_value = ""
    match_type = "contains"

    try:
        root = ET.fromstring(oval_xml)
    except ET.ParseError:
        return commands, expected_value, match_type

    for obj in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#independent}textfilecontent54_object'):
        filepath_el = obj.find('ind:filepath', OVAL_NS)
        pattern_el = obj.find('ind:pattern', OVAL_NS)
        if filepath_el is not None and pattern_el is not None:
            filepath = replace_jinja_vars(filepath_el.text or "")
            pattern_text = pattern_el.text or ""
            pattern_text_clean = pattern_text.replace("'", "'\\''")
            cmd = f"grep -P '{pattern_text_clean}' {filepath} 2>/dev/null || echo 'NOT_FOUND'"
            commands.append(cmd)
            expected_value = pattern_text
            match_type = "regex"

    for obj in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#independent}textfilecontent_object'):
        filepath_el = obj.find('ind:filepath', OVAL_NS)
        pattern_el = obj.find('ind:pattern', OVAL_NS)
        if filepath_el is not None and pattern_el is not None:
            filepath = replace_jinja_vars(filepath_el.text or "")
            pattern_text = pattern_el.text or ""
            pattern_text_clean = pattern_text.replace("'", "'\\''")
            cmd = f"grep -P '{pattern_text_clean}' {filepath} 2>/dev/null || echo 'NOT_FOUND'"
            commands.append(cmd)
            expected_value = pattern_text
            match_type = "regex"

    for obj in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}file_object'):
        path_el = obj.find('unix:path', OVAL_NS)
        filename_el = obj.find('unix:filename', OVAL_NS)
        if path_el is not None and filename_el is not None:
            path = replace_jinja_vars(path_el.text or "")
            filename = replace_jinja_vars(filename_el.text or "")
            filepath = os.path.join(path, filename)
            cmd = f"stat -c '%a %U %G' '{filepath}' 2>/dev/null || echo 'NOT_FOUND'"
            commands.append(cmd)
            match_type = "contains"

    for obj in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}symlink_object'):
        path_el = obj.find('unix:filepath', OVAL_NS)
        if path_el is not None:
            filepath = replace_jinja_vars(path_el.text or "")
            cmd = f"readlink -f '{filepath}' 2>/dev/null || echo 'NOT_FOUND'"
            commands.append(cmd)
            match_type = "contains"

    for obj in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#linux}partition_object'):
        mp_el = obj.find('linux:mount_point', OVAL_NS)
        if mp_el is not None:
            mp = replace_jinja_vars(mp_el.text or "")
            cmd = f"mount | grep ' {mp} ' 2>/dev/null || echo 'NOT_FOUND'"
            commands.append(cmd)
            match_type = "contains"

    for obj in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}sysctl_object'):
        name_el = obj.find('unix:name', OVAL_NS)
        if name_el is not None:
            name = replace_jinja_vars(name_el.text or "")
            cmd = f"sysctl -n {name} 2>/dev/null || echo 'NOT_FOUND'"
            commands.append(cmd)
            match_type = "contains"

    for _ in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}inetlisteningservers_test'):
        cmd = "ss -tlnp 2>/dev/null || netstat -tlnp 2>/dev/null"
        if cmd not in commands:
            commands.append(cmd)
        match_type = "contains"

    for _ in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}process_test'):
        cmd = "ps -ef 2>/dev/null"
        if cmd not in commands:
            commands.append(cmd)
        match_type = "contains"

    for _ in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#linux}rpmverifyfile_test'):
        cmd = "rpm -Va 2>/dev/null | head -50"
        if cmd not in commands:
            commands.append(cmd)
        match_type = "contains"

    for obj in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#linux}systemdunitproperty_object'):
        unit_el = obj.find('linux:unit', OVAL_NS)
        prop_el = obj.find('linux:property', OVAL_NS)
        if unit_el is not None and prop_el is not None:
            unit = replace_jinja_vars(unit_el.text or "")
            prop = replace_jinja_vars(prop_el.text or "")
            cmd = f"systemctl show -p {prop} {unit}.service 2>/dev/null || echo 'NOT_FOUND'"
            commands.append(cmd)
            match_type = "contains"

    for obj in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#independent}environmentvariable_object'):
        name_el = obj.find('ind:name', OVAL_NS)
        if name_el is not None:
            name = replace_jinja_vars(name_el.text or "")
            cmd = f"echo \\${name} 2>/dev/null || echo 'NOT_FOUND'"
            commands.append(cmd)
            match_type = "contains"

    for obj in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}shadow_object'):
        username_el = obj.find('unix:username', OVAL_NS)
        if username_el is not None:
            username = replace_jinja_vars(username_el.text or "")
            cmd = f"grep '^{username}:' /etc/shadow 2>/dev/null || echo 'NOT_FOUND'"
            commands.append(cmd)
            match_type = "contains"

    for _ in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}runlevel_test'):
        cmd = "systemctl get-default 2>/dev/null || runlevel 2>/dev/null || echo 'NOT_FOUND'"
        if cmd not in commands:
            commands.append(cmd)
        match_type = "contains"

    for _ in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}password_test'):
        cmd = "cat /etc/passwd 2>/dev/null | awk -F: '{print $1\":\"$2}' | head -20"
        if cmd not in commands:
            commands.append(cmd)
        match_type = "contains"

    for obj in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#linux}rpminfo_object'):
        name_el = obj.find('linux:name', OVAL_NS)
        if name_el is not None:
            name = replace_jinja_vars(name_el.text or "")
            cmd = f"rpm -q --queryformat '%{VERSION}' {name} 2>/dev/null || echo 'NOT_INSTALLED'"
            commands.append(cmd)
            match_type = "contains"

    for obj in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#linux}selinux_object'):
        for state in ['enforce_enforcing', 'enforce_permissive', 'current_mode']:
            pass
        cmd = "getenforce 2>/dev/null || echo 'NOT_FOUND'"
        if cmd not in commands:
            commands.append(cmd)
        match_type = "contains"

    for _ in root.iter('{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}iflisteners_test'):
        cmd = "ss -tlnp 2>/dev/null || netstat -tlnp 2>/dev/null"
        if cmd not in commands:
            commands.append(cmd)
        match_type = "contains"

    return commands, expected_value, match_type


TEMPLATE_COMMANDS = {
    'kernel_module_disabled': {
        'cmd': "lsmod | grep -E '^{mod}$' 2>/dev/null || echo 'NOT_LOADED'",
        'expected': "NOT_LOADED", 'match': "contains",
        'var_map': {'mod': 'kernmodule'},
    },
    'package_removed': {
        'cmd': "rpm -q {pkg} 2>/dev/null || dpkg -l {pkg} 2>/dev/null || echo 'NOT_INSTALLED'",
        'expected': "NOT_INSTALLED", 'match': "contains",
        'var_map': {'pkg': 'pkgname'},
    },
    'package_installed': {
        'cmd': "rpm -q {pkg} 2>/dev/null || dpkg -l {pkg} 2>/dev/null || echo 'NOT_INSTALLED'",
        'expected': "{pkg}", 'match': "contains",
        'var_map': {'pkg': 'pkgname'},
    },
    'service_enabled': {
        'cmd': "systemctl is-enabled {service}.service 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {'service': 'servicename'},
    },
    'service_disabled': {
        'cmd': "systemctl is-enabled {service}.service 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "disabled", 'match': "contains",
        'var_map': {'service': 'servicename'},
    },
    'service_running': {
        'cmd': "systemctl is-active {service}.service 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "active", 'match': "contains",
        'var_map': {'service': 'servicename'},
    },
    'sshd_lineinfile': {
        'cmd': "sshd -T 2>/dev/null | grep -E '^{param}\\s' || grep -E '^{param}\\s+' /etc/ssh/sshd_config 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{param} {value}", 'match': "contains",
        'var_map': {'param': 'parameter', 'value': 'value'},
    },
    'sysctl': {
        'cmd': "sysctl -n {name} 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{value}", 'match': "contains",
        'var_map': {'name': 'name', 'value': 'value'},
    },
    'mount_option': {
        'cmd': "mount | grep ' {mp} ' 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{opt}", 'match': "contains",
        'var_map': {'mp': 'mount_point', 'opt': 'mount_option'},
    },
    'file_owner': {
        'cmd': "find {filepath} -maxdepth 1 -regextype posix-egrep -regex '{file_regex}' ! -user {uid} 2>/dev/null | head -5 | awk '{print} END {if(NR==0) print \"ALL_OK\"}'",
        'expected': "ALL_OK", 'match': "contains",
        'var_map': {'filepath': 'filepath', 'file_regex': 'file_regex', 'uid': 'uid_or_name'},
    },
    'file_permissions': {
        'cmd': "find {filepath} -maxdepth 1 -regextype posix-egrep -regex '{file_regex}' ! -perm {perms} 2>/dev/null | head -5 | awk '{print} END {if(NR==0) print \"ALL_OK\"}'",
        'expected': "ALL_OK", 'match': "contains",
        'var_map': {'filepath': 'filepath', 'file_regex': 'file_regex', 'perms': 'file_permissions'},
    },
    'file_groupowner': {
        'cmd': "find {filepath} -maxdepth 1 -regextype posix-egrep -regex '{file_regex}' ! -group {gid} 2>/dev/null | head -5 | awk '{print} END {if(NR==0) print \"ALL_OK\"}'",
        'expected': "ALL_OK", 'match': "contains",
        'var_map': {'filepath': 'filepath', 'file_regex': 'file_regex', 'gid': 'group_or_name'},
    },
    'sebool': {
        'cmd': "getsebool {bool} 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{value}", 'match': "contains",
        'var_map': {'bool': 'name', 'value': 'value'},
    },
    'kernel_build_config': {
        'cmd': "grep -E '^{config}=' /boot/config-$(uname -r) 2>/dev/null || grep -E '^{config}=' /proc/config.gz 2>/dev/null | zcat 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{config}=y", 'match': "contains",
        'var_map': {'config': 'config'},
    },
    'audit_rules_dac_modification': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'chmod|chown|fchmod|fchown|fremovexattr|fsetxattr|lchown|lremovexattr|lsetxattr|removexattr|setxattr' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "chmod", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_path_syscall': {
        'cmd': "auditctl -l 2>/dev/null | grep -E '{syscall}' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "{syscall}", 'match': "contains",
        'var_map': {'syscall': 'syscall', 'path': 'path'},
    },
    'mount': {
        'cmd': "mount | grep -E '{mount_point}' 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{mount_point}", 'match': "contains",
        'var_map': {'mount_point': 'mount_point'},
    },
    'coreos_kernel_option': {
        'cmd': "grep -E '{arg}' /proc/cmdline 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{arg}", 'match': "contains",
        'var_map': {'arg': 'arg'},
    },
    'zipl_bls_entries_option': {
        'cmd': "grep -E '{arg}' /boot/loader/entries/*.conf 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{arg}", 'match': "contains",
        'var_map': {'arg': 'arg'},
    },
    'auditd_lineinfile': {
        'cmd': "grep -E '{param}' /etc/audit/auditd.conf 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{param}", 'match': "contains",
        'var_map': {'param': 'parameter', 'value': 'value'},
    },
    'pam_account_password_faillock': {
        'cmd': "grep -E 'pam_faillock\\.so' /etc/pam.d/password-auth /etc/pam.d/system-auth 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "pam_faillock.so", 'match': "contains",
        'var_map': {},
    },
    'cis_banner': {
        'cmd': "cat /etc/issue 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "NOT_FOUND", 'match': "not_contains",
        'var_map': {},
    },
    'audit_rules_unsuccessful_file_modification_o_creat': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'O_CREAT' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "O_CREAT", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_unsuccessful_file_modification_o_trunc_write': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'O_TRUNC_WRITE|O_WRONLY' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "O_TRUNC_WRITE", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_unsuccessful_file_modification_rule_order': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'open.*trunc|open.*write' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "open", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'always,exit.*-F.*perm=x' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "perm=x", 'match': "contains",
        'var_map': {},
    },
    'accounts_password': {
        'cmd': "grep -E '^minlen|^dcredit|^ucredit|^lcredit|^ocredit|^minclass' /etc/security/pwquality.conf 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "minlen", 'match': "contains",
        'var_map': {},
    },
    'accounts_umask': {
        'cmd': "grep -E '^UMASK' /etc/login.defs 2>/dev/null || grep -E 'umask' /etc/profile /etc/bashrc 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "UMASK", 'match': "contains",
        'var_map': {},
    },
    'accounts_tmout': {
        'cmd': "grep -E '^TMOUT' /etc/profile /etc/profile.d/*.sh 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "TMOUT", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_chmod': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/chmod' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/chmod", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_chown': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/chown' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/chown", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_rm': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/rm' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/rm", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_su': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/su' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/su", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_sudo': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/sudo' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/sudo", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_mount': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/mount' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/mount", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_umount': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/umount' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/umount", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_ssh': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/ssh' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/ssh", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_passwd': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/passwd' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/passwd", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_crontab': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/crontab' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/crontab", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_pam': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/unix_chkpwd' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/unix_chkpwd", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_newgrp': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/newgrp' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/newgrp", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_gpasswd': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/gpasswd' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/gpasswd", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_chage': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/chage' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/chage", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_userhelper': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/userhelper' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/userhelper", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_usernetctl': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/usernetctl' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/usernetctl", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_postdrop': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/postdrop' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/postdrop", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_postqueue': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/postqueue' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/postqueue", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_ssh_agent': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/ssh-agent' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/ssh-agent", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_pam_timestamp': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/pam_timestamp_check' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/pam_timestamp_check", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_kmod': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/kmod' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/kmod", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_unsuccessful_file_modification': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'exit=-EACCES|exit=-EPERM' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "exit=-EACCES", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_file_deletion_events': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'S.*rmdir|S.*unlink' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "unlink", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_media_exports': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'mount|umount' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "mount", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_kernel_module_loading': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'init_module|finit_module|delete_module' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "init_module", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_login_events': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'faillog|lastlog|tallylog' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "faillog", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_time_watch': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'clock_settime|adjtimex|settimeofday' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "clock_settime", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_mac_modification': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'setxattr|lsetxattr|removexattr' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "setxattr", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_networkconfig_modification': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'sethostname|setdomainname|issue|hosts' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "sethostname", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_system_shutdown': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'reboot|halt|poweroff|shutdown' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "reboot", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_usergroup_modification': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'useradd|usermod|userdel|groupadd|groupmod|groupdel' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "useradd", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_session_events': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'loginuid|session' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "loginuid", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_suid_events': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'suid|sgid' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "suid", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_immutable': {
        'cmd': "grep -E '^-e 2' /etc/audit/audit.rules /etc/audit/rules.d/*.rules 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "-e 2", 'match': "contains",
        'var_map': {},
    },
    'grub2_bootloader_password': {
        'cmd': "grep -E '^password_pbkdf2|^GRUB2_PASSWORD' /boot/grub2/user.cfg /boot/grub2/grub.cfg 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "password_pbkdf2", 'match': "contains",
        'var_map': {},
    },
    'selinux_state': {
        'cmd': "getenforce 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "Enforcing", 'match': "contains",
        'var_map': {},
    },
    'selinux_policy': {
        'cmd': "sestatus 2>/dev/null | grep 'Loaded policy name' || grep '^SELINUXTYPE' /etc/selinux/config 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "targeted", 'match': "contains",
        'var_map': {},
    },
    'sshd_set_keepalive': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'clientaliveinterval\\|clientalivecountmax' || echo 'NOT_FOUND'",
        'expected': "ClientAliveInterval", 'match': "contains",
        'var_map': {},
    },
    'sshd_disable_root_login': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'permitrootlogin' || echo 'NOT_FOUND'",
        'expected': "PermitRootLogin no", 'match': "contains",
        'var_map': {},
    },
    'sshd_disable_empty_passwords': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'permitemptypasswords' || echo 'NOT_FOUND'",
        'expected': "PermitEmptyPasswords no", 'match': "contains",
        'var_map': {},
    },
    'sshd_allow_only_protocol2': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'protocol' || grep -i '^Protocol' /etc/ssh/sshd_config 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "Protocol 2", 'match': "contains",
        'var_map': {},
    },
    'sshd_use_approved_ciphers': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'ciphers' || echo 'NOT_FOUND'",
        'expected': "Ciphers", 'match': "contains",
        'var_map': {},
    },
    'sshd_use_approved_macs': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'macs' || echo 'NOT_FOUND'",
        'expected': "MACs", 'match': "contains",
        'var_map': {},
    },
    'sshd_use_approved_kex': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'kexalgorithms' || echo 'NOT_FOUND'",
        'expected': "KexAlgorithms", 'match': "contains",
        'var_map': {},
    },
    'sshd_print_last_log': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'printlastlog' || echo 'NOT_FOUND'",
        'expected': "PrintLastLog yes", 'match': "contains",
        'var_map': {},
    },
    'sshd_set_loglevel': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'loglevel' || echo 'NOT_FOUND'",
        'expected': "LogLevel VERBOSE", 'match': "contains",
        'var_map': {},
    },
    'sshd_set_max_auth_tries': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'maxauthtries' || echo 'NOT_FOUND'",
        'expected': "MaxAuthTries", 'match': "contains",
        'var_map': {},
    },
    'sshd_set_max_sessions': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'maxsessions' || echo 'NOT_FOUND'",
        'expected': "MaxSessions", 'match': "contains",
        'var_map': {},
    },
    'sshd_set_max_startups': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'maxstartups' || echo 'NOT_FOUND'",
        'expected': "MaxStartups", 'match': "contains",
        'var_map': {},
    },
    'sshd_idle_timeout': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'clientaliveinterval' || echo 'NOT_FOUND'",
        'expected': "ClientAliveInterval", 'match': "contains",
        'var_map': {},
    },
    'sshd_enable_warning_banner': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'banner' || echo 'NOT_FOUND'",
        'expected': "Banner", 'match': "contains",
        'var_map': {},
    },
    'sshd_disable_gssapi_auth': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'gssapiauthentication' || echo 'NOT_FOUND'",
        'expected': "GSSAPIAuthentication no", 'match': "contains",
        'var_map': {},
    },
    'sshd_disable_kerberos_auth': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'kerberosauthentication' || echo 'NOT_FOUND'",
        'expected': "KerberosAuthentication no", 'match': "contains",
        'var_map': {},
    },
    'sshd_disable_rhosts': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'rhosts' || grep -i 'ignoreRhosts' /etc/ssh/sshd_config 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "IgnoreRhosts yes", 'match': "contains",
        'var_map': {},
    },
    'sshd_disable_rhosts_rsa': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'rhostsrsaauthentication' || echo 'NOT_FOUND'",
        'expected': "RhostsRSAAuthentication no", 'match': "contains",
        'var_map': {},
    },
    'sshd_disable_hostbased_auth': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'hostbasedauthentication' || echo 'NOT_FOUND'",
        'expected': "HostbasedAuthentication no", 'match': "contains",
        'var_map': {},
    },
    'sshd_disable_compression': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'compression' || echo 'NOT_FOUND'",
        'expected': "Compression no", 'match': "contains",
        'var_map': {},
    },
    'sshd_disable_tcp_forwarding': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'tcpforwarding' || echo 'NOT_FOUND'",
        'expected': "AllowTcpForwarding no", 'match': "contains",
        'var_map': {},
    },
    'sshd_disable_x11_forwarding': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'x11forwarding' || echo 'NOT_FOUND'",
        'expected': "X11Forwarding no", 'match': "contains",
        'var_map': {},
    },
    'sshd_disable_agent_forwarding': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'allowagentforwarding' || echo 'NOT_FOUND'",
        'expected': "AllowAgentForwarding no", 'match': "contains",
        'var_map': {},
    },
    'sshd_enable_strictmodes': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'strictmodes' || echo 'NOT_FOUND'",
        'expected': "StrictModes yes", 'match': "contains",
        'var_map': {},
    },
    'sshd_enable_pam': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'usepam' || echo 'NOT_FOUND'",
        'expected': "UsePAM yes", 'match': "contains",
        'var_map': {},
    },
    'sshd_set_idle_timeout': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'clientaliveinterval' || echo 'NOT_FOUND'",
        'expected': "ClientAliveInterval 300", 'match': "contains",
        'var_map': {},
    },
    'sshd_enable_x11_forwarding': {
        'cmd': "sshd -T 2>/dev/null | grep -i 'x11forwarding' || echo 'NOT_FOUND'",
        'expected': "X11Forwarding yes", 'match': "contains",
        'var_map': {},
    },
    'accounts_maximum_age_login_defs': {
        'cmd': "grep -E '^PASS_MAX_DAYS' /etc/login.defs 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "PASS_MAX_DAYS", 'match': "contains",
        'var_map': {},
    },
    'accounts_minimum_age_login_defs': {
        'cmd': "grep -E '^PASS_MIN_DAYS' /etc/login.defs 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "PASS_MIN_DAYS", 'match': "contains",
        'var_map': {},
    },
    'accounts_password_warn_age_login_defs': {
        'cmd': "grep -E '^PASS_WARN_AGE' /etc/login.defs 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "PASS_WARN_AGE", 'match': "contains",
        'var_map': {},
    },
    'accounts_minlen_login_defs': {
        'cmd': "grep -E '^PASS_MIN_LEN' /etc/login.defs 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "PASS_MIN_LEN", 'match': "contains",
        'var_map': {},
    },
    'no_empty_passwords': {
        'cmd': "grep -E '^[^:]+::' /etc/shadow 2>/dev/null || echo 'ALL_ACCOUNTS_HAVE_PASSWORDS'",
        'expected': "ALL_ACCOUNTS_HAVE_PASSWORDS", 'match': "contains",
        'var_map': {},
    },
    'no_legacy_plus_entries': {
        'cmd': "grep '^\\+' /etc/passwd /etc/shadow /etc/group 2>/dev/null || echo 'NO_LEGACY_ENTRIES'",
        'expected': "NO_LEGACY_ENTRIES", 'match': "contains",
        'var_map': {},
    },
    'no_shelllogin_for_systemaccounts': {
        'cmd': "awk -F: '($3<1000){print $1\":\"$7}' /etc/passwd | grep -v 'nologin\\|/bin/false' 2>/dev/null || echo 'ALL_SYSTEM_ACCOUNTS_LOCKED'",
        'expected': "ALL_SYSTEM_ACCOUNTS_LOCKED", 'match': "contains",
        'var_map': {},
    },
    'root_unique_uid': {
        'cmd': "awk -F: '($3==0){print $1}' /etc/passwd 2>/dev/null",
        'expected': "root", 'match': "contains",
        'var_map': {},
    },
    'ensure_logrotate': {
        'cmd': "ls /etc/logrotate.d/ 2>/dev/null | head -20 || echo 'NO_LOGROTATE_CONFIG'",
        'expected': "NO_LOGROTATE_CONFIG", 'match': "not_contains",
        'var_map': {},
    },
    'ensure_rsyslog': {
        'cmd': "systemctl is-enabled rsyslog 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_auditd': {
        'cmd': "systemctl is-enabled auditd 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_crond': {
        'cmd': "systemctl is-enabled crond 2>/dev/null || systemctl is-enabled cron 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_ntp': {
        'cmd': "systemctl is-enabled ntpd 2>/dev/null || systemctl is-enabled chronyd 2>/dev/null || systemctl is-enabled systemd-timesyncd 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_firewalld': {
        'cmd': "systemctl is-enabled firewalld 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_iptables': {
        'cmd': "systemctl is-enabled iptables 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_nftables': {
        'cmd': "systemctl is-enabled nftables 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_syslog_ng': {
        'cmd': "systemctl is-enabled syslog-ng 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_postfix': {
        'cmd': "systemctl is-enabled postfix 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_httpd': {
        'cmd': "systemctl is-enabled httpd 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_nginx': {
        'cmd': "systemctl is-enabled nginx 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_dhcp': {
        'cmd': "systemctl is-enabled dhcpd 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_named': {
        'cmd': "systemctl is-enabled named 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_snmpd': {
        'cmd': "systemctl is-enabled snmpd 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_vsftpd': {
        'cmd': "systemctl is-enabled vsftpd 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_samba': {
        'cmd': "systemctl is-enabled smb 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_squid': {
        'cmd': "systemctl is-enabled squid 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_dovecot': {
        'cmd': "systemctl is-enabled dovecot 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_bind': {
        'cmd': "systemctl is-enabled named 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_nfs': {
        'cmd': "systemctl is-enabled nfs-server 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_rpcbind': {
        'cmd': "systemctl is-enabled rpcbind 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_avahi': {
        'cmd': "systemctl is-enabled avahi-daemon 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_cups': {
        'cmd': "systemctl is-enabled cups 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_telnet': {
        'cmd': "systemctl is-enabled telnet 2>/dev/null || systemctl is-enabled telnet.socket 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_rsh': {
        'cmd': "systemctl is-enabled rsh 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_tftp': {
        'cmd': "systemctl is-enabled tftp 2>/dev/null || systemctl is-enabled tftp.socket 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_ypserv': {
        'cmd': "systemctl is-enabled ypserv 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'ensure_autofs': {
        'cmd': "systemctl is-enabled autofs 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_watch': {
        'cmd': "auditctl -l 2>/dev/null | grep -E '{path}' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "{path}", 'match': "contains",
        'var_map': {'path': 'path'},
    },
    'grub2_bootloader_argument': {
        'cmd': "grep -E '^{arg}' /etc/default/grub 2>/dev/null | grep -E '{value}' || grep -E '\\b{arg}\\b' /boot/grub2/grub.cfg 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{arg}", 'match': "contains",
        'var_map': {'arg': 'arg_name', 'value': 'arg_value'},
    },
    'grub2_bootloader_argument_absent': {
        'cmd': "grep -E '^{arg}' /etc/default/grub 2>/dev/null | grep -E '{value}' && echo 'FOUND' || echo 'NOT_FOUND'",
        'expected': "NOT_FOUND", 'match': "contains",
        'var_map': {'arg': 'arg_name', 'value': 'arg_value'},
    },
    'pam_options': {
        'cmd': "grep -E '^{module}' /etc/pam.d/{file} 2>/dev/null | grep -E '{option}' || echo 'NOT_CONFIGURED'",
        'expected': "{option}", 'match': "contains",
        'var_map': {'module': 'module_name', 'file': 'pam_file', 'option': 'option_name'},
    },
    'lineinfile': {
        'cmd': "grep -E '{text}' {path} 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{text}", 'match': "contains",
        'var_map': {'text': 'text', 'path': 'path'},
    },
    'file_existence': {
        'cmd': "ls -la {path} 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{path}", 'match': "contains",
        'var_map': {'path': 'path'},
    },
    'mount_option_remote_filesystems': {
        'cmd': "mount | grep '{filesystem}' 2>/dev/null | grep '{option}' || echo 'NOT_FOUND'",
        'expected': "{option}", 'match': "contains",
        'var_map': {'filesystem': 'filesystem', 'option': 'mount_option'},
    },
    'mount_option_removable_partitions': {
        'cmd': "mount | grep -E 'sd[b-z]|hd[b-z]' 2>/dev/null | grep '{option}' || echo 'NOT_FOUND'",
        'expected': "{option}", 'match': "contains",
        'var_map': {'option': 'mount_option'},
    },
    'sudo_defaults_option': {
        'cmd': "grep -E '^Defaults\\s+{option}' /etc/sudoers /etc/sudoers.d/* 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{option}", 'match': "contains",
        'var_map': {'option': 'option_name'},
    },
    'rsyslog_logfiles_attributes_modify': {
        'cmd': "grep -E '{attribute}' /etc/rsyslog.conf /etc/rsyslog.d/*.conf 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{attribute}", 'match': "contains",
        'var_map': {'attribute': 'attribute'},
    },
    'dconf_ini_file': {
        'cmd': "grep -E '{key}' /etc/dconf/db/local.d/* 2>/dev/null || grep -E '{key}' /etc/dconf/db/*.d/* 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{key}", 'match': "contains",
        'var_map': {'key': 'key'},
    },
    'socket_disabled': {
        'cmd': "systemctl is-enabled {socket}.socket 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "disabled", 'match': "contains",
        'var_map': {'socket': 'servicename'},
    },
    'shell_lineinfile': {
        'cmd': "grep -E '{text}' {path} 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{text}", 'match': "contains",
        'var_map': {'text': 'text', 'path': 'path'},
    },
    'systemd_dropin_configuration': {
        'cmd': "ls /etc/systemd/system/{service}.service.d/ 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{service}", 'match': "contains",
        'var_map': {'service': 'servicename'},
    },
    'key_value_pair_in_file': {
        'cmd': "grep -E '^{key}\\s*=\\s*{value}' {path} 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{key}={value}", 'match': "contains",
        'var_map': {'key': 'key', 'value': 'value', 'path': 'path'},
    },
    'systemd_mount_enabled': {
        'cmd': "systemctl is-enabled {mount}.mount 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {'mount': 'mountpoint'},
    },
    'timer_enabled': {
        'cmd': "systemctl is-enabled {timer}.timer 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "enabled", 'match': "contains",
        'var_map': {'timer': 'timername'},
    },
    'sebool': {
        'cmd': "getsebool {bool} 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{value}", 'match': "contains",
        'var_map': {'bool': 'name', 'value': 'value'},
    },
    'kernel_build_config': {
        'cmd': "grep -E '^{config}=' /boot/config-$(uname -r) 2>/dev/null || grep -E '^{config}=' /proc/config.gz 2>/dev/null | zcat 2>/dev/null || echo 'NOT_FOUND'",
        'expected': "{config}=y", 'match': "contains",
        'var_map': {'config': 'config'},
    },
    'audit_rules_path_syscall': {
        'cmd': "auditctl -l 2>/dev/null | grep -E '{path}' | grep -E '{syscall}' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "{path}", 'match': "contains",
        'var_map': {'path': 'path', 'syscall': 'syscall'},
    },
    'audit_rules_dac_modification': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'chmod|chown|fchmod|fchown|fchmodat|fchownat|lchown|setxattr|lsetxattr|fsetxattr|removexattr' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "chmod", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_chsh': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/chsh' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/chsh", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_chfn': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/chfn' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/chfn", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_at': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/at' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/at", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_ping': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/ping' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/ping", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_traceroute': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/traceroute' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/traceroute", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_netstat': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/netstat' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/netstat", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_find': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/find' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/find", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_tcpdump': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/tcpdump' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/tcpdump", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_w': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/w' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/w", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_wall': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/wall' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/wall", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_write': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/write' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/write", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_dmesg': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/dmesg' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/dmesg", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_fdisk': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/fdisk' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/fdisk", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_halt': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/halt' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/halt", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_reboot': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/reboot' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/reboot", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_shutdown': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/shutdown' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/shutdown", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_iptables': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/iptables' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/iptables", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_ifconfig': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/ifconfig' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/ifconfig", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_route': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/route' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/route", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_kill': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/kill' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/kill", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_quota': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/quota' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/quota", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_rsync': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/rsync' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/rsync", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_strace': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/strace' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/strace", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_sysctl': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/sysctl' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/sysctl", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_crontab': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/crontab' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/crontab", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_lsmod': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/lsmod' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/lsmod", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_rmmod': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/rmmod' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/rmmod", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_nice': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/nice' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/nice", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_uname': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/uname' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/uname", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_zip': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/zip' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/zip", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_logrotate': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/logrotate' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/logrotate", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_last': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/last' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/last", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_ps': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/ps' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/ps", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_ss': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/ss' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/ss", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_yum': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/yum' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/yum", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_dnf': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/dnf' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/dnf", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_apt': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/apt' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/apt", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_screen': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/screen' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/screen", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_xauth': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/xauth' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/xauth", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_sudoedit': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/sudoedit' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/sudoedit", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_ssh_keysign': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/libexec/openssh/ssh-keysign' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/libexec/openssh/ssh-keysign", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_mount_nfs': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/mount.nfs' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/mount.nfs", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_umount_nfs': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/umount.nfs' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/umount.nfs", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_ping6': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/ping6' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/ping6", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_traceroute6': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/traceroute6' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/traceroute6", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_arping': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/arping' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/arping", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_clock': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/clock' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/clock", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_newgidmap': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/newgidmap' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/newgidmap", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_newuidmap': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/newuidmap' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/newuidmap", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_setfiles': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/setfiles' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/setfiles", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_semanage': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/semanage' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/semanage", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_unix_chkpwd': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/unix_chkpwd' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/unix_chkpwd", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_pam_console_apply': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/pam_console_apply' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/pam_console_apply", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_netreport': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/netreport' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/netreport", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_su_l10n': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/su' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/su", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_pt_chown': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/pt_chown' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/pt_chown", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_usermod': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/usermod' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/usermod", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_groupmod': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/sbin/groupmod' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/sbin/groupmod", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_x86_64': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/x86_64' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/x86_64", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_xmms': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/xmms' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/xmms", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_privileged_commands_xterm': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'path=/usr/bin/xterm' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "path=/usr/bin/xterm", 'match': "contains",
        'var_map': {},
    },
    'accounts_password_pam_unix': {
        'cmd': "grep -E '^password.*pam_unix\\.so' /etc/pam.d/password-auth /etc/pam.d/system-auth 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "pam_unix.so", 'match': "contains",
        'var_map': {},
    },
    'accounts_password_pam_pwhistory': {
        'cmd': "grep -E '^password.*pam_pwhistory\\.so' /etc/pam.d/password-auth /etc/pam.d/system-auth 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "pam_pwhistory.so", 'match': "contains",
        'var_map': {},
    },
    'accounts_password_pam_minlen': {
        'cmd': "grep -E '^minlen' /etc/security/pwquality.conf 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "minlen", 'match': "contains",
        'var_map': {},
    },
    'accounts_password_pam_dcredit': {
        'cmd': "grep -E '^dcredit' /etc/security/pwquality.conf 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "dcredit", 'match': "contains",
        'var_map': {},
    },
    'accounts_password_pam_ucredit': {
        'cmd': "grep -E '^ucredit' /etc/security/pwquality.conf 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "ucredit", 'match': "contains",
        'var_map': {},
    },
    'accounts_password_pam_lcredit': {
        'cmd': "grep -E '^lcredit' /etc/security/pwquality.conf 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "lcredit", 'match': "contains",
        'var_map': {},
    },
    'accounts_password_pam_ocredit': {
        'cmd': "grep -E '^ocredit' /etc/security/pwquality.conf 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "ocredit", 'match': "contains",
        'var_map': {},
    },
    'accounts_password_pam_minclass': {
        'cmd': "grep -E '^minclass' /etc/security/pwquality.conf 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "minclass", 'match': "contains",
        'var_map': {},
    },
    'accounts_password_pam_retry': {
        'cmd': "grep -E '^retry' /etc/security/pwquality.conf 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "retry", 'match': "contains",
        'var_map': {},
    },
    'accounts_password_pam_faillock': {
        'cmd': "grep -E 'pam_faillock\\.so' /etc/pam.d/password-auth /etc/pam.d/system-auth 2>/dev/null || echo 'NOT_CONFIGURED'",
        'expected': "pam_faillock.so", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_time_adjtimex': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'adjtimex' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "adjtimex", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_time_clock_settime': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'clock_settime' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "clock_settime", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_time_stime': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'stime' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "stime", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_time_watch_localtime': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'localtime' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "localtime", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_media_export': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'mount' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "mount", 'match': "contains",
        'var_map': {},
    },
    'audit_rules_sysadmin_actions': {
        'cmd': "auditctl -l 2>/dev/null | grep -E 'sudo' | head -10 || echo 'NO_AUDIT_RULES'",
        'expected': "sudo", 'match': "contains",
        'var_map': {},
    },
}


def extract_ocil_commands(ocil_text):
    """从OCIL文本中提取shell命令"""
    if not ocil_text:
        return [], "", "contains"
    
    text = clean_text(ocil_text)
    commands = []
    
    # 尝试从<pre>标签中提取命令
    pre_cmds = re.findall(r'<pre[^>]*>(.*?)</pre>', ocil_text, re.DOTALL)
    for pre in pre_cmds:
        pre = pre.strip()
        if pre and any(cmd in pre for cmd in ['grep', 'cat', 'ls', 'find', 'stat', 'sysctl', 'systemctl', 'rpm', 'dpkg', 'auditctl', 'ss', 'netstat', 'mount', 'ps', 'getenforce', 'sestatus', 'readlink', 'file', 'awk', 'sed', 'echo', 'printf', 'test', 'if', 'which', 'whereis', 'chkconfig', 'df', 'du', 'sha256sum', 'md5sum']):
            pre_clean = pre.replace('\n', ' ').replace('&gt;', '>').replace('&lt;', '<')
            commands.append(pre_clean)
    
    # 尝试匹配命令模式
    cmd_pattern = r'(?:^|[.;])\s*(sudo\s+)?(grep|cat|ls|find|stat|sysctl|systemctl|rpm|dpkg|auditctl|ss|netstat|getenforce|sestatus|mount|ps|df|du|chkconfig|readlink|file|sha256sum|md5sum|awk|sed|echo|printf|test|if|for|while|which|whereis|type|command|chkconfig)\s+[-/\w].*?(?=[.;]|$)'
    cmd_matches = re.findall(cmd_pattern, text, re.IGNORECASE)
    for sudo_prefix, cmd_name in cmd_matches:
        full_match = re.search(
            r'(sudo\s+)?' + re.escape(cmd_name) + r'\s+[-/\w].*?(?=[.;]|$)',
            text, re.IGNORECASE
        )
        if full_match:
            cmd = full_match.group(0).strip()
            if cmd not in commands:
                commands.append(cmd)
    
    # 尝试匹配反引号中的命令
    backtick_cmds = re.findall(r'`([^`]+)`', text)
    for cmd in backtick_cmds:
        cmd = cmd.strip()
        if any(c in cmd for c in ['grep', 'cat ', 'ls ', 'find ', 'stat ', 'sysctl', 'systemctl', 'rpm ', 'dpkg ', 'auditctl', 'ss ', 'netstat', 'mount', 'ps ', 'getenforce', 'sestatus', 'readlink', 'file ', 'awk', 'sed', 'echo', 'printf', 'test', 'if ', 'which', 'whereis']):
            if cmd not in commands:
                commands.append(cmd)
    
    if commands:
        return commands[:3], commands[0], "contains"
    
    return [], "", "contains"


def extract_bash_commands(bash_content):
    """从bash修复脚本中提取检查命令"""
    if not bash_content:
        return [], "", "contains"
    
    commands = []
    lines = bash_content.split('\n')
    
    for line in lines:
        line = line.strip()
        if line.startswith('#') or line.startswith('else') or line.startswith('fi') or line.startswith('done'):
            continue
        if any(cmd in line for cmd in ['grep ', 'find ', 'cat ', 'ls ', 'stat ', 'sysctl ', 'systemctl ', 'rpm ', 'dpkg ', 'test ', '[ ', 'auditctl ', 'ss ', 'mount ', 'ps ', 'getenforce ', 'sestatus ', 'readlink ', 'file ', 'awk ', 'sed ', 'echo ', 'printf ', 'which ', 'whereis ']):
            line_clean = line.strip(';').strip()
            if line_clean and line_clean not in commands:
                commands.append(line_clean)
    
    if commands:
        return commands[:3], commands[0], "contains"
    
    return [], "", "contains"


def get_risk_from_severity(severity):
    return SEVERITY_MAP.get(severity, 4)


def get_category_from_prodtype(prodtype):
    if not prodtype:
        return DEFAULT_CATEGORY
    prodtype_lower = prodtype.lower()
    if 'ubuntu' in prodtype_lower or 'debian' in prodtype_lower:
        return 1
    if 'rhel' in prodtype_lower or 'centos' in prodtype_lower or 'fedora' in prodtype_lower:
        return 1
    if 'sle' in prodtype_lower or 'suse' in prodtype_lower:
        return 1
    if 'win' in prodtype_lower:
        return 2
    return DEFAULT_CATEGORY


def get_os_type(prodtype):
    if not prodtype:
        return 1
    prodtype_lower = prodtype.lower()
    if 'win' in prodtype_lower:
        return 2
    if 'ubuntu' in prodtype_lower or 'debian' in prodtype_lower:
        return 1
    if 'rhel' in prodtype_lower or 'centos' in prodtype_lower or 'fedora' in prodtype_lower:
        return 1
    if 'sle' in prodtype_lower or 'suse' in prodtype_lower:
        return 1
    if 'openeuler' in prodtype_lower or 'kylin' in prodtype_lower:
        return 3
    if 'embedded' in prodtype_lower:
        return 4
    return 1


def get_category_from_rule_id(rule_id):
    rule_lower = rule_id.lower()
    if any(k in rule_lower for k in ['account', 'password', 'pam', 'login']):
        return 1
    if any(k in rule_lower for k in ['user', 'group', 'sudo', 'umask']):
        return 2
    if any(k in rule_lower for k in ['firewall', 'iptables', 'nftables', 'firewalld']):
        return 3
    if any(k in rule_lower for k in ['selinux', 'kernel', 'boot', 'grub', 'apparmor']):
        return 4
    if any(k in rule_lower for k in ['file', 'perm', 'mount', 'partition']):
        return 5
    if any(k in rule_lower for k in ['audit', 'log', 'rsyslog', 'journald']):
        return 6
    if any(k in rule_lower for k in ['network', 'ssh', 'sshd', 'ip', 'tcp', 'dns']):
        return 7
    if any(k in rule_lower for k in ['update', 'package', 'rpm', 'dnf', 'yum', 'apt']):
        return 8
    if any(k in rule_lower for k in ['ssh', 'sshd']):
        return 9
    return DEFAULT_CATEGORY


def get_category_from_oval(oval_xml):
    if not oval_xml:
        return None
    try:
        root = ET.fromstring(oval_xml)
        for obj in root.iter():
            tag = obj.tag
            if 'audit' in tag.lower():
                return 6
            if 'file' in tag.lower() and 'perm' in tag.lower():
                return 5
            if 'sysctl' in tag.lower() or 'kernel' in tag.lower():
                return 4
            if 'account' in tag.lower() or 'password' in tag.lower() or 'shadow' in tag.lower():
                return 1
            if 'network' in tag.lower() or 'listen' in tag.lower():
                return 7
            if 'package' in tag.lower() or 'rpm' in tag.lower():
                return 8
    except:
        pass
    return None


def main():
    if len(sys.argv) < 2:
        print("Usage: python convert_compliance_rules.py <content_dir> [output_file]")
        sys.exit(1)
    
    content_dir = sys.argv[1]
    output_file = sys.argv[2] if len(sys.argv) > 2 else 'data/baseline/compliance_rules.json'
    
    guide_dir = os.path.join(content_dir, 'linux_os', 'guide')
    if not os.path.isdir(guide_dir):
        print(f"Error: guide directory not found: {guide_dir}")
        sys.exit(1)
    
    rules = []
    total_processed = 0
    total_converted = 0
    total_skipped_oval = 0
    total_skipped_template = 0
    total_skipped_ocil = 0
    total_skipped_bash = 0
    total_skipped_unknown = 0
    total_yaml_errors = 0
    
    print(f"Scanning rules in: {guide_dir}")
    
    for root, dirs, files in os.walk(guide_dir):
        if 'rule.yml' not in files:
            continue
        
        total_processed += 1
        rule_path = os.path.join(root, 'rule.yml')
        rule_id = os.path.basename(root)
        
        if total_processed % 100 == 0:
            print(f"  Processed {total_processed} rules, converted: {total_converted}...")
        
        try:
            with open(rule_path, 'r', encoding='utf-8') as f:
                raw_content = f.read()
        except Exception as e:
            total_yaml_errors += 1
            continue
        
        # 预处理YAML，移除Jinja2语法
        preprocessed = preprocess_yaml(raw_content)
        
        try:
            data = yaml.safe_load(preprocessed)
        except Exception as e:
            total_yaml_errors += 1
            continue
        
        if not data or not isinstance(data, dict):
            total_yaml_errors += 1
            continue
        
        title = data.get('title', '')
        if not title:
            continue
        
        description = data.get('description', '') or title
        description = clean_text(description)
        
        severity = data.get('severity', 'unknown')
        risk = get_risk_from_severity(severity)
        
        prodtype = data.get('prodtype', '')
        os_type = get_os_type(prodtype)
        
        # 确定分类
        category = get_category_from_rule_id(rule_id)
        
        # 检查OVAL
        oval_dir = os.path.join(root, 'oval')
        has_oval = os.path.isdir(oval_dir) and any(f.endswith('.xml') for f in os.listdir(oval_dir))
        
        # 检查模板
        tpl = data.get('template', {})
        has_template = bool(tpl and isinstance(tpl, dict) and tpl.get('name'))
        
        # 检查OCIL
        ocil = data.get('ocil', '')
        has_ocil = bool(ocil)
        
        # 检查bash修复
        bash_dir = os.path.join(root, 'bash')
        has_bash = os.path.isdir(bash_dir) and any(f.endswith('.sh') for f in os.listdir(bash_dir))
        
        commands = []
        expected_value = ""
        match_type = "contains"
        
        # 策略1: 尝试OVAL
        if has_oval:
            for fname in os.listdir(oval_dir):
                if fname.endswith('.xml'):
                    oval_path = os.path.join(oval_dir, fname)
                    break
            try:
                with open(oval_path, 'r', encoding='utf-8') as f:
                    oval_content = f.read()
                oval_content = replace_jinja_vars(oval_content)
                cmds, exp, mt = oval_to_shell_commands(oval_content)
                if cmds:
                    commands = cmds
                    expected_value = exp or title
                    match_type = mt
                    # 尝试从OVAL确定分类
                    oval_cat = get_category_from_oval(oval_content)
                    if oval_cat is not None:
                        category = oval_cat
            except Exception as e:
                pass
        
        # 策略2: 尝试模板
        if not commands and has_template:
            tpl_name = tpl.get('name', '')
            tpl_vars = tpl.get('vars', {}) or {}
            if tpl_name in TEMPLATE_COMMANDS:
                tpl_entry = TEMPLATE_COMMANDS[tpl_name]
                cmd_template = tpl_entry['cmd']
                var_map = tpl_entry.get('var_map', {})
                
                # 替换模板变量
                cmd = cmd_template
                for var_key, var_source in var_map.items():
                    if var_source in tpl_vars:
                        cmd = cmd.replace('{' + var_key + '}', str(tpl_vars[var_source]))
                    else:
                        cmd = cmd.replace('{' + var_key + '}', '.*')
                
                expected_template = tpl_entry['expected']
                for var_key, var_source in var_map.items():
                    if var_source in tpl_vars:
                        expected_template = expected_template.replace('{' + var_key + '}', str(tpl_vars[var_source]))
                
                commands = [cmd]
                expected_value = expected_template
                match_type = tpl_entry.get('match', 'contains')
        
        # 策略3: 尝试OCIL
        if not commands and has_ocil:
            cmds, exp, mt = extract_ocil_commands(ocil)
            if cmds:
                commands = cmds
                expected_value = exp or title
                match_type = mt
        
        # 策略4: 尝试bash修复
        if not commands and has_bash:
            for fname in os.listdir(bash_dir):
                if fname.endswith('.sh'):
                    bash_path = os.path.join(bash_dir, fname)
                    break
            try:
                with open(bash_path, 'r', encoding='utf-8') as f:
                    bash_content = f.read()
                cmds, exp, mt = extract_bash_commands(bash_content)
                if cmds:
                    commands = cmds
                    expected_value = exp or title
                    match_type = mt
            except:
                pass
        
        if not commands:
            # 记录跳过原因
            if has_oval:
                total_skipped_oval += 1
            elif has_template:
                total_skipped_template += 1
            elif has_ocil:
                total_skipped_ocil += 1
            elif has_bash:
                total_skipped_bash += 1
            else:
                total_skipped_unknown += 1
            continue
        
        # 质量过滤：检查命令是否包含无效文本
        valid = True
        for cmd in commands:
            cmd_lower = cmd.lower()
            if any(w in cmd_lower for w in ['to the ', 'for the ', 'of the ', 'run the following ', 'if the ', 'in the ', 'on the ', 'with the ', 'is the ', 'by the ', 'at the ']):
                valid = False
                break
        
        if not valid:
            total_skipped_ocil += 1
            continue
        
        # 生成修复建议
        fix_suggestion = "请参考系统安全最佳实践进行修复"
        if category == 1:
            fix_suggestion = "请修改密码策略配置，确保符合安全要求"
        elif category == 2:
            fix_suggestion = "请检查用户权限配置，确保最小权限原则"
        elif category == 3:
            fix_suggestion = "请检查防火墙规则配置，确保网络访问控制"
        elif category == 4:
            fix_suggestion = "请检查内核安全配置，确保系统加固"
        elif category == 5:
            fix_suggestion = "请检查文件权限配置，确保访问控制"
        elif category == 6:
            fix_suggestion = "请检查审计日志配置，确保日志记录完整"
        elif category == 7:
            fix_suggestion = "请检查网络服务配置，确保服务安全"
        elif category == 8:
            fix_suggestion = "请检查系统更新配置，确保及时更新"
        elif category == 9:
            fix_suggestion = "请检查SSH配置，确保远程访问安全"
        
        # 生成风险描述
        risk_description = f"检测到系统存在安全配置问题: {title}"
        if risk == 0:
            risk_description = f"严重风险: {title}。请立即修复此安全问题。"
        elif risk == 1:
            risk_description = f"高危风险: {title}。建议尽快修复此安全问题。"
        elif risk == 2:
            risk_description = f"中危风险: {title}。建议在适当时候修复此安全问题。"
        elif risk == 3:
            risk_description = f"低危风险: {title}。建议关注此安全问题。"
        
        # 生成唯一ID
        rule_id_num = abs(hash(rule_id)) % 1000000
        
        rule = {
            "id": rule_id_num,
            "name": title[:200],
            "description": description[:500] if description else title[:200],
            "category": category,
            "risk": risk,
            "osType": os_type,
            "commands": commands[:3],
            "expectedValue": expected_value[:200] if expected_value else title[:200],
            "matchType": match_type,
            "fixSuggestion": fix_suggestion,
            "riskDescription": risk_description[:500],
        }
        
        rules.append(rule)
        total_converted += 1
    
    print(f"\n=== 转换完成 ===")
    print(f"总处理规则数: {total_processed}")
    print(f"成功转换: {total_converted}")
    print(f"YAML解析错误: {total_yaml_errors}")
    print(f"跳过-OVAL不支持: {total_skipped_oval}")
    print(f"跳过-未知模板: {total_skipped_template}")
    print(f"跳过-OCIL无命令: {total_skipped_ocil}")
    print(f"跳过-bash无检查命令: {total_skipped_bash}")
    print(f"跳过-未知原因: {total_skipped_unknown}")
    
    # 保存结果
    os.makedirs(os.path.dirname(output_file), exist_ok=True)
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(rules, f, ensure_ascii=False, indent=2)
    
    print(f"\n规则已保存到: {output_file}")
    print(f"总规则数: {len(rules)}")
    
    # 统计分布
    from collections import Counter
    cat_counter = Counter(r['category'] for r in rules)
    os_counter = Counter(r['osType'] for r in rules)
    risk_counter = Counter(r['risk'] for r in rules)
    
    print(f"\n分类分布:")
    for cat, count in sorted(cat_counter.items()):
        print(f"  分类{cat}: {count}")
    
    print(f"\n操作系统类型分布:")
    for os_t, count in sorted(os_counter.items()):
        print(f"  OS类型{os_t}: {count}")
    
    print(f"\n风险等级分布:")
    for risk, count in sorted(risk_counter.items()):
        print(f"  风险{risk}: {count}")


if __name__ == '__main__':
    main()