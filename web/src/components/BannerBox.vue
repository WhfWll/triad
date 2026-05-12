<template>
    <div class="content_title  tastBtnCont"  > 
            <div class="search-box"  :class="tips == ''?'padding':''">
                <div class="operationbutton">  
                    <slot/>
                </div>
                <div class="tipsCont" v-if="tips">
                    <label class="topTips">提示：</label> 
                    <span  >{{tips}}</span>
                </div>
            </div>
        </div>
</template>
<style lang="less" scoped>
.padding{
    padding: 16px 24px !important;
}
.content_title{
    background: url(../assets/images/banner@2x.png) no-repeat;
    background-size: contain;
    background-position: right;
    background-color: #fff; 
    box-shadow: 0px 2px 4px 0px rgba(76,122,227,0.12);
}
.tastBtnCont {
        font-size: 13px;
        .search-box {
            padding: 24px;
            border-radius: 4px;
            min-height: 0;
            box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
            .operationbox {
                label {
                    margin-left: 24px;
                    color: #4c7ae3;
                    font-size: 13px;
                }
                span {
                    margin-left: 4px;
                    font-size: 13px;
                
                }
            }
        }
        .tipsCont {
            float: left; 
            line-height: 31px;
            margin-left: 24px;
            .topTips {
                color: #4c7ae3;
                width: 40px;
            }
            span{
                color:rgba(72,72,102,0.64); 
                font-size: 13px;
            }
        }
    }
</style>
<script>
export default{
    name:'',
    props:{
        tips:{}
    },
    data(){
        return{
            
        }
    },
    mounted(){

    },
    methods:{
        
    }
}
</script>