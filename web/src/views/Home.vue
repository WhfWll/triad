<template>
  <div class="app-shell">
    <div class="sidebar" :class="{ collapsed: isCollapsed }">
      <div class="sidebar-header">
        <div class="logo-area">
          <svg viewBox="0 0 48 48" width="32" height="32" class="logo-icon">
            <path d="M24 4L8 12v10c0 12 7 22 16 24 9-2 16-12 16-24V12L24 4z" fill="none" stroke="#00d4aa" stroke-width="2.5"/>
            <path d="M18 24l4 4 8-8" fill="none" stroke="#00d4aa" stroke-width="2.5" stroke-linecap="round"/>
          </svg>
          <span class="logo-text" v-show="!isCollapsed">SecGuard</span>
        </div>
      </div>

      <div class="sidebar-menu">
        <el-scrollbar style="height: 100%">
          <el-menu
            :default-active="sidebarActivePath"
            :collapse="isCollapsed"
            :collapse-transition="false"
            router
            background-color="transparent"
            text-color="rgba(148,163,184,0.7)"
            active-text-color="#00d4aa"
          >
            <el-menu-item index="/index" v-if="role != 3">
              <i class="el-icon-data-analysis"></i>
              <span slot="title">仪表盘</span>
            </el-menu-item>

            <el-submenu index="/task" v-if="role != 3">
              <template slot="title">
                <i class="el-icon-s-order"></i>
                <span>任务中心</span>
              </template>
              <el-menu-item index="/task">渗透测试</el-menu-item>
              <el-menu-item index="/vulScanTask">漏洞扫描</el-menu-item>
              <el-menu-item index="/logicvuln">逻辑漏洞</el-menu-item>
              <el-menu-item index="/taskgroup">任务组</el-menu-item>
            </el-submenu>

            <el-submenu index="submenu-hostsec" v-if="role != 3">
              <template slot="title">
                <i class="el-icon-monitor"></i>
                <span>主机安全检查</span>
              </template>
              <el-menu-item index="/hostsec/tasks">任务管理</el-menu-item>
              <el-menu-item index="/hostsec/rules">配置核查规则</el-menu-item>
              <el-menu-item index="/hostsec/vuln-rules">漏洞扫描规则</el-menu-item>
              <el-menu-item index="/hostsec/malware-rules">病毒库规则</el-menu-item>
            </el-submenu>

            <el-submenu index="submenu-appsec" v-if="role != 3">
              <template slot="title">
                <i class="el-icon-mobile-phone"></i>
                <span>应用安全检查</span>
              </template>
              <el-menu-item index="/appsec/tasks">任务管理</el-menu-item>
              <el-menu-item index="/appsec/rules">检测规则</el-menu-item>
            </el-submenu>

            <el-submenu index="submenu-datasec" v-if="role != 3">
              <template slot="title">
                <i class="el-icon-document"></i>
                <span>数据安全检查</span>
              </template>
              <el-menu-item index="/datasec/tasks">任务管理</el-menu-item>
              <el-menu-item index="/datasec/rules">检测规则</el-menu-item>
            </el-submenu>

            <el-submenu index="/report">
              <template slot="title">
                <i class="el-icon-document"></i>
                <span>报告中心</span>
              </template>
              <el-menu-item index="/reportlist">报告清单</el-menu-item>
              <el-menu-item index="/createreport" v-if="false">生成报告</el-menu-item>
            </el-submenu>

            <el-submenu index="/experienceSet" v-if="role != 3">
              <template slot="title">
                <i class="el-icon-setting"></i>
                <span>场景管理</span>
              </template>
              <el-menu-item index="/taskscenario">任务场景</el-menu-item>
            </el-submenu>

            <el-submenu index="/toolmanagement" v-if="role != 3">
              <template slot="title">
                <i class="el-icon-tools"></i>
                <span>工具管理</span>
              </template>
              <el-menu-item index="/vulnerability">漏洞库</el-menu-item>
              <el-menu-item index="/fingerprint">指纹库</el-menu-item>
              <el-menu-item index="/dictionary">字典库</el-menu-item>
              <el-menu-item index="/auxiliarytool">辅助工具</el-menu-item>
            </el-submenu>

            <el-submenu index="/assetview" v-if="role != 3">
              <template slot="title">
                <i class="el-icon-monitor"></i>
                <span>资产中心</span>
              </template>
              <el-menu-item index="/assetview">资产概览</el-menu-item>
              <el-menu-item index="/assettree">资产树</el-menu-item>
            </el-submenu>

            <el-menu-item index="/riskvuln" v-if="role != 3">
              <i class="el-icon-warning"></i>
              <span slot="title">漏洞中心</span>
            </el-menu-item>



            <el-submenu index="/usermanagement" v-if="role == 2 || role == 4">
              <template slot="title">
                <i class="el-icon-user"></i>
                <span>用户管理</span>
              </template>
              <el-menu-item index="/usermanagement">用户管理</el-menu-item>
              <el-menu-item index="/usergroup">用户组管理</el-menu-item>
            </el-submenu>

            <el-submenu index="/system" v-if="role !== 1 && role !== 2 && role !== 3">
              <template slot="title">
                <i class="el-icon-s-tools"></i>
                <span>系统管理</span>
              </template>
              <el-menu-item index="/systemsetting">系统配置</el-menu-item>
              <el-menu-item index="/node">节点管理</el-menu-item>
            </el-submenu>

            <el-submenu index="/log" v-if="role !== 1 && role !== 2">
              <template slot="title">
                <i class="el-icon-document-copy"></i>
                <span>日志管理</span>
              </template>
              <el-menu-item index="/log">日志审计</el-menu-item>
              <el-menu-item index="/logconfig">日志配置</el-menu-item>
            </el-submenu>
          </el-menu>
        </el-scrollbar>
      </div>

      <div class="sidebar-footer">
        <div class="collapse-btn" @click="isCollapsed = !isCollapsed">
          <i :class="isCollapsed ? 'el-icon-s-unfold' : 'el-icon-s-fold'"></i>
        </div>
      </div>
    </div>

    <div class="main-area" :class="{ expanded: isCollapsed }">
      <header class="top-header">
        <div class="header-left">
          <div class="breadcrumb">
            <span class="page-title">{{ currentPageTitle }}</span>
          </div>
        </div>
        <div class="header-right">
          <div class="header-actions">
            <el-dropdown trigger="click" @command="handleCommand">
              <span class="user-info">
                <span class="avatar">
                  <i class="el-icon-user-solid"></i>
                </span>
                <span class="username">{{ sysUserName }}</span>
                <i class="el-icon-arrow-down"></i>
              </span>
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item command="updatepwd">
                  <i class="el-icon-key"></i> 修改密码
                </el-dropdown-item>
                <el-dropdown-item command="updateedit">
                  <i class="el-icon-edit"></i> 编辑用户
                </el-dropdown-item>
                <el-dropdown-item command="logout" divided>
                  <i class="el-icon-switch-button"></i> 退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
          </div>
        </div>
      </header>

      <main class="content-area">
        <router-view />
      </main>
    </div>

    <el-dialog title="修改密码" :visible.sync="dialogFormVisible" width="520px" :close-on-click-modal="false" class="theme-dialog">
      <el-form :model="pwddata" :rules="rules1" ref="ruleFormPWD" label-width="80px">
        <el-form-item label="原密码" prop="oldpwd">
          <el-input type="password" v-model="pwddata.oldpwd" placeholder="请输入原密码"></el-input>
        </el-form-item>
        <el-form-item label="新密码" prop="onepwd">
          <el-input type="password" v-model="pwddata.onepwd" placeholder="请输入新密码"></el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="twopwd">
          <el-input type="password" v-model="pwddata.twopwd" placeholder="请输入确认密码"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="cancelform">取 消</el-button>
        <el-button type="primary" @click="submitForm">确 定</el-button>
      </span>
    </el-dialog>

    <el-dialog title="提示" :visible.sync="logoutDialogVisible" width="400px" :close-on-click-modal="false" class="theme-dialog">
      <div class="logout-confirm">
        <i class="el-icon-warning-outline"></i>
        <span>确定退出系统吗？</span>
      </div>
      <span slot="footer">
        <el-button @click="logoutDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="fnlogout">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<style scoped lang="less">
