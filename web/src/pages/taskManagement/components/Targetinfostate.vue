<template>
    <!-- 不同目标状态：目标总数、存活目标、不存活目标、未知状态、存活检测中目标、未检测目标 -->
    <el-col :span="6" class="target_showbox_col">
        <div :class="{[activelayout]:activelayoutflag ==flag}"> 
            <div :class="'msg_icon icon' + flag">
                <i :class="'iconfont '+iconfont"></i>
            </div> 
            <div class="msg_label">{{label}}
                <el-tooltip v-if="content" placement="top">
                    <div slot="content">{{content}}</div> 
                    <i class="iconfont icontishi icontsstyle" style=" margin-left: 4px;vertical-align: middle; width: 16px; height: 16px;"  ></i> 
                </el-tooltip>
            </div>
            <div class="msg_count"   >
                <label>{{targetNumber}}</label>  
            </div> 
        </div>
    </el-col>
</template>

<script>
export default {
    name: "targetinfostate",
    props:{
        label:{},
        content:{},
        className:{},
        flag:{},
        activelayoutflag:{},
        icon:{},
        iconfont:{},
        targetNumber:{},
        targetdata:{}
    },
    data: () => ({
        activelayout:'', 
    }), 
    methods: {
       
    }
};
</script>

<style lang="less" scoped>

.target_showbox{
    background:#F7F7FB;
    // min-height: 200px;
    border-radius: 4px;
    padding: 24px;
    box-sizing: border-box;
    .target_showbox_col{
        >div{
            height: 48px;
            display: flex;
            align-items: center;
            background: #fff;
            border-radius: 4px;
            /* padding: 12px 24px 14px 24px; */
            padding: 14px 24px;
            font-size: 0;
            box-sizing: border-box;
            box-shadow: 0px 2px 4px 0px rgba(76,122,227,0.12);
            >div{
                display: inline-block;
            }
            .msg_icon{
                width: 32px;
                /* height: 32px; */
                border-radius: 4px;
                /* line-height: 32px; */
                text-align: center;
            }
            .msg_label{
                margin-left: 16px;
                font-size: 13px;
                width: 100px;
                /* height: 32px;  */
                color: rgba(72, 72, 102, 0.64);
                font-weight: 500;
            }
            .msg_count{
                width: calc(100% - 148px);
                text-align: right; 
                label{
                    display: inline-block;
                    font-size: 20px;
                    text-align: right;
                    color: rgba(72, 72, 102, 0.87);
                    margin-right:16px; 
                    font-weight: 500;
                }
                i{ 
                    color: rgba(72, 72, 102, 0.32); 
                }
                .icondisabled{
                    color: rgba(76, 122, 227, 0.12);
                    cursor: not-allowed
                } 
            }
            .nocount{
                cursor: default;
                label{
                    cursor:  default;
                }
                i{
                    cursor:  not-allowed;
                }
            }
        }
    }  
}
.target_active{
    .msg_icon{ 
        // background: rgba(255, 255, 255, 0.3);
        color: #fff;
      
    }
    .msg_label{
        color:#fff !important;
    }
    .msg_count{
        label{
            color: #fff !important;
        }
        i{
            color: #fff !important;
        } 
    } 
    .icontsstyle{
        color: rgba(255, 255, 255, 0.7) !important;;
    }
} 
@targeticon1:#4C7AE3;
@targeticon2:#66C782;
@targeticon3:#F35F28;
@targeticon4:#A142CF;
@targeticon5:#24B5F7;
@targeticon6:#F9B640;

@targeticon7:#F87D7D;
@targeticon8: #F9B640;
@targeticon9: #4C7AE3;
@targeticon10: #66C681;

.target_active1{
    background:@targeticon1 !important; 
}
.target_active2{
    background:@targeticon2 !important; 
}
.target_active3{
    background:@targeticon3 !important; 
}
.target_active4{
    background:@targeticon4 !important; 
}
.target_active5{
    background:@targeticon5 !important; 
}
.target_active6{
    background:@targeticon6 !important; 
}
.target_active7 {
    background: @targeticon7  !important;
}
.target_active8 {
    background: @targeticon8  !important;
}
.target_active9 {
    background: @targeticon9  !important;
}
.target_active10 {
    background: @targeticon10  !important;
}
.target_active7,.target_active8,.target_active9,.target_active10{
    .msg_icon{
        i{
            background: #fff;
        }
    }
}

.icon1{ 
    // background: rgba(76, 122, 227, 0.16);
    color: @targeticon1;
}
.icon2{  
    // background: rgba(21, 197, 61, 0.16);
    color: @targeticon2;
}
.icon3{   
    // background: rgba(243, 95, 40, 0.16);
    color: @targeticon3;
}
.icon4{    
    // background: rgba(161, 66, 207, 0.16);
    color: @targeticon4;
}
.icon5{    
    // background: rgba(36, 181, 247, 0.16);
    color: @targeticon5;
}
.icon6{    
    // background: rgba(249, 182, 64, 0.16);;
    color: @targeticon6;
}

.icon7,.icon8,.icon9,.icon10{
    i{
        display: inline-block;
        width: 8px;
        height: 8px;
        border-radius: 50%;  
    }
}
.icon7 i{
    background: @targeticon7;
}
.icon8 i {
    background: @targeticon8;
}
.icon9 i {
    background: @targeticon9;
}
.icon10 i {
    background: @targeticon10;
}
</style>