
<template>
   <div>
        <div class="main-title" @click="$refs.fileInput.click()" >字典库</div>
        <input ref="fileInput" style="display:none" type="file" id="data2" @change="upload()">
        <div class="tollboxlist context_box_bg">
            <div class="search-box">
              
                <!-- selectPlaceholder="字典类型"
                :selectList="dictionariesType" 原来的下拉框-->
                <Operation 
                    one="新建" 
                    two="删除"
                    three="搜索"
                    four="重置"
                    :delList="delList"
                    @handleOneClick="handleOperationUpdate"
                    @handleTwoClick="handleOperationDel"
                    @handleThreeClick="handleOperationSearch"
                    @handleFourClick="handleOperationReset"
                    @handleSelectChange="handleSelectChange"
                ></Operation>
                
  <el-table :data="tableData" style="width: 98%;align: center" @cell-click="tsetClick" @header-click="headerTest">
    <el-table-column :label="item.propName" :property="item.prop" v-for="item in tableColumnList" :key="item.prop" align="center">
      <template slot-scope="scope">
        <span @click="getIndex(scope.$index)">{{scope.row[scope.column.property]}}</span>
      </template>
    </el-table-column>
    <el-table-column label="." width="20" prop="addTableHeaderName" align="center"/>
  </el-table>
    <el-dialog :visible.sync="dialogForHeader" title="修改项目流程名称" width="800">
      <el-form ref="form" :model="tableHeader" label-width="80px">
        <el-form-item label="表头名称">
          <el-input v-model="tableHeader.tableHeaderName" placeholder="请输入表头名称" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitForm">确 定</el-button>
        <el-button @click="cancel">取 消</el-button>
      </div>
    </el-dialog>
 
    <el-dialog :visible.sync="dialogForTable" title="修改项目流程内容" width="800">
      <el-form ref="form" :model="tableCell" label-width="120px">
        <el-form-item label="流程内容名称">
          <el-input v-model="tableCell.tableCellData" placeholder="流程内容名称" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitForm1">确 定</el-button>
        <el-button @click="cancel1">取 消</el-button>
      </div>
    </el-dialog>

                <el-pagination
                    :page-size="page_size" 
                    background
                    layout="total,  prev, pager, next, sizes,jumper"
                    :total="total"
                    :current-page="page"
                    @size-change="handleSizeChange"
                    @current-change="handleCurrentChange"
                    >
                </el-pagination> 
            </div>
        </div>
        <el-dialog class="zidiandia" title="新建字典" :visible.sync="createDictionaries" v-if="createDictionaries"  :before-close="cancelform" width='1184px' :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">  
                <el-button size="small" @click="handleCreateSave">保存</el-button>
                <el-button size="small" @click="handleSaveCancel">取消</el-button>
            </div>
            <Edit
                :exhibitionList="exhibitionList"
                :serviceType="serviceType"
                :dictionariesType="dictionariesType"
                :pageType = "pageType"
                ref="createEditRef"
            ></Edit>
        </el-dialog>

         <el-dialog title="详情"  :visible.sync="editDictionaries"  :before-close="cancelform" width='1184px' :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">  
                <el-button size="small" @click="handleClickEdit" v-if="isUpdate">编辑</el-button>
                <el-button size="small" v-else >编辑中</el-button>
                <el-button size="small" @click="handleEditSave" v-if="!isUpdate">保存</el-button>
                <el-button size="small" @click="handleEditCancel">取消</el-button>
            </div>
            <Edit
                v-if="editShow"
                :show="isUpdate"
                :exhibitionList="exhibitionList"
                :serviceType="serviceType"
                :dictionariesType="dictionariesType"
                :editData="editData"
                :pageType = "pageType"
                ref="cancelEditRef"
            ></Edit>
        </el-dialog>
    </div>
<!-- 6666666666666 -->

</template>
 
