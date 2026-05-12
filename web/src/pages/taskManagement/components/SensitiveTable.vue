<template>
    <div>
        <el-table
            :data="sensitiveTableData"
            tooltip-effect="dark"
            style="width: 100%">
            <el-table-column
                label="文件名称"
                prop="filename">
            </el-table-column>
            <el-table-column
                label="文件大小"
                prop="size">
            </el-table-column>
            <el-table-column
                label="更新时间"
                prop="update_time">
            </el-table-column>
            <el-table-column
                prop=""
                label="操作"
            >
                <template slot-scope="scope" > 
                    <div >
                        <el-link :underline="false"  @click="downfile(scope.row)">下载</el-link>
                        <!-- <el-link 
                            class="link_danger" 
                            :underline="false"   >
                                <el-popover
                                    placement="bottom"
                                    width="170"   
                                    :visible-arrow="false"
                                    :ref="`popover_id-${scope.$index}`"
                                    popper-class="delButton_popper" >
                                    <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                    <div style="text-align: right; margin: 0">
                                        <el-button size="mini" class="delCancel" @click="scope._self.$refs[`popover_id-${scope.$index}`].doClose()">取消</el-button>
                                        <el-button size="mini" type="primary" @click="targetDelete(scope,'yes')">确定</el-button>
                                    </div>  
                                    <span  slot="reference"  >删除</span>
                                </el-popover>
                        </el-link> -->
                    </div>
                </template>
            </el-table-column>
        </el-table>
    </div>
</template>
<style lang="less" scoped>

</style>
<script>
import jsFileDownload from 'js-file-download'
export default ({
    name:'sensitiveTable',
    props: {
        tableData: Array
    },
    data(){
        return{
            sensitiveTableData:[
                // {'file_name':'xxx','file_size':'xxx','update_time':'2020-1-2'}
            ]
        }
    },
    mounted() {
        this.getData();
    },
    methods: {
        getData() {
            this.sensitiveTableData = this.tableData;
        },
        downfile(row){
            jsFileDownload(row.content, row.filename);
        }
    }
})
</script>
