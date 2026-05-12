<!-- 渗透任务-信息收集页面 -->
<template>
  <div>
   <div class="templatelist context_box_bg" style="border-top: 1px solid #E8E8F5; ">
     <div class="search-box">
       <div class="operationbutton">  
         <el-button type="primary"    v-show="false">新建</el-button> 
         <del-button
           :width="170"
           @click="btnMultiDelete"
           :disabled="!multipleSelection.length"
           style="margin-left: 8px"
         >
         </del-button> 
       </div>
       <div class="serach-condition">
         <div class="search-text">
           <el-input
             placeholder="请输入关键字"
             @keydown.enter.native="handlesearch"
             v-model="search_item.search"
             class="input-with-select"
             
             clearable
           >
           </el-input>
           <el-button type="primary"   @click="handlesearch">搜索</el-button> 
         </div>
         <div>
           <el-button type="primary"   @click="handleReset">重置</el-button> 
         </div>
       </div>
     </div>
     <el-table
       ref="multipleTable"
       :data="tableData" 
       v-loading="Loading"
       style="width: 100%"
       @selection-change="handleSelectionChange"
       @cell-mouse-enter="mouseenter"
       @cell-mouse-leave="mouseleave"
      height="calc(100% - 102px)"
     >
     <!--  -->
       <el-table-column type="selection" width="55" :selectable="checkboxT">
       </el-table-column>
       <!-- <el-table-column prop="ip" label="IP">
         <template #default="scope">
           <el-link @click="fnDetails(scope.row)">{{
             scope.row.ip
           }}</el-link>
         </template> 
       </el-table-column> -->
       <el-table-column prop="url" label="URL" :show-overflow-tooltip = "true"> </el-table-column>
       <el-table-column prop="method" label="方法" :show-overflow-tooltip = "true"> </el-table-column>
       <el-table-column prop="status" label="状态码" :show-overflow-tooltip = "true"> </el-table-column>
       <el-table-column prop="title" label="标题" :show-overflow-tooltip = "true"> </el-table-column>
       <el-table-column prop="param" label="参数" :show-overflow-tooltip = "true"> </el-table-column>
       <el-table-column prop="riskType" label="风险类型" :show-overflow-tooltip = "true"> </el-table-column>
       <el-table-column prop="bodyLength" label="body长度" :show-overflow-tooltip = "true"> </el-table-column>
       <el-table-column prop="source" label="来源" :show-overflow-tooltip = "true"> </el-table-column>
       <!-- <el-table-column prop="isDefault" label="默认场景">
         <template #default="scope">
           <span v-if="scope.row.isDefault == 1">是</span>
           <span v-else>否</span>
         </template>
       </el-table-column> -->
       <!-- <el-table-column prop="createTime" label="创建时间"> </el-table-column>
       <el-table-column prop="userName" label="提交者">
         <template #default="scope">
           <div
             v-if="showOperateButton && rowId == scope.row.id"
           >  
             <el-popover
                 placement="bottom"
                 width="200"  
                 popper-class="delButton_popper" 
                 trigger="click"  
                 :ref="`popover_id-${scope.row.id}`"
               >
                 <p class="delText">
                   <i class="el-icon-warning"></i>确定设置为默认场景吗？
                 </p>
                 <div style="text-align: right; margin: 0">
                   <el-button
                     size="small"
                     class="delCancel"
                     @click="fncancel1(scope) "
                     >取消1
                   </el-button>
                   <el-button
                     size="small"
                     type="primary"
                     @click="btnDefault(scope)"
                     >确定1
                   </el-button>
                 </div>
                 <template #reference>
                   <span   class="link_info linkafter" style="padding:0;cursor: pointer;">默认</span>
                 </template>
               </el-popover>
             <el-link
                   :underline="false"
                   class="link_primary"
                   style="vertical-align: initial;"
                   @click="copyScene(scope.row)"
                   >复制</el-link
                 >
             <el-popover
               v-if="scope.row.isDefault !=1"
               placement="bottom"
               width="170"
               :visible-arrow="false"
               :ref="`popover-${scope.row.id}`"
               popper-class="delButton_popper"
             >
               <p class="delText">
                 <i class="el-icon-warning"></i>确定删除吗？
               </p>
               <div style="text-align: right; margin: 0">
                 <el-button
                   size="small"
                   class="delCancel"
                   @click="fncancel(scope)"
                   >取消</el-button
                 >
                 <el-button
                   size="small"
                   type="primary"
                   @click="btnDel(scope)"
                   >确定</el-button
                 >
               </div>
               <span slot="reference" class="link_danger linkafter2" style="cursor: pointer;">删除</span>
             </el-popover>
              
           </div>
           <div v-else>
             <span>{{ scope.row.userName }}</span>
           </div>
         </template>
       </el-table-column>  -->
     </el-table>
     <el-pagination 
       background
       @size-change="handleSizeChange"
       @current-change="handleCurrentChange"
       :current-page="currentPage"
       :page-size="pageSize"
       layout=" total,  prev, pager, next, sizes,jumper"
       :total="total"
     >
     </el-pagination>
   </div>
 </div>
