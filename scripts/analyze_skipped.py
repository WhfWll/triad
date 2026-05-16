import os, yaml, json, sys, xml.etree.ElementTree as ET

guide_dir = 'linux_os/guide'
total = 0
converted = 0
skipped_oval = 0
skipped_template = 0
skipped_ocil = 0
skipped_bash = 0
skipped_unknown = 0

# Known template names from the conversion script
KNOWN_TEMPLATES = [
    'kernel_module_disabled', 'package_removed', 'package_installed',
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
]

# Supported OVAL object types
OVAL_OBJECTS = {
    'textfilecontent54': '{http://oval.mitre.org/XMLSchema/oval-definitions-5#independent}textfilecontent54_object',
    'file': '{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}file_object',
    'symlink': '{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}symlink_object',
    'partition': '{http://oval.mitre.org/XMLSchema/oval-definitions-5#linux}partition_object',
    'sysctl': '{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}sysctl_object',
    'inetlisteningservers': '{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}inetlisteningservers_test',
    'process': '{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}process_test',
    'rpmverifyfile': '{http://oval.mitre.org/XMLSchema/oval-definitions-5#linux}rpmverifyfile_test',
    'systemdunitproperty': '{http://oval.mitre.org/XMLSchema/oval-definitions-5#linux}systemdunitproperty_object',
    'environmentvariable': '{http://oval.mitre.org/XMLSchema/oval-definitions-5#independent}environmentvariable_object',
    'shadow': '{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}shadow_object',
    'runlevel': '{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}runlevel_test',
    'password': '{http://oval.mitre.org/XMLSchema/oval-definitions-5#unix}password_test',
}

for root, dirs, files in os.walk(guide_dir):
    if 'rule.yml' not in files:
        continue
    total += 1
    if total > 2100:
        break
    
    try:
        with open(os.path.join(root, 'rule.yml'), encoding='utf-8') as f:
            data = yaml.safe_load(f)
    except:
        continue
    
    if not data or not isinstance(data, dict):
        continue
    if not data.get('title'):
        continue
    
    # Check what content types this rule has
    has_oval = os.path.isdir(os.path.join(root, 'oval')) and any(f.endswith('.xml') for f in os.listdir(os.path.join(root, 'oval')))
    tpl = data.get('template', {})
    has_template = bool(tpl and isinstance(tpl, dict) and tpl.get('name'))
    has_ocil = bool(data.get('ocil'))
    has_bash = os.path.isdir(os.path.join(root, 'bash')) and any(f.endswith('.sh') for f in os.listdir(os.path.join(root, 'bash')))
    
    if not (has_oval or has_template or has_ocil or has_bash):
        continue
    
    # Try OVAL
    if has_oval:
        oval_dir = os.path.join(root, 'oval')
        for fname in os.listdir(oval_dir):
            if fname.endswith('.xml'):
                oval_path = os.path.join(oval_dir, fname)
                break
        with open(oval_path, encoding='utf-8') as f:
            oval_content = f.read()
        try:
            root_xml = ET.fromstring(oval_content)
            found = False
            for obj_type, ns in OVAL_OBJECTS.items():
                if root_xml.findall('.//' + ns):
                    found = True
                    break
            if not found:
                skipped_oval += 1
                has_oval = False
        except:
            skipped_oval += 1
            has_oval = False
    
    if has_oval:
        converted += 1
        continue
    
    # Try template
    if has_template:
        tpl_name = tpl.get('name', '')
        if tpl_name in KNOWN_TEMPLATES:
            converted += 1
            continue
        else:
            skipped_template += 1
            has_template = False
    
    # Try OCIL
    if has_ocil:
        ocil = data.get('ocil', '')
        import re
        cmd_match = re.search(
            r'(?:^|[.;])\s*(sudo\s+)?(grep|cat|ls|find|stat|sysctl|systemctl|rpm|dpkg|auditctl|ss|netstat|getenforce|sestatus|mount|ps|df|du|chkconfig|readlink|file|sha256sum|md5sum|awk|sed)\s+[-/\w].*?(?=[.;]|$)',
            ocil
        )
        if cmd_match:
            converted += 1
            continue
        else:
            skipped_ocil += 1
            has_ocil = False
    
    # Try bash
    if has_bash:
        bash_dir = os.path.join(root, 'bash')
        for fname in os.listdir(bash_dir):
            if fname.endswith('.sh'):
                bash_path = os.path.join(bash_dir, fname)
                break
        with open(bash_path, encoding='utf-8') as f:
            bash_content = f.read()
        check_cmds = ['grep ', 'find ', 'cat ', 'ls ', 'stat ', 'sysctl ',
                      'systemctl ', 'rpm ', 'dpkg ', 'test ', '[ ',
                      'auditctl ', 'ss ', 'mount ', 'ps ', 'getenforce ']
        has_check = any(cmd in bash_content for cmd in check_cmds)
        if has_check:
            converted += 1
            continue
        else:
            skipped_bash += 1
            has_bash = False
    
    skipped_unknown += 1

print(f'Total rules with content: {total}')
print(f'Converted: {converted}')
print(f'Skipped - OVAL unsupported type: {skipped_oval}')
print(f'Skipped - unknown template: {skipped_template}')
print(f'Skipped - OCIL no cmd: {skipped_ocil}')
print(f'Skipped - bash no check cmd: {skipped_bash}')
print(f'Skipped - unknown reason: {skipped_unknown}')
print(f'Potential to add: {total - converted}')