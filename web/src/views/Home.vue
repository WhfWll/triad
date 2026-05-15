&lt;template&gt;
  &lt;div class="app-shell"&gt;
    &lt;div class="sidebar" :class="{ collapsed: isCollapsed }"&gt;
      &lt;div class="sidebar-header"&gt;
        &lt;div class="logo-area"&gt;
          &lt;svg viewBox="0 0 48 48" width="32" height="32" class="logo-icon"&gt;
            &lt;path d="M24 4L8 12v10c0 12 7 22 16 24 9-2 16-12 16-24V12L24 4z" fill="none" stroke="#00d4aa" stroke-width="2.5"/&gt;
            &lt;path d="M18 24l4 4 8-8" fill="none" stroke="#00d4aa" stroke-width="2.5" stroke-linecap="round"/&gt;
          &lt;/svg&gt;
          &lt;span class="logo-text" v-show="!isCollapsed"&gt;SecGuard&lt;/span&gt;
        &lt;/div&gt;
      &lt;/div&gt;

      &lt;div class="sidebar-menu"&gt;
        &lt;el-scrollbar style="height: 100%"&gt;
          &lt;el-menu
            :default-active="navselected"
            :collapse="isCollapsed"
            :collapse-transition="false"
            router
            background-color="transparent"
            text-color="rgba(148,163,184,0.7)"
            active-text-color="#00d4aa"
          &gt;
            &lt;el-menu-item index="/index" v-if="role != 3"&gt;
              &lt;i class="el-icon-data-analysis"&gt;&lt;/i&gt;
              &lt;span slot="title"&gt;仪表盘&lt;/span&gt;
            &lt;/el-menu-item&gt;

            &lt;el-submenu index="/task" v-if="role != 3"&gt;
              &lt;template slot="title"&gt;
                &lt;i class="el-icon-s-order"&gt;&lt;/i&gt;
                &lt;span&gt;任务中心&lt;/span&gt;
              &lt;/template&gt;
              &lt;el-menu-item index="/task"&gt;渗透测试&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/vulScanTask"&gt;漏洞扫描&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/logicvuln"&gt;逻辑漏洞&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/taskgroup"&gt;任务组&lt;/el-menu-item&gt;
            &lt;/el-submenu&gt;

            &lt;el-submenu index="/security" v-if="role != 3"&gt;
              &lt;template slot="title"&gt;
                &lt;i class="el-icon-shield"&gt;&lt;/i&gt;
                &lt;span&gt;安全检查&lt;/span&gt;
              &lt;/template&gt;
              &lt;el-menu-item index="/bastask"&gt;基线核查&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/malware"&gt;恶意代码检测&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/dbcheck"&gt;数据库安全&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/sensitive"&gt;敏感数据发现&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/appspecific"&gt;专项应用检测&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/dynamicscan"&gt;动态扫描&lt;/el-menu-item&gt;
            &lt;/el-submenu&gt;

            &lt;el-submenu index="/report"&gt;
              &lt;template slot="title"&gt;
                &lt;i class="el-icon-document"&gt;&lt;/i&gt;
                &lt;span&gt;报告中心&lt;/span&gt;
              &lt;/template&gt;
              &lt;el-menu-item index="/reportlist"&gt;报告清单&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/createreport" v-if="false"&gt;生成报告&lt;/el-menu-item&gt;
            &lt;/el-submenu&gt;

            &lt;el-submenu index="/experienceSet" v-if="role != 3"&gt;
              &lt;template slot="title"&gt;
                &lt;i class="el-icon-setting"&gt;&lt;/i&gt;
                &lt;span&gt;场景管理&lt;/span&gt;
              &lt;/template&gt;
              &lt;el-menu-item index="/taskscenario"&gt;任务场景&lt;/el-menu-item&gt;
            &lt;/el-submenu&gt;

            &lt;el-submenu index="/toolmanagement" v-if="role != 3"&gt;
              &lt;template slot="title"&gt;
                &lt;i class="el-icon-tools"&gt;&lt;/i&gt;
                &lt;span&gt;工具管理&lt;/span&gt;
              &lt;/template&gt;
              &lt;el-menu-item index="/vulnerability"&gt;漏洞库&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/fingerprint"&gt;指纹库&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/dictionary"&gt;字典库&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/auxiliarytool"&gt;辅助工具&lt;/el-menu-item&gt;
            &lt;/el-submenu&gt;

            &lt;el-submenu index="/assetview" v-if="role != 3"&gt;
              &lt;template slot="title"&gt;
                &lt;i class="el-icon-monitor"&gt;&lt;/i&gt;
                &lt;span&gt;资产中心&lt;/span&gt;
              &lt;/template&gt;
              &lt;el-menu-item index="/assetview"&gt;资产概览&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/assettree"&gt;资产树&lt;/el-menu-item&gt;
            &lt;/el-submenu&gt;

            &lt;el-menu-item index="/riskvuln" v-if="role != 3"&gt;
              &lt;i class="el-icon-warning"&gt;&lt;/i&gt;
              &lt;span slot="title"&gt;漏洞中心&lt;/span&gt;
            &lt;/el-menu-item&gt;



            &lt;el-submenu index="/usermanagement" v-if="role == 2 || role == 4"&gt;
              &lt;template slot="title"&gt;
                &lt;i class="el-icon-user"&gt;&lt;/i&gt;
                &lt;span&gt;用户管理&lt;/span&gt;
              &lt;/template&gt;
              &lt;el-menu-item index="/usermanagement"&gt;用户管理&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/usergroup"&gt;用户组管理&lt;/el-menu-item&gt;
            &lt;/el-submenu&gt;

            &lt;el-submenu index="/system" v-if="role !== 1 &amp;&amp; role !== 2 &amp;&amp; role !== 3"&gt;
              &lt;template slot="title"&gt;
                &lt;i class="el-icon-s-tools"&gt;&lt;/i&gt;
                &lt;span&gt;系统管理&lt;/span&gt;
              &lt;/template&gt;
              &lt;el-menu-item index="/systemsetting"&gt;系统配置&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/node"&gt;节点管理&lt;/el-menu-item&gt;
            &lt;/el-submenu&gt;

            &lt;el-submenu index="/log" v-if="role !== 1 &amp;&amp; role !== 2"&gt;
              &lt;template slot="title"&gt;
                &lt;i class="el-icon-document-copy"&gt;&lt;/i&gt;
                &lt;span&gt;日志管理&lt;/span&gt;
              &lt;/template&gt;
              &lt;el-menu-item index="/log"&gt;日志审计&lt;/el-menu-item&gt;
              &lt;el-menu-item index="/logconfig"&gt;日志配置&lt;/el-menu-item&gt;
            &lt;/el-submenu&gt;
          &lt;/el-menu&gt;
        &lt;/el-scrollbar&gt;
      &lt;/div&gt;

      &lt;div class="sidebar-footer"&gt;
        &lt;div class="collapse-btn" @click="isCollapsed = !isCollapsed"&gt;
          &lt;i :class="isCollapsed ? 'el-icon-s-unfold' : 'el-icon-s-fold'"&gt;&lt;/i&gt;
        &lt;/div&gt;
      &lt;/div&gt;
    &lt;/div&gt;

    &lt;div class="main-area" :class="{ expanded: isCollapsed }"&gt;
      &lt;header class="top-header"&gt;
        &lt;div class="header-left"&gt;
          &lt;div class="breadcrumb"&gt;
            &lt;span class="page-title"&gt;{{ currentPageTitle }}&lt;/span&gt;
          &lt;/div&gt;
        &lt;/div&gt;
        &lt;div class="header-right"&gt;
          &lt;div class="header-actions"&gt;
            &lt;el-dropdown trigger="click" @command="handleCommand"&gt;
              &lt;span class="user-info"&gt;
                &lt;span class="avatar"&gt;
                  &lt;i class="el-icon-user-solid"&gt;&lt;/i&gt;
                &lt;/span&gt;
                &lt;span class="username"&gt;{{ sysUserName }}&lt;/span&gt;
                &lt;i class="el-icon-arrow-down"&gt;&lt;/i&gt;
              &lt;/span&gt;
              &lt;el-dropdown-menu slot="dropdown"&gt;
                &lt;el-dropdown-item command="updatepwd"&gt;
                  &lt;i class="el-icon-key"&gt;&lt;/i&gt; 修改密码
                &lt;/el-dropdown-item&gt;
                &lt;el-dropdown-item command="updateedit"&gt;
                  &lt;i class="el-icon-edit"&gt;&lt;/i&gt; 编辑用户
                &lt;/el-dropdown-item&gt;
                &lt;el-dropdown-item command="logout" divided&gt;
                  &lt;i class="el-icon-switch-button"&gt;&lt;/i&gt; 退出登录
                &lt;/el-dropdown-item&gt;
              &lt;/el-dropdown-menu&gt;
            &lt;/el-dropdown&gt;
          &lt;/div&gt;
        &lt;/div&gt;
      &lt;/header&gt;

      &lt;main class="content-area"&gt;
        &lt;router-view /&gt;
      &lt;/main&gt;
    &lt;/div&gt;

    &lt;el-dialog title="修改密码" :visible.sync="dialogFormVisible" width="520px" :close-on-click-modal="false" class="theme-dialog"&gt;
      &lt;el-form :model="pwddata" :rules="rules1" ref="ruleFormPWD" label-width="80px"&gt;
        &lt;el-form-item label="原密码" prop="oldpwd"&gt;
          &lt;el-input type="password" v-model="pwddata.oldpwd" placeholder="请输入原密码"&gt;&lt;/el-input&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="新密码" prop="onepwd"&gt;
          &lt;el-input type="password" v-model="pwddata.onepwd" placeholder="请输入新密码"&gt;&lt;/el-input&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="确认密码" prop="twopwd"&gt;
          &lt;el-input type="password" v-model="pwddata.twopwd" placeholder="请输入确认密码"&gt;&lt;/el-input&gt;
        &lt;/el-form-item&gt;
      &lt;/el-form&gt;
      &lt;span slot="footer"&gt;
        &lt;el-button @click="cancelform"&gt;取 消&lt;/el-button&gt;
        &lt;el-button type="primary" @click="submitForm"&gt;确 定&lt;/el-button&gt;
      &lt;/span&gt;
    &lt;/el-dialog&gt;

    &lt;el-dialog title="提示" :visible.sync="logoutDialogVisible" width="400px" :close-on-click-modal="false" class="theme-dialog"&gt;
      &lt;div class="logout-confirm"&gt;
        &lt;i class="el-icon-warning-outline"&gt;&lt;/i&gt;
        &lt;span&gt;确定退出系统吗？&lt;/span&gt;
      &lt;/div&gt;
      &lt;span slot="footer"&gt;
        &lt;el-button @click="logoutDialogVisible = false"&gt;取 消&lt;/el-button&gt;
        &lt;el-button type="primary" @click="fnlogout"&gt;确 定&lt;/el-button&gt;
      &lt;/span&gt;
    &lt;/el-dialog&gt;
  &lt;/div&gt;
