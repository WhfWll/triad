<template>
  <div class="login-page">
    <div class="login-bg">
      <div class="login-particles">
        <div v-for="i in 30" :key="i" class="particle" :style="getParticleStyle(i)"></div>
      </div>
      <div class="login-content">
        <div class="login-brand">
          <div class="brand-icon">
            <svg viewBox="0 0 48 48" width="48" height="48">
              <path d="M24 4L8 12v10c0 12 7 22 16 24 9-2 16-12 16-24V12L24 4z" fill="none" stroke="#00d4aa" stroke-width="2.5"/>
              <path d="M18 24l4 4 8-8" fill="none" stroke="#00d4aa" stroke-width="2.5" stroke-linecap="round"/>
            </svg>
          </div>
          <h1 class="brand-title">SecGuard</h1>
          <p class="brand-subtitle">智能安全检测与分析平台</p>
        </div>
        <div class="login-card">
          <div class="card-header">
            <h2>欢迎回来</h2>
            <p>请登录您的账户</p>
          </div>
          <el-alert class="error-msg" :title="errorbox.errmsg" type="error" :closable="false" v-if="errorbox.errshow" show-icon></el-alert>
          <el-form :model="ruleForm2" :rules="rules2" ref="ruleForm2" class="login-form">
            <el-form-item prop="username">
              <div class="input-group">
                <span class="input-icon">
                  <svg viewBox="0 0 24 24" width="18" height="18"><path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" fill="none" stroke="currentColor" stroke-width="2"/><circle cx="12" cy="7" r="4" fill="none" stroke="currentColor" stroke-width="2"/></svg>
                </span>
                <el-input type="text" v-model="ruleForm2.username" placeholder="用户名" class="custom-input"></el-input>
              </div>
            </el-form-item>
            <el-form-item prop="password">
              <div class="input-group">
                <span class="input-icon">
                  <svg viewBox="0 0 24 24" width="18" height="18"><rect x="3" y="11" width="18" height="11" rx="2" fill="none" stroke="currentColor" stroke-width="2"/><path d="M7 11V7a5 5 0 0110 0v4" fill="none" stroke="currentColor" stroke-width="2"/></svg>
                </span>
                <el-input type="password" v-model="ruleForm2.password" placeholder="密码" class="custom-input" @keyup.enter.native="handleSubmit"></el-input>
              </div>
            </el-form-item>
            <el-form-item prop="vcode">
              <div class="input-group">
                <span class="input-icon">
                  <svg viewBox="0 0 24 24" width="18" height="18"><rect x="3" y="3" width="18" height="18" rx="2" fill="none" stroke="currentColor" stroke-width="2"/><path d="M9 12l2 2 4-4" fill="none" stroke="currentColor" stroke-width="2"/></svg>
                </span>
                <el-input type="text" v-model="ruleForm2.vcode" placeholder="验证码" class="custom-input code-input" @keyup.enter.native="handleSubmit"></el-input>
                <img :src="check_code" class="code-img" @click="changevCode" alt="验证码" />
              </div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" class="login-btn" @click.prevent="handleSubmit" :loading="loading">
                <span v-if="!loading">登 录</span>
                <span v-else>登录中...</span>
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">
@import './css/login.less';
</style>

<script>
import router from '@/router';
import common from '@/utils/common.js'
import login from '@/api/login.js'
import { system } from '@/api/system.js'

