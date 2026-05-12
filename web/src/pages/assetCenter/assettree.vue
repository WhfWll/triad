<template> 
  <div >
    <!-- style="background-color: #fff; " -->
     <div class="main-title  ">  
        
        <i class="nav_icon"></i>
      <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i> 
      <label class="taskSceneBtn" >资产树</label>
    </div> 
    <div class="box_Style" style="height: calc(100% - 39px);" >
        <!-- 左侧资产树 -->
        <RightZiChan style="width: 280px;" />
        <!-- 右边资产列表 -->
        <ZiChan style="width: calc(100% - 290px);" />
    </div>
  </div>
</template>
<style lang="less" scoped>
.box_Style{
    background-color: #fff;display:flex;  padding-top: 12px; box-sizing: border-box;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);border-radius: 4px;
}
</style>
<script>
import  RightZiChan from './component/RightZiChan.vue'
import  ZiChan from './component/ZiChan.vue'
export default {
    components:{
        RightZiChan,
        ZiChan,
    },
    created(){
        this.$store.state.activefirstMenu="/assettree";  
    },
}
</script>
