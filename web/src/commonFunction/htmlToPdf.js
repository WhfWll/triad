/* 导出pdf文档 */
import html2Canvas from "html2canvas";
import JsPDF from "jspdf";

export default {
    
    install(Vue, options) {
        Vue.prototype.getPdf = function (id, title) {
            console.log(11);
            const loading = Vue.prototype.$loading({
                fullscreen: true,
                lock: true,
                text: 'Loading',
                spinner: 'el-icon-loading',
                background: 'rgba(0, 0, 0, 0.7)'
            });
            let shareContent = document.getElementById(id), //需要截图的包裹的（原生的）DOM 对象
                width = shareContent.clientWidth, //获取dom 宽度
                height = shareContent.clientHeight, //获取dom 高度
                canvas = document.createElement("canvas"), //创建一个canvas节点
                scale = 1; //定义任意放大倍数 支持小数
            canvas.width = width * scale; //定义canvas 宽度 * 缩放
            canvas.height = height * scale; //定义canvas高度 *缩放
            canvas.style.width = shareContent.clientWidth * scale + "px";
            canvas.style.height = shareContent.clientHeight * scale + "px";
            canvas.getContext("2d").scale(scale, scale); //获取context,设置scale
            let opts = {
                scale: scale, // 添加的scale 参数
                canvas: canvas, //自定义 canvas
                logging: false, //日志开关，便于查看html2canvas的内部执行流程
                width: width, //dom 原始宽度
                height: height,
                useCORS: true, // 【重要】开启跨域配置
            }
            html2Canvas(shareContent, opts).then(() => {
                var contentWidth = canvas.width;
                var contentHeight = canvas.height;
                //一页pdf显示html页面生成的canvas高度;
                var pageHeight = (contentWidth / 592.28) * 841.89;
                //未生成pdf的html页面高度
                var leftHeight = contentHeight;
                //页面偏移
                var position = 0;
                //a4纸的尺寸[595.28,841.89]，html页面生成的canvas在pdf中图片的宽高
                var imgWidth = 595.28;
                var imgHeight = (592.28 / contentWidth) * contentHeight;
                var pageData = canvas.toDataURL("image/jpeg", 1.0);
                var PDF = new JsPDF("", "pt", "a4");
                if (leftHeight < pageHeight) {
                    PDF.addImage(pageData, "JPEG", 0, 0, imgWidth, imgHeight);
                } else {
                    while (leftHeight > 0) {
                        PDF.addImage(pageData, "JPEG", 0, position, imgWidth, imgHeight);
                        leftHeight -= pageHeight;
                        position -= 841.89;
                        if (leftHeight > 0) {
                            PDF.addPage();
                        }
                    }
                }
                PDF.save(title + ".pdf"); // 这里是导出的文件名
                loading.close();
                this.$router.go(-1)
            });
        };
    }
};