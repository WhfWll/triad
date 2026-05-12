<template>
    <div>
        <div class="main-title">
             
             <i class="nav_icon"></i>
            <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i> 
            <label >钓鱼目标管理</label>
        </div>
        <div class="gophishbox context_box_bg"> 
             <div class="search-box">
                <div class="operationbutton">
                    <xz-button type="primary" @click="AddGroup" size="small">新建目标组</xz-button>
                   
                </div> 
                <div class="serach-condition"> 
                    <div class="search-text">
                        <el-input placeholder="请输入关键字" @keydown.enter.native="handlesearch" v-model="formData.search" class="input-with-select" size="small"
                        >
                        </el-input>
                        <el-button type="primary" size="small" @click="handlesearch">搜索</el-button>
                    </div>
                    <div>
                        <el-button type="primary" size="small" @click="handleReset">重置</el-button>
                    </div>
                </div>
            </div>
             <el-table :data="tableData" style="width: 100%"  @selection-change="handleSelectionChange"  height="calc(100% - 102px)"
                @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave"> 
                <el-table-column prop="name" label="名称">
                </el-table-column>  
                <el-table-column prop="num_targets" label="邮箱数量">
                </el-table-column>   
                <el-table-column prop="last_time" label="修改时间" width="200">
                    <template slot-scope="scope">
                        <div v-if="showOperateButton && rowId == scope.row.id"> 
                            <el-link class="link_primary" :underline="false"
                                @click="updateGroup(scope.row)" > 编辑</el-link>  
                            <el-popover placement="bottom" width="170" :ref="`popover-${scope.row.id}`"
                                popper-class="delButton_popper">
                                <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                <div style="text-align: right; margin: 0">
                                    <el-button size="mini" class="delCancel" @click="fncancel(scope)">取消</el-button>
                                    <el-button size="mini" type="primary" @click="fnDel(scope)">确定</el-button>
                                </div>
                                <el-link :underline="false" class="link_primary" style="padding:0"
                                    slot="reference">删除</el-link>
                            </el-popover>  
                        </div>
                        <div v-else>
                            <span class="">{{ scope.row.modified_date }}</span> 
                        </div>
                    </template>
                </el-table-column>
            </el-table>
            <el-pagination :page-size="10" background layout=" total,  prev, pager, next,sizes, jumper" :total="totalpage"
                :current-page="currentpage" @current-change="currentchange" @size-change="handleSizeChange">
            </el-pagination>
        </div>

        <el-dialog :title="groupid!=0?'编辑钓鱼目标组':'创建钓鱼目标组'" :visible.sync="dialogaddFormVisible" :before-close="cancelform" width="900px"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false" class="newUserDialog">
            <div class="dialog_b_btn">
                <el-button size="small" @click="submitForm">确定</el-button>
                <el-button size="small" @click="cancelform">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form :model="groupform" label-width="0" status-icon ref="ruleFormadduser"  >
                    <el-form-item label=" " prop="username">
                        <label class="dialog_item_label">目标组名称</label>
                        <el-input v-model="groupform.name" size="small" style="width:320px" placeholder="请输入目标组名称"
                            maxlength="50"></el-input>
                        <div style="margin-left: 114px;">
                            <el-button type="primary" size="small" style="vertical-align: middle; margin-right: 16px"
                            @click="clickupload()">批量导入邮箱账号</el-button>
                            <input type="file" class="btnUploadID" ref="upload" @change="changeuploaID($event)"
                                style="display: none" id="input-file-ID" />
                            <el-link  class="link_primary" style="vertical-align: middle;" @click="downfile()">下载模板</el-link> 
                        </div> 
                    </el-form-item>
                    
                     <div>
                        <div>
                              <div class="dialog_item_label">邮箱内容</div> 
                        <el-button type="primary" size="mini" style="vertical-align: middle;float: right; "  
                                @click="clickadd()">新增</el-button>
                        </div>  
                         <div class="div_width"   style="margin-top:16px;margin-bottom:16px; ">
                            <el-table :data="login_conf" size="small" style="width: 100%">
                                <el-table-column  prop="target" label="姓名">
                                    <template slot-scope="scope">
                                        <span v-if="!scope.row.dataShow">{{ scope.row.first_name }}</span>    
                                        <el-input v-else v-model="scope.row.first_name" size="small" ></el-input>    
                                    </template>
                                </el-table-column>
                                <el-table-column prop="protocol"  label="账号邮箱"> 
                                   <template slot-scope="scope">
                                        <span v-if="!scope.row.dataShow">{{ scope.row.email }}</span>    
                                        <el-input v-else v-model="scope.row.email" size="small" ></el-input>    
                                    </template>
                                </el-table-column>
                                <el-table-column  prop="voucher" label="职位"> 
                                    <template slot-scope="scope">
                                        <span v-if="!scope.row.dataShow">{{ scope.row.position }}</span>    
                                        <el-input  type="textarea" :rows="1" v-else v-model="scope.row.position" size="small" resize="none"></el-input>    
                                    
                                    </template>
                                </el-table-column> 
                                <el-table-column  label="操作"  width="150">
                                    <template slot-scope="scope"  >
                                            <el-link :underline="false" @click="tbSave(scope)" 
                                                v-if="scope.row.dataShow" > 保存 </el-link>
                                            <el-link :underline="false" @click="tbUpdate(scope)" > 编辑 </el-link>
                                            <el-link :underline="false" @click="tbDelete(scope)" > 删除 </el-link>
                                        </template>
                                </el-table-column>
                            </el-table>
                        </div> 
                    </div>

                </el-form>
            </div>
        </el-dialog>


    </div>
