<template>
  <div id="app"> 
    <router-view /> 
  </div>
</template>

<script> 

export default {
  name: 'App',
  data(){
    return{}
  },
}
</script>
<style  scoped>
@import '@/assets/less/style.less';
</style>
<style lang="less" scoped>   
:deep(.el-table__expanded-cell) {
  border-bottom: 0 !important;
}
:deep(.el-table__empty-text)  {
  width: 100%;
  padding-left: 32px;
  text-align: left;
}
:deep( .el-button--small) {
  border-radius: 2px;
  font-size: 14px;
}
:deep( .el-input__inner) {
  // border: 1px solid #e8e8f5;
  border-radius: 2px;
  color: rgba(72, 72, 102, 0.64);
  font-weight: 500;
  transition: border-color 0.2s cubic-bezier(76, 122, 227, 1) !important;
  caret-color: @theme-BaseColor;
  // height: 31px !important;
  line-height: 31px !important;
}
:deep(.el-input:not(.is-disabled) .el-input__inner:hover ){
  border-color: @theme-BaseColor !important;
}
:deep( .el-textarea__inner) {
  border: 1px solid #e8e8f5;
  color: rgba(72, 72, 102, 0.64);
  border-radius: 2px;
  caret-color: @theme-BaseColor;
}
:deep( .el-textarea__inner:hover ){
  border-color: @theme-BaseColor !important;
}
:deep( .el-dialog__body) {
  color: rgba(72, 72, 102, 0.87);
}
:deep( .el-date-editor .el-range-input) {
  color: rgba(72, 72, 102, 0.64);
  // font-weight: 100 !important;
}
:deep(.el-range-editor.el-input__inner) {
  padding: 2px 10px;
}
:deep( .el-table th ){
  font-weight: 500;
  font-size: 13px;
  color: rgba(72, 72, 102, 0.87);
  background: #f7f7fb !important;
  padding: 12px 0;
}
:deep( .el-table--enable-row-hover .el-table__body tr:hover > td ){
  background: #f7f7fb !important;
}
:deep( .el-table td ){
  font-weight: 500;
  font-size: 13px;
  padding: 12px 0;
  color: rgba(72, 72, 102, 0.64);
}
:deep( .el-tooltip .el-link.el-link--default) {
  // color: rgba(72,72,102,.64) ;
  width: 100%;
  span {
    width: 100%;
    display: inline-block;
    white-space: nowrap;
    word-break: break-all;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}
:deep( .el-link.el-link--default ){
  color: var(--baseFontColor) ;
  line-height: 14px;
}
:deep( .el-table-column--selection .cell ){
  padding-left: 24px !important;
}
:deep( .el-table .cell ){
  line-height: 16px;
}

:deep( .btn-transparentBlue) {
  border-color: @theme-BaseColor;
  color: @theme-BaseColor;
}
:deep( .el-checkbox__input.is-checked + .el-checkbox__label ){
  color: @theme-BaseColor;
} 
:deep( .el-link ){
  &.el-link--primary,
  &.el-link--default:after,
  &.el-link--primary.is-underline:hover:after,
  &.el-link--primary:after {
    border-color: @theme-BaseColor;
  }

}
:deep( .el-checkbox__inner:hover ){
  /* border-color: #DCDFE6!important; */
  border-color: @theme-BaseColor !important;
}
:deep( .el-radio__inner:hover ){
  border-color: @theme-BaseColor !important;
}
/* //修改选中框默认颜色 */
:deep( .el-checkbox__input.is-focus .el-checkbox__inner ){
  border-color: @theme-BaseColor !important;
}
// deep() .el-checkbox__input.is-checked .el-checkbox__inner,
// deep() .el-checkbox__input.is-indeterminate .el-checkbox__inner {
//   background-color: @theme-BaseColor !important;
//   border-color: @theme-BaseColor !important;
// }
:deep( .el-checkbox__input){
  &.is-checked .el-checkbox__inner,
  &.is-indeterminate .el-checkbox__inner{
    background-color: @theme-BaseColor !important;
    border-color: @theme-BaseColor !important;
  }
}
/* 修改页码条里的箭头位置 */
:deep( .el-input__prefix){
  top: -1px;
}
:deep( .el-input__suffix) {
  top: -1px;
}
:deep( .el-tabs__active-bar ){
  background-color: @theme-BaseColor;
}
// 
:deep(.el-pagination.is-background) {
   .btn-next.is-active, 
  .btn-prev.is-active, 
  .el-pager li.is-active{
    background-color: @theme-BaseColor;
  }
}

// 按钮
:deep(.el-button){
 
  padding: 6px 12px;
  border-radius: 3px;
  vertical-align: middle;
  height: 31px; 
  &.el-button--primary{
    background-color: @theme-BaseColor;
    border-color: @theme-BaseColor;
  }
  &.el-button--primary:hover{
    background-color: @theme-BaseColor;
    border-color: @theme-BaseColor;
  }
  &.el-button--primary:focus{
    background-color: @theme-BaseColor;
    border-color: @theme-BaseColor;
  }
}
:deep(.el-button.el-button--warning.is-disabled), 
:deep(.el-button.el-button--warning.is-disabled:focus), 
:deep(.el-button.el-button--warning.is-disabled:hover){
  color: rgba(243,95,40,.5);
  background-color: #fff;
  border-color: rgba(243,95,40,.5);
}
.el-button--primary{

}
:deep(.el-button+.el-button){
  margin-left: 8px;
} 
// 新版表格样式
/deep/ .el-table::before{
  height: 0;
}
/deep/ .el-table td.el-table__cell, 
/deep/.el-table th.el-table__cell.is-leaf{
  border-bottom: 0;
}
/deep/ .el-table__fixed-right::before, 
/deep/ .el-table__fixed::before{
  height: 0;
}
</style>
