<template>
    <div>
        <el-dialog
            title="资产导入"
            :visible.sync="dialogVisible"
            width="40%"
            :close-on-click-modal="false"
            :show-close="false"
            >
            <div class="dialog_b_btn"> 
                <el-button size="small" @click="handleCancel">取消</el-button>
            </div>
            <div style="padding: 24px">
                <div style="margin-bottom:24px">
                    <el-button type="primary"   size="small" @click="downLoadTemp()" >下载模板</el-button>
                </div>
                
                <div>
                    <input type="file" class="btnUploadID" ref="upload"
                        @change="changeuploaID($event)"
                        style="display: none" id="input-file-ID" />
                    <el-button type="primary"   size="small" @click="clickupload()">资产上传</el-button>
                </div> 
            </div>
        </el-dialog>
    </div>
</template>
<style lang="less" scoped>

</style>
<script>
import jsFileDownload from 'js-file-download'
import { traffic } from '@/api/assetManagement.js'
export default {
    name:'',
    props:{
        value: {}, 
    },
    data(){
        return{
            dialogVisible:false,
            uploadfileName:'',
            uploadfile:null,
            addDevice:{}
        }
    },
    watch: {
        value(newVal, oldVal) {
            this.dialogVisible = newVal;  
            
        }
    },
    methods:{
        clickupload(){
            document.querySelector(".btnUploadID").click();
        },
        changeuploaID(e){
            let deviceFile = e.target.files;
            for(let i=0;i<deviceFile.length;i++){
                this.uploadfileName = deviceFile[i].name;
                this.uploadfile =  deviceFile[i];
                e.target.value='';
            } 
            console.log(this.uploadfileName )
            const formDate = new FormData();  
            formDate.append('file',this.uploadfile);
            this.$ajax({
                url:"/smart/asset/import",
                method:"POST",
                data: formDate
            })
            .then(res =>{
                if(res.data.code == 200){
                    this.$message.success('导入资产成功');
                    this.$emit("importCancel");
                }
            
            })
            .catch(e=>{})
        },
        async downLoadTemp(){ //下载模板  
            const res = await traffic.downloadasset({
                getTempToken:true
            })
            if(res.code ==200){ 
                let url = '/smart/asset/templatedownload?temp_token='+res.data.tempToken;
                if (process.env.NODE_ENV === 'development') {
                    url = '/api'+url
                } else {
                    url = ''+url
                }
                window.open(url, '_blank')
            }
           
        },
        handleCancel(){
            this.dialogVisible = false; 
            this.$emit("importCancel");
        },
    }
}
</script>