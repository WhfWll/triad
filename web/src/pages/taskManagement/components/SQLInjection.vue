<template>
    <div>
        <el-table
            :data="sqlTableDatas"
            tooltip-effect="dark"
            style="width: 100%">
            <!-- 行列均不固定，遍历列来实现表格显示 -->
            <div 
                v-for="(item,key) in sqlTableDatas[0]"
                :key="key">
                 <el-table-column
                    :label="key"
                    :prop="key">
                </el-table-column>
            </div>
           
            <!-- <el-table-column
                label="Password"
                prop="password">
            </el-table-column>
            <el-table-column
                label="Name"
                prop="name"> 
            </el-table-column>
            <el-table-column
                label="Sex"
                prop="sex">
            </el-table-column>
            <el-table-column
                label="Number"
                prop="number">
            </el-table-column> -->
        </el-table>
    </div>
</template>
<style lang="less" scoped>

</style>
<script>

export default ({
    name:'SQLInjection',
    props:{
        tableData:Array
    },
    data(){
        return{
            sqlTableDatas:[
                // {'user':'xxx','password':'xxx','name':'Bob','sex':'女','number':'123'},
                // {'user':'ccc','password':'ccc','name':'ccc','sex':'女','number':'ccc'},
            ]
        }
    },
    mounted(){
        this.getData();
    },
    methods:{
        getData(){
            this.sqlTableDatas = this.tableData;
        }
    }
})
</script>
