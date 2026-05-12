<template>
  <div
    style="position: relative; overflow: hidden; height: 100%"
    id="target_box_id"
  >
    <div class="main-title">
      <router-link :underline="false" class="classA" :to="{ path: '/task' }"
        >渗透任务</router-link
      >
      <router-link
        :underline="false"
        class="classB"
        :to="{
          path: '/taskdetail',
          query: { id: task_Id, name: task_name, risk_level: risk_level },
        }"
      >
        <el-tooltip
          class="item"
          effect="dark"
          :content="task_name"
          placement="bottom"
        >
          <span> {{ task_name }}</span> </el-tooltip
        >任务
      </router-link>
      <label class="currentpagetitle">
        <el-tooltip
          class="item"
          effect="dark"
          :content="target_name"
          placement="bottom"
        >
          <span> {{ target_name }}</span> </el-tooltip
        >目标
      </label>
    </div>
    <div class="targetlist context_box_bg">
      <div class="iconlist">
        <ul>
          <!-- <li class="attack" @click="AddAttack" title="新增攻击面">
            <i class="iconfont iconxinzenggongjimian"></i>
          </li>
          <li class="loop" @click="AddLoop" title="新增漏洞执行">
            <i class="iconfont iconzhihangloudong"></i>
          </li> -->
        </ul>
      </div>
      <el-tabs v-model="activeName" @tab-click="handleClick">
        <el-tab-pane label="渗透路径" name="tabs1">
          <Seepagepath />
        </el-tab-pane>
      </el-tabs>
    </div>
    <!-- 新增攻击面 -->
    <div class="node_pop" id="addattack">
      <div class="test_dialog">
        <div class="close" @click="closeDialogattack">
          <i class="iconfont iconquxiao"></i>
        </div>
        <div class="el-dialog__header">
          <div class="dialog__title clearfix">
            <span class="title_name"> 新增攻击面 </span>
          </div>
        </div>
        <div class="dialog_b_btn">
          <el-button size="small" @click="submitattackForm">测试</el-button>
        </div>
        <div class="el-dialog__body">
          <el-form :model="attackform" :rules="rules" ref="attack">
            <el-form-item label="" prop="attack_type">
              <label class="dialog_item_label"
                >攻击面类型<i class="is-required" style="float: right"
                  >*</i
                ></label
              >
              <el-select
                v-model="attackform.attack_type"
                size="small"
                placeholder="请选择"
                style="width: 320px"
              >
                <el-option
                  v-for="(item, index) in attack_type_list"
                  :key="index"
                  :label="item.label"
                  :value="item.value"
                ></el-option>
              </el-select>
            </el-form-item>
            <!-- 端口 -->
            <div v-if="attackform.attack_type == 1">
                <el-form-item prop="target">
                  <label class="dialog_item_label" for="">目标<i class="is-required" style="float: right">*</i></label
                  >
                  <el-input v-model="attackform.target" style="width: 320px" disabled>
                  </el-input>
              </el-form-item>
                <el-form-item prop="port">
                  <label class="dialog_item_label" for=""
                    >端口<i class="is-required" style="float: right">*</i></label
                  >
                  <el-input v-model="attackform.port" style="width: 320px">
                  </el-input>
                </el-form-item> 
            </div>
            <!-- 敏感路径 -->
            <div v-if="attackform.attack_type == 2">
                <el-form-item prop="path_type">
                    <label class="dialog_item_label" for="">敏感路径类型<i class="is-required"
                            style="float: right">*</i></label>
                    <el-select v-model="attackform.path_type" size="small" placeholder="请选择" style="width: 320px">
                        <el-option v-for="(item, index) in path_type_list" :key="index" :label="item.name"
                            :value="item.id"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item prop="target">
                    <label class="dialog_item_label" for="">URL<i class="is-required"
                            style="float: right">*</i></label>
                    <el-input v-model="attackform.url" style="width: 320px">
                    </el-input>
                </el-form-item> 
            </div>
            <!-- 登录凭证 -->
            <div v-if="attackform.attack_type == 3">
              <el-form-item prop="target">
                  <label class="dialog_item_label" for="">目标IP<i class="is-required"
                          style="float: right">*</i></label>
                  <el-input v-model="attackform.ip" style="width: 320px">
                  </el-input>
              </el-form-item>
              <el-form-item prop="port">
                  <label class="dialog_item_label" for="">端口<i class="is-required"
                          style="float: right">*</i></label>
                  <el-input v-model="attackform.port" style="width: 320px">
                  </el-input>
              </el-form-item>
              <el-form-item prop="scheme">
                  <label class="dialog_item_label" for="">协议</label>
                  <el-select v-model="attackform.scheme" size="small" style="width: 320px">
                      <el-option v-for="item in scheme_list" :key="item.value" :label="item.label"
                          :value="item.value">
                      </el-option>
                  </el-select>
              </el-form-item>
              <el-form-item prop="auth_type">
                  <label class="dialog_item_label" for="">认证方式<i class="is-required"
                          style="float: right">*</i></label>
                  <el-select v-model="attackform.auth_type" size="small" placeholder="请选择" style="width: 320px">
                      <el-option v-for="(item, index) in auth_type_list" :key="index" :label="item.name"
                          :value="item.id"></el-option>
                  </el-select>
              </el-form-item>
              <el-form-item prop="auth_value" v-if="attackform.auth_type == 1">
                  <label class="dialog_item_label" for="">Cookie</label>
                  <el-input v-model="attackform.auth_value" type="textarea" :rows="4" size="small"
                      style="width: 320px">
                  </el-input>
              </el-form-item>
              <el-form-item prop="auth_value" v-else>
                  <label class="dialog_item_label" for="">Headers</label>
                  <el-input v-model="attackform.auth_value" type="textarea" :rows="4" size="small"
                      style="width: 320px">
                  </el-input>
              </el-form-item>

            </div>
          </el-form>
        </div>
      </div>
    </div>
    <!-- 新增漏洞执行 -->
    <div class="node_pop" id="vulnlistdialog">
      <div class="test_dialog">
        <div class="close" @click="closeDialogloop">
          <i class="iconfont iconquxiao"></i>
        </div>
        <div class="el-dialog__header">
          <div class="dialog__title clearfix">
            <span class="title_name"> 新增漏洞执行 </span>
          </div>
        </div>
        <div class="el-dialog__body">
          <div class="serach-condition">
            <div class="search-text">
              <el-input
                placeholder="搜索漏洞名称"
                @keydown.enter.native="handlloopesearch"
                v-model="loopformData.loopsearch"
                class="input-with-select"
                size="small"
                clearable
              >
              </el-input>
              <el-button type="primary" size="small" @click="handlloopesearch"
                >搜索</el-button
              >
            </div>
          </div>
          <el-table
            :data="looptableData"
            style="width: 100%"
            highlight-current-row
            ref="table"
            @row-click="handleRowClick"
          >
            <el-table-column
              prop="name"
              label="漏洞名称"
              :show-overflow-tooltip="true"
            >
            </el-table-column>
            <el-table-column prop="vul_risk_label" label="漏洞风险" width="100">
              <template slot-scope="scope">
                <span
                  :class="[
                    { 'riskstyle risk_hight': scope.row.risk == 1 },
                    { 'riskstyle risk_middle': scope.row.risk == 2 },
                    { 'riskstyle risk_low': scope.row.risk == 3 },
                    { 'riskstyle risk_nofind': scope.row.risk == 4 },
                  ]"
                >
                  <i></i>{{ scope.row.riskName }}</span
                >
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            :page-size="looppageSize"
            :pager-count="5"
            background
            layout=" total,  prev, pager, next, jumper"
            :total="looptotalpage"
            :current-page="loopcurrentpage"
            @current-change="handleloopcurrentchange"
            @size-change="handleloopSizeChange"
          >
          </el-pagination>
        </div>
      </div>
    </div>

    <!-- 漏洞执行 -->
    <div class="node_pop" id="addloopexecute">
      <div class="test_dialog">
        <div class="el-dialog__header">
          <div class="dialog__title clearfix">
            <span class="title_name">
              {{ loopexecute_dialog_title }}
            </span>
          </div>
        </div>
        <div class="dialog_b_btn">
          <el-button size="small" @click="submitFormloop">测试</el-button>
        </div>
        <div class="el-dialog__body"  v-loading="loading"> 
          <el-form :model="testform" :rules="looprules" ref="testform">
             
            <el-form-item prop="path">
              <label class="dialog_item_label" for="">路径 <i class="is-required" style="float: right">*</i></label>
                <el-input v-model="testform.path" style="width: 320px">
                </el-input>
            </el-form-item> 
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>
  
  <style scoped>
