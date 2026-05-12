<template>
  <div class="mainbox">
    <div class="main-title">
      <router-link
        :underline="false"
        class="classA"
        :to="{ path: '/usergroup' }"
        >用户组管理</router-link
      >
      <label class="currentpagetitle">
        <span>用户管理</span>
      </label>
    </div>
    <div class="operationbox whitebg">
      <xzbutton type="primary" @click="addExperienceToSet" size="small"
        >保存至用户组</xzbutton
      >
      <label for="">提示：</label>
      <span>点击“保存至用户组”按钮，则添加用户成功，且退出用户管理页面。</span>
    </div>
    <div class="addexperiencebox clearfix">
      <div class="content-r whitebg">
        <div class="search-box">
          <div class="serach-condition">
            <div class="search-text">
              <el-input
                placeholder="请输入关键字"
                @keydown.enter.native="handlesearch"
                v-model="systemsearch_field"
                class="input-with-select"
                size="small"
                clearable
              >
              </el-input>
              <xzbutton type="primary" @click="handlesearch" size="small"
                >搜索</xzbutton
              >
            </div>
            <div>
              <xzbutton type="primary" @click="handleReset" size="small"
                >重置</xzbutton
              >
            </div>
          </div>
          <div style="float: left">
            <xzbutton
              type="primary"
              @click="addToRight"
              size="small"
              style="margin-right: 8px"
              >添加至组</xzbutton
            >
            <el-button
              type="warning"
              size="small"
              :width="170"
              :disabled="!SelNoToRight.length"
              @click="clearSelNoToRight"
              >清除</el-button
            >
          </div>
        </div>

        <div class="tablebox">
          <el-table
            v-loading="Loading"
            class="myTable"
            ref="multipleTable"
            :data="systemtableData"
            tooltip-effect="dark"
            style="width: 100%"
            @selection-change="handlesystemSelectionChange"
            @select="handleSelect"
            @select-all="handleSelectAll"
          >
            <el-table-column type="selection" :selectable="checkboxT">
            </el-table-column>
            <el-table-column
              prop="name"
              label="用户名"
              :show-overflow-tooltip="true"
            >
            </el-table-column>
            <el-table-column prop="role_str" label="用户角色">
            </el-table-column>
            <el-table-column prop="group_list" label="已归属组">
            </el-table-column>
            <el-table-column prop="status_str" label="状态"> </el-table-column>
          </el-table>
        </div>
        <div class="box-footer">
          <el-pagination
            :page-size="pageSize"
            background
            :pager-count="5"
            layout="total,  prev, pager, next,sizes,  jumper"
            :total="totalpage"
            :current-page="currentpage"
            @current-change="currentchange"
            @size-change="handleSizeChange"
          >
          </el-pagination>
        </div>
      </div>
      <div class="menu-l whitebg" style="padding: 24px">
        <el-button
          type="warning"
          size="small"
          style="margin-bottom: 16px"
          :disabled="!selectmultipleSelection.length"
          @click="removeSelect"
          >移除</el-button
        >
        <!-- <delbutton 
                size="small"
                style="margin-bottom:16px"
                @click="removeSelect"  
                :disabled="!selectmultipleSelection.length">移除22</delbutton>   -->
        <div style="overflow-y: auto; height: calc(100% - 40px)">
          <!-- height="600" -->
          <!-- <el-table
            :data="rightData"
            class="infinite-list"
            tooltip-effect="dark"
            ref="selectmultipleTable"
            v-infinite-scroll="load"
            :infinite-scroll-distance="15"
             :infinite-scroll-delay="600"
            @selection-change="SelectedSelectionChange"
             style="width: 100%;height:100%"
          >
            <el-table-column type="selection" width="55"> </el-table-column>
            <el-table-column
              prop="name"
              label="全选"
              :show-overflow-tooltip="true"
            >
            </el-table-column>
          </el-table> -->
          <el-table :data="rightData" class="rightTable" tooltip-effect="dark" ref="selectmultipleTable"
                        v-el-table-infinite-scroll="load" :infinite-scroll-distance="15" :infinite-scroll-delay="600"
                        @selection-change="SelectedSelectionChange" style="width: 100%;height:100%">
                        <el-table-column type="selection" width="55">
                        </el-table-column>
                        <el-table-column prop="name" label="全选" :show-overflow-tooltip="true">
                        </el-table-column>

                    </el-table>
        </div>
      </div>
    </div>
  </div>
