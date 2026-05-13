const path = require('path');
const { defineConfig } = require('@vue/cli-service')
// const FileManagerPlugin = require('filemanager-webpack-plugin')
module.exports = defineConfig({
	transpileDependencies: true,
	publicPath: '/',
    assetsDir: 'static',
    productionSourceMap: false,
    devServer: {
        // 前端请求的链接
        host: '127.0.0.1',
        // 前端请求的端口
        port: 8101,
        open: false,
        // process:true, //显示打包进度条
        // compress: true, //启用gzip压缩
        // 代理
        proxy: { 
            '/api': {  
                // target: 'http://172.16.102.72:8011', 
                target: 'http://192.168.0.70:8011', 
                changeOrigin: true,
                pathRewrite: {
                    '^/api': '/'
                }
            }
        }
    },
    parallel: false,
    lintOnSave: false, 
    pluginOptions: {
        "style-resources-loader": {
            preProcessor: "less",
            patterns: [
               // 存放less变量文件的路径
                path.resolve(__dirname, "./src/assets/less/common.less")
            ]
        },
        electronBuilder:{
            nodeIntegration:true
        }
    },
    configureWebpack:{
        externals: [{
         './cptable': 'var cptable'
        }],
        resolve: { fallback: { fs: false } },
        // plugins: [
        //     new FileManagerPlugin({  //初始化 filemanager-webpack-plugin 插件实例
        //         events: {
        //             onEnd: {
        //                 // delete: [   //首先需要删除项目根目录下的dist.zip
        //                 //     './dist/basic.zip',
        //                 // ],
        //                 archive: [ //然后我们选择dist文件夹将之打包成dist.zip并放在根目录
        //                     { source: './dist/', destination: './basic.zip' },
        //                 ]
        //             }
        //         }
        //     }), 
        // ],
    }, 
})