.packtxt textarea {
  resize: none !important;
}
</style>
  <style lang="less" scoped>
.is-required {
  margin-right: 4px;
  color: #f56c6c;
  font-size: 12px;
}
.iconlist {
  position: absolute;
  right: 24px;
  top: 14px;
  z-index: 9;

  ul {
    list-style: none;

    li {
      width: 20px;
      height: 20px;
      line-height: 20px;
      display: inline-block;
      margin-left: 20px;

      i {
        display: inline-block;
        cursor: pointer;
      }

      &.line {
        span {
          display: inline-block;
          height: 16px;
          width: 2px;
          background-color: #d8d8d8;
          vertical-align: middle;
        }
      }

      &.run {
        i {
          color: #66c681;
        }
      }

      &.status {
        width: 52px;
        vertical-align: middle;

        span {
          display: block;
        }
      }

      &.pasuse,
      &.resume {
        i {
          color: #f9b640;
        }
      }

      &.rungrey {
        i {
          color: rgba(72, 72, 102, 0.32);
          cursor: not-allowed;
        }
      }

      &.stop {
        i {
          color: #f87d7d;
        }
      }

      &.target {
        margin-left: 4px;

        i {
          color: #4c7ae3;
        }
      }

      &.attack {
        i {
          color: #f9b640;
        }
      }

      &.loop {
        i {
          color: #f87d7d;
        }
      }
    }
  }
}

