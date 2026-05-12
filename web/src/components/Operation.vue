<!--
    dialog组件
-->
<template>
    <div>
        <div class="search-box">
            <el-button
                v-if="one"
                type="primary"
                @click="handleOne"
                size="small"
                >{{ one || "" }}</el-button
            >
            <!-- <el-button
                v-if="two"
                type="warning"
                @click="handleTwo"
                size="small"
                plain
                :disabled="delList.length < 0"
                >{{ two || "" }}</el-button
            >  -->
            <el-popover
                v-if="two"
                popper-class="delButton_popper"
                placement="bottom-start"
                width="170"
                style="padding-left:8px"
                trigger="click" 
                :visible-arrow="false"
                v-model="alldelvisible" >
                <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                <div style="text-align: right; margin: 0" class="" >
                    <el-button size="mini" class="delCancel" @click="alldelvisible = false" >取消</el-button>
                    <el-button size="mini" type="primary"  @click="handleTwo">确定</el-button>
                </div>  
                <el-button type="warning"  size="small"  slot="reference" :disabled="!delList.length">删除</el-button> 
            </el-popover>  
            <div class="serach-condition">
                <el-select 
                    v-model="value" 
                    :placeholder="selectPlaceholder" 
                    v-if="selectList.length > 0"
                    @change="handleSelectChange"
                    clearable
                >
                    <el-option
                        v-for="item in selectList"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value">
                    </el-option>
                </el-select>
                <div class="search-text">
                    <el-input
                        placeholder="请输入关键字"
                        @keydown.enter.native="handleThree" 
                        v-model="search"
                        class="input-with-select"
                        size="small"
                        clearable
                    >
                    </el-input>
                    <el-button
                        type="primary"
                        size="small"
                        @click="handleThree"
                        >{{ three || "" }}</el-button
                    >
                </div>
                <div>
                    <el-button
                        type="primary"
                        size="small"
                        @click="handleFour"
                        >{{ four || "" }}</el-button
                    >
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    props: {
        one: String,
        two: String,
        three: String,
        four: String,
        selectList:{
            type: Array,
            default: function () {
                return []
            }
        },
        selectPlaceholder:{
            type:String,
            default:"请选择"
        },
        delList:{
            type:Array,
            default:()=>{
                return []
            }
        }
    },
    data:()=>({
        search:"",
        value:"",
        alldelvisible:false
    }),
    methods: {
        handlesearch() {},
        handleOne(){
            this.$emit('handleOneClick')
        },
        handleTwo(){
            if(this.delList.length < 1){
                return this.$message.error('请先勾选需要删除的数据！')
            }
            this.alldelvisible = false;
            this.$emit('handleTwoClick')
        },
        handleThree(){
            this.$emit('handleThreeClick',this.search)
        },
        handleFour(){
            if (this.four === '重置') {
                this.search = ''
            }
            this.$emit('handleFourClick')
        },
        handleSelectChange(select){
            this.$emit('handleSelectChange',select)
        }

    },
};
</script>

<style lang="less" scoped>
.el-select{
    width: 140px;
}
</style>