</template>
<style scoped lang="less">
.search-box .search-text .input-with-select {
  vertical-align: inherit;
}

.search-box {
  margin-bottom: 0;
  padding-bottom: 0;
}

.operationbox {
  padding: 24px;
  margin-bottom: 15px;

  label {
    margin-left: 24px;
    color: #4c7ae3;
    font-size: 13px;
  }

  span {
    margin-left: 4px;
    font-size: 13px;
    color: rgba(72, 72, 102, 0.64);
  }
}

.serach-condition > div {
  margin-bottom: 16px;
}

.dialog_b_btn {
  position: absolute;
  top: 15px;
  right: 24px;
  font-size: 14px;

  button {
    color: #4c7ae3;
  }
}

@media (max-width: 1440px) {
  /deep/ .el-dialog {
    height: calc(100% - 96px);
  }
}

@media (min-width: 1440px) {
  /deep/ .el-dialog {
    height: calc(100% - 176px);
  }
}

.dialog_item_label {
  font-size: 14px;
  border-left: 3px solid #4c7ae3;
  padding-left: 8px;
  font-weight: 500;
  width: 113px;
  display: inline-block;
  line-height: 16px;
  box-sizing: border-box;
}

.search-filed > div {
  margin-right: 10px;
}

.mainbox {
  height: 100%;
  box-sizing: border-box;
}

.addexperiencebox {
  height: calc(100% - 130px);
  box-sizing: border-box;
  margin-bottom: 30px;
}

.whitebg {
  background: #fff;
  border-radius: 4px;
  box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);

  /deep/ .rightTable {
    .el-table__body-wrapper {
      height: calc(100% - 56px) !important;
      /* height: auto!important; */
    }
  }

  /deep/ .is-checked.is-disabled {
    .el-checkbox__inner {
      background-color: #f2f6fc !important;
      border-color: #d1d4dc !important;
    }
  }
}

.boxtop {
  padding: 24px 24px 4px;
}

.addexperiencemodelbox {
  overflow: hidden;
  min-height: calc(100% - 54px);
  box-sizing: border-box;
}

.addexperiencefoot {
  height: 80px;
  width: 100%;
  border-top: 1px solid #e4e7ed;
  padding: 20px;
  box-sizing: border-box;
}

.menu-l {
  float: right;
  width: 300px;
  height: 100%;
  box-sizing: border-box;
  margin-left: -1px;
  background: #fff;
  box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
  border-radius: 4px;
}

