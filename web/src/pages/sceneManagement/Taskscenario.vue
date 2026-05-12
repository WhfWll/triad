<template>
  <div>
    <div class="main-title">
      <i class="nav_icon"></i>
      <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i>  
      <label for="" class="taskSceneBtn">{{ title }}</label>
    </div>
    <div class="scenlist context_box_bg">
      <div class="search-box">
        <div class="operationbutton">
          <el-button type="primary" @click="createScene">新建</el-button>

          <del-button
            :width="170"
            @click="btnMultiDelete"
            :disabled="!multipleSelection.length"
            style="margin-left: 8px"
          >
          </del-button>
        </div>
        <div class="serach-condition">
          <div class="search-text">
            <el-input
              placeholder="请输入关键字"
              @keydown.enter.native="handlesearch"
              v-model="search_item.search"
              class="input-with-select"
              clearable
            >
            </el-input>
            <el-button type="primary" @click="handlesearch">搜索</el-button>
          </div>
          <div>
            <el-button type="primary" @click="handleReset">重置</el-button>
          </div>
        </div>
      </div>
      <el-table
        ref="multipleTable"
        :data="tableData"
        v-loading="Loading"
        style="width: 100%"
        @selection-change="handleSelectionChange"
        @cell-mouse-enter="mouseenter"
        @cell-mouse-leave="mouseleave"
         height="calc(100% - 102px)"
      >
        <!--  -->
        <el-table-column type="selection" width="55" :selectable="checkboxT">
        </el-table-column>
        <el-table-column
          prop="templateName"
          label="场景名称"
          :show-overflow-tooltip="true"
        >
          <template #default="scope">
            <el-link @click="fnDetails(scope.row)">{{
              scope.row.templateName
            }}</el-link>
          </template>
        </el-table-column>
        <el-table-column
          prop="describe"
          label="场景描述"
          :show-overflow-tooltip="true"
        >
        </el-table-column>
        <el-table-column prop="isDefault" label="默认场景">
          <template #default="scope">
            <span v-if="scope.row.isDefault == 1">是</span>
            <span v-else>否</span>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间"> </el-table-column>
        <el-table-column prop="userName" label="提交者">
          <template #default="scope">
            <div v-if="showOperateButton && rowId == scope.row.id">
              <el-popover
                placement="bottom"
                width="200"
                popper-class="delButton_popper"
                trigger="click"
                :ref="`popover_id-${scope.row.id}`"
              >
                <p class="delText">
                  <i class="el-icon-warning"></i>确定设置为默认场景吗？
                </p>
                <div style="text-align: right; margin: 0">
                  <el-button
                    size="small"
                    class="delCancel"
                    @click="fncancel1(scope)"
                    >取消
                  </el-button>
                  <el-button
                    size="small"
                    type="primary"
                    @click="btnDefault(scope)"
                    >确定
                  </el-button>
                </div>
                <template #reference>
                  <span
                    class="link_info linkafter"
                    style="padding: 0; cursor: pointer"
                    >默认</span
                  >
                </template>
              </el-popover>
              <el-link
                :underline="false"
                class="link_primary"
                style="vertical-align: initial"
                @click="copyScene(scope.row)"
                >复制</el-link
              >
              <el-link
                :underline="false"
                class="link_primary"
                style="vertical-align: initial"
                @click="fnKnowledge(scope.row.id, scope.row.templateName)"
                >图谱</el-link
              >
              <el-popover
                v-if="scope.row.isDefault != 1"
                placement="bottom"
                width="170"
                :visible-arrow="false"
                :ref="`popover-${scope.row.id}`"
                popper-class="delButton_popper"
              >
                <p class="delText">
                  <i class="el-icon-warning"></i>确定删除吗？
                </p>
                <div style="text-align: right; margin: 0">
                  <el-button
                    size="small"
                    class="delCancel"
                    @click="fncancel(scope)"
                    >取消</el-button
                  >
                  <el-button size="small" type="primary" @click="btnDel(scope)"
                    >确定</el-button
                  >
                </div>
                <span
                  slot="reference"
                  class="link_danger linkafter2"
                  style="cursor: pointer"
                  >删除</span
                >
              </el-popover>
            </div>
            <div v-else>
              <span>{{ scope.row.userName }}</span>
            </div>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        background
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="currentPage"
        :page-size="pageSize"
        layout=" total,  prev, pager, next, sizes,jumper"
        :total="total"
      >
      </el-pagination>
    </div>
  </div>
