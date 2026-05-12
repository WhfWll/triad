 <!--   资产树组件 -->
 <template> 
  <div v-show="isshow" style=" padding: 0 16px;   box-sizing: border-box;"> 
    <!-- 操作按钮和搜索框 -->
    <div  class="operation-bar">  
      <el-button
         type="primary"
        size="small" 
        @click="handleDelete"
      >删除</el-button>
      <el-input
        v-model="searchQuery" 
        size="small"
        placeholder="搜索"
        style="width: 200px; margin-left:8px;height:28px;margin-right: 8px;" 
        @input="handleSearch"
      >
        <template #append> 
           <i @click="handleSearch" class="el-icon-search"></i> 
        </template>
      </el-input>
         <el-button size="small" type="primary" @click="openAddDialog" style="height:28px"> 
             <i class="el-icon-plus"></i>
        </el-button>
    </div>

    <div  style="height: calc(100vh - 200px); overflow: auto">
        <el-button size="small" @click="toggleTree">{{ isExpanded ? '折叠全部' : '展开全部' }}</el-button> 
        <el-tree
            ref="treeRef"
            :data="filteredData"
            show-checkbox
            node-key="id"
            :props="defaultProps"
            :expand-on-click-node="false" >
            <template #default="{ node, data }" >
                <div class="tree-node" @mouseenter="node.hover = true" @mouseleave="node.hover = false">
                    <div style="color: #000;position: relative;width:180px;" @click="sss(node, data)" >{{ node.label }}
                        <div style="position:absolute;top:0;left:100px;"  class="tree-node-actions"  >
                            <span
                                class="tree-node-button"
                                type="text"
                                @click="() => openEditDialog(data)" >
                                编辑
                            </span>
                            <span
                                class="tree-node-button"
                                type="text"
                                @click="() => remove(node, data)" >
                                删除
                            </span>
                        </div>
                    </div > 
                </div>
            </template>
        </el-tree>
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog width="800px" :title="dialogTitle" :visible.sync="dialogVisible"> 
        <el-form ref="formRef" :rules="rules" :model="form" label-width="120px" style="padding: 20px;">
            <el-form-item prop="name" label="资产组名称：">
                <el-input v-model="form.name" style="width: 400px;"></el-input>
            </el-form-item> 
            <el-form-item label="上级资产组：">
                <el-select
                  
                    style="width: 400px;"
                    v-model="selectedNode"
                    placeholder="请选择资产组" 
                    @visible-change="handleVisibleChange"> 
                    <el-tree
                        ref="treeRefs"
                        :data="filteredData"
                        node-key="id"
                        :props="defaultProps"
                        :expand-on-click-node="false"
                        :default-expanded-keys="defaultExpandedKeys"
                        
                        class="tree-select" >
                        <template #default="{ node, data }">
                            <el-option
                                :key="data.id"
                                :label="node.label"
                                :value="data.id"
                                @click="handleNodeClick(node, data)" >
                                <span class="tree-node">{{ node.label }}</span>
                            </el-option>
                        </template>
                    </el-tree>
                </el-select>
                <!-- <el-input
                    v-model="selectedLabel"
                    placeholder="请选择"
                    readonly
                    @focus="popoverVisible = true"
                    ref="input"
                    style="width: 400px;"
                    />
                <el-popover
                    ref="treePopover"
                    placement="bottom"
                    width="400"
                    trigger="manual"
                    v-model="popoverVisible" 
                    @hide="onPopoverHide">
                    <el-tree
                        :data="filteredData"
                        :props="defaultProps"
                        @node-click="handleNodeClick1"
                        :highlight-current="true"
                        default-expand-all
                    ></el-tree>
                </el-popover> -->


            </el-form-item>
            <el-form-item label="备注：">
                <el-input
                    type="textarea"
                    v-model="form.remark"
                    style="width: 400px;"
                ></el-input>
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="dialogVisible = false">取消</el-button>
            <el-button type="primary" @click="handleSubmit">确定</el-button>
        </template>
    </el-dialog>
  </div>
</template>

<style scoped>
 
/* 操作栏样式 */
.operation-bar {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

/* 基础样式 */
.el-tree {
  font-size: 14px;
  color: #333;
}

/* 自定义树节点样式 */
.custom-tree-node .tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 10px;
  transition: color 0.3s ease;
}