</template>
<style scoped lang="less">
.gophishbox {
    padding: 24px;
    background: #fff;
    height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
</style>
<script>
import XzButton from "../../components/XzButton.vue";
import DelButton from "../../components/DelButton.vue";
import gophish from "@/api/gophish.js";
import _ from 'lodash'
var XLSX = require('xlsx');
export default ({
    name: 'usermanagement',
    components: {
        XzButton,
        DelButton,
    },
    data() {
        return{
            multipleSelection: [],
            currentpage: 1,
            totalpage: 0,
            page_num: 1,
            pageSize:10,
            tableData: [],
            isMore: false,
            showOperateButton: false,
            formData:{
                search:'',
            },
            dialogaddFormVisible:false,
            groupform:{
                name:'',
            },
            login_conf:[],
            groupid:0,
            rowId:0,
        }
    },
    created(){
         this.$store.state.activefirstMenu = '/group';
    },
    mounted(){
        this.getData(); //获得列表数据
    },
    methods:{
        AddGroup(){
            this.dialogaddFormVisible = true;
        },
        async updateGroup(row){ //编辑
            const dt = await gophish.groupdetail({
                id:row.id
            })
            if (dt.code === 200) {
                this.groupid =row.id;
                this.dialogaddFormVisible = true;
                this.groupform.name = dt.data.name; 
                const newArray = dt.data.targets.map(item => ({
                    ...item,
                    dataShow: false
                }));
                this.login_conf = newArray;

            } else {
                this.$message({
                    message: dt.msg,
                    type: 'error'
                });
            } 
        },
        // 下载模板
        downfile(){
             window.open('/group_template.csv');
        },
        handlesearch(){
            this.getData();
        },
        handleReset(){
            this.formData.search = '';
            this.getData();
        },
        async getData() { 
            const dt = await gophish.groupall({
               page:this.page_num,
                size:this.pageSize,
                search:this.formData.search,
            });
            if (dt.code === 200) {
                this.tableData = dt.data.groups;
                this.totalpage = dt.data.total;
            } else {
                this.$message({
                    message: dt.msg,
                    type: 'error'
                });
            } 
        },
        handleSelectionChange(val) {
            this.multipleSelection = val;
            this.selectedID = []
            this.multipleSelection.forEach(item => {
                this.selectedID.push(item.role)
            })
        },
        currentchange(t) {
            this.page_num = t;
            this.getData();
            this.currentpage = t;
        },
        handleSizeChange(t) {
            this.page_num = 1;
            this.pageSize = t;
            this.getData();
        },
        mouseenter(row, colum, cell, event) {
            this.showOperateButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断
        },
        mouseleave(row, colum, cell, event) {
            if (!this.$refs['popover-' + row.id]) {
                this.showOperateButton = false;
                this.rowId = "";
                return;
            } else {
                let isShow = this.$refs['popover-' + row.id].showPopper;
                if (!isShow) {
                    this.showOperateButton = false;
                    this.rowId = "";
                }

            } 
        },
        fncancel(scope) {
            scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
            scope._self.$refs[`popover-${scope.row.id}`].doClose()
        },
        async fnDel(scope) { //单个删除
           
            const data = await gophish.groupdelete({
                id:scope.row.id
            });
            if (data.code === 200) {
                this.$message({
                    message: data.msg || '删除成功！' ,
                    type: 'success'
                });
               
                this.getData();
            } else {
                this.$message({
                    message: data.msg || '错误，请重试',
                    type: 'error'
                });
            }

        },
        async submitForm(){ 
            if(this.groupid!=0){ //编辑
                const data = await gophish.groupupdate({
                    id:Number(this.groupid),
                    name:this.groupform.name,
                    targets:this.login_conf
                });
                if (data.code === 200) {
                    this.$message({
                        message:'编辑保存成功',
                        type: 'success'
                    });
                    this.dialogaddFormVisible = false;
                    this.getData();
                    this.cancelform();
                } else {
                    this.$message({
                        message: data.msg,
                        type: 'error'
                    });
                }
            }else{ //新增
                const data = await gophish.groupcreate({
                    name:this.groupform.name,
                    targets:this.login_conf
                });
                if (data.code === 200) {
                    this.$message({
                        message:'创建保存成功',
                        type: 'success'
                    });
                    this.dialogaddFormVisible = false;
                    this.getData();
                    this.cancelform();
                } else {
                    this.$message({
                        message: data.msg,
                        type: 'error'
                    });
                }
            }

           
        },
        cancelform(){
            this.dialogaddFormVisible = false;
            this.groupform.name='';
            this.login_conf=[];
            this.groupid = 0;
        },
        clickadd(){
             this.login_conf.push({
                dataShow: true,
            })
        },
        tbUpdate(scope){
            console.log(scope.$index)
            this.login_conf[scope.$index].dataShow = true;
        },
        tbDelete(scope){
            this.login_conf.splice(scope.$index, 1);
        },
        tbSave(scope){
            this.login_conf[scope.$index].dataShow = false;
        },
        clickupload(){
            document.querySelector(".btnUploadID").click();
        },
        changeuploaID(e){
            var that = this;
            var f = e.target.files[0];
            if (!f) return; 
            let fileSuffix = f.name.substr(f.name.lastIndexOf(".") + 1);
            if (fileSuffix.indexOf("csv") != -1) {
                
                var reader = new FileReader();
                reader.onload = function (e) { 
                    const csvData = e.target.result; // 这是字符串
                    try {
                        // 使用 XLSX 解析 CSV 字符串
                        const workbook = XLSX.read(csvData, { type: 'string', bookType: 'csv' });
                        const sheetName = workbook.SheetNames[0];
                        const worksheet = workbook.Sheets[sheetName];

                        // 转为 JSON，header: 1 表示按数组形式返回
                        const jsonData = XLSX.utils.sheet_to_json(worksheet, { header: 1 });

                        const arr = [];
                        // 假设第一行是表头，从第二行开始
                        for (let i = 1; i < jsonData.length; i++) {
                            const row = jsonData[i];
                            arr.push({
                                first_name: row[0] || "",
                                email: row[2] || "",
                                position: row[3] || "",
                                dataShow:false,
                            });
                        }

                        // that.login_conf = arr;
                        const merged = [... that.login_conf, ...arr];
                        const uniqueByMail = Array.from(
                        new Map(
                            merged.map(item => [item.email, item]) // 以 email 为 key，item 为 value
                        ).values()
                        );

                        that.login_conf = uniqueByMail;

                        console.log("解析结果：", arr);

                    } catch (error) {
                        console.error("解析 CSV 失败：", error);
                        alert("CSV 文件解析失败，请检查格式是否正确。");
                    }
 
                };
                 // ✅ 明确以 UTF-8 文本方式读取
                reader.readAsText(f, "UTF-8");
                // if (that.rABS) {
                //     reader.readAsArrayBuffer(f);
                // } else {
                //     reader.readAsBinaryString(f);
                // }
            } else{
                
            }

            e.target.value = "";
        },
    }
})
</script>