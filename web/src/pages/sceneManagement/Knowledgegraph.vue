<template>
  <!-- 知识图谱 -->
  <div class="knowledgeboxdiv">
    <div class="main-title">
      <router-link
        :underline="false"
        class="classA"
        :to="{ path: '/taskscenario' }"
        >任务场景</router-link
      >
      <label class="currentpagetitle">
        <el-tooltip class="item" effect="dark" placement="bottom">
          <span> {{ name }}</span>
        </el-tooltip>
        的知识图谱
      </label>
    </div>
    <div class="knowledgebox" id="knowledge" v-loading="loading"></div>
  </div>
</template>
<style >
.knowledgeboxdiv {
  height: 100%;
}
.knowledgebox {
  padding: 24px;
  background: #fff;
  min-height: calc(100% - 39px);
  box-sizing: border-box;
  box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
</style>
<style lang="less" scoped>
.main-title .currentpagetitle span {
  vertical-align: sub;
}
</style>
<script>  
var echarts = require('echarts')
export default ({
  name: 'knowledgegraph',
  data () {
    return {
      id: this.$route.query.id,
      name: this.$route.query.name,
      loading: false,
    }
  },
  created: function () {
    this.$store.state.activefirstMenu = "/tasktemplate"
  },
  mounted: function () {
    this.getData()
  },
  methods: {
    goBack () {
      this.$router.push({
        path: `/experienceSet`,
        query: {
          // page_num: this.fpage_num, 
        }
      })
    },
    getData () {
      this.loading = true
      this.$ajax.get('/smart/scene/graph', {
        params: {
          taskTemplateId: this.id,
          // name:this.name, 
        }
      })
        .then((data) => {
          var dt = data.data
          if (dt.code == 200) {
            this.loading = false
            var xdt = dt.data.nodes
            var line = dt.data.links
            xdt.forEach(element => {
              console.log(typeof element.id,element);
            });
            this.chartKnowledge(xdt, line)

          } else {
            this.$message({
              message: dt.msg || '获取数据失败',
              type: 'error'
            })
          }
        })
        .catch(function (error) {
          console.log(error)
        })
    },
    chartKnowledge: function (dt, line) {
      try {
        line.forEach(element => {
          element.source = element.source.toString()
          element.target = element.target.toString()
        });
           console.log(dt, line, 'dt,line888888')
      // 基于准备好的dom，初始化echarts实例
      var myChart = echarts.init(document.getElementById('knowledge'))
      // 绘制图表
      myChart.setOption({
        animationDuration: 1000,
        animationEasingUpdate: 'quinticInOut',
        color: ['#67c23a', '#4c7ae3', '#4cbae3', 'yellow', 'green', 'red'],
        series: [{
          name: '知识图谱',
          type: 'graph',
          layout: 'force',
          force: {
            repulsion: 500,
            edgeLength: 120
          },

          data: dt,
          links: line,
     


          // 这里得类型代表节点得颜色
          categories: [{
            'name': '初始化'
          }, {
            'name': '信息收集'
          }, {
            'name': '验证攻击'
          },],
          emphasis: { focus: 'adjacency' },
          roam: true,
          label: {
            show: true,
            position: 'top',
          },
          lineStyle: {
            width: 1,
            color: 'source',
            curveness: 0,
            type: "solid"
            // color: 'red',
            // curveness: 1,
            // type: "solid",

          },
        }]
      })
      } catch (error) {
        console.log(error);
      }
    },

  }
})

</script>
