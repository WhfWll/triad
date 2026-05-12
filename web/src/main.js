import Vue from 'vue'
import App from './App.vue'
import router from './router/index' 
import store from './store'
import './assets/icon/iconfont.css'
import axios from './axios/axios'
import ElementUI from 'element-ui'
import './assets/theme/index.css'
import './assets/theme/ElementUI.less'
import './assets/Xtheme/index.less'
import qs from 'qs'
// import JsonViewer from 'vue-json-viewer'
import common from './utils/common.js'
// import fullscreen from 'vue-fullscreen'
// import VueTypedJs from 'vue-typed-js'
import dataV from '@jiaminghi/data-view'
import common2 from './commonFunction/common.js'


import JsonViewer from 'vue-json-viewer'
import {
  Editable,
  EditableColumn
} from 'vue-element-extends'
import 'vue-element-extends/lib/index.css'
import VueDOMPurifyHTML from 'vue-dompurify-html'


// import CKEditor from 'ckeditor4-vue'
// Vue.use(CKEditor)  // 全局注册后可在任意组件使用 <ckeditor> 标签

Vue.use(VueDOMPurifyHTML)
Vue.use(Editable)
Vue.use(EditableColumn)
console.log('当前的环境是：',process.env.NODE_ENV === 'development' ? '开发环境' : '生产环境')
if(process.env.NODE_ENV === 'development'){
  // 开发环境
  console.log('看看输出：',process.env); // 输出为123

}else{
  // 生产环境
  console.log('当前的环境是：',process.env.NODE_ENV === 'development' ? '开发环境' : '生产环境')
}
  
Vue.config.productionTip = false 
Vue.prototype.$ajax = axios
Vue.prototype.qs = qs
Vue.prototype.store = store
Vue.prototype.$commonjs = common
Vue.prototype.commonjs = common2
// Vue.use(VueTypedJs)
Vue.use(dataV)
Vue.use(ElementUI)
Vue.use(JsonViewer)
// Vue.use(JsonViewer)
// Vue.use(fullscreen)
new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app')
