<template>
     <el-form class="dynamicCrawlerForm" :model="webCrawler" label-width="0" status-icon ref="dynamicCrawlerFormRef" >
        <el-form-item class="formItem1" label="">
            <div class="infobox" style="display:inline-block;">  
                <el-form-item prop = 'scope' label=" " > 
                    <label class="dialog_item_label_m topline"  >爬取范围</label>
                    <el-select v-model="webCrawler.scanRange"  class="frame_width" :class="flag == 2 ? 'frame_width_m' : ''" :disabled="!isUpdate"
                        clearable placeholder="爬取范围"  size="small" ref="vulSelect">  
                        <el-option
                            v-for="(item, i) in scopelist"
                            :key="i"
                            :label="item.label"
                            :value="item.value"> 
                        </el-option>
                    </el-select>  
                </el-form-item>
                <el-form-item
                    prop = 'depth'
                    label=" "  > 
                    <label class="dialog_item_label_m topline">爬取深度</label>
                        <el-select v-model="webCrawler.maxDepth"   class="frame_width" :class="flag == 2 ? 'frame_width_m' : ''" clearable placeholder="爬取深度"  size="small" ref="vulSelect" :disabled="!isUpdate">  
                            <el-option
                                v-for="(item, i) in depthlist"
                                :key="i"
                                :label="item.label"
                                :value="item.value"> 
                            </el-option>
                        </el-select> 
                </el-form-item> 
                <el-form-item
                    prop = 'speed'
                    label=" "  > 
                    <label class="dialog_item_label_m topline">最大链接数</label>
                    <el-select v-model="webCrawler.maxUrl"   class="frame_width" :class="flag == 2 ? 'frame_width_m' : ''" clearable  size="small"  :disabled="!isUpdate">  
                        <el-option
                            v-for="(item, i) in maxUrllist"
                            :key="i"
                            :label="item.label"
                            :value="item.value"> 
                        </el-option>
                    </el-select>  
                </el-form-item> 
                <el-form-item
                    prop = 'url'
                    label=" " > 
                    <label class="dialog_item_label_m topline">URL去重</label>
                    <el-select v-model="webCrawler.scanRepeat"   class="frame_width" :class="flag == 2 ? 'frame_width_m' : ''" clearable placeholder="URL去重"  size="small" ref="vulSelect" :disabled="!isUpdate">  
                        <el-option v-for="(item,i) in  urlRepeat" :key="i" :label="item.label" :value="item.value"></el-option>
                    </el-select>  
                </el-form-item> 
                <el-form-item
                    prop = 'timeout'
                    label=" "  > 
                    <label class="dialog_item_label_m topline">单链接超时</label>
                        <el-select v-model="webCrawler.timeout"   class="frame_width" :class="flag == 2 ? 'frame_width_m' : ''" clearable placeholder="单链接超时"  size="small" ref="vulSelect" :disabled="!isUpdate">   
                            <el-option
                                v-for="(item, i) in timeoutlist"
                                :key="i"
                                :label="item.label"
                                :value="item.value"> 
                            </el-option>
                        </el-select> 
                </el-form-item>
                <el-form-item
                    prop = 'timeout'
                    label=" "  > 
                    <label class="dialog_item_label_m topline">最大爬取时长</label>
                        <el-select v-model="webCrawler.fullTimeout"   class="frame_width" :class="flag == 2 ? 'frame_width_m' : ''" clearable placeholder="单链接超时"  size="small" ref="vulSelect" :disabled="!isUpdate">   
                            <el-option
                                v-for="(item, i) in fullTimeout"
                                :key="i"
                                :label="item.label"
                                :value="item.value"> 
                            </el-option>
                        </el-select> 
                </el-form-item> 
                <el-form-item prop = 'blackkey' label=" " > 
                    <label class="dialog_item_label_m">后缀过滤</label>
                    <el-input
                            class="infoinput frame_width"  :class="flag == 2 ? 'frame_width_m' : ''" 
                            v-model="webCrawler.suffixFilter" :disabled="!isUpdate"
                            style="    vertical-align: middle;" 
                            size="small"
                            placeholder=""
                            maxlength="50"  
                        ></el-input>
                </el-form-item> 
                <el-form-item
                    prop = 'blackkey'
                    label=" " > 
                    <label class="dialog_item_label_m">关键字白名单</label>
                        <el-input
                            class="infoinput frame_width"  :class="flag == 2 ? 'frame_width_m' : ''" 
                            v-model="webCrawler.whiteList" :disabled="!isUpdate"
                            style="    vertical-align: middle;"
                            type="textarea" resize="none"
                            size="small"
                            placeholder=""
                            maxlength="50"  
                        ></el-input>
                </el-form-item>  
                <el-form-item
                    prop = 'whiteurl'
                    label=" "  > 
                    <label class="dialog_item_label_m">关键字黑名单</label>
                        <el-input
                            class="infoinput frame_width"  :disabled="!isUpdate" :class="flag == 2 ? 'frame_width_m' : ''"
                            v-model="webCrawler.blackList"
                            type="textarea"  resize="none"
                            style="    vertical-align: middle;"
                            size="small"
                            placeholder=""
                            maxlength="50"  
                        ></el-input>

                </el-form-item> 
                <el-form-item prop = ' '
                    label=" " > 
                    <label class="dialog_item_label_m">http请求头</label>
                    <!-- <el-radio-group v-model="webCrawler.headerType" class="radiobox" :disabled="!isUpdate" @change="changeHeader">
                        <el-radio :label="1">默认header</el-radio>
                        <el-radio :label="2">自定义header</el-radio>
                    </el-radio-group> -->
                    <div style="width:730px;    display: inline-block; vertical-align: top;">
                        <el-button type="primary" size="mini" style="vertical-align: middle;margin-right: 27px"  :disabled="!isUpdate"
                                    @click="addhttpheader()">新增</el-button>
                    
                        <div class="div_width"   style="margin-top:16px;margin-bottom:16px;">
                            <el-table :data="header_conf" size="small" style="width: 100%">
                                <el-table-column  prop="key" label="key">
                                    <template slot-scope="scope">
                                        <span v-if="!scope.row.dataShow">{{ scope.row.key }}</span>    
                                        <el-input v-else v-model="scope.row.key" size="small" ></el-input>    
                                    </template>
                                </el-table-column> 
                                <el-table-column  prop="value" label="value"> 
                                    <template slot-scope="scope">
                                        <span v-if="!scope.row.dataShow">{{ scope.row.value }}</span>    
                                        <el-input  :rows="1" v-else v-model="scope.row.value" size="small" resize="none"></el-input>    
                                    
                                    </template>
                                </el-table-column> 
                                <el-table-column  label="操作"  width="150">
                                    <template slot-scope="scope"  > 
                                            <el-link :underline="false" @click="headerDelete(scope)" > 删除 </el-link>
                                    </template>
                                </el-table-column>
                            </el-table>
                        </div> 
                    </div>
                </el-form-item>   
            
            </div>   
        </el-form-item>
    </el-form>