.openAnimate {
  right: 0px;
  box-shadow: 0px 8px 32px 0px rgba(76, 122, 227, 0.12);
  animation: openAnimate 0.5s;
  -webkit-animation: openAnimate 0.5s;
  animation-fill-mode: forwards;
}

.closeAnimate {
  right: -640px;
  animation: closeAnimate 0.5s;
  -webkit-animation: closeAnimate 0.5s;
  animation-fill-mode: forwards;
}

@keyframes openAnimate {
  0% {
    right: -635px;
  }

  100% {
    right: 0px;
  }
}

@keyframes closeAnimate {
  0% {
    right: 0px;
  }

  100% {
    right: -635px;
  }
}

.openAnimate2 {
  right: 580px;
  box-shadow: 0px 8px 32px 0px rgba(76, 122, 227, 0.12);
  animation: openAnimate2 0.5s;
  -webkit-animation: openAnimate2 0.5s;
  animation-fill-mode: forwards;
}

@keyframes openAnimate2 {
  0% {
    right: 0;
  }

  100% {
    right: 580px;
  }
}

.closeAnimate2 {
  right: 0;
  animation: closeAnimate2 0.5s;
  -webkit-animation: closeAnimate2 0.5s;
  animation-fill-mode: forwards;
}

@keyframes closeAnimate2 {
  0% {
    right: 580px;
  }

  100% {
    right: 0;
  }
}

#addattack {
  top: 20px;
  z-index: 99;
}

#vulnlistdialog {
  top: 20px;
  z-index: 99;
}

#addloopexecute {
  top: 20px;
  z-index: 90;
}

// 节点详情
.node_pop {
  position: absolute;
  // right: 0;
  right: -635px;
  top: 76px;
  width: 580px;
  height: 100%;
  background: #fff;
  z-index: 88;
  box-shadow: 0px 8px 32px 1px rgba(76, 122, 227, 0.12);
  border-radius: 4px 0px 0px 0px;

  .test_dialog {
    width: 100%;
    height: 100%;
    position: relative;

    .close {
      cursor: pointer;
      width: 54px;
      height: 54px;
      text-align: center;
      line-height: 54px;
      position: absolute;
      left: -54px;
      top: 70px;
      background: #4c7ae3;
      color: #fff;

      i {
        font-size: 22px;
      }
    }

    /deep/ .el-dialog__header {
      border-radius: 4px 0px 0px 0px;
    }

    .dialog__title {
      font-size: 14px;

      box-sizing: border-box;

      .title_name {
        display: inline-block;
        border-left: 2px solid #fff;
        background: hsla(0, 0%, 100%, 0.12);
        color: #fff;
        height: 32px;
        line-height: 32px;
        text-align: center;
        padding: 0 24px;
      }

      .title_item {
        // float: left;
        // width: 200px;
        margin: 8px 24px 8px 0;

        .micon {
          color: rgba(255, 255, 255, 0.7);
          vertical-align: top;
          height: 20px;
          line-height: 19px;
          display: inline-block;
        }

        .lbname {
          display: inline-block;
          // width: 54px;
          margin-right: 18px;
          margin-left: 8px;
          color: rgba(255, 255, 255, 0.7);
          height: 20px;
          line-height: 20px;
          vertical-align: middle;
        }

        .spvalue {
          // width: calc(100% - 100px) !important;
          display: inline-block;
          overflow: hidden;
          white-space: nowrap;
          text-overflow: ellipsis;
          color: #fff;
          vertical-align: text-top;
          height: 20px;
          line-height: 20px;
          position: relative;
          vertical-align: middle;
        }

        .name {
          display: block;
          font-size: 14px;
          width: 100%;
          margin-bottom: 10px;
          border-left: 3px solid #4c7ae3;
          padding-left: 10px;
          padding-right: 10px;
          box-sizing: border-box;
          color: #fff;

          span {
            display: block;
            float: right;
            font-size: 12px;
            border-radius: 10px;
            color: #fff;
            padding: 2px 8px;
          }

          .spcor1 {
            background: #09c1f7;
          }

          .spcor2 {
            background: #15c53d;
          }

          .spcor3 {
            background: #f35f28;
          }
        }
      }
    }

    /deep/ .el-dialog__body {
      padding: 14px;
      overflow-y: auto;
      height: calc(100% - 150px);
      .desc {
        padding: 8px;
        font-size: 14px;
      }

      .msg_list {
        div {
          margin-top: 10px;

          label {
            display: inline-block;
            min-width: 60px;
            font-size: 12px;
            padding-left: 10px;
            height: 14px;
            line-height: 14px;
            color: rgba(72, 72, 102, 0.64);
            // vertical-align: sub;
          }

          span {
            font-size: 12px;
            display: inline-block;
            // width: calc(100% - 130px);
            color: rgba(72, 72, 102, 0.64);
            vertical-align: text-top;
          }
        }
      }

      .serach-condition {
        width: 100%;
        margin-bottom: 20px;

        .search-text {
          width: 100%;

          .input-with-select {
            // width: calc(100% - 540px);
            width: 480px;
            vertical-align: bottom;
          }

          .input-with-select input {
            border-radius: 2px 0px 0px 2px !important;
          }

          button {
            border-radius: 0px 2px 2px 0px;
          }
        }
      }

      .serach-condition > div {
        display: inline-block;
        margin-right: 8px;
      }

      .serach-condition > div:last-child {
        margin-right: 0;
      }

      .serach-condition > div > label {
        color: #606266;
        font-size: 14px;
      }
    }
  }

  /deep/ .el-dialog__body {
    padding: 0;
  }

  /deep/ .el-tabs__header {
    box-shadow: none;
  }

  /deep/ .el-tabs__content {
    padding: 0 16px;
  }
}