<script>
import Operation from "../../components/Operation";
import Edit from "../../components/edit";
  export default {
    data() {
      return {
        createDictionaries: false,
        editDictionaries: false,
        editShow: false,
        isUpdate: true,
        pageType: 1,
        page: 1,
        page_size: 10,
        
        delList:[],
        total:123,
        tableCellIndex:"",
        columnName:"",
        tableCell: { tableCellData:""},
        dialogForTable:false,
        num:7,
        tableHeader: { tableHeaderName:"",property:"" },
        dialogForHeader:false,
        // 这里为了简便我就没有调用后台接口获取数据，直接写的假数据  你要用的话可以调用后台接口获取tableColumnList，注意数据格式
        tableColumnList:
          [
          {prop: '0', propName: '编号8888'},
          {prop: '2', propName: '名字'},
          {prop: '3', propName: '保质期'},
          {prop: '4', propName: '特点1'},
          
          ],
 
        tableData: [{
          0: '2016-05-01',
          2: '王小虎1',
          3: '上海市普陀区金沙江路 1518 弄',
          4: '2016-05-02',
        
        }, {
          0: '2016-05-02',
          2: '王小虎2',
          3: '上海市普陀区金沙江路 1518 弄',
          4: '2016-05-02',
         
        }, {
          0: '2016-05-03',
          2: '王小虎3',
          3: '上海市普陀区金沙江路 1518 弄',
          4: '2016-05-02',
         
        }, {
          0: '2016-05-04',
          2: '王小虎4',
          3: '上海市普陀区金沙江路 1518 弄',
          4: '2016-05-02',
      
        }, {
          0: '2016-05-05',
          2: '王小虎5',
          3: '上海市普陀区金沙江路 1518 弄',
          4: '2016-05-06',
        
        }, 
        {
          0: '2016-05-07',
          2: '王小虎6',
          3: '上海市普陀区金沙江路 1518 弄',
          4: '2016-05-02',
        
        },
        {
          0: '2016-05-07',
          2: '王小虎6',
          3: '上海市普陀区金沙江路 1518 弄',
          4: '2016-05-02',
        
        },
        {
          0: '2016-05-07',
          2: '王小虎6',
          3: '上海市普陀区金沙江路 1518 弄',
          4: '2016-05-02',
        
        },
        {
          0: '2016-05-07',
          2: '王小虎6',
          3: '上海市普陀区金沙江路 1518 弄',
          4: '2016-05-02',
        
        },
        {
          0: '2016-05-07',
          2: '王小虎6',
          3: '上海市普陀区金沙江路 1518 弄',
          4: '2016-05-02',
        
        }
        ],
      }
    },
    components:{
        Operation,
        Edit,
    },
    // watch: {
    //   'tableData':{
    //     handler: function(newVal, oldVal) {
    //       // TO DO
    //       console.log("newVal:", newVal);
    //       console.log("oldVal:", oldVal);
    //     },
    //     deep: true,
    //     immediate: true
    //   },
    // },
    methods: {
  
      upload(e){
        let that = this;
        console.log("点击了上传按钮");
        let da = window.event || e  // change事件获取到的数据
            var file = da.target.files[0]
            console.log(file, 'file===');
            if (true) {
              var reader = new FileReader()
              reader.readAsText(file, "utf-8")//gbk编码 还有其他方式比如gpk
              reader.onload = function (da) {
               console.log(da.target.result, 'da',Object.prototype.toString.call(da.target.result));
               //查找字符串中第三个‘[’并截取，分别保存到str1和str2中
            let str1 = "";
            let str2 = "";
                var str = da.target.result;

                // 计数器，用于记录'['的出现次数
                var count = 0;

                // 保存第三个'['字符的索引
                var startIndex = -1;

                // 遍历字符串
                for (var i = 0; i < str.length; i++) {
                  // 如果找到了第2个'['字符
                  if (str[i] === '[') {
                    // 计数器加1
                    count++;
                    // 如果计数器为3，则保存索引并跳出循环
                    if (count === 2) {
                      startIndex = i;
                      break;
                    }
                  }
                }

                // 如果找到了第三个'['字符
                if (startIndex !== -1) {
                  // 使用字符串的截取函数substring将第三个'['之后的内容截取到str1变量
                   str1 = str.substring(0,startIndex);

                   str2 = str.substring(startIndex, str.length );

                  console.log("str1: " + str1);
                  console.log("str2: " + str2);
                } else {
                  console.log("未找到第三个'['字符");
                }



                // console.log(str2, 'str2',Object.prototype.toString.call(str2));
                that.daoRuJson = eval('(' + str1 + ')')
                let tabelData = eval('(' + str2 + ')')
                // that.daoRuJson = eval('(' + da.target.result + ')')
                that.tableColumnList = that.daoRuJson   
                that.tableData =  tabelData 
                console.log("表头数据",that.tableColumnList);
                console.log("表格数据",that.tableData);
                // console.log(that.daoRuJson, 'that.daoRuJson',Object.prototype.toString.call(that.daoRuJson));
              }
            }
      },
      handleEdit(index, row) {
        console.log(index, row);
      },
      handleDelete(index, row) {
        console.log(index, row);
      },
      getIndex(index){
        console.log("index",index);
        this.tableCellIndex = "";
        this.tableCellIndex = index;
      },
      tsetClick(row, column, cell, event){
        this.dialogForTable = true;
        this.columnName = "";
        this.columnName = column.property;
        this.tableCell.tableCellData = "";
        this.tableCell.tableCellData = row[this.columnName];
        console.log("这里是表格内容单击事件",row);
        console.log("这里是表格内容单击事件",column)
        // console.log("这里是表格双击事件",cell)
        // console.log("这里是表格双击事件",event)
      },
      //添加表头，修改表头
      headerTest(val){
        if(val.property == "addTableHeaderName"){
          console.log("这里是表格tou部点击,添加头部事件",val)
          this.tableColumnList.push({prop: this.num.toString(), propName: '点击编辑项目流程名称'})
          for (let i = 0; i < this.tableData.length; i++) {
            //this.tableData[i][parseInt(this.num)] = "请添加内容";
            this.$set(this.tableData[i],[parseInt(this.num)],"请添加内容");
          }
          this.num = this.num+1;
        }else {
          console.log("这里是表格tou部点击，修改头部事件",val)
          this.tableHeader.tableHeaderName = "";
          this.tableHeader.property = "";
          this.tableHeader.tableHeaderName = val.label;
          this.tableHeader.property = val.property;
          console.log(this.tableHeader.tableHeaderName)
          this.dialogForHeader = true;
        }
      },
      //表头编辑提交
      submitForm(){
        console.log("点击提交按钮");
        for (let i = 0; i < this.tableColumnList.length; i++) {
          console.log("this.tableHeader.property",this.tableHeader.property)
          console.log("this.tableColumnList[i].prop",this.tableColumnList[i].prop)
          if(this.tableColumnList[i].prop === this.tableHeader.property){
            this.tableColumnList[i].propName = this.tableHeader.tableHeaderName;
          }
        }
        console.log(this.tableColumnList)
        this.dialogForHeader = false;
      },
      //表格内容编辑提交
      submitForm1(){
        //console.log("点击提交按钮");
        if(this.tableCellIndex === "" || this.tableCellIndex === null){
          alert("亲，请精准点击表格中的字进行修改！")
        }else {
            console.log("this.columnName",this.columnName)
            console.log("this.tableData[this.tableCellIndex][this.columnName]",this.tableData[this.tableCellIndex][this.columnName])
            console.log("this.tableData[this.tableCellIndex]",this.tableData[this.tableCellIndex])
          this.tableData[this.tableCellIndex][this.columnName] = this.tableCell.tableCellData;
          console.log(this.tableData);
          this.rush();
        }
        this.dialogForTable = false;
      },
      //强制刷新数据
      rush(){
        this.$set(this.tableData);
      },
      //取消表头编辑
      cancel(){
        //console.log("点击取消按钮");
        this.dialogForHeader = false;
      },
      //取消表格内容编辑
      cancel1(){
        //console.log("点击取消按钮");
        this.dialogForTable = false;
      },
          cancelform(){},
      handleOperationDel(){},
      handleOperationUpdate(){},
      handleOperationSearch(){},
      handleOperationReset(){},
      handleSelectChange(){},
      handleSizeChange(){},
      handleCurrentChange(){},
      handleCreateSave(){},
      handleSaveCancel(){},
      handleClickEdit(){},
      handleEditSave(){},
      handleEditCancel(){},
    },
  }