@import './css/home.less';
</style>

<script>
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
    sidebarActivePath() {
      const p = this.$route.path
      if (p.startsWith('/hostsec')) return p
      if (p === '/hostbaseline' || p === '/malware' || p === '/hostsecurityhub') return '/hostsec/tasks'
      if (p.startsWith('/appsec')) return p
      if (p === '/dynamicscan' || p === '/appspecific') return '/appsec/tasks'
      if (p.startsWith('/datasec')) return p
      if (p === '/dbcheck' || p === '/sensitive') return '/datasec/tasks'
      return p
    },
    currentPageTitle() {
      const map = {
        '/index': '仪表盘',
        '/task': '渗透测试',
        '/vulScanTask': '漏洞扫描',
        '/logicvuln': '逻辑漏洞',
        '/taskgroup': '任务组',
        '/hostsec/tasks': '主机安全检查 · 任务管理',
        '/hostsec/rules': '主机安全检查 · 配置核查规则',
        '/hostsec/vuln-rules': '主机安全检查 · 漏洞扫描规则',
        '/hostsec/malware-rules': '主机安全检查 · 病毒库规则',
        '/appsec/tasks': '应用安全检查 · 任务管理',
        '/appsec/rules': '应用安全检查 · 检测规则',
        '/datasec/tasks': '数据安全检查 · 任务管理',
        '/datasec/rules': '数据安全检查 · 检测规则',
        '/hostsecurityhub': '主机安全检查 · 任务管理',
        '/baseline': '安全配置核查',
        '/hostbaseline': '安全配置核查',
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
      this.$refs.ruleFormPWD.validate(valid => {
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
        }).then(res => {
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
</script>