.loglist {
  list-style: none;
  margin-top: 24px;

  li {
    color: rgba(72, 72, 102, 0.64);
    font-size: 13px;
  }
}

.network i {
  display: inline-block;
  // width: 6px;
  // height: 6px;
  // background-color: #ebeef5;
  // background: #15C53D;
  margin: 0 2px;
}

.network i.wifig {
  background: #fff;
}

.network {
  i {
    position: absolute;
    width: 4px;
    background: rgba(255, 255, 255, 0.7);

    &:nth-child(1) {
      top: 10px;
      left: 0;
      height: 4px;
    }

    &:nth-child(2) {
      top: 8px;
      left: 6px;
      height: 6px;
    }

    &:nth-child(3) {
      top: 6px;
      left: 12px;
      height: 8px;
    }

    &:nth-child(4) {
      top: 4px;
      left: 18px;
      height: 10px;
    }

    &:nth-child(5) {
      top: 2px;
      left: 24px;
      height: 12px;
    }
  }
}

/deep/ .el-switch.is-checked .el-switch__core {
  border-color: #fff;
  background-color: #fff;
}

/deep/ .el-switch__core:after {
  background-color: #4c7ae3;
}

#svgdiv {
  height: calc(100% - 180px);
}

.vulnstatistics {
  font-size: 13px;
  display: inline-block;
  margin-right: 54px;

  i {
    margin-right: 5px;
    vertical-align: baseline;
  }

  label {
    color: rgba(72, 72, 102, 0.64);
    margin-right: 16px;
  }

  span {
    color: rgba(72, 72, 102, 0.87);
  }
}

.port {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 4px;
}

.port1 {
  background: #f87d7d;
}

.port2 {
  background: #f9b640;
}

.port3 {
  background: #4c7ae3;
}

.port4 {
  background: #15c53d;
}

.port1,
.port2,
.port3,
.port4 {
  vertical-align: super;
}

// .iconduocijiance{
//     vertical-align: bottom;
// }
.tag_status {
  width: auto;
  padding: 0 8px;
}

// .selstatus /deep/ .el-input__inner{
//     height:40px !important;
// }
.updatestatus {
  /deep/ .el-dialog__body {
    height: 192px !important;
  }

  /deep/ .el-dialog {
    height: auto !important;
  }

  /deep/ .el-dialog__body {
    padding: 72px 152px !important;
  }
}

.rotateZ {
  transform: rotateZ(180deg);
}

/deep/ .el-dropdown-menu__item:hover {
  background: #f7f7fb !important;
}

.notarget {
  color: rgba(72, 72, 102, 0.32);
  cursor: default;
  min-width: 160px;
  box-sizing: border-box;
}

.notarget:hover {
  color: rgba(72, 72, 102, 0.32);
  background: #fff !important;
}

.targetitem {
  min-width: 200px;

  span {
    i {
      display: inline-block;
      width: 8px;
      height: 8px;
      border-radius: 50%;
    }

    margin-right: 16px;
    width: 55px;
    display: inline-block;
    height: 17px;
    line-height: 17px;
  }

  a {
    text-decoration: none;
    color: rgba(72, 72, 102, 0.64);
    display: inline-block;
    // min-width:100%;
  }

  a:hover {
    background: #f7f7fb !important;
  }

  label {
    display: inline-block;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    width: 110px;
    height: 17px;
    line-height: 17px;
    vertical-align: sub;
    cursor: pointer;
  }
}

.targetitem:hover {
  a {
    color: #4c7ae3 !important;
  }

  background: #f7f7fb !important;
}

