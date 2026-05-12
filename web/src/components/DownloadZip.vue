/* 
    
    1.安装jszip和file-saver
    2.将后端数组数据的每一条下载为一个word文件，然后打包成zip文件下载

 */
<template>
  <div>
    <el-button type="primary" @click="downloadData">zip下载</el-button>
  </div>
</template>

<script>

import Zip from 'jszip';
import {saveAs} from 'file-saver';
export default {
  data () {
    return {}
  },
  methods: {
    downloadData () {
      let downloadUrl = ['1111', '22222', 'xxx'] //数组数据，诶一条代表一个word文件
      const zip = new Zip()
      for (let item of downloadUrl) {
          let fileName = `${item}.doc` //文件名需要做处理 必须加后缀
          zip.file(fileName, item, { binary: true })
      }

        if (Object.keys(zip.files).length > 0) {
          zip.generateAsync({ type: 'blob' }).then((blob) => {
            saveAs(blob, '批量文件.zip')
            console.log('批量下载成功')
          })
        } else {
          console.log('下载全部失败')
        }

    },
  },
}
</script>

<style></style>
