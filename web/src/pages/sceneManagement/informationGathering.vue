<!-- 渗透任务-信息收集页面 -->
<template>
  <div>
    <el-tabs v-model="subactiveName" @tab-click="handleClick" class="info_tabs"  > 
     <el-tab-pane label="服务" name="subtabs1">
      <serviceC v-if="subactiveName == 'subtabs1'"></serviceC>
     </el-tab-pane>
     <el-tab-pane label="站点" name="subtabs2" >
      <siteC v-if="subactiveName == 'subtabs2'"></siteC>
     </el-tab-pane>
     <el-tab-pane label="URL" name="subtabs3" > 
      <UrlC v-if="subactiveName == 'subtabs3'"></UrlC>
     </el-tab-pane> 
     <el-tab-pane label="子域名" name="subtabs4" > 
      <SubdomainC v-if="subactiveName == 'subtabs4'"></SubdomainC>
     </el-tab-pane> 
 
     </el-tabs>
 </div>
</template>
<style lang="less" >
@import './css/taskscenario.less';
</style>
<script>
import scene from '@/api/scene.js'
import DelButton from "@/components/DelButton.vue"
import serviceC from "@/pages/sceneManagement/infomationCle/service.vue"
import siteC from "@/pages/sceneManagement/infomationCle/site.vue"
import SubdomainC from "@/pages/sceneManagement/infomationCle/Subdomain.vue"
import UrlC from "@/pages/sceneManagement/infomationCle/Url.vue"
import { string } from 'jszip2/lib/support'

export default ({
   name:'infomationCle',
   components:{
       DelButton,
       serviceC,
       siteC,
       SubdomainC,
       UrlC
   },
   data(){
       return{
        subactiveName:'subtabs1',
           title:'任务场景',
           alldelvisible:false,
           showOperateButton:false,
           tableData:[],
           multipleSelection:[],
           total:0,
           pageSize:10,
           currentPage:1,
           search_item: {
               search: '',
               page: 1,
           },
           Loading:false,
           rowId:'',
           visible:false,
           alldelvisible1:false,
       }
   }, 
   created() {
       this.$store.state.activefirstMenu = '/task';
       this.userid = this.$commonjs.decryptCBC( localStorage.getItem('user_id'),this.$commonjs.myKey)
  
   },
   mounted(){ 
   },
   methods:{
      handleClick(){}, 
   }
})
</script>