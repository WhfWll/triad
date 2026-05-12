//系统设置---系统工具---安全设置 页面代码
<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="24">
        <!-- <div class="totalbox"> -->
        <div class="totalbox">
          <el-form
            ref="securityform"
            :model="securityform"
            label-width="250px"
            class="securityform"
            :rules="safety"
          >
            <!-- <div class="clearfix"> -->
            <el-form-item label="密码强度：" prop="passWordlevel">
              <el-select
                v-model="securityform.passWordlevel"
                clearable
                placeholder="密码强度"
                size="small"
              >
                <el-option
                  v-for="(item, i) in passWordlevel"
                  :key="i"
                  :label="item.label"
                  :value="item.value"
                >
                </el-option>
              </el-select>
              <el-tooltip placement="bottom">
                <div slot="content">
                  低：长度不低于8位，可为全字母或数字
                  <br />中：长度不低于8位，包含字母、符号、数字中的两类<br />高：长度不低于8位，包含字母、符号、数字中的三类
                </div>
                <i
                  class="iconfont icontishi icontsstyle"
                  style="vertical-align: middle"
                ></i>
              </el-tooltip>
            </el-form-item>
            <el-form-item label="密码有效期：" prop="user_limit">
              <el-select
                v-model="securityform.passWorddate"
                clearable
                placeholder="密码有效期"
                size="small"
              >
                <el-option
                  v-for="(item, i) in passWorddate"
                  :key="i"
                  :label="item.label"
                  :value="item.value"
                >
                </el-option> </el-select
              ><span class="month">月</span>
              <el-tooltip placement="bottom">
                <div slot="content">
                  密码有效期密码是指账号密码在有效期内未<br />进行密码修改，账号将被禁用
                </div>
                <i
                  class="iconfont icontishi icontsstyle"
                  style="vertical-align: middle"
                ></i>
              </el-tooltip>
            </el-form-item>
            <!-- </div> -->

            <el-form-item
              label="用户长时间未操作自动退出系统："
              prop="user_limit"
            >
              <el-select
                v-model="securityform.leavesystem"
                clearable
                placeholder=""
                size="small"
              >
                <el-option
                  v-for="(item, i) in leavesystem"
                  :key="i"
                  :label="item.label"
                  :value="item.value"
                >
                </el-option> </el-select
              ><span class="month">分钟</span>
              <el-tooltip placement="bottom">
                <div slot="content">
                  用户在设定的时间内未对系统执行任何操作，<br />登录账号将退出
                </div>
                <i
                  class="iconfont icontishi icontsstyle"
                  style="vertical-align: middle"
                ></i>
              </el-tooltip>
            </el-form-item>

            <el-form-item
              label="账号长时间未登录系统禁用账号："
              prop="login_timeout"
            >
              <el-select
                v-model="securityform.login_timeout"
                clearable
                placeholder=""
                size="small"
              >
                <el-option
                  v-for="(item, i) in login_timeout"
                  :key="i"
                  :label="item.label"
                  :value="item.value"
                >
                </el-option> </el-select
              ><span class="month">月</span>
              <el-tooltip placement="bottom">
                <div slot="content">
                  账号在设定的时间内没有登录过系统，<br />账号将被禁用
                </div>
                <i
                  class="iconfont icontishi icontsstyle"
                  style="vertical-align: middle"
                ></i>
              </el-tooltip>
            </el-form-item>
            <div>
              <el-form-item
                label="普通账号连续登录出错禁用账号："
                prop="user_limit"
                style="display: inline-block; margin-right: 50px"
              >
                <el-select
                  v-model="securityform.user_limit"
                  clearable
                  placeholder=""
                  size="small"
                >
                  <el-option
                    v-for="(item, i) in user_limit"
                    :key="i"
                    :label="item.label"
                    :value="item.value"
                  >
                  </el-option> </el-select
                ><span class="month">次</span>
              </el-form-item>
              <!-- <el-form-item label="禁用账号：" prop="outuser" label-width="110px" style="display:inline-block;">
                                     <el-input v-model="securityform.outuser">
                                     </el-input>
                                </el-form-item>   -->
            </div>

            <div>
              <el-form-item
                label="管理/审计员连续登录出错限制："
                prop="user_limit"
                style="display: inline-block; margin-right: 50px"
              >
                <el-select
                  v-model="securityform.admin_limit"
                  clearable
                  placeholder=""
                  size="small"
                >
                  <el-option
                    v-for="(item, i) in admin_limit"
                    :key="i"
                    :label="item.label"
                    :value="item.value"
                  >
                  </el-option> </el-select
                ><span class="month">次</span>
              </el-form-item>
              <el-form-item
                label="禁止登录时长："
                label-width="110px"
                style="display: inline-block"
              >
                <el-select
                  v-model="securityform.ban_time"
                  clearable
                  placeholder=""
                  size="small"
                >
                  <el-option
                    v-for="(item, i) in ban_time"
                    :key="i"
                    :label="item.label"
                    :value="item.value"
                  >
                  </el-option> </el-select
                ><span class="month">小时</span>
              </el-form-item>
            </div>
          </el-form>
          <xzbutton type="primary" @click="benSavesecurity" size="small"
            >保存</xzbutton
          >
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script> 
import xzbutton from "../../components/XzButton.vue"
import { system } from '@/api/system.js'
export default ({
  name: 'safetyset',
  components: {
    xzbutton,
  },
  data () {
    //  安全设置相关
    // let cycle = (rule, value, callback)=>{
    //     if(value < 0 || value > 36) callback(new Error('只能输入0-36范围之内'));
    //     else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
    //     else callback()
    // }
    // let timeout = (rule, value, callback)=>{
    //     if(value < 0 || value > 120) callback(new Error('只能输入0-120范围之内'));
    //     else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
    //     else callback()
    // }
    // let user = (rule, value, callback)=>{
    //     if(value < 0 || value > 10) callback(new Error('只能输入0-10范围之内'));
    //     else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
    //     else callback()
    // }
    // let admin = (rule, value, callback)=>{
    //     if(value < 0 || value > 10) callback(new Error('只能输入0-10范围之内'));
    //     else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
    //     else callback()
    // } 
    return {
      safeid: '',
      securityform: {
        passWordlevel: '',
        passWorddate: '',
        leavesystem: '',
        login_timeout: '',
        user_limit: '',
        admin_limit: '',
        ban_time: '',
      },
      passWordlevel: [
        {
          value: '1',
          label: '低'
        }, {
          value: '2',
          label: '中'
        },
        {
          value: '3',
          label: '高'
        },
      ],
      passWorddate: [
        {
          value: '1',
          label: '1'
        }, {
          value: '2',
          label: '2'
        },
        {
          value: '3',
          label: '3'
        },
        {
          value: '4',
          label: '4'
        }, {
          value: '5',
          label: '5'
        },
        {
          value: '6',
          label: '6'
        },
        {
          value: '7',
          label: '7'
        }, {
          value: '8',
          label: '8'
        },
        {
          value: '9',
          label: '9'
        },
        {
          value: '10',
          label: '10'
        }, {
          value: '11',
          label: '11'
        },
        {
          value: '12',
          label: '12'
        },
        {
          value: '0',
          label: '不限'
        },
      ],
      login_timeout: [
        {
          value: '1',
          label: '1'
        }, {
          value: '2',
          label: '2'
        },
        {
          value: '3',
          label: '3'
        },
        {
          value: '4',
          label: '4'
        }, {
          value: '5',
          label: '5'
        },
        {
          value: '6',
          label: '6'
        },
        {
          value: '7',
          label: '7'
        }, {
          value: '8',
          label: '8'
        },
        {
          value: '9',
          label: '9'
        },
        {
          value: '10',
          label: '10'
        }, {
          value: '11',
          label: '11'
        },
        {
          value: '12',
          label: '12'
        },
        {
          value: '18',
          label: '18'
        },
        {
          value: '24',
          label: '24'
        },
        {
          value: '0',
          label: '不限'
        },
      ],
      leavesystem: [
        {
          value: '10',
          label: '10'
        }, {
          value: '20',
          label: '20'
        },
        {
          value: '30',
          label: '30'
        },
        {
          value: '40',
          label: '40'
        }, {
          value: '50',
          label: '50'
        },
        {
          value: '60',
          label: '60'
        },
        {
          value: '80',
          label: '80'
        }, {
          value: '100',
          label: '100'
        },
        {
          value: '120',
          label: '120'
        },
        {
          value: '0',
          label: '不限'
        },
      ],
      admin_limit: [
        {
          value: '1',
          label: '1'
        }, {
          value: '2',
          label: '2'
        },
        {
          value: '3',
          label: '3'
        },
        {
          value: '4',
          label: '4'
        }, {
          value: '5',
          label: '5'
        },
        {
          value: '6',
          label: '6'
        },
        {
          value: '0',
          label: '不限'
        },
      ],
      user_limit: [
        {
          value: '1',
          label: '1'
        }, {
          value: '2',
          label: '2'
        },
        {
          value: '3',
          label: '3'
        },
        {
          value: '4',
          label: '4'
        }, {
          value: '5',
          label: '5'
        },
        {
          value: '6',
          label: '6'
        },
        {
          value: '不限',
          label: '不限'
        },
      ],
      ban_time: [
        {
          value: '0.5',
          label: '0.5'
        },
        {
          value: '1',
          label: '1'
        }, {
          value: '2',
          label: '2'
        },
        {
          value: '3',
          label: '3'
        },
        {
          value: '4',
          label: '4'
        }, {
          value: '5',
          label: '5'
        },
        {
          value: '6',
          label: '6'
        },
        {
          value: '12',
          label: '12'
        },
        {
          value: '24',
          label: '24'
        },
        {
          value: '0',
          label: '不禁止'
        },
      ],

      safety: {
        // 密码修改周期
        passWordlevel: [
          { required: true, message: '', trigger: [] },
        ],
        // 系统登录超时
        login_timeout: [
          { required: true, message: '', trigger: [] },
        ],
        user_limit: [
          { required: true, message: '', trigger: [] },
        ],
        admin_limit: [
          { required: true, message: '', trigger: [] },
        ],
      },
    }
  },
  created: function () {
    this.$store.state.activefirstMenu = "/systemsetting"
  },
  mounted: function () {
    this.getsafetyData()
  },
  beforeDestroy () {
  },
  methods: {
    async getsafetyData () {//进来先获取配置
      const dt = await system.getSafeinfo()
      if (dt.code === 200) {
        this.securityform.passWordlevel = dt.data.password_rank.toString()
        this.securityform.passWorddate = dt.data.password_cycle.toString()
        this.securityform.login_timeout = dt.data.account_no_login_disable.toString()
        this.securityform.leavesystem = dt.data.login_timeout.toString()
        this.securityform.admin_limit = dt.data.admin_limit.toString()
        this.securityform.user_limit = dt.data.user_limit.toString()
        this.securityform.ban_time = dt.data.ban_time.toString()
        this.safeid = dt.data.id
      } else {
        this.$message({
          message: dt.msg,
          type: 'error'
        })
      }
    },
    benSavesecurity () { //保存安全设置
      this.$refs.securityform.validate(async valid => {
        if (!valid) return

        let params = {
          id: this.safeid,
          password_rank: this.securityform.passWordlevel,
          password_cycle: this.securityform.passWorddate,
          login_timeout: this.securityform.leavesystem,
          account_no_login_disable: this.securityform.login_timeout,
          user_limit: this.securityform.user_limit,
          admin_limit: this.securityform.admin_limit,
          BanTime: this.securityform.ban_time,
        }
        const dt = await system.saveSecurity(params)
        if (dt.code === 200) {
          this.$message.success(dt.msg)
        } else {
          this.$message({
            message: dt.msg,
            type: 'error'
          })
        }
      })
    },

  }
})

</script>

<style scoped lang="less">
.selectbox {
  width: 90px;
  margin: 0 10px 0 40px;
}
.month {
  display: inline-block;
  color: rgba(72, 72, 102, 0.87);
  margin-left: 10px;
}
.totalbox {
  padding: 24px 24px;
  background: #fff;
  box-shadow: 0px 2px 4px 1px rgba(76, 122, 227, 0.12);
  border-radius: 4px 4px 4px 4px;
  /deep/ .el-table {
    .cell {
      height: 15px;
      line-height: 15px;
    }
  }
}
.tabsbox {
  border-radius: 4px;
}
/deep/ .el-input--small .el-input__icon {
  line-height: 38px !important;
}
/deep/ .el-form-item {
  margin-bottom: 14px;
}
</style>

