<template>
    <el-button
        type="button"
        class="xz-button"
        :class="[`xz-button--${type}`, {
        'is-plain': plain,
        'is-round': round,
        'is-circle': circle,
        'is-disabled': disabled
        }]"
        :disabled="disabled"
        :size="size"
        @click="handleClick" 
        >
            <slot/>
        </el-button>
</template>  
<script>

export default { 
    name: "xzbutton",
    props:{
        size:String, 
        type: {
            // 数据类型： 字符串
            type: String,
            default: 'default'
        },
        plain: {
            type: Boolean,
            default: false
        },
        round: {
            type: Boolean,
            default: false
        },
        circle: {
            type: Boolean,
            default: false
        },
        icon: {
            type: String,
            default: ''
        },
        disabled: {
            type: Boolean,
            default: false
        } 
    },
    data(){
         
        return {
        }
    },
    
    created(){
        
    },
    methods:{
        handleClick (e) {
            this.$emit('click', e)
        }
    }
};
</script>

<style lang="less" scoped>

.xz-button {
  display: inline-block;
  line-height: 1;
  white-space: nowrap;
  cursor: pointer;
  background: #fff;
  border: 1px solid #dcdfe6;
  color: #606266;
  -webkit-appearance: none;
  text-align: center;
  box-sizing: border-box;
  outline: none;
  margin: 0;
  transition: 0.1s;
  font-weight: 500;
  // 禁止元素的文字被选中
  -moz-user-select: none;
  -webkit-user-select: none;
  -ms-user-select: none; 
  font-size: 14px;
  border-radius: 4px;
  &:hover,
  &:focus {
    color: #4c7ae3;
    border-color: #c6e2ff;
    background-color: #ecf5ff;
  }
  /* margin-right: 8px; */
} 
.xz-button--primary {
  color: #fff;
  background-color: #4c7ae3;
  border-color: #4c7ae3;

  &:hover,
  &:focus {
    background: #4c7ae3;
    border-color: #4c7ae3;
    color: #fff;
  }
}
.xz-button--warning {
  color: rgba(243, 95, 40, 0.5);
  background-color: #fff;
  border-color: rgba(243, 95, 40, 0.5);
} 

.xz-button--warning {
  color: #f35f28 ;
  background-color: #fff;
  border-color: #f35f28;

   &:hover,
  &:focus{
    background: #f35f28;
    border-color: #f35f28;
      color: #fff;
  }
} 
.el-button--primary.is-disabled {
    background-color: rgba(76, 122, 227, .5) !important;
    border-color: rgba(76, 122, 227, .2) !important;
}
.xz-button--warning.is-disabled, 
.xz-button--warning.is-disabled:active, 
.xz-button--warning.is-disabled:focus, 
.xz-button--warning.is-disabled:hover {
    color: rgba(243, 95, 40, 0.5);
    background-color: #ffffff;
    border-color: rgba(243, 95, 40, 0.5);
}
.xz-button--blueBorderWhiteBg {
  color: #4C7AE3 ;
  background-color: #fff;
  border-color: #4C7AE3;
  &:hover,
  &:focus{
    background: #ecf5ff;
    border-color: #4C7AE3;
    color: #4C7AE3;
  }
}

.el-button--small, 
.el-button--small.is-round {
    padding: 8px 12px;
    vertical-align: middle;
}
.el-button--small {
    border-radius: 2px;
    font-size: 14px;
}
.el-button--mini, 
.el-button--small {
    font-size: 12px;
    border-radius: 3px;
}
/deep/ .el-button--mini, .el-button--small {
    font-size: 14px;
    border-radius: 3px;
}
</style>