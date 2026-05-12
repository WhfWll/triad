/* 
  资产树组件


 */
<template>
  <div style="height: calc(100vh - 200px);overflow:auto">
    <!-- 添加按钮组 -->
    <!-- 添加按钮组下方添加进度条 -->
    <div class="button-group" style="margin-bottom: 15px;">
      <el-button 
        size="small" 
        type="primary" 
        @click="$emit('updateAssetGroup', '')">
        新增资产组
      </el-button>
      <el-button 
        size="small" 
        type="primary" 
        @click="$emit('updateAsset', '')">
        新增资产
      </el-button>
      <el-button 
        size="small" 
        type="primary" 
        @click="handleBatchDelete"
        :disabled="!hasCheckedNodes">
        批量删除
      </el-button>
      <el-button 
        size="small" 
        type="primary" 
        @click="handleBatchPenetrationTest"
        :disabled="!hasCheckedNodes">
        渗透测试
      </el-button>
    </div>
    <div class="progress-section" v-if="showProgress">
      <div class="progress-container">
        <span class="progress-title">资产发现进度</span>
        <el-progress :percentage="progressPercent"  style="width: 100%; margin-left: 15px;">
        </el-progress>
      </div>
    </div>
    <el-tree 
      ref="tree"
      :data="data"
      show-checkbox
      node-key="id"
      default-expand-all 
      :props="defaultProps"
      :expand-on-click-node="false"
      @check-change="handleCheckChange">
      <span class="custom-tree-node" slot-scope="{ node, data }">
        <span @click="sss(node, data)">{{ node.label }}</span>
        <span class="tree-node-actions" >
          <el-button
          icon="el-icon-edit"
            type="text"
            size="mini"
            @click="() => Update(data)">
          </el-button>
          <el-button
            icon="el-icon-delete"
            type="text"
            size="mini"
            @click="() => remove(node, data)">
            
          </el-button>
          <el-button
            icon="el-icon-cpu"
            type="text"
            size="mini"
            @click="() => penetrationTask(node, data)">
            
          </el-button>
        </span>
      </span>
    </el-tree>  
  </div>
</template>

