<template>
    <el-form  :model="portconfig" label-width="0" status-icon ref="portScanFormRef"  :isUpdate = isUpdate >
        <el-form-item
            prop ='tcpscan'
            label=" " > 
            <label class="dialog_item_label_m">扫描方式</label>
            <el-radio-group v-model="portconfig.tcpScanType" :disabled="!isUpdate"> 
                <el-radio v-for="(item,i) in tcp_scan_type" :key="i" :label="item.value" value="item.value">{{item.label}}</el-radio>
            </el-radio-group>
        </el-form-item>
        <!-- <el-form-item prop="udp_scan" style="margin-bottom: 0" class="udpsan">
            <label class="dialog_item_label_m">UDP扫描</label> 
            <el-switch v-model="portconfig.isUdp" class="elSwitch" :disabled="!isUpdate">
            </el-switch> 
        </el-form-item> -->
        <!-- <el-form-item> 
            <label class="dialog_item_label_m"  >智能端口扫描</label> 
            <el-switch v-model="portconfig.intelligencePort" class="elSwitch" :disabled="!isUpdate">
            </el-switch> 
        </el-form-item> -->
        <el-form-item label="" prop="port_scan_type">
            <label class="dialog_item_label_m">端口范围</label>
            <el-select v-model="portconfig.portScanType" size="small" placeholder="请选择" class="frame_width"  @change="changePortSacn" :disabled="!isUpdate">
                <el-option v-for="(item, index) in portRange_list" :key="index" :label="item.label" :value="item.value"></el-option>
            </el-select>
        </el-form-item>
        <el-form-item label="" prop="scan_port" >
            <label class="dialog_item_label_m">扫描端口</label>
            <el-input v-model="portconfig.scanPort" size="small" :disabled="!isUpdate" class="frame_width" style=" vertical-align: text-top;" rows="6" autocomplete="off" type="textarea" resize="none" placeholder="请输入扫描端口"  ></el-input>
        </el-form-item>
        <el-form-item label="" prop="timeout">
            <label class="dialog_item_label_m">扫描超时</label>
            <el-select v-model="portconfig.timeout" size="small" placeholder="请选择" class="frame_width" :disabled="!isUpdate">
                <el-option v-for="(item, index) in timeout_list" :key="'timeout-' + index" :label="item.label" :value="item.value"></el-option>
            </el-select>
        </el-form-item>
        <el-form-item label="" prop="concurrent">
            <label class="dialog_item_label_m">扫描并发</label>
            <el-select v-model="portconfig.concurrent" size="small" placeholder="请选择" class="frame_width" :disabled="!isUpdate">
                <el-option v-for="(item, index) in concurrent_list" :key="'concurrent-' + index" :label="item.label" :value="item.value"></el-option>
            </el-select>
        </el-form-item>
    </el-form>
</template>
<style lang="less" scoped>
    .frame_width{
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
export default {
    name: "port_module",
    data(){
        return{
            portconfig:{
                portScanType: 0,
                scanPort: '',
                tcpScanType: 1,
                timeout: 10,
                concurrent: 10,
                intelligencePort: false,
                // isUdp: false,
            },
            portRange_list: [ ], // 端口范围下拉列表
            isUpdate:false,
            portScan:{},
            portRangeValue:[],
            tcp_scan_type:[],
            timeout_list: [],
            concurrent_list: [],
            update_flag:0,
            old_custom_scanPort:'', //编辑时，自定义端口
            
        }
    },
    props: {  
        isupdate: {},  
    },
    watch:{
        is_Update(curVal,oldVal){
            if(curVal){
               is_Update =  curVal
            }
        }
    },
    created() { 
        this.isUpdate = this.is_Update;   
    },
    mounted(){ 
       
    },
    methods: {
        getEnum(portScan,flag){
            this.portScan = portScan;   
            this.tcp_scan_type =  this.portScan.tcpScanType;//列表   
            this.portRangeValue = this.portScan.portRangeValue;
            this.portRange_list = this.portScan.portRange;  
            this.timeout_list = this.portScan.timeout || [];
            this.concurrent_list = this.portScan.concurrent || [];
            // this.portRangeValue.forEach(item =>{
            //     if(item.value == this.portconfig.portScanType){    
            //         this.portconfig.scanPort = item.label;   
            //     }
            // }) 
            if(flag == 1){ 
                this.tcp_scan_type.forEach(item => {
                    if(item.isDefault == true){
                        this.portconfig.tcpScanType = item.value;   
                    }
                });
                this.portRange_list.forEach(item => {
                    if(item.isDefault == true){
                        this.portconfig.portScanType = item.value;  
                        this.changePortSacn();
                    }
                });
                this.timeout_list.forEach(item => {
                    if(item.isDefault == true){
                        this.portconfig.timeout = item.value;
                    }
                });
                this.concurrent_list.forEach(item => {
                    if(item.isDefault == true){
                        this.portconfig.concurrent = item.value;
                    }
                });


            }
        },
        getIsUpdate(flag){
            this.isUpdate = flag;
        },
        getConifg(_config){
            this.update_flag = 1;
            this.portconfig = Object.assign({
                portScanType: 0,
                scanPort: '',
                tcpScanType: 1,
                timeout: 10,
                concurrent: 10,
                intelligencePort: false,
            }, _config); 
            if( _config.portScanType == 0){
                this.old_custom_scanPort = _config.scanPort;  
            } 
            this.portRangeValue.forEach(item =>{
                if(item.value == this.portconfig.portScanType){    
                    if(item.value == 0){
                        this.portconfig.scanPort = this.old_custom_scanPort;
                    }else{
                        this.portconfig.scanPort = item.label;  
                    } 
                }
            })
            if(!this.portconfig.timeout && this.timeout_list.length > 0){
                let defaultTimeout = this.timeout_list.find(item => item.isDefault === true);
                this.portconfig.timeout = defaultTimeout ? defaultTimeout.value : this.timeout_list[0].value;
            }
            if(!this.portconfig.concurrent && this.concurrent_list.length > 0){
                let defaultConcurrent = this.concurrent_list.find(item => item.isDefault === true);
                this.portconfig.concurrent = defaultConcurrent ? defaultConcurrent.value : this.concurrent_list[0].value;
            }
        },
        getAllData(){
            return this.portconfig;
        },
        changePortSacn() { //端口范围选择事件  
            this.portRangeValue.forEach(item =>{  
                if(item.value == this.portconfig.portScanType){  
                    if(item.value ==0 && this.update_flag == 1){
                        this.portconfig.scanPort = this.old_custom_scanPort; 
                    }else{
                        this.portconfig.scanPort = item.label;
                       
                    } 
                }
            })
            
        },
      
    }
}
</script>
