<template>
  	<div class="systemtip" >  
          <div class="toplogo" >
            <!-- <img src="../../assets/logo.png" style="height:23px">  -->
            <!-- <span style="padding-left:12px;font-weight: 600;">自动化渗透测试平台</span> -->
          </div> 
          <div class="centerpart">
            <el-progress type="circle" :percentage="percentage" :stroke-width="10" :width="200" class="centerbar"></el-progress>
            <!-- <p class="tip" >系统更新过程持续时间10分钟左右，升级过程中不能访问系统...</p> -->
            <p class="tip" v-if="percentage !== 100">系统正在更新中，请终止对本系统的操作，以防更新失败...</p>
            <p class="tip" v-if="percentage == 100">系统更新已完成，正在重启中...</p>
          </div>
    </div>
</template>
<script> 
import { mapState } from 'vuex'
// import { delUserInfo } from '@/utils'
import { logout } from '@/lib'
 export default({
    name:'systemsetting',
    data(){
      return {
        updateStatus: '',
        // percentage: 0,
        arr: [],
        timer: null,
        timer2: null
      }
    },
    watch: {
      async percentage (val) {
        this.$store.commit('setPercentage', val)
        if (val === 100) {
          clearInterval(this.timer)
          let userID = localStorage.getItem('user_id');
          const res = await logout({userId: userID});
            localStorage.removeItem('user'); 
            localStorage.removeItem('user_id'); 
            localStorage.removeItem('role'); 
            this.$router.push('/login');
        }
      }
    },
    computed: {
    ...mapState({
    }),
    percentage: {
      get () {
        return this.$store.state.uploadPercentage
      },
      set (val) {
        this.$store.state.uploadPercentage = val
      }
    },
  },
    created () {
      // this.getOfflineStatus();
      this.timer = setInterval(() => {
        if (this.percentage < 100) {
          this.percentage = this.percentage + 1
        }
      }, 1000)
    },
    methods: {
      // getOfflineStatus (time) {
      //       this.$ajax.get('/systems/interfaces/off_update/status/',{
      //         params:{
      //           system_uid: this.$store.state.system_uid
      //         }
      //       }).then((res)=>{   
      //               var dt = res.data;
      //               if(dt.success){
      //                    this.updateStatus = dt.status;
      //                    switch (dt.status) {

      //                       case 'running':
      //                           setTimeout(() => {
      //                             this.getOfflineStatus() // 持续请求该路由，直到‘finish'结束请求
      //                           }, 1000)                              
      //                           break
      //                       case 'finish':
      //                       this.percentage = 100;
      //                       clearInterval(this.timer)
      //                           break
      //                       case 'fail':
      //                           this.$message.error('更新失败，请联系管理员')
      //                           // this.newWin.close() // 直接结束更新
      //                           this.$router.push({
      //                                   path: "/systemsetting",
      //                                   query: {}
      //                               });
      //                           break
      //                   }
      //               }else{
      //                 this.$message.error(dt.error)
      //                 this.$router.push('/systemsetting')
      //               }
      //           })
      //           .catch(err=>{
      //               setTimeout(() => {
      //                   this.getOfflineStatus()
      //               }, 1000)
      //           })
      //   },

    },
    beforeDestroy () {
      clearInterval(this.timer)
      clearTimeout(this.timer2)
    }
  })
 
</script>
<style scoped lang="less">
  .systemtip{
    width: 100%;
    height: 100%;
    background-color: #FFF;
    display: flex;
    flex-direction: column;
    .toplogo{
        width: 100%;
        height: 48px;
        line-height:60px;
        background: #4C7AE3;
        box-shadow: 0px 2px 6px 1px rgba(0, 21, 41, 0.11999999731779099);
        border-radius: 0px 0px 0px 0px;
    }
    .centerpart{
      flex: 1;
      padding-top: 132px;
      text-align: center;
      .tip{
        font-size: 14px;
        font-weight: 400;
        color: #4C7AE3;
        margin-top: 64px;
      }
    }
  }
    
    
</style>