.textarea /deep/ textarea {
  resize: none !important;
}

.controlinput {
  width: 90% !important;
  box-sizing: border-box;

  /deep/ .el-input__inner {
    border: none !important;
    padding-left: 0;
  }
}

.controlbox {
  margin-top: 16px;
  color: rgba(72, 72, 102, 0.64);

  .cmdresult {
    padding: 16px 0;
    word-wrap: break-word;
    word-break: normal;
  }
}

.useinput {
  width: 90% !important;
  box-sizing: border-box;

  /deep/ .el-input__inner {
    border: none !important;
    padding-left: 0;
  }
}

/deep/ .el-progress-bar {
  width: 110px !important;
}

/deep/ .el-table td:not(.el-table-column--selection):first-child .cell,
/deep/ .el-table th:not(.el-table-column--selection):first-child .cell {
  padding-left: 32px !important;
}

.bugbasicinfo /deep/ .el-table--enable-row-transition .el-table__body td {
  transition: none;
}

.bugbasicinfo /deep/ .el-table--enable-row-hover .el-table__body tr:hover > td {
  background: none !important;
}

.iconcircular {
  position: absolute;
  border-radius: 50%;
}

.circular1 {
  width: 30px;
  height: 30px;
  background: #3768d6;
  box-shadow: 0px -2px 2px 0px #3157af;
  bottom: -19px;
  left: 131px;
}

.circular2 {
  width: 86px;
  height: 86px;
  box-shadow: 0px 2px 2px 0px #2e5dc8;
  border: 20px solid #3768d6;
  top: -82px;
  left: 50%;
}

.circular3 {
  width: 67px;
  height: 67px;
  background: #3768d6;
  box-shadow: 0px -2px 2px 0px #3157af;
  bottom: -44px;
  right: 100px;
}

/deep/ .el-progress__text {
  color: #fff;
  font-size: 13px !important;
}

/deep/ .el-progress-bar__outer {
  background: rgba(255, 255, 255, 0.3);
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

/deep/ .el-tabs__item {
  height: 48px;
  line-height: 48px;
  padding: 0 24px;
}

/deep/ .el-tabs__item.is-active {
  color: #4c7ae3;
  font-weight: 500;
}

/deep/ .el-tabs__nav-wrap {
  padding: 0 24px;
}

/deep/ .el-tabs__nav-wrap::after {
  background: #e8e8f5;
  height: 1px;
}

/deep/ .el-tabs__header {
  margin: 0 0 24px;
}
.svgCanvas /deep/ g.node rect {
  padding: 5px;
  fill: #dbe4f9;
  stroke: #4c7ae3;
  stroke-width: 1px;
}

.svgCanvas /deep/ g.node text {
  color: #4c7ae3;
  font: 14px sans-serif;
  /*font-weight:700;*/
  font-size: 12px;
  cursor: pointer;
}

#svgCanvas {
  height: 100%;
  min-height: 590px;
  width: 100%;
}

.svgCanvas .node text {
  font-weight: bold;
  font-family: "Microsoft YaHei";
  font-size: 12px;
  pointer-events: none;
  text-anchor: middle;
  fill: white;
}

.svgCanvas .label g {
  transform: translate(0, -13px);
}

.svgCanvas /deep/ .node rect {
  fill: white;
  stroke-width: 0px;
  color: white;
}

.svgCanvas g > rect {
  stroke: #dee6e8;
  stroke-width: 1px;
}

.svgCanvas g text {
  fill: #687386 !important;
  font-size: 14px !important;
}

.svgCanvas /deep/ .edgePath path {
  stroke: #ccc;
  stroke-width: 1px;
}

/deep/ .el-tabs__header {
  margin: 0 0 15px;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
  border: none;
}

.targetlist {
  min-height: calc(100% - 39px);
  box-sizing: border-box;
  // background: #fff;
  // box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
  position: relative;
}

.context_box_bg {
  background: none;
}

.target_box {
  // padding:  24px 24px;
  box-sizing: border-box;
  position: relative;
  height: 100%;
  background: #fff;
  box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
  border-radius: 4px;
  margin-bottom: 10px;

  .el-table__body-wrapper {
    height: calc(100% - 54px);
  }
}

.networkpopover_1 {
  display: inline-block;
}

.networkpopover_1 > label {
  margin-bottom: 8px;
  padding-left: 13px;
  font-size: 13px;
  color: rgba(72, 72, 102, 0.64);
  margin-right: 24px;
}

.targettypelist {
  position: absolute;
  right: 24px;
  top: 8px;
  z-index: 9;

  .el-tabs {
    height: 100%;
    position: absolute;

    .el-tabs__content {
      height: calc(100% - 70px);

      .el-tab-pane {
        height: 100%;
      }
    }
  }
}

.el-tabs__content {
  overflow: auto !important;
}