export default {
  name: 'Login',
  data() {
    return {
      loading: false,
      code: '',
      baseUrl: '/api',
      ruleForm2: {
        username: '',
        password: '',
        vcode: '',
      },
      rules2: {
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
      },
      checked: true,
      errorbox: { errshow: false, errmsg: '' },
      check_code: '',
      timestamp: new Date().getTime(),
      captcha_id: '',
    };
  },
  created() {
    this.getCaptcha();
    if (process.env.VUE_APP_FREELOGIN == 1) {
      let obj = this.getUrlParams(window.location.href);
      if (obj && obj.sso) {
        let loginParams = { userName: obj.username, sso_token: obj.sso };
        this.$ajax({ url: "/smart/user/loginfreepassb", method: "POST", data: loginParams })
          .then(res => {
            var dt = res.data;
            if (dt.code === 200) {
              var data = dt.data;
              this.errorbox.errshow = false;
              let role = Number(data.role);
              localStorage.setItem('user_id', this.$commonjs.encryptCBC(data.uid, this.$commonjs.myKey));
              localStorage.setItem('user_id-par', data.uid);
              localStorage.setItem('token', data.token);
              localStorage.setItem('role', this.$commonjs.encryptCBC(role, this.$commonjs.myKey));
              localStorage.setItem('user', this.$commonjs.encryptCBC("admin", this.$commonjs.myKey));
              this.$router.push({ path: role == 3 ? '/log' : '/index' });
            }
          });
      }
    } else if (process.env.VUE_APP_FREELOGIN == 3) {
      let obj = this.getUrlParams(window.location.href);
      var next = obj.next.split("#");
      if (next.length < 2) return;
      next = next[1];
      if (obj && obj.authkey) {
        let loginParams = { token: obj.authkey };
        this.$ajax({ url: "/smart/user/loginfreepassd", method: "POST", data: loginParams })
          .then(res => {
            var dt = res.data;
            if (dt.code === 200) {
              var data = dt.data;
              this.errorbox.errshow = false;
              let role = Number(data.role);
              localStorage.setItem('user_id', this.$commonjs.encryptCBC(data.uid, this.$commonjs.myKey));
              localStorage.setItem('user_id-par', data.uid);
              localStorage.setItem('token', data.token);
              localStorage.setItem('role', this.$commonjs.encryptCBC(role, this.$commonjs.myKey));
              localStorage.setItem('user', this.$commonjs.encryptCBC(loginParams.userName || loginParams.token || 'user', this.$commonjs.myKey));
              this.$router.push({ path: '/index' });
            }
          });
      }
    }
  },
  methods: {
    getParticleStyle(i) {
      const size = Math.random() * 4 + 2;
      return {
        width: size + 'px',
        height: size + 'px',
        left: Math.random() * 100 + '%',
        top: Math.random() * 100 + '%',
        animationDelay: Math.random() * 5 + 's',
        animationDuration: (Math.random() * 3 + 2) + 's',
        opacity: Math.random() * 0.5 + 0.2,
      };
    },
    changevCode() {
      this.getCaptcha();
    },
    getCaptcha() {
      this.$ajax({
        url: '/smart/user/logincaptcha',
        method: 'GET',
        params: { time: new Date().getTime() },
      }).then(res => {
        if (res.data.code === 200) {
          this.check_code = res.data.data.data;
          this.captcha_id = res.data.data.captchaId;
        }
      });
    },
    handleSubmit() {
      this.$refs.ruleForm2.validate(valid => {
        if (!valid) return;
        this.loading = true;
        this.errorbox.errshow = false;
        var loginParams = {
          username: this.ruleForm2.username,
          password: this.$commonjs.encryptCBC(this.ruleForm2.password, this.$commonjs.myKey),
          checkCode: this.ruleForm2.vcode,
          captchaId: this.captcha_id,
        };
        this.$ajax({
          url: '/smart/user/login',
          method: 'POST',
          data: loginParams,
        })
          .then(res => {
            var dt = res.data;
            if (dt.code === 200) {
              var data = dt.data;
              let role = Number(data.role);
              localStorage.setItem('user_id', this.$commonjs.encryptCBC(data.uid, this.$commonjs.myKey));
              localStorage.setItem('user_id-par', data.uid);
              localStorage.setItem('token', data.token);
              localStorage.setItem('role', this.$commonjs.encryptCBC(role, this.$commonjs.myKey));
              localStorage.setItem('user', this.$commonjs.encryptCBC(this.ruleForm2.username, this.$commonjs.myKey));
              if (role == 3) {
                this.$router.push({ path: '/log' });
              } else {
                this.$router.push({ path: '/index' });
              }
            } else {
              this.errorbox.errshow = true;
              this.errorbox.errmsg = dt.msg || '登录失败';
              this.changevCode();
            }
            this.loading = false;
          })
          .catch(() => {
            this.errorbox.errshow = true;
            this.errorbox.errmsg = '网络错误';
            this.loading = false;
            this.changevCode();
          });
      });
    },
    getUrlParams(url) {
      var params = {};
      url.replace(/[?&]+([^=&]+)=([^&]*)/gi, function (m, key, value) { params[key] = value; });
      return params;
    },
  },
};
</script>
