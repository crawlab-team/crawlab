<template>
  <div id="network-chart"></div>
</template>

<script>
import echarts from 'echarts'

export default {
  name: 'NodeNetwork',
  props: {
    activeTab: {
      type: String
    }
  },
  watch: {
    activeTab () {
      setTimeout(() => {
        this.render()
      }, 0)
    }
  },
  data () {
    return {
      chart: undefined
    }
  },
  computed: {
    masterNode () {
      const nodes = this.$store.state.node.nodeList
      for (let i = 0; i < nodes.length; i++) {
        if (nodes[i].is_master) {
          return nodes[i]
        }
      }
      return {}
    },
    nodes () {
      let nodes = this.$store.state.node.nodeList
      nodes = nodes
        .filter(d => d.status !== 'offline')
        .map(d => {
          d.id = d._id
          d.x = Math.floor(100 * Math.random())
          d.y = Math.floor(100 * Math.random())
          d.itemStyle = {
            color: d.is_master ? '#409EFF' : '#e6a23c'
          }
          return d
        })

      // mongodb
      nodes.push({
        id: 'mongodb',
        name: 'MongoDB',
        x: Math.floor(100 * Math.random()),
        y: Math.floor(100 * Math.random()),
        itemStyle: {
          color: '#67c23a'
        }
      })

      // redis
      nodes.push({
        id: 'redis',
        name: 'Redis',
        x: Math.floor(100 * Math.random()),
        y: Math.floor(100 * Math.random()),
        itemStyle: {
          color: '#f56c6c'
        }
      })

      return nodes
    },
    links () {
      const links = []
      for (let i = 0; i < this.nodes.length; i++) {
        if (this.nodes[i].status === 'offline') continue
        if (['redis', 'mongodb'].includes(this.nodes[i].id)) continue
        // mongodb
        links.push({
          source: this.nodes[i].id,
          target: 'mongodb',
          value: 10,
          lineStyle: {
            color: '#67c23a'
          }
        })

        // redis
        links.push({
          source: this.nodes[i].id,
          target: 'redis',
          value: 10,
          lineStyle: {
            color: '#f56c6c'
          }
        })

        if (this.masterNode.id === this.nodes[i].id) continue

        // master
        links.push({
          source: this.masterNode.id,
          target: this.nodes[i].id,
          value: 0.5,
          lineStyle: {
            color: '#409EFF'
          }
        })
      }
      return links
    }
  },
  methods: {
    render () {
      const option = {
        title: {
          text: this.$t('Node Network')
        },
        tooltip: {
          formatter: params => {
            let str = '<span style="margin-right:5px;display:inline-block;height:12px;width:12px;border-radius:6px;background:' + params.color + '"></span>'
            if (params.data.name) str += '<span>' + params.data.name + '</span><br>'
            if (params.data.ip) str += '<span>IP: ' + params.data.ip + '</span><br>'
            if (params.data.mac) str += '<span>MAC: ' + params.data.mac + '</span><br>'
            return str
          }
        },
        animationDurationUpdate: 1500,
        animationEasingUpdate: 'quinticInOut',
        series: [
          {
            type: 'graph',
            layout: 'force',
            symbolSize: 50,
            roam: true,
            label: {
              normal: {
                show: true
              }
            },
            edgeSymbol: ['circle', 'arrow'],
            edgeSymbolSize: [4, 10],
            edgeLabel: {
              normal: {
                textStyle: {
                  fontSize: 20
                }
              }
            },
            focusOneNodeAdjacency: true,
            force: {
              initLayout: 'force',
              repulsion: 30,
              gravity: 0.001,
              edgeLength: 30
            },
            draggable: true,
            data: this.nodes,
            links: this.links,
            lineStyle: {
              normal: {
                opacity: 0.9,
                width: 2,
                curveness: 0
              }
            }
          }
        ]
      }
      this.chart = echarts.init(this.$el)
      this.chart.setOption(option)
      this.chart.resize()
    }
  },
  mounted () {
    this.render()
  }
}
</script>

<style scoped>
  #network-chart {
    height: 480px;
    width: 100%;
  }
</style>
