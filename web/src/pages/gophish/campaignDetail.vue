<template>
    <div> 
        <div class="main-title">
            <i class="nav_icon"></i>
            <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i> 
          
            <router-link :underline="false" class="classA" :to="{ path: '/campaign' }" >钓鱼活动管理
            </router-link>  
            <label class="currentpagetitle">
                <el-tooltip class="item" effect="dark"    placement="bottom">
                    <span>钓鱼活动详情</span>
                </el-tooltip>
            </label>
        </div> 
        <div class="gophishbox context_box_bg"> 
            <el-table :data="tableData" style="width: 100%"   @selection-change="handleSelectionChange"
                @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave">  
                <el-table-column prop="email" label="邮箱">
                </el-table-column>  
                <el-table-column prop="first_name" label="姓名">
                </el-table-column> 
                 <el-table-column prop="position" label="职位">
                </el-table-column> 
                <el-table-column prop="status" label="状态" width="200">
                    <template slot-scope="scope">
                        <div v-if="showOperateButton && rowId == scope.row.id">
                            <el-link class="link_primary" :underline="false" @click="handleinfo(scope.row)" >
                                详情</el-link> 
                        </div>
                        <div v-else>
                            <span class="">{{ scope.row.status }}</span>

                        </div>
                    </template>
                </el-table-column>
            </el-table> 
 
        </div>
        <el-dialog  title="详情" 
            :visible.sync="dialogaddFormVisible" 
            :before-close="cancelform" width="1000px"
            :close-on-click-modal="false" 
            :validate-on-rule-change="false" 
            :show-close="false"  
            class="newUserDialog">
            <div class="dialog_b_btn"> 
                <el-button size="small" @click="cancelform">关闭</el-button>
            </div>
            <div style="padding:24px">
                <div>
                    <label for="">钓鱼目标名称：</label>
                    <el-input v-model="target_name" size="small" style="width:500px" disabled ></el-input>
                </div>
                <div class="timeline"> 
                    <div v-for="(item,i) in timelinedt" :key="i">   
                        <div v-if="item.message == 'Campaign Created'">
                            <i class="c_0"></i> 
                            <label for="">活动已创建</label>
                            <span>{{item.time}}</span>
                        </div>
                        <div v-if="item.message == 'Email Sent'">
                            <i class="c_1"></i> 
                            <label for="">邮件已发送</label>
                            <span>{{item.time}}</span>
                        </div>
                        <div v-if="item.message == 'Email Opened'">
                            <i class="c_2"></i> 
                            <label for="">邮件已被打开</label>
                           <span>{{item.time}}</span>
                        </div>
                        <div v-if="item.message == 'Clicked Link'">
                            <i class="c_3"></i> 
                            <label for="">已点击钓鱼链接</label>
                            <span>{{item.time}}</span>
                            <div class="detail">
                                Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36
                            </div>
                        </div>
                        <div v-if="item.message == 'Submitted Data'">
                            <i class="c_4"></i> 
                            <label for="">已获取钓鱼数据</label>
                             <span>{{item.time}}</span>
                            <el-table
                                :data="item.details.payload"
                                style="width: 100%;margin-top: 10px;">
                                <el-table-column
                                    prop="key"
                                    label="参数名" >
                                </el-table-column>
                                <el-table-column
                                    prop="value"
                                    label="值" >
                                </el-table-column> 
                            </el-table>
                        </div>
                    </div>

                    <!-- <div>
                        <label for="">邮件已发送</label>
                        <span>2025-10-10-11</span>
                    </div>
                    <div>
                        <label for="">邮件已被打开</label>
                        <span>2025-10-10-11</span>
                    </div>
                    <div>
                        <label for="">已点击钓鱼链接</label>
                        <span>2025-10-10-11</span>
                        <div class="detail">
                            Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36
                        </div>
                    </div>
                    <div>
                        <label for="">以获取钓鱼数据</label>
                        <span>2025-10-10-11</span>
                           <el-table
                            :data="tableData"
                            style="width: 100%">
                            <el-table-column
                                prop="date"
                                label="参数名" >
                            </el-table-column>
                            <el-table-column
                                prop="name"
                                label="值" >
                            </el-table-column> 
                        </el-table>
                    </div> -->
                </div>
            </div>
        </el-dialog> 

    </div>