</script>
<style lang="less" scoped>
.tollboxlist {
    padding: 24px;
    background: #fff;
    min-height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
    border-radius: 4px;
}
::v-deep .el-dialog{
    background: #F7F7FB;
    border-radius: 4px;
}
.el-table .el-link::after{
    content: "";
    padding: 0;
    width: 0;
}
::v-deep .el-popconfirm__main{
    padding: 8px 0;
}
.el-popper{
    margin: 0;
}
::v-deep .popconfirm{
    margin: 220px !important;
}
// .el-link.el-link--danger {
//     color: #F35F28;
// }
/deep/ .zidiandia .el-form-item__label{
    text-align: left!important;;
}
/deep/ .zidiandia .el-form-item__label:after{
    left:-6px;
}
/deep/ .el-input__inner{
    height: 32px !important;
    line-height: 32px  !important;
}
/deep/ .el-input__suffix{
    top: -1px;
}
/deep/ .myTable{
        thead {
            .spanActive{
                color:#4C7AE3;
                i{
                    color: #4C7AE3;
                }
            }
            .cursorPointer{
                cursor: pointer;
            }
            .cell{
                line-height: 15px;
                >span{
                    position: absolute;
                }
            }
            .iconfont{
                color:rgba(72,72,102,0.32);
                margin-left:5px;
            }
            .el-select{
                height: 0;
                visibility: hidden;
                .el-input, .el-input__inner{
                    height: 0!important;
                }
            }
        }
    }
    .linkafter  {
        display: inline-block;
        border-right: 1px solid #E8E8F5;
        padding-right: 10px !important;
        height: 14px;
        line-height: 16px;
        padding-left: 0px;
    }
    .linkafter2{
        display: inline-block;
        border-right: 0;
        padding-right: 10px !important;
        height: 14px;
        line-height: 16px;
        padding-left: 10px;
    }

</style>