</template>
<style lang="less">
@import "./css/taskscenario.less";
</style>
<script>
import scene from '@/api/scene.js'
import DelButton from "@/components/DelButton.vue"
export default ({
  name: 'taskscenario',
  components: {
    DelButton
  },
  data () {
    return {
      title: '任务场景',
      alldelvisible: false,
      showOperateButton: false,
      tableData: [],
      multipleSelection: [],
      total: 0,
      pageSize: 10,
      currentPage: 1,
      search_item: {
        search: '',
        page: 1,
      },
      Loading: false,
      rowId: '',
      visible: false,
      alldelvisible1: false,
    }
  },
  created () {
    this.$store.state.activefirstMenu = '/taskscenario'
    this.userid = this.$commonjs.decryptCBC(localStorage.getItem('user_id'), this.$commonjs.myKey)
  },
  mounted () {
    this.getData()
  },
  methods: {
    fnKnowledge (_id, _name) { //知识图谱
      this.$router.push({
        path: `/knowledgegraphsm`,
        query: {
          id: _id,
          name: _name,
        }
      })
    },
    getData () {
      let params = {
        search: this.search_item.search,
        page: this.search_item.page,
        size: this.pageSize,
      }
      scene.getData(params).then(res => {
        if (res.code == 200) {
          this.tableData = res.data.list
          this.total = res.data.total
        } else {
          this.$ElMessage({
            message: res.msg,
            type: 'error'
          })
        }
      }).catch(err => {

      })
    },
    checkboxT (row, index) {
      if (row.isDefault == 1) {
        return 0
      } else {
        return 1
      }
    },
    createScene () { //新建场景
      this.$router.push({
        path: `/createscene`,
        query: {
          isAdd: 1,
        }
      })
    },
    handlesearch () {
      this.search_item.page = 1
      this.currentPage = 1
      this.getData()
    },
    handleReset () {
      this.search_item.page = 1
      this.search_item.search = ''
      this.pageSize = 10
      this.currentPage = 1
      this.getData()
    },
    handleSizeChange (t) {
      this.search_item.page = 1
      this.pageSize = t
      this.getData()
    },
    handleCurrentChange (t) {
      this.search_item.page = t
      this.getData()
    },
    fncancel (scope) {
      scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
      scope._self.$refs[`popover-${scope.row.id}`].doClose()
    },
    fncancel1 (scope) {

      scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
      scope._self.$refs[`popover-${scope.row.id}`].doClose()
      // this.clickEmpty()//wumeng处理两个popovar下下策
    },
    // 单个删除
    async btnDel (scope) {
      let _ids = scope.row.id

      const res = await scene.delScene({
        taskTemplateIds: _ids
      })
      if (res.code == 200) {
        this.$message({
          message: '删除场景成功',
          type: 'success'
        })
        scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
        // this.deldialogVisible = false;
        this.getData()

      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    async copyScene (row) { //复制场景
      const res = await scene.copyScene({
        taskTemplateId: row.id,
        userId: this.userid,
      })
      if (res.code == 200) {
        this.$message({
          message: '复制场景成功',
          type: 'success'
        })
        this.getData()
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    fnDetails (row) { //详情
      this.$router.push({
        path: `/createscene`,
        query: {
          isAdd: 0,
          sceneid: row.id,
          scene_name: row.templateName,
        }
      })
    },
    async btnDefault (scope) {
      let params = {
        taskTemplateId: scope.row.id,
      }
      const res = await scene.defaultScene(params)
      if (res.code == 200) {
        this.$message({
          message: '设置默认场景成功',
          type: 'success'
        })
        scope._self.$refs[`popover-${scope.row.id}`].doClose()
        this.getData()

      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    // 批量删除
    async btnMultiDelete () {
      if (this.multipleSelection.length == 0) return

      var _ids = this.multipleSelection.map(item => item.id)

      let params = {
        taskTemplateIds: _ids.join(',')
      }

      const res = await scene.delScene(params)
      if (res.code == 200) {
        this.$message({
          message: '批量删除成功',
          type: 'success'
        })
        this.alldelvisible = false
        this.getData()
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    handleSelectionChange (val) {
      this.multipleSelection = val
    },
    mouseenter (row, column, cell, event) {
      this.showOperateButton = true
      let _id = row.id
      this.rowId = _id//赋值行id，便于页面判断
    },
    mouseleave (row, colum, cell, event) {
      let t = this.$refs['popover_id-' + row.id] && this.$refs['popover_id-' + row.id].showPopper
      let t2 = this.$refs['popover-' + row.id] && this.$refs['popover-' + row.id].showPopper

      if (!t && !t2) {
        this.showOperateButton = false
        this.rowId = ""
      }

    },
  }
})
</script>