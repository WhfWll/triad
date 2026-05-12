<template>
    <div class="basic">
        <div v-for="itme in aboutList" :key="itme.id">
            <strong>{{ itme.title }}</strong>
            <span>{{ itme.content }}</span>
        </div>
    </div>
</template>

<script>
export default {
    name: "about",
    data: () => ({
        aboutList: [
            {
                id: 1,
                title: "介绍：",
                content:
                    "渗透资源模块是渗透过程中，当不能对被测目标上传文件，但可以下载文件时，可以让被测目标从渗透资源模块下载指定文件。",
            },
            {
                id: 2,
                title: "使用帮助：",
                content:
                    "编写渗透脚本的功能中涉及需要下载渗透资源文件时，可复制渗透资源文件的地址到脚本中，脚本执行时就能够在被测目标实现远程下载渗透资源文件。",
            },
        ],
    }),
};
</script>

<style lang="less" scoped>
.basic {
    background: #f7f7fb !important;
    border: 1px solid #e2e5ed;
    border-left: 2px solid #4c7ae3;
    padding: 5px 10px;
}
.basic > div {
    color: rgba(72, 72, 102, 0.64);
    font-weight: 500;
    margin: 10px 8px;
    font-size: 13px;
}
.basic > div strong {
    display: inline-block;
    width: 80px;
    color: #4c7ae3;
}
</style>