.el-dropdown-link {
  background: rgba(76, 122, 227, 0.12);
  display: inline-block;
  height: 32px;
  width: 160px;
  line-height: 32px;
  padding-left: 12px;
  box-sizing: border-box;
  font-weight: 500;
  color: rgba(76, 122, 227, 1);
  cursor: pointer;
  outline: none;

  i {
    float: right;
    margin-right: 12px;
    // margin-top: 10px;
  }
}

.attack_tabs > div {
  float: left;
  height: 32px;
  line-height: 32px;
  border-radius: 2px;
  border: 1px solid rgba(232, 232, 245, 1);
  padding: 0 24px;
  font-weight: 500;
  color: rgba(72, 72, 102, 0.64);
  font-size: 14px;
  margin-left: -1px;
  -moz-user-select: none;
  -khtml-user-select: none;
  user-select: none;
}

.attack_tabs > div.active {
  background: #4c7ae3;
  color: #fff;
  border: 1px solid #4c7ae3;
}

.attack_tabs_content {
  margin-top: 24px;
}

.target_box /deep/ .el-tabs__nav-wrap {
  padding: 0;
}

.target_box /deep/ .el-tabs--card > .el-tabs__header .el-tabs__item {
  height: 32px;
  line-height: 32px;
  padding: 0 24px;
  border: 1px solid #e8e8f5;
  margin-left: -1px;
  transition: none;
}

.target_box /deep/ .el-tabs--card > .el-tabs__header {
  border: 0;
}

.target_box /deep/ .el-tabs--card > .el-tabs__header .el-tabs__item.is-active {
  background: #4c7ae3;
  color: #fff;
  border: 1px solid #4c7ae3;
}

.target_box /deep/ .el-tabs--card > .el-tabs__header .el-tabs__nav {
  border-left: 1px solid #e8e8f5;
  border-top: none;
  border-right: none;
}

.linkafter {
  display: inline-block;
  border-right: 1px solid #e8e8f5;
  padding-right: 10px !important;
  height: 14px;
  line-height: 16px;
  padding-left: 10px;
}

.linkbefore {
  display: inline-block;
  border-left: 1px solid #e8e8f5;
  padding-left: 10px !important;
  height: 14px;
  line-height: 16px;
  padding-left: 10px;
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
.bugotherinfo {
  margin-top: 32px;

  .part_title {
    margin-bottom: 8px;
  }

  .content {
    background: rgba(255, 255, 255, 1);
    border-radius: 2px;
    border: 1px solid rgba(232, 232, 245, 1);
    padding: 12px 16px;
    color: rgba(72, 72, 102, 0.64);
    font-size: 13px;
  }
}

.message > div {
  //   margin-bottom: 24px;
  background: #f7f7fb;
  border-radius: 4px;
  border: 1px solid #e8e8f5;
  padding: 16px;
  box-sizing: border-box;
}

.requestpack > div {
  background: #fff !important;
  padding: 0 !important;
}

.packbtn {
  height: 65px;
  box-sizing: border-box;
  padding: 16px;
  text-align: left;
  padding-left: 0;
}
.message > label {
  display: inline-block;
  width: 80px;
  text-align: center;
  height: 26px;
  line-height: 26px;
  color: #fff;
  background-color: #4c7ae3;
  font-weight: bold;
  font-size: 12px;
}

.message > div {
  height: 253px;
  overflow-y: auto;
}

.delButton_popper {
  padding: 16px !important;

  .el-button--mini {
    padding: 5px 10px;
    border-radius: 2px;
  }
}

.delText {
  margin-bottom: 16px;
  color: rgba(72, 72, 102, 0.64);

  i {
    color: #f9b640;
    margin-right: 10px;
  }
}

.filetree {
  padding: 16px;
  float: left;
  width: 310px;
  border: 1px solid #e8e8f5;
  box-sizing: border-box;
  border-radius: 4px;

  // display: flex;
  // align-items: stretch;
  /deep/ .el-button--small {
    padding: 9px;
  }

  /deep/ .el-button--small.is-round {
    padding: 9px;
  }

  .el-button + .el-button {
    margin-left: 8px;
  }
}

.filelist {
  float: left;
  width: calc(100% - 334px);
  margin-left: 24px;
}

.equipmentinfo {
  margin-bottom: 24px;
}

.targetNofind {
  display: inline-block;
  background: rgba(255, 255, 255, 0.12);
  border-radius: 2px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.7);
  height: 32px;
  line-height: 32px;
  padding: 0 17px;

  i {
    margin-right: 16px;
  }
} 
</style>
  <script>
