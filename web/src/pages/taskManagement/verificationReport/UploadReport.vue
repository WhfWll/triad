<template>
    <div class="outContainer">
        <div class="main-title  ">   
            <router-link :underline="false" class="classA" :to="{ path: '/verificationReport' }" >验证报告</router-link>  
            <label class="currentpagetitle"> 
                 <!-- <el-tooltip class="item" > -->
                    <span> 上传报告</span>
                <!-- </el-tooltip> -->
            </label>
	  	</div> 
		<div class="taskinfolist  tastBtnCont">    
              <div class="search-box" >
				<div class="operationbutton"  >
					<el-button type="primary"  size="small" @click="startVerification" >开始验证</el-button>  <!--:disabled="startDisabled"-->
				</div> 
                <div class="tipsCont operationbox">
                    <label class="topTips">提示：</label>
                    <span >点击导入报告添加要验证的报告，支持报告格式为html/zip，支持厂商包括绿盟、天镜。</span>
					<span style="color:#F35F28">本页面刷新后已导入的报告将被清空！</span>
                </div>
			</div>
        </div>
        <div class="tasklist context_box_bg mt16"> 
            <div class="search-box" >

            <!-- 报告验证 -->
            <el-form :model="importReoprtForm" :rules="reportRule" ref="reportsForm" label-width="100px">
				<div>
					<label class="dialog_item_label" style=" vertical-align: top; margin-top: 10px; width: 103px;">报告名称</label>
					<el-form-item label=" " prop="name"  label-width="10px" style="display: inline-block;margin-right:0">
						<el-input type="text" v-model="importReoprtForm.name" autocomplete="off" placeholder="名称仅限大小写字母、数字、下划线且小于100字 "></el-input>
					</el-form-item>
				</div>
				<div>
					<label class="dialog_item_label" style=" vertical-align: top; margin-top: 10px; width: 103px;">报告厂商</label>
					<el-form-item label="" prop="manufacturer"  label-width="10px" style="display: inline-block;margin-right:0">
						<el-select v-model="importReoprtForm.manufacturer" placeholder="请选择">
							<el-option
								v-for="item in manufacturerOptions"
								:key="item.value"
								:label="item.label"
								:value="item.value">
							</el-option>
						</el-select>
					</el-form-item>
				</div>
				<div>
					<label class="dialog_item_label" style=" vertical-align: top; margin-top: 10px; width: 103px;">执行方式</label>
					<el-form-item label=" " prop="type" label-width="10px" style="display: inline-block;margin-right:0">
						<el-select v-model="importReoprtForm.type" placeholder="请选择">
							<el-option
								v-for="item in typeOptions"
								:key="item.value"
								:label="item.label"
								:value="item.value">
							</el-option>
						</el-select>
						<div class="margin_left_width" v-if="importReoprtForm.type == 2">
							<label class="execute_type_label">计划时间</label>
							<el-date-picker
								:clearable="false"
								v-model="importReoprtForm.time"
								type="datetime"
								size="small"
								format="yyyy-MM-dd HH:mm:ss"
								value-format="yyyy-MM-dd HH:mm:ss"
								placeholder="选择日期时间">
							</el-date-picker>
						</div>
					</el-form-item>
				</div>
                <div class="reportElItem">
					<label class="dialog_item_label" style=" vertical-align: top; margin-top: 10px; width: 103px;">扫描报告</label>
					<el-form-item label=" " prop="reports" label-width="10px" style="display: inline-block;margin-right:0">
						<input type="file"  class="fileInput" id='fileUpload' @change="changeFiles($event)" multiple="multiple" style="display:none" accept="text/html,application/zip"/>
						<el-button size="small" type="primary" @click="importFileClick">导入报告</el-button>
						<!-- <el-upload
							class="upload-demo"
							accept="text/html,application/zip"
							ref="reportRef"
							action=""
							:on-remove="handleRemove"
							:before-remove="beforeRemove"
							:before-upload="beforeUpload"
							:on-change="fileChange"
							:http-request="httpRequest"
							multiple
							:limit="20"
							:on-exceed="handleExceed"
							:file-list="fileList">
							<el-button size="small" type="primary">导入报告</el-button>
						</el-upload> -->
						<div v-for="(item,index) in fileStatusData" :key="index" :class="['progressCont',{'default':item.success},{'errorProgre':item.success == false}]">
							<div class="fileCont">
								<i class="iconfont iconfujian"></i>
								<span class="fileNameText" >{{ item.name }}</span>
								<!-- <i class="iconfont iconquxiao" @click="deleteFile(item.name,index)"></i> -->
							</div>
							<!-- <el-progress :stroke-width="2" :percentage="item.percentage" :format="progressText"></el-progress> -->
							<span class="uploadStatus successColor"  v-if="item.success">上传成功</span>
							<span class="uploadStatus errorColor"  e-else>{{ item.error }}</span>
						</div>
					</el-form-item>
				</div>
				
			</el-form>
            </div>
        </div>
    </div>