</template>
<style lang="less" scoped>
    .frame_width{
        width: 720px;
    }
    .frame_width_m{
        width: 400px;
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
    /deep/ .el-form-item {
        margin-bottom: 8px;
    }
    /deep/ .el-table::before{
        height: 1 !important;
    }
    /deep/ .el-table td.el-table__cell, 
    /deep/.el-table th.el-table__cell.is-leaf{
        border-bottom: 1px solid #ebeef5  ;
    }
    /deep/ .el-form-item__content{
        line-height: 30px;
    }
</style>
<script>
export default {
    name:'',
    data(){
        return{
             webCrawler: {
                "maxDepth":0,
                "maxUrl":0,
                "scanRange":0, //爬取范围
                "timeout":0,
                "fullTimeout":0,
                "scanRepeat":0,
                "blackList":"",
                suffixFilter:'',
                "whiteList":"",
                "headers":[]   
            },
            defalulHeader: "{'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; WOW64; rv:68.0) Gecko/20100101 Firefox/68.0', 'Accept': '*/*', 'Accept-Language': 'zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2', 'Accept-Encoding': 'gzip, deflate', 'Connection': 'close'}",
            timeoutlist: [ ], //单链接超时 
            depthlist: [ ], //爬取深度
            scopelist: [ ], //爬取范围
            maxUrllist:[], //最大连接数
            fullTimeout:[], //最大爬取时长
            urlRepeat:[],
            isUpdate:false,
            crawler:{},
            header_conf:[],
        }
    },
    props: {
        // web_crawler: {}, 
        flag: {}, 
    },
    created() {
        // this.webCrawler = this.web_crawler;
    },
    methods: {
        getEnum(crawler,flag){ 
            this.crawler = crawler;
            this.scopelist = this.crawler.scanRange; //爬取范围
            this.depthlist = this.crawler.maxDeep;  //爬取深度
            this.timeoutlist = this.crawler.timeout; //单链接超时 
            this.urlRepeat = this.crawler.scanRepeat;
            this.maxUrllist = this.crawler.maxUrl; //最大连接数
            this.fullTimeout = this.crawler.fullTimeout;
            if(flag == 1){
                this.scopelist.forEach(item =>{
                    if( item.isDefault==true){
                        this.webCrawler.scanRange = item.value
                    }
                })
                this.depthlist.forEach(item =>{
                    if( item.isDefault==true){
                        this.webCrawler.maxDepth = item.value
                    }
                })
                this.timeoutlist.forEach(item =>{
                    if(item.isDefault==true){
                        this.webCrawler.timeout = item.value
                    }
                })
                this.urlRepeat.forEach(item =>{
                    if(item.isDefault==true){
                        this.webCrawler.scanRepeat = item.value
                    }
                })
                this.maxUrllist.forEach(item =>{
                    if(item.isDefault==true){
                        this.webCrawler.maxUrl = item.value
                    }
                })
                this.fullTimeout.forEach(item =>{
                    if(item.isDefault==true){
                        this.webCrawler.fullTimeout = item.value
                    }
                })

                this.webCrawler.blackList = this.crawler.blackList;
                this.webCrawler.whiteList = this.crawler.whiteList
            }

        },
        getIsUpdate(isUpdate) {
            this.isUpdate = isUpdate;
        },
        getConifg(_config) { 
            this.webCrawler = _config; 
            this.header_conf =  this.webCrawler.headers;
        },
        getAllData() {
           
            this.webCrawler.headers = this.header_conf;
            return this.webCrawler;
        },
        // changeHeader(){ //切换是否自定义header 
        //     if(this.webCrawler.header_type == 2){ //自定义
        //         this.webCrawler.crawler_header = '';
        //     }
        // },
        addhttpheader(){
            this.header_conf.push({
                dataShow: true,
            })
        },
        headerDelete(scope){
            this.header_conf.splice(scope.$index, 1);
        },
        
    },
}
</script>