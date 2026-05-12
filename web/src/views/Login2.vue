<template>
    <div class="login-bg" style="width: 100%;height:100%;" v-loading.fullscreen.lock="loading" >
  
    </div>
  </template>
  
  <script >
import router from '@/router';
import common from '@/utils/common.js'
import login from '@/api/login.js' 
import { system } from '@/api/system.js'

export default {
  name: 'Login',
  data() {
    return{
        loading:false,
    }
  },
  beforeCreate(){
      //航天运载靶场免密登录
    console.log("beforeCreate");
    let obj = getUrlParams(window.location.href); 
    
    console.log( obj )

    if(obj && obj.token){ 
        let loginParams={
            token:obj.token, 
        }
        this.$ajax({
            url:"/smart/user/loginfreepasse",
            method:"POST",
            data: loginParams
        })
        .then(res =>{
            console.log("beforeCreate33333333333333333");

            console.log(res);
            console.log(res.data);
            var dt = res.data;
            if (dt.code === 200) { 
                var data = dt.data
                localStorage.setItem('test', "11111111111111111111");
                //this.errorbox.errshow = false;
                let role = Number(data.role)
                console.log("444444444444444");
                localStorage.setItem('user_id', this.$commonjs.encryptCBC(data.uid, this.$commonjs.myKey));
                localStorage.setItem('user_id-par',data.uid);
                localStorage.setItem('token', data.token);
                console.log("55555555555555");
                localStorage.setItem('role', this.$commonjs.encryptCBC(role, this.$commonjs.myKey));  //1  普通用户   2  管理员   3  审核员
                localStorage.setItem('user', this.$commonjs.encryptCBC(data.username, this.$commonjs.myKey));
                console.log("6666666666666666");


                this.loading = false;
                if (role == 3) {
                    this.$router.push({ path: '/log' });
                } else {
                    this.$router.push({ path: '/index' });
                } 
            }  
        })
    }else if(obj && obj.code){ 
        let loginParams={
            code:obj.code, 
        }
        this.$ajax({
            url:"/smart/user/loginfreepassa",
            method:"POST",
            data: loginParams
        })
        .then(res =>{
            console.log("beforeCreate33333333333333333");

            console.log(res);
            console.log(res.data);
            var dt = res.data;
            if (dt.code === 200) { 
                var data = dt.data
                localStorage.setItem('test', "11111111111111111111");
                //this.errorbox.errshow = false;
                let role = Number(data.role)
                console.log("444444444444444");
                localStorage.setItem('user_id', this.$commonjs.encryptCBC(data.uid, this.$commonjs.myKey));
                localStorage.setItem('user_id-par',data.uid);
                localStorage.setItem('token', data.token);
                console.log("55555555555555");
                localStorage.setItem('role', this.$commonjs.encryptCBC(role, this.$commonjs.myKey));  //1  普通用户   2  管理员   3  审核员
                localStorage.setItem('user', this.$commonjs.encryptCBC(data.username, this.$commonjs.myKey));
                console.log("6666666666666666");


                this.loading = false;
                if (role == 3) {
                    this.$router.push({ path: '/log' });
                } else {
                    this.$router.push({ path: '/index' });
                } 
            }  
        })
    } else if(obj && obj.Authorization){ 
        let loginParams={
            token:obj.Authorization, 
        }
        this.$ajax({
            url:"/smart/user/loginfreepassf",
            method:"POST",
            data: loginParams
        })
        .then(res =>{
            console.log("beforeCreate33333333333333333");

            console.log(res);
            console.log(res.data);
            var dt = res.data;
            if (dt.code === 200) { 
                var data = dt.data
                localStorage.setItem('test', "11111111111111111111");
                //this.errorbox.errshow = false;
                let role = Number(data.role)
                console.log("444444444444444");
                localStorage.setItem('user_id', this.$commonjs.encryptCBC(data.uid, this.$commonjs.myKey));
                localStorage.setItem('user_id-par',data.uid);
                localStorage.setItem('token', data.token);
                console.log("55555555555555");
                localStorage.setItem('role', this.$commonjs.encryptCBC(role, this.$commonjs.myKey));  //1  普通用户   2  管理员   3  审核员
                localStorage.setItem('user', this.$commonjs.encryptCBC(data.username, this.$commonjs.myKey));
                console.log("6666666666666666");


                this.loading = false;
                if (role == 3) {
                    this.$router.push({ path: '/log' });
                } else {
                    this.$router.push({ path: '/index' });
                } 
            }  
        })
    }

    function  getUrlParams(url) {
        // 通过 ? 分割获取后面的参数字符串
        let urlStr = url.split('?')[1];
        if(!urlStr) return;
        // 创建空对象存储参数
        let obj = {};
        // 再通过 & 将每一个参数单独分割出来
        let paramsArr = urlStr && urlStr.split('&');  
        for(let i = 0,len = paramsArr.length;i < len;i++){
            // 再通过 = 将每一个参数分割为 key:value 的形式
            let arr = paramsArr[i].split('=');
            obj[arr[0]] = arr[1];
        }
        return obj;
    }
  },
  watch:{
    '$route': function (to, from) { // 路由改变时执行     
      
      if(to.path == '/login'){
        let obj = getUrlParams(to.fullPath);

        console.log(obj);

        if(obj && obj.token){ 
                let loginParams={
                    token:obj.token, 
                }
                this.$ajax({
                    url:"/smart/user/loginfreepasse",
                    method:"POST",
                    data: loginParams
                })
                .then(res =>{
                    console.log("watch33333333333333333");

                    console.log(res);
                    console.log(res.data);
                    var dt = res.data;
                    if (dt.code === 200) { 
                        var data = dt.data
                        localStorage.setItem('test', "11111111111111111111");
                        //this.errorbox.errshow = false;
                        let role = Number(data.role)
                        console.log("444444444444444");
                        localStorage.setItem('user_id', this.$commonjs.encryptCBC(data.uid, this.$commonjs.myKey));
                        localStorage.setItem('user_id-par',data.uid);
                        localStorage.setItem('token', data.token);
                        console.log("55555555555555");
                        localStorage.setItem('role', this.$commonjs.encryptCBC(role, this.$commonjs.myKey));  //1  普通用户   2  管理员   3  审核员
                        localStorage.setItem('user', this.$commonjs.encryptCBC(data.username, this.$commonjs.myKey));
                        console.log("6666666666666666");


                        this.loading = false;
                        if (role == 3) {
                            this.$router.push({ path: '/log' });
                        } else {
                            this.$router.push({ path: '/index' });
                        } 
                    }  
                })
            }

        function  getUrlParams(url) {
            // 通过 ? 分割获取后面的参数字符串
            let urlStr = url.split('?')[1];
            if(!urlStr) return;
            // 创建空对象存储参数
            let obj = {};
            // 再通过 & 将每一个参数单独分割出来
            let paramsArr = urlStr && urlStr.split('&');  
            for(let i = 0,len = paramsArr.length;i < len;i++){
                // 再通过 = 将每一个参数分割为 key:value 的形式
                let arr = paramsArr[i].split('=');
                obj[arr[0]] = arr[1];
            }
            return obj;
        }

      }
    },
  },
  created:function(){
    //  //航天运载靶场免密登录
    // console.log("created");
    // let obj = getUrlParams(window.location.href); 
    
    // console.log( obj )

    return;
    if(obj && obj.token){ 
        let loginParams={
            token:obj.token, 
        }
        this.$ajax({
            url:"/smart/user/loginfreepasse",
            method:"POST",
            data: loginParams
        })
        .then(res =>{
            console.log("33333333333333333");
            var dt = res.data;
            if (dt.code === 200) { 
                var data = dt.data
                localStorage.setItem('test', "11111111111111111111");
                //this.errorbox.errshow = false;
                let role = Number(data.role)
                console.log("444444444444444");
                localStorage.setItem('user_id', this.$commonjs.encryptCBC(data.uid, this.$commonjs.myKey));
                localStorage.setItem('user_id-par',data.uid);
                localStorage.setItem('token', data.token);
                console.log("55555555555555");
                localStorage.setItem('role', this.$commonjs.encryptCBC(role, this.$commonjs.myKey));  //1  普通用户   2  管理员   3  审核员
                localStorage.setItem('user', this.$commonjs.encryptCBC(data.username, this.$commonjs.myKey));
                console.log("6666666666666666");


                this.loading = false;
                if (role == 3) {
                    this.$router.push({ path: '/log' });
                } else {
                    this.$router.push({ path: '/index' });
                } 
            }  
        })
    }



    function  getUrlParams(url) {
        // 通过 ? 分割获取后面的参数字符串
        let urlStr = url.split('?')[1];
        if(!urlStr) return;
        // 创建空对象存储参数
        let obj = {};
        // 再通过 & 将每一个参数单独分割出来
        let paramsArr = urlStr && urlStr.split('&');  
        for(let i = 0,len = paramsArr.length;i < len;i++){
            // 再通过 = 将每一个参数分割为 key:value 的形式
            let arr = paramsArr[i].split('=');
            obj[arr[0]] = arr[1];
        }
        return obj;
    }
  },
  methods:{
    getUrlParams(url) {
        // 通过 ? 分割获取后面的参数字符串
        let urlStr = url.split('?')[1];
        if(!urlStr) return;
        // 创建空对象存储参数
        let obj = {};
        // 再通过 & 将每一个参数单独分割出来
        let paramsArr = urlStr && urlStr.split('&');  
        for(let i = 0,len = paramsArr.length;i < len;i++){
            // 再通过 = 将每一个参数分割为 key:value 的形式
            let arr = paramsArr[i].split('=');
            obj[arr[0]] = arr[1];
        }
        return obj;
    },
  }
}  
  </script>
  