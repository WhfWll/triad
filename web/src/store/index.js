import Vue from 'vue'
import Vuex from 'vuex'
import individual from '@/commonFunction/individual' //引入定义的属性
const { themeDark, themeLight, themeWhite,themeMR} = individual  //4种不同主题

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    topNavState: 'home',
    leftNavState: 'home',
    activefirstMenu: "index",
    breadcrumb: '1',
    activesecondMenu: '',
    uploadPercentage: 0,
    system_uid:'',
    groupID:[],
    copyArr:[],//任务-复制
    sideTheme: JSON.parse(localStorage.getItem('theme'))||'themeDark',
    sideThemeNav: localStorage.getItem('themeNav')||'LeftNav',
    AttributeTheme: localStorage.getItem('AttributeTheme')||'light',
    sideThemeLOGO: localStorage.getItem('themeLOGO')||'false',
    systemAuthorized: false
  },
  mutations: {
    setPercentage (state, data) {
      state.uploadPercentage = data
    },
    setSystemAuthorized (state, data) {
      state.systemAuthorized = data
    },
    setgroupID (state, data) {
      state.groupID = data
    },
    updateCopyArr(state, newArray) {
      state.copyArr = newArray; 
    },
    CHANGE_SETTING: (state, data) => {

      console.log(state, data,'CHANGE_SETTING');
          if (data === 'themeDark') {
            state.sideTheme= themeDark
            state.AttributeTheme= 'dark'
            localStorage.setItem('theme', JSON.stringify(themeDark))
            localStorage.setItem('AttributeTheme', 'dark')
            // 切黑色
            window.document.documentElement.setAttribute('theme', 'dark');
          }
          else if (data === 'themeLight') {
            state.sideTheme= themeLight
            localStorage.setItem('theme', JSON.stringify(themeLight))
          }
          else if (data === 'themeWhite') {
            state.sideTheme= themeWhite
            localStorage.setItem('theme', JSON.stringify(themeWhite))
          }
          else if (data === 'themeMR') {
            state.sideTheme= themeMR
            state.AttributeTheme= 'light'
            localStorage.setItem('theme', JSON.stringify(themeMR))
             // 切白天
             window.document.documentElement.setAttribute('theme', 'light');
            localStorage.setItem('AttributeTheme', 'light')

          }
  },
  CHANGE_SETTING_NAV: (state, data) => {
      console.log(state, data,'CHANGE_SETTING_NAV');
          if (data === 'LeftNav') {
            state.sideThemeNav= data
            localStorage.setItem('themeNav', data)
          }
          else if (data === 'topNav') {
            state.sideThemeNav= data
            localStorage.setItem('themeNav', data)
          }
          
  },
  CHANGE_SETTING_LOGO: (state, data) => {
      console.log(state, data,'CHANGE_SETTING_LOGO');
          if (data == 'false') {
            state.sideThemeLOGO= data
            localStorage.setItem('themeLOGO', data)
          }
          else if (data == 'true') {
            
            state.sideThemeLOGO= data
            localStorage.setItem('themeLOGO', data)
          }
          
  },

},

  modules: {
  }
})