.menu-l > strong {
  display: block;
  padding: 20px;
  border-bottom: 1px solid #e4e7ed;
  font-size: 14px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.menu-l > ul {
  margin: 0;
  padding: 0;
  list-style: none;
  padding: 20px;
  font-size: 13px;
  color: #494b4f;
  height: calc(100% - 60px);
  overflow-y: auto;
  box-sizing: border-box;
}

.menu-l > ul > li {
  margin: 8px 0;
}

.menu-l > ul > li span {
  display: inline-block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  width: 250px;
  cursor: pointer;
}

.menu-l > ul > li i {
  cursor: pointer;
  float: right;
  margin-top: 3px;
}

.content-r {
  float: left;
  width: calc(100% - 315px);
  height: 100%;
  padding: 24px;
  box-sizing: border-box;
  background: #fff;
  box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
  border-radius: 4px;
}

.addexperiencebox .el-tabs {
  height: 100%;
}

.addexperiencebox .el-tabs__content {
  height: calc(100% - 55px);
}

.addexperiencebox .el-tabs__content > div {
  height: 100%;
}

.addexperiencebox .el-tabs__item {
  font-size: 16px;
  line-height: 22px;
}

.addexperiencebox .tablebox {
  height: calc(100% - 115px);
  overflow-y: auto;
}

.selectlable {
  margin-left: 10px;
  font-size: 12px;
}

.selectlable span {
  color: #4c7ae3;
}

.introduce {
  background-color: #edf2fc;
  padding: 10px 20px;
  font-size: 12px;
  color: #909399;
}

.introduce > div {
  margin: 10px 0;
}

.txtareacontent textarea {
  resize: none;
}

.formaddbox {
  padding: 10px 0;
  height: calc(100% - 174px);
  overflow-y: auto;
}

/* /deep/ .el-checkbox__input.is-checked .el-checkbox__inner, .el-checkbox__input.is-indeterminate .el-checkbox__inner[data-v-fae5bece] {
        background-color: #edf2fc !important;
        border-color: #c8ccd5 !important;
} */
/deep/ .el-button + .el-button {
  margin-left: 0;
}

/deep/ .myTable {
  thead {
    .cursorPointer {
      cursor: pointer;

      &.active {
        color: #4c7ae3;

        i {
          color: #4c7ae3;
        }
      }
    }

    .cell {
      line-height: 15px;

      > span {
        position: absolute;
      }
    }

    .iconfont {
      color: rgba(72, 72, 102, 0.32);
      margin-left: 5px;
    }

    .el-select {
      height: 0;
      visibility: hidden;

      .el-input,
      .el-input__inner {
        height: 0 !important;
      }
    }
  }
}
</style>
<script>
import { encryptCBC, decryptCBC } from '@/commonFunction/des.js'
import elTableInfiniteScroll from 'el-table-infinite-scroll';
import xzbutton from "@/components/XzButton.vue"
import delbutton from "@/components/DelButton.vue"
import API from '@/api/user.js'
export default ({
  name: 'groupmanagement',
  components: {
    xzbutton,
    delbutton,
  },
  directives: {
      'el-table-infinite-scroll': elTableInfiniteScroll
  },
  data () {
    var valiFile = (rule, value, callback) => {
      if (this.addform.hasOwnProperty(rule.field)) {
        if (!this.addform[rule.field]) {
          callback(new Error('请上传脚本文件'))
        } else {
          callback()
        }
      }

    }
    return {
      id: this.$route.params.id,
      name: this.$route.query.name,
      fpage_num: this.$route.query.page_num,
      experiencelist: [],
      experiencelistold: [],
      activeName: 'first',
      systemsource: '',
      systemsearch_field: '',
      systemSelectedCount: 0,
      systemtableData: [],
      status: [
        {
          id: 1,
          name: '系统提供'
        },
        {
          id: 2,
          name: '个人添加'
        }
      ],
      formdt: {
        name: '',
        content: '',
        remark: '',
      },
      addform: {
        name: '',
        type: '',
        file: null,
        filename: '',
        rank: '',
        object: '',
        desc: '',
        user: '',
        status: '',
      },
      rules: {
        name: [
          { required: true, message: '工具名称不能为空', trigger: 'blur' },
        ],
        type: [
          { required: true, message: '请选择工具类型', trigger: 'change' },
        ],
        rank: [
          { required: true, message: '请选择测试强度', trigger: 'change' },
        ],
        file: [
          { required: true, validator: valiFile, trigger: "change" },
        ]
      },
      multipleSelection: [],
      Loading: false,
      downfile: '', //下载文件
      page_num: 1,
      totalpage: 0,
      currentpage: 1,
      dialogVisible: false,
      selectDataTemp: [],
      firstIN: 1,
      myKey: '4dogs.cn',
      source: [],
      toolstype: [],
      objectlist: [],
      toolstatus: [],
      riskllist: [],
      role_str: '全部',
      rank: '',
      pageSize: 10,
      tableIds: [],
      SelNoToRight: [],
      selectmultipleSelection: [],
      curId: null,
      rightPage: 1,
      rightPageSize: 50,
      rightData: [],
      rightAllDataCount: 0, // 右侧总数据数
      level: []
    }
  },
  created: function () {
    this.$store.state.activefirstMenu = "/usergroup"
    this.pageSize = this.commonjs.pageSize
    this.curId = this.$route.params.id
  },
  mounted () {
    // this.getparams();





    this.getWaitUserList()
    this.getSelectedUserList();
    // this.rightData = [
    //   {
    //     "id": 31,
    //     "name": "admin123@"
    //   },
    //   {
    //     "id": 32,
    //     "name": "ces01"
    //   },
    //   {
    //     "id": 33,
    //     "name": "zs02"
    //   },
    //   {
    //     "id": 34,
    //     "name": "zhang"
    //   },
 

    // ]
  },
  methods: {
    goBack () {
      // this.$router.go(-1);
      this.$router.push({
        path: `/usergroup`,
        query: {
          page_num: this.fpage_num,
        }
      })
    },
    async getdownfile () {
      let params = {}
      const res = await API.getdownfile(params)
      // this.$ajax.get('/experience/write_specification/',{
      //  params: {}
      //})
      //.then((data) => { 
      //var dt = data.data;  
      // console.log(dt);  
      if (res.success) {
        this.downfile = res.specification_files
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    // 获取左侧工具列表
    async getWaitUserList () {
      console.log(1, this.$route)
      this.Loading = true
      let params = {
        id: this.curId,
        keyword: this.systemsearch_field,
        page: this.page_num,
        size: this.pageSize
      }
      const res = await API.getWaitUserList(params)
      this.Loading = false
      this.systemtableData = res.data.list
      this.tableIds.length = 0
      let defaultSelected = this.systemtableData.filter(item => item.selected)
      if (defaultSelected.length > 0) {
        defaultSelected.forEach(item => {
          this.tableIds.push(item.id)
        })
        console.log('defaultSelected', defaultSelected)
      }
      this.totalpage = res.data.total
      this.$nextTick(() => {
        this.toggleSelection(defaultSelected)
      })
    },
    // 默认选中项
    toggleSelection (rows) {
      if (rows) {
        rows.forEach(row => {
          this.$refs.multipleTable.toggleRowSelection(row, true)
        })
      } else {
        this.$refs.multipleTable.clearSelection()
      }
    },
    // 获取右侧列表数据
    async getSelectedUserList () {
      let params = {
        id: this.curId
      }
      const res = await API.getSelectedUserList(params)
      var dt2 = res.data
      this.rightData = []
      this.rightAllDataCount = dt2.total
      if (dt2.list) {
        this.rightData = this.rightData.concat(dt2.list)
        this.rightData.forEach(item => {
          item.hadSaved = true
        })
        this.experiencelist = this.experiencelist.concat(dt2.list)
        dt2.list.forEach(item => {
          this.experiencelistold.push(item.id)
        })
      }

    },
    // 右侧翻页
    load () {
        // this.rightAllDataCount > 50
      if (this.rightData.length>3) {
        let random = parseInt(Math.random() * 10000)
        console.log('下一页')
        // this.rightPage++
        // this.getSelectedUserList()
        this.rightData.push({
          "id": random,
          "name": random
        },)
      }

    },
    async getSystemExperience () { //工具列表
      this.Loading = true
      let params = {
        // name:this.systemsearch_field,
        // object:this.object, 
        // page:this.page_num,
        // type:this.type, 
        // level:this.level,
        // page_size:this.pageSize
        scheme_id: this.curId,
        page: this.page_num,
        page_size: this.pageSize
      }
      const res = await API.getWaitUserList(params)
      this.Loading = false
      this.systemtableData = res.results
      this.totalpage = res.total
      this.tableIds.length = 0
      let defaultSelected = this.systemtableData.filter(item => item.selected)
      defaultSelected.forEach(item => {
        this.tableIds.push(item.id)
      })
      this.$nextTick(() => {
        this.toggleSelection(defaultSelected)
      })
    },
    handleSelect (val, row) {
      if (this.firstIN === 1) { // 设置第一次进来才回触发事件
        /* 1 => add ; 0 => remove*/
        let flag = 0
        for (const i in val) {
          if (row.vul_id === val[i].vul_id) {
            flag = 1
            break
          }
        }
        if (flag === 1) {
          // 如果判断当前为添加则将当前勾选数据push到指定数组中
          this.experiencelist.push({
            id: row.id,
            name: row.name
          })
        } else {
          // 否则从数组中删除当前行数据
          for (const i in this.experiencelist) {
            if (this.experiencelist[i].id === row.id) {
              this.experiencelist.splice(i, 1)
            }
          }
          for (const j in this.SelNoToRight) {
            if (this.SelNoToRight[j].id === row.id) {
              this.SelNoToRight.splice(j, 1)
            }
          }
          // this.SelNoToRight = [];
        }
      }
      // console.log(this.rightData);
    },
    handleSelectAll (val) {
      if (this.firstIN === 1) { // 意思第一次点击不会执行hangleSelectAll里面的方法
        var v = this
        // remove
        if (val.length > 0) {
          this.tableIds.length = 0
          for (const n in val) {
            this.tableIds.push(val[n].id)
          }
        }
        if (val.length === 0) {
          this.SelNoToRight = []
          let numI = v.experiencelist.length
          for (var i = 0; i < numI; i++) {
            for (const j in this.tableIds) {
              if (v.experiencelist[i].id === this.tableIds[j]) {
                // console.log(i)
                v.experiencelist.splice(i--, 1)
                break
              }
            }
          }

        }
        if (v.experiencelist.length === 0) {
          for (const i in val) {
            v.experiencelist.push({
              id: val[i].id,
              name: val[i].name
            })
          }
        } else {
          for (const i in val) {
            let flag = false
            for (const j in v.experiencelist) {
              if (v.experiencelist[j].id === val[i].id) {
                flag = true
                break
              }
            }
            if (!flag) {
              v.experiencelist.push({
                id: val[i].id,
                name: val[i].name
              })
            }
          }
        }
        // if(val.length != this.pageSize){
        //     this.SelNoToRight = [];
        // }
      }
    },
    checkboxT (row, index) {
      if (row.selected) {
        return 0
      } else {
        return 1
      }
      // var flag = false;
      // for(var i=0;i<this.rightData.length;i++){
      //     if(row.id == this.rightData[i].id){
      //         flag = true; 
      //     }
      // }
      // if(flag){
      //     return 0;
      // }else{
      //     return 1;
      // }
    },
    async addToRight () { //添加到右边
      let userids = []
      this.experiencelist.forEach(item => {
        if (!item.hadSaved && !item.hadAddRight) {
          this.rightData.unshift(item)
          userids.push(item.id)
          let objIndex = this.systemtableData.findIndex(item2 => item2.id === item.id)
          if (objIndex > -1) {
            this.systemtableData[objIndex].selected = true
          }
        }
      })
      this.experiencelist.forEach(item => {
        item.hadAddRight = true
      })
      // this.rightData = this.experiencelist.concat();
      this.SelNoToRight = []
    },
    handlesystemSelectionChange: function (val) {
      this.multipleSelection = val
      // for(var i=0;i<this.multipleSelection.length;i++){
      //     var item = this.multipleSelection[i]; 
      //     let flag = false;
      //     for(var j = 0;j<this.rightData.length;j++){
      //         if(item.id == this.rightData[j].id){
      //             flag = true;
      //             break;
      //         }
      //     }
      //     if(!flag){
      //         this.SelNoToRight.push(item);
      //     } 
      // } 
      this.SelNoToRight = val
    },
    SelectedSelectionChange (val) {
      this.selectmultipleSelection = val
    },
    clearSelNoToRight () { //清除 选择 但未移入到 右边的
      console.log('this.experiencelist', this.experiencelist)
      let arr1 = [] // 已保存
      let arr2 = [] // 未保存
      this.experiencelist.forEach(item => {
        if (item.hadSaved || item.hadAddRight) { // 已保存或已添加到右侧
          arr1.push(item)
        } else {
          arr2.push(item)
        }
      })
      this.experiencelist = arr1
      let waitList = arr2
      let waitIdArr = []
      waitList.forEach(item => {
        waitIdArr.push(item.id)
      })
      this.systemtableData.forEach(item => {
        if (waitIdArr.includes(item.id)) {
          this.$refs.multipleTable.toggleRowSelection(item, false)
        }
      })
      console.log('this.experiencelist2', this.experiencelist)
      //    for(var x = 0;x <  this.SelNoToRight.length;x++){
      //     //    var id = this.SelNoToRight[x].id;
      // //        for(var i = 0;i<this.systemtableData.length;i++){
      // //             if(this.systemtableData[i].id == id){
      // //                 this.$refs.multipleTable.toggleRowSelection(this.systemtableData[i], false);
      // //             }
      // //         } 
      //    }


      // this.multipleSelection = [];
      // this.SelNoToRight = [];
      // this.experiencelist = this.rightData.concat();

    },
    handleReset () {
      this.systemsearch_field = ''
      this.page_num = 1
      this.currentpage = 1
      this.pageSize = 10
      this.getWaitUserList()
      this.getSelectedUserList()
    },
    handlesearch () {
      this.page_num = 1
      this.currentpage = 1
      this.getWaitUserList()
    },
    currentchange (t) {
      this.page_num = t
      this.getWaitUserList()
      this.currentpage = t
    },
    handleSizeChange (t) {
      this.page_num = 1
      this.pageSize = t
      this.getWaitUserList()
    },
    // 保存至方案
    async addExperienceToSet (flag) { //添加经验到经验集 
      // if(this.multipleSelection.length == 0) return;
      // var ids = [];
      // for (var i = 0; i < this.rightData.length; i++) {
      //     ids.push(this.rightData[i].id); 
      // }
      let add_tools_id = []
      let delete_tools_id = []
      let lastOldIds = []
      let user_ids = []
      this.rightData.forEach(item => {
        user_ids.push(item.id)
        //     if (!item.hadSaved) {
        //         // 新增
        //         add_tools_id.push(item.id)
        //     } else {
        //         // 剩余保存过的项
        //         lastOldIds.push(item.id)
        //     }
      })
      // 对比旧的原保存过的数据与剩余的原保存过的数据
      // this.experiencelistold.forEach(itemID => {
      //     if (!lastOldIds.includes(itemID)) {
      //         delete_tools_id.push(itemID)
      //     }
      // })
      // console.log('新增的', add_tools_id)
      // console.log('删除的', delete_tools_id)
      let params = {
        group_id: this.curId,
        // add_ids:add_tools_id.join(','),
        // delete_ids:delete_tools_id.join(','),
        user_ids: user_ids.join(',')
      }
      const res = await API.addGroupUser(params)
      if (res.code === 200) {
        this.$message({
          message: '添加用户组成功',
          type: 'success'
        })
        this.experiencelistold = []
        this.rightData.forEach(item => {
          item.hadSaved = true
          this.experiencelistold.push(item.id)
        })
        // this.getWaitUserList(); 
        // this.$refs.multipleTable.clearSelection();
        // this.getSystemExperience();
        // // /experienceSet
        this.$router.push({ path: '/usergroup' })
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    clearExperienceToSet: function (flag) {
      this.$refs.multipleTable.clearSelection()
    },
    removeExperience: function (id) { //从现有经验集中移除   
      for (var i = 0; i < this.systemtableData.length; i++) {
        if (this.systemtableData[i].id == id) {
          this.$refs.multipleTable.toggleRowSelection(this.systemtableData[i], false)
        }
      }
      for (var i = 0; i < this.experiencelist.length; i++) {

        if (this.experiencelist[i].id == id) {
          this.experiencelist.splice(i, 1)
        }

      }
    },
    removeSelect () {
      for (var x = 0; x < this.selectmultipleSelection.length; x++) {
        var id = this.selectmultipleSelection[x].id
        for (var i = 0; i < this.systemtableData.length; i++) {
          let objIndex = this.experiencelist.findIndex(item => item.id === id)
          if (objIndex > -1) {
            this.experiencelist.splice(objIndex, 1)
          }
          if (this.systemtableData[i].id == id) {
            this.systemtableData[i].selected = false
            this.$refs.multipleTable.toggleRowSelection(this.systemtableData[i], false)
            break
          }
        }
        for (var i = 0; i < this.rightData.length; i++) {
          if (this.rightData[i].id == id) {
            this.rightData.splice(i, 1)
            this.selectmultipleSelection.splice(x, 1)
            x--
            break
          }
        }
        this.$refs.selectmultipleTable.toggleRowSelection(this.rightData[x], false)
        // this.load();
      }
      if (this.rightData.length === 0 && this.rightPage < Math.ceil(this.rightAllDataCount / 50)) {
        this.rightPage++
        this.getSelectedUserList()
      }
      // this.experiencelist = this.rightData.concat();
    },
    cancelAddExperience: function () { //取消添加自有经验
      this.formdt.name = ''
      this.formdt.content = ''
      this.formdt.remark = ''
      this.dialogVisible = false
    },
    clickupload () {
      document.querySelector('.btnUploadID').click()
    },
    changeuploaID (e) {
      let deviceFile = e.target.files
      for (let i = 0; i < deviceFile.length; i++) {
        this.addform.file = deviceFile[i]
        this.addform.filename = deviceFile[i].name
      }
    },
    addNewrules () {
      this.dialogVisible = true
    },

  }
})

</script>