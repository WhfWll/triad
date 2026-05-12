
<template>
    <div class="textareaIndex">
        <textarea id="info" ></textarea>
    </div>
</template>
 
<script>
export default {
    name: 'ckeditor',
    components: {
    },
    props: ['content'],

    data() {
        return {
            editorDom: null,
            editorData: null,
        }
    },
    methods: {
        getData(){
            return this.editorDom.getData()
        },
        setData(dt){
            this.editorData = dt
        }
    }, 
    watch: {
        'content':function(){
            this.editorDom.setData(this.content) 
        }
    },
    mounted() {  
        this.editorDom = CKEDITOR.replace('info',{
            language: 'zh-CN',  // 中文界面
            // CKEDITOR配置 
            height:300,
            // 配置工具栏
            toolbar:[
                 
                ['Source','-','Save','NewPage','Preview','-','Templates'],
                ['Cut','Copy','Paste','PasteText','PasteFromWord','-','Print', 'SpellChecker', 'Scayt'],
                ['Undo','Redo','-','Find','Replace','-','SelectAll','RemoveFormat'],
                ['Form', 'Checkbox', 'Radio', 'TextField', 'Textarea', 'Select', 'Button', 'ImageButton', 'HiddenField'],
                
                ['Bold','Italic','Underline','Strike','-','Subscript','Superscript'],
                ['NumberedList','BulletedList','-','Outdent','Indent','Blockquote'],
                ['JustifyLeft','JustifyCenter','JustifyRight','JustifyBlock'],
                ['Link','Unlink','Anchor'],
                ['Image','Flash','Table','HorizontalRule','Smiley','SpecialChar','PageBreak'],
                
                ['Styles','Format','Font','FontSize'],
                ['TextColor','BGColor']       
            ], 
            removePlugins: ['about'], 
        });
    
        // 监听图片上传事件
        this.editorDom.on('fileUploadRequest', (evt) => {
        
        });
    
        // 监听ckeditor创建完成事件
        this.editorDom.on('instanceReady', (evt) => {
            this.editorDom.setData(this.content)
        })
    
    
        // 监听内容改变事件
        this.editorDom.on('change', () => {  
            this.editorData = this.editorDom.getData()  
            // 子组件给父组件传值，可以在失去焦点的时候在传值
            this.$emit(`change`, this.editorData)
        })
        
        // 监听失去焦点事件
        this.editorDom.on("blur", () => {
            this.editorData = this.editorDom.getData();
            this.$emit(`change`, this.editorData)
        });
  },
  beforeDestroy() {
    if (this.editorDom) {
      this.editorDom.destroy() // 销毁编辑器，防止内存泄露
      this.editorDom = null;
    }
  },
}
</script>
 
<style scoped>
.textareaIndex{
    width: 100%; 
    justify-content: center;
}
</style> 