</template>
<style scoped lang="less">
.gophishbox {
    padding: 24px;
    background: #fff;
    min-height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
.timeline{
    margin-top: 10px;
    >div{
        margin-bottom: 20px;
        label{
            display: inline-block;
            width: 200px;
            font-size: 14px;
            color: rgba(72, 72, 102, 0.87);
            vertical-align: middle;
        }
        span{
            font-size: 14px;
            color: rgba(72, 72, 102, 0.32);
        }
        .detail{
            margin-top: 10px;
        }
        i{
            display: inline-block;
            width: 36px;
            height: 36px; 
            vertical-align: middle;
            margin-right: 10px;
        }
        .c_0{
            background:url(../../assets/images/c_0.png);
            background-size: contain;
        }
        .c_1{
            background:url(../../assets/images/c_1.png);
            background-size: contain;
        }
        .c_2{
            background:url(../../assets/images/c_2.png);
            background-size: contain;
        }
        .c_3{
            background:url(../../assets/images/c_3.png);
            background-size: contain;
        }
        .c_4{
            background:url(../../assets/images/c_4.png);
            background-size: contain;
        }
    }
}
</style>
<script> 
import user from "@/api/user.js";
import _ from 'lodash'
import gophish from "@/api/gophish.js";
export default ({
    name: 'usermanagement', 
    data() { 
        return { 
            currentpage: 1,
            totalpage: 0,
            page_num: 1,
            tableData: [], 
            pageSize: 10,
            showOperateButton: false,
            id: this.$route.query.id,  
            name:this.$route.query.name,
            rowId:0,
            dialogaddFormVisible:false,
            timeline:[],
            target_name:'',
            timelinedt:[],
        }
    },
    watch: { 
    },
    created: function () {
        this.$store.state.activefirstMenu = "/campaign"; 
    },
    mounted: function () {
        this.getinfo(); //获得列表数据
    },
    methods: { 
        async getinfo() {
            this.Loading = false;
            const dt = await gophish.campaignresult({
               id:this.id
            });
            if (dt.code === 200) {
                this.tableData = dt.data.results; 
                this.timeline = dt.data.timeline;
            } else {
                this.$message({
                    message: dt.msg,
                    type: 'error'
                });
            } 
        }, 
        async handleinfo(row){
            this.dialogaddFormVisible = true;
            this.target_name = row.email;
            //根据email筛选
            const filtered = this.timeline.filter(item => 
                item.email === row.email || item.email === ""
            );
            //把Detail字符串解析成对象
            const parsedTimeline = filtered.map(item => {
                if (item.details) {
                    try {
                    return {
                        ...item,
                        details: JSON.parse(item.details)  // 替换为解析后的对象
                    };
                    } catch (e) {
                    console.error("解析失败:", e);
                    return item;
                    }
                }
                return item;
            });
            //把Detail中payload 转换成数组
            const result = parsedTimeline.map(item => {
                if (item.details?.payload) {
                    // 将 payload 对象转为数组 [{ key, value }]
                    const payloadArray = Object.entries(item.details.payload)
                    .map(([key, value]) => ({ key, value }));

                    item.details.payload = payloadArray;
                }
                return item;
            });

            this.timelinedt = result;

            console.log(this.timelinedt)

        },
        fncancel(scope) {
            scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
            scope._self.$refs[`popover-${scope.row.id}`].doClose()
        },
        async fnDel(scope) { //单个删除 
            const data = await gophish.campaigndelete({
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
                    message: data.msg,
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
        cancelform() { 
            this.dialogaddFormVisible = false; 
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

    }
})

</script>