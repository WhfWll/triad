<template>
  	<div class="systemtip" >  
          <div class="toplogo" >
            <!-- <img src="../../assets/logo.png" style="height:23px">  -->
            <span style="padding-left:12px;font-weight: 600;">现场漏洞验证工具集
            </span>
          </div> 
           <div class="infos">
             <p>产品：现场漏洞验证工具集
            </p>
             <p v-for="(item, index) in infos" :key="index" :class="item === '' ? 'noContent' : ''">{{item}}</p>
           </div>
    </div>
</template>
<script> 
import  {current_version_info}   from '@/api/system.js' 
 export default({
    name:'nowsystip',
    data(){
      return {
        infos:'',
      }
    },
    created () {
      this.getOfflineStatus();
    },
    methods: {
      async getOfflineStatus (time) { 
            const res = await current_version_info();
            if(res.success){
                this.infos = res.result.split('\n') 
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            }
        }, 
    },
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
        line-height:48px;
        background: #4C7AE3;
        box-shadow: 0px 2px 6px 1px rgba(0, 21, 41, 0.11999999731779099);
        border-radius: 0px 0px 0px 0px;
        span{
          color: #FFF;
          padding-left: 12px;
          font-weight: 600;
        }
    }
    .infos{
      padding: 40px 60px;
      box-sizing: border-box;
      p{
        line-height: 40px;
        font-size:16px;
        &.noContent + p{
          font-weight: bold;
          margin-top: 20px;
        }
        &:nth-child(1), &:nth-child(2), &:nth-child(3){
          font-weight: bold;
        }
      }
    }
  }
    
    
</style>

