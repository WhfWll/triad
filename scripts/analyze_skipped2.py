import os, yaml, json, sys, xml.etree.ElementTree as ET
from collections import Counter

guide_dir = 'linux_os/guide'

# Count unsupported OVAL types
oval_types_counter = Counter()
template_counter = Counter()
ocil_no_cmd_examples = []
bash_no_cmd_examples = []

for root, dirs, files in os.walk(guide_dir):
    if 'rule.yml' not in files:
        continue
    
    try:
        with open(os.path.join(root, 'rule.yml'), encoding='utf-8') as f:
            data = yaml.safe_load(f)
    except:
        continue
    
    if not data or not isinstance(data, dict):
        continue
    if not data.get('title'):
        continue
    
    title = data.get('title', '')
    
    # Check OVAL
    oval_dir = os.path.join(root, 'oval')
    has_oval = os.path.isdir(oval_dir) and any(f.endswith('.xml') for f in os.listdir(oval_dir))
    
    if has_oval:
        for fname in os.listdir(oval_dir):
            if fname.endswith('.xml'):
                oval_path = os.path.join(oval_dir, fname)
                break
        with open(oval_path, encoding='utf-8') as f:
            oval_content = f.read()
        try:
            root_xml = ET.fromstring(oval_content)
            # Find all test types
            for el in root_xml.iter():
                tag = el.tag
                if '_test' in tag and 'object' not in tag and 'state' not in tag:
                    # Extract the type name
                    if '}' in tag:
                        tag = tag.split('}')[1]
                    oval_types_counter[tag] += 1
        except:
            pass
    
    # Check template
    tpl = data.get('template', {})
    if tpl and isinstance(tpl, dict) and tpl.get('name'):
        tpl_name = tpl['name']
        known = ['kernel_module_disabled', 'package_removed', 'package_installed',
                 'service_enabled', 'service_disabled', 'service_running',
                 'sshd_lineinfile', 'sysctl', 'mount_option', 'file_owner',
                 'file_permissions', 'file_groupowner',
                 'audit_rules_dac_modification', 'audit_rules_privileged_commands',
                 'audit_rules_unsuccessful_file_modification', 'audit_rules_file_deletion_events',
                 'audit_rules_time_adjtimex', 'audit_rules_time_clock_settime',
                 'audit_rules_time_stime', 'audit_rules_time_watch_localtime',
                 'audit_rules_media_export', 'audit_rules_kernel_module_loading',
                 'audit_rules_login_events', 'audit_rules_networkconfig_modification',
                 'audit_rules_session_events', 'audit_rules_sysadmin_actions',
                 'audit_rules_usergroup_modification',
                 'accounts_password_pam_unix', 'accounts_password_pam_pwhistory',
                 'accounts_password_pam_minlen', 'accounts_password_pam_dcredit',
                 'accounts_password_pam_ucredit', 'accounts_password_pam_lcredit',
                 'accounts_password_pam_ocredit', 'accounts_password_pam_minclass',
                 'accounts_password_pam_retry', 'accounts_password_pam_faillock',
                 'audit_rules_privileged_commands_chmod', 'audit_rules_privileged_commands_rm',
                 'audit_rules_privileged_commands_su', 'audit_rules_privileged_commands_sudo',
                 'audit_rules_privileged_commands_mount', 'audit_rules_privileged_commands_umount',
                 'audit_rules_privileged_commands_ssh', 'audit_rules_privileged_commands_passwd',
                 'audit_rules_privileged_commands_crontab', 'audit_rules_privileged_commands_chsh',
                 'audit_rules_privileged_commands_chfn', 'audit_rules_privileged_commands_usermod',
                 'audit_rules_privileged_commands_groupmod', 'audit_rules_privileged_commands_newgrp',
                 'audit_rules_privileged_commands_pam_timestamp_check',
                 'audit_rules_privileged_commands_postdrop', 'audit_rules_privileged_commands_postqueue',
                 'audit_rules_privileged_commands_ssh_agent', 'audit_rules_privileged_commands_pt_chown',
                 'audit_rules_privileged_commands_newgidmap', 'audit_rules_privileged_commands_newuidmap',
                 'audit_rules_privileged_commands_setfiles', 'audit_rules_privileged_commands_semanage',
                 'audit_rules_privileged_commands_userhelper', 'audit_rules_privileged_commands_unix_chkpwd',
                 'audit_rules_privileged_commands_gpasswd', 'audit_rules_privileged_commands_at',
                 'audit_rules_privileged_commands_pam_console_apply',
                 'audit_rules_privileged_commands_pam_timestamp_check',
                 'audit_rules_privileged_commands_netreport',
                 'audit_rules_privileged_commands_pam_timestamp_check',
                 'audit_rules_privileged_commands_ssh_keysign',
                 'audit_rules_privileged_commands_su_l10n',
                 'audit_rules_privileged_commands_sudoedit',
                 'audit_rules_privileged_commands_xauth',
                 'audit_rules_privileged_commands_screen',
                 'audit_rules_privileged_commands_mount_nfs',
                 'audit_rules_privileged_commands_umount_nfs',
                 'audit_rules_privileged_commands_ping',
                 'audit_rules_privileged_commands_ping6',
                 'audit_rules_privileged_commands_traceroute',
                 'audit_rules_privileged_commands_traceroute6',
                 'audit_rules_privileged_commands_netstat',
                 'audit_rules_privileged_commands_arping',
                 'audit_rules_privileged_commands_clock',
                 'audit_rules_privileged_commands_crontab',
                 'audit_rules_privileged_commands_dmesg',
                 'audit_rules_privileged_commands_fdisk',
                 'audit_rules_privileged_commands_find',
                 'audit_rules_privileged_commands_halt',
                 'audit_rules_privileged_commands_ifconfig',
                 'audit_rules_privileged_commands_iptables',
                 'audit_rules_privileged_commands_kill',
                 'audit_rules_privileged_commands_kmod',
                 'audit_rules_privileged_commands_last',
                 'audit_rules_privileged_commands_logrotate',
                 'audit_rules_privileged_commands_lsmod',
                 'audit_rules_privileged_commands_mount',
                 'audit_rules_privileged_commands_netstat',
                 'audit_rules_privileged_commands_nice',
                 'audit_rules_privileged_commands_ping',
                 'audit_rules_privileged_commands_ping6',
                 'audit_rules_privileged_commands_ps',
                 'audit_rules_privileged_commands_quota',
                 'audit_rules_privileged_commands_reboot',
                 'audit_rules_privileged_commands_rmmod',
                 'audit_rules_privileged_commands_route',
                 'audit_rules_privileged_commands_rsync',
                 'audit_rules_privileged_commands_shutdown',
                 'audit_rules_privileged_commands_ss',
                 'audit_rules_privileged_commands_ssh',
                 'audit_rules_privileged_commands_strace',
                 'audit_rules_privileged_commands_su',
                 'audit_rules_privileged_commands_sudo',
                 'audit_rules_privileged_commands_sysctl',
                 'audit_rules_privileged_commands_tcpdump',
                 'audit_rules_privileged_commands_traceroute',
                 'audit_rules_privileged_commands_umount',
                 'audit_rules_privileged_commands_uname',
                 'audit_rules_privileged_commands_userhelper',
                 'audit_rules_privileged_commands_usermod',
                 'audit_rules_privileged_commands_w',
                 'audit_rules_privileged_commands_wall',
                 'audit_rules_privileged_commands_write',
                 'audit_rules_privileged_commands_x86_64',
                 'audit_rules_privileged_commands_xauth',
                 'audit_rules_privileged_commands_xmms',
                 'audit_rules_privileged_commands_xterm',
                 'audit_rules_privileged_commands_yum',
                 'audit_rules_privileged_commands_zip',
        ]
        if tpl_name not in known:
            template_counter[tpl_name] += 1

print("=== 最常见的未支持 OVAL 类型 ===")
for t, c in oval_types_counter.most_common(20):
    print(f"  {t}: {c}")

print("\n=== 最常见的未支持模板 ===")
for t, c in template_counter.most_common(30):
    print(f"  {t}: {c}")