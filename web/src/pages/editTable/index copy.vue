
<template>
  <div>
    <el-button @click="$refs.elxEditable.insert({name: 'new1'})">新增</el-button>
    <el-button @click="$refs.elxEditable.removeSelecteds()">删除选中</el-button>
 
    <elx-editable header ref="elxEditable" :data.sync="tableData">
      <elx-editable-column type="selection" width="55"></elx-editable-column>
      <elx-editable-column type="index" width="55"></elx-editable-column>
      <elx-editable-column prop="name" label="只读"></elx-editable-column>
      <elx-editable-column prop="sex" label="下拉" :edit-render="{name: 'ElSelect', options: sexList}"></elx-editable-column>
      <elx-editable-column prop="age" label="数值" :edit-render="{name: 'ElInputNumber'}"></elx-editable-column>
      <elx-editable-column prop="date" label="日期" :edit-render="{name: 'ElDatePicker', props: {type: 'date', format: 'yyyy-MM-dd'}}"></elx-editable-column>
      <elx-editable-column prop="flag" label="开关" :edit-render="{name: 'ElSwitch', type: 'visible'}"></elx-editable-column>
      <elx-editable-column prop="remark" label="文本" :edit-render="{name: 'ElInput'}">
        <template slot="header">
          <!-- <span>{{scope.row.remark}}</span> -->
        </template>
      </elx-editable-column>
    </elx-editable>
  </div>
</template>
 
<script>
export default {
  data () {
    return {
      tableData: [{
        date: 1551322088449,
        name: '小徐',
        sex: '1',
        age: '26',
        flag: false,
        remark: '备注'
      }],
      sexList: [
        {
          'label': '男',
          'value': '1'
        },
        {
          'label': '女',
          'value': '0'
        }
      ]
    }
  }
}
</script>