/* 节点文本鼠标悬停效果 */
.custom-tree-node:hover .tree-node-label {
  color: #409eff;
}

/* 节点操作按钮样式 */
.tree-node-actions {
  display: flex;
  align-items: center;
  opacity: 0;
  transition: opacity 0.3s ease-in-out;
}

/* 悬停时显示操作按钮 */
.tree-node:hover .tree-node-actions {
  opacity: 1;
}

/* 操作按钮样式 */
.tree-node-button {
  margin-left: 8px;
  color: #c0c4cc;
  font-size: 14px;
  cursor: pointer;
  transition: color 0.3s, transform 0.3s;
}

/* 操作按钮悬停颜色变化 */
.tree-node-button:hover {
  color: #409eff;
}

/* 树节点文本样式 */
.tree-node-label {
  transition: color 0.3s ease;
}
.el-form-item__label {
  width: 120px; /* 固定label的宽度 */
}
.el-form-item__content {
  width: calc(100% - 120px); /* 调整右侧输入框的宽度，使其与label对齐 */
}
.tree-select {
  width: 400px; /* 统一下拉树组件的宽度 */
}
</style>
<script>
import asset from '@/api/asset.js'; 
export default {
    data(){
        return{
            data:[],
            selectedNode:'',
            filteredData:[],
            defaultExpandedKeys:[],/// 存放要默认展开的节点ID
            defaultProps:{
                children: 'children',
                label: 'label',
            },
            allids:[],
            groupips:[],
            assetips:[],
            searchQuery:'',
            multipleSelection:[],
            dialogVisible:false,
            dialogTitle:'',
            fatherNodeID:'',
            fatherNodelabel:'',
            form:{
                id: '',
                name: '',
                pid: '',
                remark: ''
            },
            rules:{
                name: [{ required: true, message: '资产组名称不能为空', trigger: 'blur' }],
                remark: [{ required: true, message: '备注不能为空', trigger: 'blur' }],
                fatherNode: [{ required: false, message: '归属资产组不能为空', trigger: 'click' }]
            },
            enumOptions:[],
            isExpanded:false,
            allChecked:false,
            isshow:true,

            nodeDataMap: {} ,
            popoverVisible: false,
            selectedValue:'', 
            selectedLabel: '',
        }
    },
    created(){

    },
    mounted(){
        this.getZCList();
        localStorage.setItem('groupID','');
    },
    methods:{
        getAllNodeIds(data){
            let ids = [];
            for (let item of data) {
                ids.push(item.id);
                if (item.children) {
                    ids = ids.concat(this.getAllNodeIds(item.children));
                }
            }
            return ids;
        },
        handleNodeClick(node, data){
            console.log('资产树点击')
            this.selectedNode = data.id;
            this.selectedValue = data.id
            console.log(node)
             console.log(data)
        },
        sss(node, data){ 
            this.fatherNodelabel = node.parent.label
            let id = data.id.split('_')[1];

            localStorage.setItem('groupID', id);
            // 手动触发自定义事件
            const event = new CustomEvent('groupIDChanged', { detail: id });
            window.dispatchEvent(event);
        },
        fatherNode(node, data){
            let id = data.id.split('_')[1];
            console.log(data, '------------');
            this.fatherNodeID = id
            this.fatherNodelabel = data.label
        },
        transformData(data, level = 1){
            if (!Array.isArray(data)) {
                return []; // 返回空数组作为兜底
            }
            return data.map(item => ({
                id: typeof item === 'object' && item.type && item.id ? `${item.type}_${item.id}` : `unknown_${id++}`,
                label: typeof item === 'object' && item.name ? item.name : '未命名',
                type: typeof item === 'object' && item.type ? item.type : '未知类型',
                level: level,
                children: item.items && Array.isArray(item.items) ? this.transformData(item.items, level + 1) : [],
                hover: false, // 用于控制按钮显示
            }));
        },
        getTreeAssetID(list){
             for (var i = 0; i < list.length; i++) {
                this.allids.push(list[i].id);
                if (list[i].children && list[i].children.length > 0) {
                    this.getTreeAssetID(list[i].children);
                }
            }
        },
        async getZCList(){
             const res = await asset.trafficlistinfodel();
            if (res.code == 200) {
                this.data = this.transformData(res.data.list);
                this.filteredData = this.data;

                this.defaultExpandedKeys = this.getAllNodeIds(this.data); // 初始化默认展开的节点
                this.getTreeAssetID(this.data);
            } else { 
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }
        },
        openAddDialog(){
            // 点击新增的时候清空选择的父节点
            this.fatherNodelabel = ''
            this.dialogTitle = '新增资产组';
            Object.assign(this.form, {
                id: '',
                name: '',
                pid: '',
                remark: ''
            });
            // await fetchEnumOptions();
            this.dialogVisible = true;

            this.buildNodeMap(this.filteredData)

        },
        async openEditDialog(data){ 
            this.dialogTitle = '编辑资产组';
 
            let _id = data.id.split('_')[1];
            this.buildNodeMap(this.filteredData)
            const res = await asset.assetinfo({
                id:_id
            });
            if(res.code == 200){
                Object.assign(this.form, {
                    id: _id,
                    name: res.data.name,
                    pid: res.data.pid,
                    remark: res.data.remark
                });

                const parentNode = this.findParentNode(this.filteredData, data.id);
                console.log(parentNode)
                if (parentNode) {
                    this.selectedNode = parentNode.id;
                    this.selectedValue = parentNode.id;
                    this.selectedLabel = parentNode.label;
                } else {
                    this.selectedNode = '无'; // 如果没有父节点，设置为当前节点
                    this.selectedValue = '';
                     this.selectedLabel = '无'
                }

                this.dialogVisible = true;
            }
        },
        // 查找父节点的方法
        findParentNode(data, id){
            for (let item of data) {
                if (item.children && item.children.length > 0) {
                for (let child of item.children) {
                    if (child.id === id) {
                    return item;
                    }
                }
                const parent = this.findParentNode(item.children, id);
                if (parent) {
                    return parent;
                }
                }
            }
            return null;
        },
        handleSubmit(){ 
            this.$refs.formRef.validate(async(valid) => { 
                if (valid) { 
                    if (this.dialogTitle === '新增资产组') {
                        await this.addAsset();
                    } else {
                        await this.editAsset();
                    }
                    this.dialogVisible = false;
                    this.replaceData();
                } else { 
                    this.$message({
                        message:'您有必填项未填写！',
                        type: 'error'
                    });
                    return false;
                }
            });
        },
        async addAsset(){
            const res = await asset.addApi({
                name: this.form.name,
                // pid: Number(this.fatherNodeID),
                pid:this.selectedValue.split('_')[1],
                remark: this.form.remark
            });
            if (res.code == 200) { 
                this.$message({
                    message:'新增资产组成功',
                    type: 'success'
                });
                this.getZCList(); 
            } else { 
                this.$message({
                    message:'新增资产组失败',
                    type: 'error'
                });
            }
        },
        async editAsset(){
            const res = await asset.editApi({
                id: Number(this.form.id),
                name: this.form.name,
                // pid:this.fatherNodeID == '无'? 0:Number(this.fatherNodeID),
                pid:this.selectedValue.split('_')[1],
                remark: this.form.remark
            });
            if (res.code == 200) { 
                 this.$message({
                    message:'编辑资产组成功',
                    type: 'success'
                });
                this.getZCList();
            } else { 
                this.$message({
                    message:'编辑资产组失败',
                    type: 'error'
                });
            }
        },
        remove(node, data){
            const groupIds = data.id.split('_')[1];

            let msg = '你确定要删除数据吗？';

            this.$confirm(msg, '提示', {
                confirmButtonText: '确认',
                cancelButtonText: '取消',
                type: 'warning',
            }).then(async () => { 
                const res = await asset.deleteApi({ groupIds });
                if (res.code == 200) { 
                    this.$message({
                        message:'删除资产组成功',
                        type: 'success'
                    });
                    this.replaceData();
                } else { 
                    this.$message({
                        message:'删除资产组失败',
                        type: 'error'
                    });
                }
            })
        },
        handleDelete(){ 
            const selectedNodes = this.$refs.treeRef.getCheckedNodes();
            const selectedGroupIds = selectedNodes.map(node => node.id.split('_')[1]).join(',');

            if(selectedGroupIds ==''){ 
                this.$message({
                    message:'请选择要删除的数据。',
                    type: 'error'
                });
                return
            }
            this.remove(null, { id: `group_${selectedGroupIds}` });
        },
        recursiveFilter(node, query){
             // 如果当前节点匹配，直接返回该节点
            if (node.label.toLowerCase().includes(query)) {
                return true;
            }
            // 如果有子节点，递归搜索子节点
            if (node.children && node.children.length) {
                // 递归检查每个子节点是否有匹配的
                const childrenMatch = node.children.some(child => this.recursiveFilter(child, query));
                if (childrenMatch) {
                // 如果任一子节点匹配，返回整个节点（包括子节点）
                    node.children = node.children.filter(child => this.recursiveFilter(child, query));
                    return true;
                }
            }
            // 当前节点和子节点都不匹配
            return false;
        },
        handleSearch(){
             const query = this.searchQuery.toLowerCase();
            // 应用递归搜索到每个顶级节点
            this.filteredData = this.data.filter(node => this.recursiveFilter(node, query));
        },
        replaceData(){
            this.getZCList();
        },
        toggleTree(){
            if (this.isExpanded) {
                this.collapseAllNodes(this.$refs.treeRef.store.root);
            } else {
                this.expandAllNodes(this.$refs.treeRef.store.root);
            }
            this.isExpanded = !this.isExpanded;
        },
        expandAllNodes(node){
            node.expanded = true;
            if (node.childNodes && node.childNodes.length > 0) {
                node.childNodes.forEach(childNode => {
                    this.expandAllNodes(childNode);
                });
            }
        },
        collapseAllNodes(node){
            node.expanded = false;
            if (node.childNodes && node.childNodes.length > 0) {
                node.childNodes.forEach(childNode => {
                    this.collapseAllNodes(childNode);
                });
            }
        },
        toggleCheck(){
            if (this.allChecked) {
                // 取消全选
                this.uncheckAllNodes(this.$refs.treeRef.store.root);
            } else {
                // 全选
                this.checkAllNodes(this.$refs.treeRef.store.root);
            }
            this.allChecked = !this.allChecked;
        },
        checkAllNodes(node){
            this.$refs.treeRef.setChecked(node.data.id, true, true);
            if (node.childNodes && node.childNodes.length > 0) {
                node.childNodes.forEach(childNode => {
                    this.checkAllNodes(childNode);
                });
            }
        },
        uncheckAllNodes(node){
            this.$refs.treeRef.setChecked(node.data.id, false, true);
            if (node.childNodes && node.childNodes.length > 0) {
                    node.childNodes.forEach(childNode => {
                    this.uncheckAllNodes(childNode);
                });
            }
        }, 
        handleCheckChange(){
            //  this.hasCheckedNodes = this.$refs.tree && this.$refs.tree.getCheckedNodes().length > 0; 

            const selectedNodes = this.$refs.treeRefs.getCheckedNodes();
            const selectedGroupIds = selectedNodes.map(node => node.id.split('_')[1]).join(',');

            console.log('selectedGroupIds:'+selectedNodes)
            console.log('selectedGroupIds:'+selectedGroupIds)
        },
        handleVisibleChange(visible) {
            this.popoverVisible = visible
        },
       
        // 点击 tree 节点
        handleNodeClick1(data, node, component) {
            this.selectedValue = data.id
            this.selectedLabel = data.label
            this.popoverVisible = false
            console.log('选中节点:', data)
        },
        buildNodeMap(data){
            data.forEach(node => {
                this.nodeDataMap[node.id] = node
                if (node.children) {
                this.buildNodeMap(node.children)
                }
            })
        },
        onPopoverHide() {
            // 失去焦点时关闭
            this.$refs.input.blur()
        }
    },
    watch:{
       if (newVal) { 
        if (typeof newVal === 'string' && newVal.includes('_')) {
          fatherNodeID.value = newVal.split('_')[1];

        } else {
          fatherNodeID.value = newVal;
        }
      }
    }
    
}
</script>