<script>
import { traffic } from '@/api/assetManagement.js'
export default {
  data() { 
    return {
      data:[],
      newdata:[],
      defaultProps: {
        children: 'children',
        label: 'label'
      },
      allids:[],
      groupips:[],
      assetips:[],
      hasCheckedNodes: false,
      showProgress: false,
      progressPercent: 0,
      progressStatus: '',
      progressTimer: null,
    }
  },
  methods: {
    sss(node, data) {
      let id = data.id.split('_')[1];
      this.$store.commit('setgroupID', [id, data.type]);


      console.log('点击tree');

      this.$emit("clicknode");

    },
    transformData(data, level = 1) {
      if (!Array.isArray(data)) {
        return []; // 返回空数组作为兜底
      }
      if (!data){
        return 
      }
      return data.map(item => ({
        id: typeof item === 'object' && item.type && item.id ? `${item.type}_${item.id}` : `unknown_${id++}`,
        label: typeof item === 'object' && item.name ? item.name : '未命名',
        type: typeof item === 'object' && item.type ? item.type : '未知类型',
        level: level,
        children: item.items && Array.isArray(item.items) ? this.transformData(item.items, level + 1) : []
      }));
    },
    async handleBatchDelete() {
      const checkedNodes = this.$refs.tree.getCheckedNodes();
      const groupIds = [];
      const assetIds = [];
      
      checkedNodes.forEach(node => {
        const id = node.id.split('_')[1];
        if (node.type === 1) {
          groupIds.push(id);
        } else {
          assetIds.push(id);
        }
      });
    
      await this.handleassetdel(groupIds, assetIds);
    },
    async handleBatchPenetrationTest() {
      const checkedNodes = this.$refs.tree.getCheckedNodes();
      const ips = checkedNodes
        .filter(node => node.type === 2)
        .map(node => node.label);
    
      if (ips.length > 0) {
        this.$router.push({
          path: `/createtask`,
          query: { 
            flag: 1,  
            type: 3,
          }
        });
        localStorage.setItem("checkedasset", ips.join(','));
      } else {
        this.$message.warning('请选择资产节点');
      }
    },
    async getZCList() {
      const res = await traffic.trafficlistinfodel()
      if (res.code == 200) { 
        this.data = this.transformData(res.data.list); 
        // this.getTreeAssetID(this.data) 
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    Update(data) { //编辑 
      let ID = data.id.split('_')[1]; 
      if(data.type==1){ //资产组   
        this.$emit("updateAssetGroup",ID);
      }else{
        this.$emit("updateAsset",ID);
      }
    },
    async remove(node, data) { 
      let ID = data.id.split('_')[1]; 
      this.groupips = [];
      this.assetips=[];
      if(data.type==1){ //资产组   
        const res = await traffic.trafficVulnInfo({
          page: 1,
          size: 10000,
          search: '',
          groupID:ID, 
        }) 
        if(res.code ==200){
          let list = res.data.subAssetGroupInfo;
          if(list&&list.length>0){
            this.groupips= list.map(item => item.groupID); 
          } 
        } 
        this.groupips.push(ID); //当前资产组ID
    
        const res1 = await traffic.trafficvulndel({
          page: 1,
          size: 10000,
          search: '',
          groupID:ID, 
        })
        if (res1.code == 200) {  
          let list = res1.data.assetsInfo;
          if(list&&list.length>0){
            this.assetips = list.map(item => item.assetID); 
          } 
         
        } 
        this.handleassetdel(this.groupips,this.assetips);
      }else{
        //资产
        this.handleassetdel([],[ID]);
      }
    }, 
    async handleassetdel(group,asset){  
      const res = await traffic.assetdelete({
        groupIds:group.length==0?'':group.join(','), 
        assetIds:asset.length==0?'':asset.join(','),
      });
      if(res.code ==200){
        this.$message({
          message: '批量删除资产组资产成功',
          type: 'success'
        })
        this.replaceData();
      }else{
        this.$message({
          message: '批量删除资产组资产失败',
          type: 'erroe'
        })
      }
    },
    async penetrationTask(node, data){ //渗透任务 
      let groupID = data.id.split('_')[1]; 
      if(data.type==1){
        const res = await traffic.trafficvulndel({
          page: 1,
          size: 10000,
          search: '',
          groupID:groupID, 
        })
        if (res.code == 200) {  
          let ips = res.data.assetsInfo.map(item => item.ip); 
          this.$router.push({
              path: `/createtask`,
              query: { 
                  flag: 1,  
                  type:3, //资产组 渗透任务
              }
          });
          localStorage.setItem("checkedasset",ips.join(','));
        } else {
          this.$message({
            message: '获取资产数据失败',
            type: "error"
          })
        }
      }
      else{
        this.$router.push({
            path: `/createtask`,
            query: { 
                flag: 1,  
                type:3, //资产组 渗透任务
            }
        });
        localStorage.setItem("checkedasset", data.label);
      }
      
    },
    checkedNodes(){ 
      return this.$refs.tree.getCheckedNodes();
    },
    cancelSelected(){ 
      this.$refs.tree.setCheckedKeys([]); 
    },
    replaceData(){ //操作后更新数据
      this.getZCList();
      // this.$refs.tree.refreshTree()
    },
    selectedAsset(){ //选择资产
      this.$refs.tree.setCheckedKeys(this.allids);
    },
    customAction(node, data) {
      // 实现自定义操作的逻辑
    },
    handleCheckChange() {
      this.hasCheckedNodes = this.$refs.tree && this.$refs.tree.getCheckedNodes().length > 0;
      // 这里可以添加任何需要在节点勾选状态变化时执行的逻辑
    },
    async getTaskProgress(findAsset=false) {
      try {
        const res = await traffic.assetDetectProgress();
        if ((res.code === 200 && res.data.progress !="100.00") || findAsset){
          this.showProgress = true;
          this.progressPercent = parseFloat(res.data.progress) || 0;
          // 只在进度为100时设置status为success，其他情况不设置status
          // this.progressStatus = this.progressPercent === 100 ? 'success' : '';
          
          if (this.progressPercent < 100) {
            this.progressTimer = setTimeout(() => {
              this.getTaskProgress();
            }, 3000);
          }
        }else if (res.code === 200 && res.data.progress === "100.00"){
          this.showProgress = false;
          clearTimeout(this.progressTimer);
          this.progressPercent = 0;
          // this.progressStatus = ''; // 重置时设置为空字符串
          this.$message({
            message: '资产检测完成',
            type: "success"
          })
          this.replaceData();
        }
      } catch (error) {
        console.error('获取任务进度失败:', error);
      }
    },
  },
  created() {
    this.getZCList();
    this.getTaskProgress(); // 添加进度获取
  },
  beforeDestroy() {
    // 清除定时器
    if (this.progressTimer) {
      clearTimeout(this.progressTimer);
    }
  }
}
</script>


<style scoped>
/* 基础样式 */
.el-tree {
  font-size: 14px;
  color: #333;
}
.button-group {
  display: flex;
  gap: 10px;
}

.button-group .el-button {
  margin-left: 0;
}
/* 自定义树节点样式 */
.custom-tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 10px;
  transition: color 0.3s ease;
}
.progress-section {
  margin-bottom: 20px;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 4px;
}

.progress-container {
  display: flex;
  align-items: center;
}

.progress-title {
  white-space: nowrap;
  color: #606266;
  font-size: 14px;
}
.progress-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
  color: #606266;
  font-size: 14px;
}
/* 节点文本鼠标悬停效果 */
.custom-tree-node:hover .tree-node-label {
  color: #409EFF; 
}

/* 节点操作按钮样式 */
.tree-node-actions {
  display: flex;
  align-items: center;
  opacity: 0;
  transition: opacity 0.3s ease-in-out;
}

/* 悬停时显示操作按钮 */
.custom-tree-node:hover .tree-node-actions {
  opacity: 1;
}

/* 操作按钮样式 */
.tree-node-actions .el-button {
  margin-left: 8px;
  color: #C0C4CC;
  font-size: 16px;
  transition: color 0.3s;
}

/* 操作按钮悬停颜色变化 */
.tree-node-actions .el-button:hover {
  color: #409EFF;
}

/* 树节点文本样式 */
.tree-node-label {
  transition: color 0.3s ease;
}

</style>


