<template>
    <el-popover
        popper-class="delButton_popper"
        placement="bottom-start"
        :width="width"
        style="padding-right:8px"
        trigger="click" 
        :visible-arrow="false"
        v-model="alldelvisible" >
        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
        <div style="text-align: right; margin: 0" class="" >
            <el-button size="mini" class="delCancel" @click="alldelvisible = false" >取消</el-button>
            <el-button size="mini" type="primary"   @click="handleClick">确认</el-button>
        </div>   
        <xz-button type="warning" slot="reference" :disabled="disabled" size="small"  >删除</xz-button> 
    </el-popover>    
</template>  
<script>
import XzButton from "./XzButton.vue";
export default { 
    name: "delbutton",
    components: {
    	XzButton,
  	},
    props:{
        width:Number, 
        txt:String,
        btnTruetxt:String,
        btnCanceltxt:String,
        disabled: {
            type: Boolean,
            default: false
        } 
    },
    data(){ 
        return {
            alldelvisible:false
        }
    }, 
    methods:{
        handleClick (e) {
            this.alldelvisible = false
            this.$emit('click', e)
        },
        handCancel(e){
            this.$emit('click', e)
        }
    }
};
</script>

<style lang="less" scoped>
.delText{
    margin-bottom: 16px ;
    color:rgba(72,72,102,0.64);
    i{
        color: #F9B640;
        margin-right: 10px;
    }
}
.delCancel{
    color: rgba(72, 72, 102, 0.64) !important;
    border: 1px solid #E8E8F5 !important;
}
.delCancel:hover{
    color: rgba(72, 72, 102, 0.64) !important; 
    border: 1px solid #E8E8F5 !important; 
    background: #fff !important; 
}
/deep/ .el-button--mini, .el-button--small {
    font-size: 14px;
    border-radius: 3px;
}
</style>