</template>
<style scoped lang="less">
.context_box_bg {
	padding:24px;
	background: #FFF;
}
.tastBtnCont{
	font-size:13px;
}
.tastBtnCont .search-box{
	padding: 24px;
	min-height: 0;
	background: #FFF;
	border-radius: 4px;
	box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
.topTips{
	color:#4C7AE3;
	width:40px;
}
.tipsCont{
	float:left;
	/* margin-left:24px; */
	line-height:31px;
}
.operationbox{
	label{
		margin-left: 24px;
		color: #4C7AE3;
		font-size: 13px;
	}
	span{
		margin-left: 4px;
		font-size: 13px;
		color: rgba(72, 72, 102, 0.64);
	}
}
.mt16{
	margin-top:16px;
	position:relative;
}
.execute_type_label {
    color: rgba(72, 72, 102, 0.64);
    font-size: 13px;
    margin-right: 16px;
}
/deep/.el-form-item__content{
	width:360px;
}
/deep/.el-select{
	width:100%;
}
/deep/.el-form-item__label{
	padding:0;
}
/deep/.el-progress-bar__outer{
	width:255px;
}
.reportElItem /deep/.el-form-item__content{
	line-height:0!important;
}
// .reportElItem /deep/.el-form-item__content{
// 	// width:500px;
// }
.reportElItem /deep/.el-form-item__content .uploadStatus {
	width:250px;
	right:-240px
}
/deep/.el-progress{
	margin-left:20px;
}
.upload-demo{
	margin-bottom:8px;
}
.progressCont{
	padding-top:16px;
	line-height:18px!important;
	position:relative;
	.uploadStatus{
		position: absolute;
		top:26px;
		right:-38px;
	}
}
.successColor{
	color:rgba(76, 122, 227, 1);
}
.errorColor{
	color:rgba(243, 95, 40, 1)!important;
}
.fileCont{
	height:18px;
	color:rgba(72, 72, 102, 0.64);
	font-size:13px;
	// overflow:hidden;
	.fileNameText{
		padding-left:6px;
		display:block;
		width:220px;	
		height:18px;
		line-height:18px;
		overflow:hidden;
		text-overflow:ellipsis;
		white-space:nowrap;
		float:left;
		color:rgba(76, 122, 227, 1);
		// color:rgba(243, 95, 40, 1);
	}
	i{
		display:block;
		float:left;
	}
	.iconquxiao{
		margin-left:16px;
		cursor:pointer;
		// z-index: 99;
	}
}
/deep/.el-date-editor.el-input, 
/deep/.el-date-editor.el-input__inner{
	width:185px!important;
}
/deep/.el-progress__text{
	display:none;
}
/deep/.el-upload-list{
	display:none;
}
.default /deep/.el-progress-bar__inner{
	background-color: rgba(76, 122, 227, 1);//rgba(243, 95, 40, 1);
}
.default /deep/span{
	color:rgba(76, 122, 227, 1)!important;
}
.errorProgre /deep/.el-progress-bar__inner{
	background-color: rgba(243, 95, 40, 1);
}
.errorProgre /deep/span{
	color: rgba(243, 95, 40, 1)!important;
}
</style>
<script>
import { verReport } from '@/api/verificationReport.js'
export default({
    name:'uploadReport',
    data(){
        var reportName = (rule, value, callback) =>{
            if (!value) {
              return callback(new Error('任务名称不能为空'));
            }else{
                const pattern = /^[a-zA-Z\d\_\u2E80-\u9FFF]{0,100}$/
				//const pattern = /^[a-zA-Z\d\_]{0,100}$/
                if(pattern.test(value)){
                    callback();
                }else{
                    return callback(new Error('任务名称格式有误'));
                }
            }
        }

		var reportFile = (rule, value, callback) =>{
            if (!value) {
              return callback(new Error('请上传扫描报告'));
            }else{
				 callback();
                // const pattern = /^[a-zA-Z\d\_\u2E80-\u9FFF]{0,8}$/
                // if(pattern.test(value)){
                //     callback();
                // }else{
                //     return callback(new Error('任务名称格式有误'));
                // }
            }
        }

		let showDate =new Date();
		let seperator ='-';
		let year = showDate.getFullYear();
		let month = showDate.getMonth() + 1;
		let day =showDate.getDate();
		let hour =showDate.getHours();
		let mins  = showDate.getMinutes();
		let sec = showDate.getSeconds();
		var strDate = showDate.getDate();
		if (month >= 1 && month <= 9) {
			month = "0" + month;
		}
		if (strDate >= 0 && strDate <= 9) {
			strDate = "0" + strDate;
		}
		let currentdate = year + seperator + month + seperator + strDate +' '+ hour +':'+mins+':'+sec;
        return {
            taskRiskList:[],
            taskstatus:[],
			reportId:this.$route.query.id,
            typeOptions:[
                // {
                //     value: '1',
                //     label: '立即执行'
                // },
                // {
                //     value: '2',
                //     label: '定时执行'
                // },
            ],
			manufacturerOptions:[
				// {
                //     value: '1',
                //     label: '启明天镜'
                // },
                // {
                //     value: '2',
                //     label: '绿盟nsfocus'
                // },
			],
            importReoprtForm:{
				name:'',
                type:'',
				reports:[],
				time:currentdate,
				manufacturer:''
			},
			fileList:[
				// {percentage:10,status:false,name:'fsjk'},
				// {percentage:80,status:true,name:'oopso'}
			],
			reportRule:{
				name:[
					{ required: true, trigger: 'blur',validator:reportName  },
				],
				manufacturer:[
					{required: true, message: '请选择厂商', trigger: 'blur'}
				],
                type:[
                    {required: true, message: '请选择执行方式', trigger: 'blur'}
                ],
				reports:[
					{required: true,  message: '请上传扫描报告',trigger: 'change',}
				]
			},
			fileStatusData:[
				// {filename:'"不是正确格式的启明天镜报告"',error:'lds',success:flase}
			],
			startDisabled:true,
			repeatRequest:'',
			repeatNames:[],
			newUploadFile:[],
        }
    },
	mounted(){
		this.getEnum();
	},
	created:function(){ 
        this.$store.state.activefirstMenu="/verificationReport";   
        
    },
    methods:{
		getEnum() { 
            this.$ajax.get('/smart/reportverify/enum').
            then((res) =>{
                var res = res.data;    
                if(res.code == 200){ 
                    this.typeOptions = res.data.executeType; 
                    this.manufacturerOptions = res.data.producerType;
                }
				else{
                    this.$message({
                        message:res.msg,
                        type: 'error'
                    });
                }
            }).catch((error) =>{
                console.log(error);
            });
            
        },
        btReportTesting(){
            this.$router.push({
                path:`uploadReport`,
                // query:{}
            })
        },
		importFileClick(){
			document.querySelector('.fileInput').click();
		},
		changeFiles(e){ 
			let _files = e.target.files;
			let _arr = [];
			this.newUploadFile = []; 
			this.fileStatusData=[];
			for(var i =0;i< _files.length;i++) {
				let item = _files[i];
				const fileType = item.name.substring(item.name.lastIndexOf(".") + 1);
				let isFileType = false ;
				if(fileType == 'html' || fileType == 'zip'){
					isFileType = true;
				}
				const isLt100M = item.size / 1024 / 1024 < 100;

				if (!isFileType) {
					this.$message.error('上传文件只能是html/zip格式!');
					return false;
				}
				if (!isLt100M) {
					this.$message.error('上传文件大小不能超过 100MB!');
					return false;
				} 
				_arr.push(item)
			} 
			_arr.forEach((i,n) =>{
				if(this.repeatNames.indexOf(i.name) == -1){
					let _len = this.fileList.length;

					if(_len>=20){
						this.$message.warning(`上传文件个数不能超过20`);//(`当前限制选择 20 个文件，本次选择了`+_len+`个文件`);
						// this.fileList.length = 0;
						// this.fileStatusData.length = 0;
						// this.importReoprtForm.reports.length = 0;
						var file = document.getElementById('fileUpload');
						file.value = ''; //file的value值只能设置为空字符串
						return false;
					}else{
						this.repeatNames.push(i.name);
						this.fileList.push(i);
						// i.success = true;
						this.fileStatusData.push(i); 

						this.newUploadFile.push(i)

						console.log(this.fileStatusData)
					}
				}
			})
			this.importReoprtForm.reports = this.fileStatusData;
 

			// let _formData = new FormData();
			// newUploadFile.forEach(function (file) {// 因为要上传多个文件，所以需要遍历
			// 	_formData.append('file', file); 
			// })
			// _formData.append('task_id',this.reportId);
			// _formData.append('product',this.importReoprtForm.manufacturer);
			// let config = {
			// 	headers: {
			// 		'Content-Type': 'multipart/form-data'
			// 	}
			// };
			// this.$ajax.post(
			// 	'/v2/task/report_test/upload_file',
			// 	_formData,
			// 	 config
			// ).
			// then((res) =>{
			// 	let dt = res.data;
			// 	if(dt.code ==200){
			// 		this.$message({
			// 			message: '上传成功',
			// 			type: 'success'
			// 		});
			// 		// this.requestData()
			// 		// if(dt.status == 'running'){
			// 		// 	repeatRequest = setInterval(this.requestData, 2000);
			// 		// }else{
			// 		// 	clearInterval(repeatRequest);
			// 		// }
					
			// 	}else{
					
			// 		this.$message({
			// 			message: dt.msg,
			// 			type: 'error'
			// 		});
			// 	}
			// }).catch((error) =>{
			// 	console.log(error)
			// })
		},
		requestData(){
			let taskId = this.reportId;
			this.$ajax.get('/reportverify/v1/files/status/',{
                params:{
                    task_id:taskId,
                }
            }).
			then((res) =>{
				let dt = res.data;
				let filesArr = dt.results;
				if(dt.success){
					this.fileStatusData.forEach((item,index) =>{
						filesArr.forEach((_item,_index) =>{
							if(item.name == _item.filename){
								item.error = _item.error;
								item.success = _item.success;
								if(_item.success == true){
									item.percentage = 100;
								}else{
									item.percentage = 80;
								}
							}
						})
					})
					this.$forceUpdate();
					if(dt.status == 'running'){
						this.repeatRequest = setInterval(this.requestData, 2000);
					}else{
						clearInterval(this.repeatRequest);
					}
					var file = document.getElementById('fileUpload');
					file.value = ''; //file的value值只能设置为空字符串
				}
			}).catch((error) =>{
				console.log(error);
			})
		},
		deleteFile(name,index){//删除
			this.$ajax({
				url:'/reportverify/v1/files/?filename='+name,
				method:'delete',
				data:{
					filename:name,
					task_id:this.reportId
				}
			}).then((res) =>{
				let dt = res.data;
				if(dt.success){
					this.fileStatusData.splice(index,1);
					this.repeatNames.splice(index,1);
					this.fileList.splice(index,1);
					// console.log(this.fileStatusData)
					this.$forceUpdate();
					var file = document.getElementById('fileUpload');
					file.value = null; //file的value值只能设置为空字符串
				}
			}).catch((error) =>{
				console.log(error)
			})
		},
        handleExceed(files, fileList) {
        	this.$message.warning(`当前限制选择 20 个文件，本次选择了 ${files.length} 个文件，共选择了 ${files.length + fileList.length} 个文件`);
        },
		// fileChange(file,files){
		// 	const fileType = file.name.substring(file.name.lastIndexOf(".") + 1);
		// 	let isFileType = false ;
		// 	if(fileType == 'html' || fileType == 'zip'){
		// 		isFileType = true;
		// 	}
		// 	const isLt100M = file.size / 1024 / 1024 < 100;

		// 	if (!isFileType) {
		// 		this.$message.error('上传文件只能是html/zip格式!');
		// 		let uid = file.uid // 关键作用代码，去除文件列表失败文件
		// 		let idx = this.$refs.reportRef.uploadFiles.findIndex(item => item.uid === uid) // 关键作用代码，去除文件列表失败文件（uploadFiles为el-upload中的ref值）
		// 		this.$refs.reportRef.uploadFiles.splice(idx, 1) // 关键作用代码，去除文件列表失败文件
		// 		return false;
		// 	}
		// 	if (!isLt100M) {
		// 		this.$message.error('上传文件大小不能超过 100MB!');
		// 		let uid = file.uid // 关键作用代码，去除文件列表失败文件
		// 		let idx = this.$refs.reportRef.uploadFiles.findIndex(item => item.uid === uid) // 关键作用代码，去除文件列表失败文件（uploadFiles为el-upload中的ref值）
		// 		this.$refs.reportRef.uploadFiles.splice(idx, 1) // 关键作用代码，去除文件列表失败文件
		// 		return false;
		// 	}

		// 	this.$refs.reportRef.uploadFiles = this.functionUnique(files).slice();//上传文件去重
		// 	this.fileList = this.$refs.reportRef.uploadFiles;
		// },
		// functionUnique(arr){
		// 	var result = [];
		// 	var obj = {};
		// 	for(var i =0; i<arr.length; i++){
		// 		if(!obj[arr[i].name]){
		// 			result.push(arr[i]);
		// 			obj[arr[i].name] = true;
		// 		}
		// 	}
		// 	return result;
		// },
		// httpRequest(param){
		// 	this.fileData.push( param.file);  // append增加数据
		// },
		startVerification(){//开始验证
			this.$refs.reportsForm.validate(async (valid) =>{
				if(valid){
					const formDate = new FormData();
					formDate.append('name',this.importReoprtForm.name);
					formDate.append('executeType',this.importReoprtForm.type); 
					formDate.append('producer',this.importReoprtForm.manufacturer);  
					this.newUploadFile.forEach(function (file) { 
						formDate.append('file', file); 
					})
					let config = {
						headers: {
							'Content-Type': 'multipart/form-data'
						}
					}; 
					this.$ajax.post(
						'/smart/reportverify/upload',
						formDate,
						config
					).then((res) =>{
						let dt = res.data;
						if(dt.code == 200){
							this.$message.success('验证报告成功');
							this.$router.push({
								path:`/verificationReport`,
							})
						}else{
							this.$message.error(dt.msg); 
						}
					}).catch((error) =>{
						console.log(error);
						
					})
				}
			})
		},
		// handlecancelReport(){
		// 	this.$refs.reportsForm.resetFields();
		// 	this.$refs.reportRef.clearFiles();
		// },
        // handlesearch(){},
        // btnTaskinfo(){
        //     this.$router.push({
        //         path:`/reportTaskInfo`,
        //         query: { 
		// 			// id: row.id, 
        //         }
        //     })
        // },
		progressText(percentage) {
			return percentage === 100 ? '' : `${percentage}%`;
		}
    }
})
</script>