</template>
<style lang="less" scoped>
@import '../css/taskscenario.less';
/deep/.el-table__empty-text{
  text-align: center;
}
</style>
<script>
import scene from '@/api/taskCenter/infomation'
import DelButton from "@/components/DelButton.vue"
export default ({
   name:'infomationCle',
   components:{
       DelButton
   },
   data(){
       return{
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
      //  this.$store.state.activefirstMenu = '/infomationCle';
      //  this.userid = this.$commonjs.decryptCBC( localStorage.getItem('user_id'),this.$commonjs.myKey)
   },
   mounted(){
       this.getData();
   },
   methods:{

       getData(){
        let that = this
           let params = {
            objType:1,//int，固定传1,必选
            subObjType:'1_3',//string ,1_1服务，1_2站点，1_3url,1_4子域名,必选
            // objId:'26',//string,任务id,必选
            objId:that.$route.query.id||'0',//string,任务id,必选
            search: that.search_item.search,//string,筛选项,可选
            page: that.search_item.page,
            size: that.pageSize, 
           }
           scene.getData(params).then(res =>{
               if(res.code == 200){
                console.log(res.data)
                that.tableData = res.data.list;
                that.total = res.data.total
               }else{
                that.$ElMessage({
                       message: res.msg,
                       type: 'error'
                   })
               }
           }).catch(err=>{

           })
       },
       checkboxT (row, index) {
           if (row.isDefault == 1  ) {
               return 0
           } else {
               return 1
           }
       },

       handlesearch(){
           this.search_item.page = 1
           this.currentPage = 1
           this.getData()
       },
       handleReset(){ 
           this.search_item.page = 1
           this.search_item.search = ''
           this.pageSize = 10
           this.currentPage = 1
           this.getData()
       },
       handleSizeChange(t){
           this.search_item.page = 1
           this.pageSize = t
           this.getData()
       },
       handleCurrentChange(t){
           this.search_item.page = t
           this.getData()
       },
       fncancel(scope){
           scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
           scope._self.$refs[`popover-${scope.row.id}`].doClose()
       },
       fncancel1 (scope) { 
            
           scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
           scope._self.$refs[`popover-${scope.row.id}`].doClose()
           // this.clickEmpty()//wumeng处理两个popovar下下策
       },
       // 单个删除
       async btnDel(scope){
        let that = this

         let _ids =  scope.row.id 

         const res = await scene.delScene({
            taskTaskResultIds :  _ids ,
            taskId :that.$route.query.id||'0', //string,信息收集的主键id,多个用英文逗号隔开，必选
         })
         if (res.code == 200) {
           this.$message({
             message: '删除场景成功',
             type: 'success'
           })
           scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
           // this.deldialogVisible = false;
           this.getData()

         } else {
           this.$message({
             message: res.msg,
             type: 'error'
           })
         }
       },
 
  
       // 批量删除
       async btnMultiDelete(){
        let that = this
         if (this.multipleSelection.length == 0) return

         var _ids = this.multipleSelection.map(item => item.id)

         let params = {
            taskTaskResultIds: _ids.join(','),
            taskId :that.$route.query.id||'0', //string,信息收集的主键id,多个用英文逗号隔开，必选
         }

         const res = await scene.delScene(params)
         if (res.code == 200) {
           this.$message({
             message: '批量删除成功',
             type: 'success'
           })
           this.alldelvisible = false
           this.getData()
         } else {
           this.$message({
             message: res.msg,
             type: 'error'
           })
         }
       },
       handleSelectionChange (val) {
           this.multipleSelection = val
       },
       mouseenter (row, column, cell, event) {
           this.showOperateButton = true;
           let _id = row.id; 
           this.rowId = _id ;//赋值行id，便于页面判断
       },
       mouseleave (row, colum, cell, event) {  
           let t = this.$refs['popover_id-' + row.id] && this.$refs['popover_id-' + row.id].showPopper;
           let t2 = this.$refs['popover-' + row.id] && this.$refs['popover-' + row.id].showPopper;

           if (!t && !t2) {
               this.showOperateButton = false;
               this.rowId = "";
           }

       },
   }
})
</script>