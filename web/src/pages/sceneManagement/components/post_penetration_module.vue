<template>
    <el-form :model="transverseForm" label-width="0" ref="transverseFormRef" >
        <div style="margin-top: 10px;">
             <el-form-item label="" prop="strategy">
                <label class="dialog_item_label_m">
                    横向策略
                    <el-tooltip :content="strategyDescription" placement="top" effect="dark">
                        <i class="el-icon-question" style="color:#909399; margin-left:2px; cursor:pointer;"></i>
                    </el-tooltip>
                </label>
                <el-select v-model="transverseForm.strategy" placeholder="请选择策略" size="small" style="width: 220px;">
                    <el-option label="同网段探测" value="same_subnet"></el-option>
                    <el-option label="邻居发现 (ARP/NetBIOS)" value="neighbor_discovery"></el-option>
                    <el-option label="排除同网段 (多网卡)" value="exclude_current_subnet"></el-option>
                    <el-option label="自定义范围" value="custom_range"></el-option>
                </el-select>
            </el-form-item>

            <el-form-item label="" prop="range" v-if="transverseForm.strategy === 'custom_range'">
                <label class="dialog_item_label_m">目标范围</label>
                 <el-input v-model="transverseForm.range" placeholder="例如: 192.168.1.0/24" size="small" style="width: 220px;"></el-input>
            </el-form-item>

            <el-form-item label="" prop="ports" v-if="transverseForm.strategy === 'custom_range'">
                <label class="dialog_item_label_m">
                    端口范围
                    <el-tooltip content="指定需要探测的端口，多个端口用逗号分隔，例如: 22,445,3389" placement="top" effect="dark">
                        <i class="el-icon-question" style="color:#909399; margin-left:2px; cursor:pointer;"></i>
                    </el-tooltip>
                </label>
                 <el-input v-model="transverseForm.ports" placeholder="例如: 22,445,3389" size="small" style="width: 220px;"></el-input>
            </el-form-item>

            <el-form-item label="" prop="targetNum">
                <label class="dialog_item_label_m">
                    目标限制
                    <el-tooltip content="限制横向移动尝试攻击的最大主机数量，建议设置在 50-100 之间以平衡覆盖率和隐蔽性。" placement="top" effect="dark">
                        <i class="el-icon-question" style="color:#909399; margin-left:2px; cursor:pointer;"></i>
                    </el-tooltip>
                </label>
                 <el-input v-model="transverseForm.targetNum" placeholder="0表示不限制" size="small" style="width: 220px;" type="number"></el-input>
            </el-form-item>

            <el-form-item label="" prop="timeout">
                <label class="dialog_item_label_m">
                    超时时间(秒)
                    <el-tooltip content="单次横向移动任务的最大执行时间，超时将自动停止。" placement="top" effect="dark">
                        <i class="el-icon-question" style="color:#909399; margin-left:2px; cursor:pointer;"></i>
                    </el-tooltip>
                </label>
                 <el-input v-model="transverseForm.timeout" placeholder="默认: 300" size="small" style="width: 220px;" type="number"></el-input>
            </el-form-item>
        </div>
    </el-form>
</template>
<style lang="less" scoped>
.frame_width {
    width: 720px;
}

.dialog_item_label_m {
    display: inline-block;
    min-width: 100px;
    font-size: 14px;
    font-weight: 500;
    color: rgba(72, 72, 102, 0.87);
    line-height: 16px;
    margin-right: 16px;
    margin-left: 10px;
}
</style>
<script>
import { number } from 'echarts';
export default {
    data() {
        return {
            transverseForm: { 
                strategy: 'same_subnet',
                range: '',
                ports: '',
                targetNum: 50,
                timeout: 600,
            },
            infiltrationHopslist: [ ],
            isUpdate:false, 
        }
    },
    props: {  
    },
    created() {
        // this.transverseForm = this.post_penetration;
        // this.isUpdate = this.is_Update;
    },
    methods: {
        getEnum(from,flag){ 
            this.infiltrationHopslist = from.jumpNum;
        },
        getIsUpdate(flag) {
            this.isUpdate = flag;
        },
        getConifg(_config) { 
            this.transverseForm = _config;
            
            if (!this.transverseForm.strategy) {
                this.$set(this.transverseForm, 'strategy', 'same_subnet');
            }

            if (!this.transverseForm.range) {
                 this.$set(this.transverseForm, 'range', '');
            }
             
            if (!this.transverseForm.ports) {
                 this.$set(this.transverseForm, 'ports', '');
            }

            if (!this.transverseForm.targetNum) {
                this.$set(this.transverseForm, 'targetNum', 50);
            }

            if (!this.transverseForm.timeout) {
                this.$set(this.transverseForm, 'timeout', 600);
            }
        },
        getAllData() {  
            return this.transverseForm;
            
        },
    },
    computed: {
        strategyDescription() {
            const map = {
                'same_subnet': '自动探测并扫描当前主机所在的整个C段网络（如 192.168.1.0/24），适合全面撒网。',
                'neighbor_discovery': '利用ARP缓存、NetBIOS广播等协议发现活跃邻居，隐蔽性高，速度快。',
                'exclude_current_subnet': '自动获取本机所有网卡信息，排除当前入侵入口所在的网段，重点探测其他网段（如内网办公段、服务器段），适合双网卡/多网卡主机跳板攻击。',
                'custom_range': '手动指定需要探测的IP范围或网段，灵活性最高。'
            };
            return map[this.transverseForm.strategy] || '';
        }
    }
}
</script>