import $ from 'jquery'
import * as d3 from 'd3'
import dagreD3 from 'dagre-d3'
import { task } from '@/api/task.js'
import scene from '@/api/scene.js'
import  Seepagepath  from "./components/seepagepath.vue";
export default ({
  name: 'targetdetail',
  components: {
    Seepagepath
  },
  data () {
    var validatorport  = (rule, value, callback) => {
            if (!value) {
                return callback(new Error("请输入端口"));
            } else {
                if (Number(value)) {
                    if(Number(value) >=0 && Number(value) <= 65535){
                        callback();
                    }else{
                        return callback(new Error('端口格式不正确'))
                    } 
                } else {
                    return callback(new Error('端口格式不正确'))
                }
            }
        };
    return {  
		loading: false,
      isSimplifyPath: false,
      fullscreen: false,
      yzloading: false,
      target_result: '',
      verify_result: [],
      pocname: '',
      target_result_id: '',
      task_Id: this.$route.query.taskId,
      task_name: this.$route.query.taskname,
      target_id: this.$route.query.id,
      target_name: this.$route.query.name,
      task_type: this.$route.query.tasktype,
      risk_level: this.$route.query.risk_level,
      tabs: this.$route.query.tabs,
      treeloading: true,
      show: false,
      activeName: 'tabs1',
      target_quality: '',
      target_delay: '',
      local_quality: '',
      local_delay: '',
      user_id: 0,
      menuflag: true, 
      flag: 0,
      currentPage: 1,
      total: 0,
      pageSize: 10, 
      logtableData: [],
      targetpage_size: 50,
      timer1: null,
      timer2: null,
      statustimer: null,
      logpage: 1,
      drawer_attack: false,
      drawer_style: '',
      gridData: [],
	    attackform: {
          task_id: '',
          attack_type: 1, 
          target: '',
          port: '',  
          auth_type:1,
          scheme:1,
          service:'',
          url:'',
          ip:'',
        },
      rules: {
        attack_type: [
          { required: true, message: '请选择攻击面类型', trigger: 'change' },
        ],
        target: [
          { required: true, message: '请输入目标', trigger: 'blur' },
        ],
        port: [
          { required: true, message: '请输入端口', trigger: 'blur' },
          {   required: true, validator: validatorport,trigger: "blur"}
        ],
         
      },
      attack_type_list: [],
      path_type_list: [{ id: 1, name: '登录入口' }, { id: 2, name: '文件上传' }],
      auth_type_list: [{ id: 1, name: 'cookie' }, { id: 2, name: 'headers' }],
      scheme_list: [{ value: 1, label: "HTTP" }, { value: 2, label: "HTTPS" }],
      services_list: [],
      looptableData: [],
      looptotalpage: 0,
      looppageSize: 10,
      loopcurrentpage: 1,
      loopformData: {
        page: 1,
        loopsearch: '',
      },
      loopexecute_dialog_title: '',
      loopexecute_dialog_desc: '',
      testform: {
          task_id: '',
          attack_type: 4, 
          path: '',
          vul_id: '', 
          pocname:'',
        },
      looprules: { 
        path: [
          { required: true, message: '路径不能为空', trigger: 'blur' },
         
        ],
      }, 
    }
  },
  created: function () {
    this.$store.state.activefirstMenu = "/task"
	this.user_id = Number(this.$commonjs.decryptCBC(localStorage.getItem('user_id'), this.$commonjs.myKey));

    this.task_Id = this.$route.query.taskId
    this.task_name = this.$route.query.taskname
    this.target_id = this.$route.query.id
    this.target_name = this.$route.query.name
    this.tabs = this.$route.query.tabs
    this.risk_level = this.$route.query.risk_level


  },

  watch: {
    '$route': function (to, from) {
      this.flag = 0
      this.task_Id = this.$route.query.taskId
      this.task_name = this.$route.query.taskname
      this.target_id = this.$route.query.id
      this.target_name = this.$route.query.name
      this.tabs = this.$route.query.tabs
      this.risk_level = this.$route.query.risk_level

      if (this.$route.query.targetTabs) {
        this.activeName = this.$route.query.targetTabs
      }
      if (this.activeName == 'tabs1') {

      }

    },

  },
  mounted: function () { 
    var _tab = !(localStorage.getItem('targetTab')) ? 'tabs1' : localStorage.getItem('targetTab')
    this.activeName = _tab

    this.total = 0 
    this.getData(this.activeName)   
  },
  beforeDestroy () {//页面销毁清除定时器
    this.timer1 ? clearInterval(this.timer1) : null
    this.timer2 ? clearInterval(this.timer2) : null
    this.statustimer ? clearInterval(this.statustimer) : null
  },
  destroyed () {
  },
  methods: {
   
    async getEnum(){
            const res = await task.taskEnum();
            if(res.code == 200){
                this.attack_type_list = res.data.attackFaceType;
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                })
            }
        }, 
    handleClick (tab) {
      this.attackactiveName = 'tabs1'
      this.getData(tab.name) 
    },
    // TODO:
    async getData (tabname) { 
    },

    closeDialog () {
      $('#node_pop').removeClass('openAnimate').addClass('closeAnimate')
    },
    closeDialogattack () {
      $('#addattack').removeClass('openAnimate').addClass('closeAnimate')
    },
    closeDialogloop () {
      $('#vulnlistdialog').removeClass('openAnimate').removeClass('openAnimate2').addClass('closeAnimate')
      $('#vulnlistdialog').removeClass('closeAnimate2')
      $('#addloopexecute').removeClass('openAnimate').addClass('closeAnimate')
    },
    AddAttack () { //新增攻击面 
		this.getEnum()
      this.attackform.target = this.target_name
      this.attackform.task_id = this.task_Id

      $('#addattack').removeClass('closeAnimate').addClass('openAnimate')
    },
    submitattackForm () {
        this.$refs.attack.validate(async (valid) => {
          if (valid) {
            let params = {

            }
            // 端口
            if(this.attackform.attack_type == 1){
              params.port = Number(this.attackform.port);
            }
            // 敏感路径
            if(this.attackform.attack_type == 2){
                if(this.attackform.path_type == 1){ //登录入口
                  params.type = 'login'; 
                }
                if(this.attackform.path_type == 1){ //文件上传
                  params.type = 'upload';
                }
                params.url=this.attackform.url;
            }
            // 登录凭证
            if(this.attackform.attack_type == 3){
              params.ip=this.attackform.ip;
              params.port = this.attackform.port;
              params.scheme = this.attackform.scheme;
              if(this.attackform.auth_type == 1){ //cookie
                params.type = 'cookie';
                params.value=this.attackform.auth_value;
              }
              if(this.attackform.auth_type == 2){
                params.type = 'header';
                params.value=this.attackform.auth_value;
              }
            } 
            const res = await task.addattackface({
                taskId:Number(this.task_Id),
                attackFaceType:this.attackform.attack_type,
                target:this.attackform.target, 
                params:JSON.stringify(params),
                userId:this.user_id,
            })
            if (res.code == 200) {
              this.$message({
                message: "添加攻击面成功",
                type: 'success'
              })
              $('#addattack').removeClass('openAnimate').addClass('closeAnimate')
              this.attackform.port = '';
            } else {
              this.$message({
                message: res.msg,
                type: 'error'
              })
            }
          }
        })
      },
    AddLoop () { //新增漏洞执行 
		this.getloopTablelist()
      this.testform.task_id = this.task_Id
      this.testform.target = this.target_name
      $('#vulnlistdialog').removeClass('closeAnimate').removeClass('closeAnimate2').removeClass('openAnimate2').addClass('openAnimate')
    },
    submitFormloop () {   
        this.$refs.testform.validate(async (valid) => { 
          if (valid) {
            this.loading = true;
            const res = await task.addvul({
                taskId:Number(this.task_Id),
                rootUrl:this.testform.path,
                pocname:this.testform.pocname,
                userId:this.user_id,
            })
            if (res.code == 200) {
				        this.loading = false;
              this.$message({
                message: '新增漏洞执行测试成功',
                type: 'success'
              })
              $('#addloopexecute').removeClass('openAnimate').addClass('closeAnimate')
              $('#vulnlistdialog').removeClass('openAnimate2').addClass('closeAnimate2')
              
              this.testform.path = '';

            } else {
              this.loading = false;
              this.$message({
                message: res.msg,
                type: 'error'
              })
  
            }
          }
        })
  
      },
    handlloopesearch () {
      this.loopformData.page = 1
      this.getloopTablelist()
      this.loopcurrentpage = 1
    },
    handleRowClick (val) {
      	this.$refs.table.clearSelection()
      	this.$refs.table.toggleRowSelection(val)

      	this.currentRow = val
        this.loopexecute_dialog_title = this.currentRow.name
        // this.loopexecute_dialog_desc = this.currentRow.vul_description
        this.testform.vul_id = this.currentRow.id
        this.testform.pocname = val.pocname;

		$('#vulnlistdialog').addClass('openAnimate2').removeClass('closeAnimate').removeClass('openAnimate').removeClass('closeAnimate2')
		$('#addloopexecute').removeClass('closeAnimate').addClass('openAnimate')
 
      	this.testform.path = '' ; 

    },
    async getloopTablelist () {
        const res = await scene.getVulnlist({
          page: this.loopformData.page,
          size: this.looppageSize,
          search: this.loopformData.loopsearch
        })
        if(res.code == 200){
            this.looptableData = res.data.list;
            this.looptotalpage = res.data.total; 
        }
        else{
            this.$message({
                message: res.msg,
                type: "error",
            });
        }
  
      },
    handleloopcurrentchange (t) {
      this.loopformData.page = t
      this.getloopTablelist()
      this.loopcurrentpage = t
    },
    handleloopSizeChange (t) {
      this.loopformData.page = 1
      this.looppageSize = t
      this.getloopTablelist()
    },
  }
})

  </script>