&lt;/template&gt;

&lt;style scoped lang="less"&gt;
@import './css/home.less';
&lt;/style&gt;

&lt;script&gt;
import store from '@/store'
import { system } from '@/api/system.js'

export default {
  name: 'Home',
  data() {
    return {
      isCollapsed: false,
      navselected: '/index',
      sysUserName: '',
      dialogFormVisible: false,
      logoutDialogVisible: false,
      dialogupdateFormVisible: false,
      pwddata: { oldpwd: '', onepwd: '', twopwd: '' },
      userform1: {},
      roleOptions: [],
      rules1: {},
      rules2: {},
      isdot: false,
      messagelist: [],
      role: null,
    };
  },
  computed: {
    currentPageTitle() {
      const map = {
        '/index': '仪表盘',
        '/task': '渗透测试',
        '/vulScanTask': '漏洞扫描',
        '/logicvuln': '逻辑漏洞',
        '/taskgroup': '任务组',
        '/baseline': '基线核查',
        '/malware': '恶意代码检测',
        '/dbcheck': '数据库安全',
        '/sensitive': '敏感数据发现',
        '/appspecific': '专项应用检测',
        '/dynamicscan': '动态扫描',
        '/reportlist': '报告清单',
        '/taskscenario': '任务场景',
        '/vulnerability': '漏洞库',
        '/fingerprint': '指纹库',
        '/dictionary': '字典库',
        '/auxiliarytool': '辅助工具',
        '/assetview': '资产概览',
        '/assettree': '资产树',
        '/riskvuln': '漏洞中心',
        '/bastask': '基线核查',
        '/evaluationScheme': '评估方案',
        '/scriptLibrary': '剧本库',
        '/agent': 'Agent管理',
        '/x-ray': 'Xray',
        '/burpsuite': 'Burpsuite',
        '/usermanagement': '用户管理',
        '/usergroup': '用户组管理',
        '/systemsetting': '系统配置',
        '/node': '节点管理',
        '/log': '日志审计',
        '/logconfig': '日志配置',
      };
      return map[this.$route.path] || '仪表盘';
    },
  },
  created() {
    this.getUserInfo();
    this.role = Number(localStorage.getItem('role') ? this.$commonjs.decryptCBC(localStorage.getItem('role'), this.$commonjs.myKey) : 0);
    this.navselected = this.$route.path;
  },
  watch: {
    '$route.path'(val) {
      this.navselected = val;
    },
  },
  methods: {
    getUserInfo() {
      this.sysUserName = localStorage.getItem('user') ? this.$commonjs.decryptCBC(localStorage.getItem('user'), this.$commonjs.myKey) : '';
    },
    handleCommand(command) {
      if (command === 'logout') {
        this.logoutDialogVisible = true;
      } else if (command === 'updatepwd') {
        this.dialogFormVisible = true;
      } else if (command === 'updateedit') {
        this.dialogupdateFormVisible = true;
      }
    },
    fnlogout() {
      localStorage.clear();
      this.$router.push({ path: '/login' });
    },
    submitForm() {
      this.$refs.ruleFormPWD.validate(valid =&gt; {
        if (!valid) return;
        var uid = localStorage.getItem('user_id-par');
        var putParams = {
          user_id: uid,
          old_password: this.pwddata.oldpwd,
          new_password: this.pwddata.onepwd,
        };
        this.$ajax({
          url: '/smart/user/updatepw',
          method: 'POST',
          data: putParams,
        }).then(res =&gt; {
          var dt = res.data;
          if (dt.code === 200) {
            this.$message({ message: '修改密码成功，请重新登录', type: 'success' });
            this.dialogFormVisible = false;
            localStorage.clear();
            this.$router.push({ path: '/login' });
          }
        });
      });
    },
    cancelform() {
      this.dialogFormVisible = false;
      this.pwddata = { oldpwd: '', onepwd: '', twopwd: '' };
    },
    updatecancelform() {
      this.dialogupdateFormVisible = false;
    },
  },
};
&lt;/script&gt;