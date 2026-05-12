<template>
  <div>
    <iframe 
      :src="iframeSrc" 
      width="100%" 
      height="100%" 
      frameborder="0" 
      title="Static Analysis Guide"
    ></iframe>
  </div>
</template>

<script>
export default {
  data() {
    return {
      iframeSrc: 'https://ssa.to/static-analysis-guide/intro',
    };
  },
  mounted() {
    // 组件挂载后会立即渲染 iframe
  }
};
</script>

<style >
iframe {
  width: 100%;
  height: 100vh; /* 全屏显示 iframe */
}
</style>
