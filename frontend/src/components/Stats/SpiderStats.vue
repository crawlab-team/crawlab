<template>
  <div class="spider-stats">
    <!--overall stats-->
    <el-row>
      <div class="metric-list">
        <metric-card :label="$t('30-Day Tasks')"
                     icon="fa fa-play"
                     :value="overviewStats.task_count"
                     type="danger"/>
        <metric-card :label="$t('30-Day Results')"
                     icon="fa fa-table"
                     :value="overviewStats.result_count"
                     type="primary"/>
        <metric-card :label="$t('Success Rate')"
                     icon="fa fa-check"
                     :value="getPercentStr(overviewStats.success_rate)"
                     type="success"/>
        <metric-card :label="$t('Avg Duration (sec)')"
                     icon="fa fa-hourglass"
                     :value="getDecimal(overviewStats.avg_duration)"
                     type="warning"/>
      </div>
    </el-row>
    <!--./overall stats-->

    <el-row>
      <el-col :span="8">
        <el-card class="chart-wrapper" style="margin-right: 20px;">
          <h4>{{$t('Tasks by Status')}}</h4>
          <div id="task-pie-status" class="chart"></div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card class="chart-wrapper">
          <h4>{{$t('Daily Tasks')}}</h4>
          <div id="task-line" class="chart"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row>
      <el-col :span="8">
        <el-card class="chart-wrapper" style="margin-right: 20px;">
          <h4>{{$t('Tasks by Node')}}</h4>
          <div id="task-pie-node" class="chart"></div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card class="chart-wrapper">
          <h4>{{$t('Daily Duration (sec)')}}</h4>
          <div id="duration-line" class="chart"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import MetricCard from './MetricCard'
import echarts from 'echarts'

export default {
  name: 'SpiderStats',
  components: { MetricCard },
  methods: {
    renderTaskPieStatus () {
      const chart = echarts.init(this.$el.querySelector('#task-pie-status'))
      const option = {
        tooltip: {
          show: true
        },
        series: [{
          name: '',
          type: 'pie',
          // radius: ['50%', '70%'],
          data: this.statusStats.map(d => {
            let color
            if (d.name === 'SUCCESS') {
              color = '#67c23a'
            } else if (d.name === 'STARTED') {
              color = '#e6a23c'
            } else if (d.name === 'FAILURE') {
              color = '#f56c6c'
            } else {
              color = 'grey'
            }
            return {
              name: this.$t(d.name),
              value: d.value,
              itemStyle: {
                color
              }
            }
          })
        }]
      }
      chart.setOption(option)
    },

    renderTaskPieNode () {
      const chart = echarts.init(this.$el.querySelector('#task-pie-node'))
      const option = {
        tooltip: {
          show: true
        },
        series: [{
          name: '',
          type: 'pie',
          // radius: ['50%', '70%'],
          data: this.nodeStats.map(d => {
            return {
              name: d.name,
              value: d.value
              // itemStyle: {
              //   color
              // }
            }
          })
        }]
      }
      chart.setOption(option)
    },

    renderTaskLine () {
      const chart = echarts.init(this.$el.querySelector('#task-line'))
      const option = {
        grid: {
          top: 20,
          bottom: 40
        },
        xAxis: {
          type: 'category',
          data: this.dailyStats.map(d => d.date)
        },
        yAxis: {
          type: 'value'
        },
        series: [{
          type: 'line',
          data: this.dailyStats.map(d => d.count),
          areaStyle: {},
          smooth: true
        }],
        tooltip: {
          trigger: 'axis',
          show: true
        }
      }
      chart.setOption(option)
    },

    renderDurationLine () {
      const chart = echarts.init(this.$el.querySelector('#duration-line'))
      const option = {
        grid: {
          top: 20,
          bottom: 40
        },
        xAxis: {
          type: 'category',
          data: this.dailyStats.map(d => d.date)
        },
        yAxis: {
          type: 'value'
        },
        series: [{
          type: 'line',
          data: this.dailyStats.map(d => d.duration),
          areaStyle: {},
          smooth: true
        }],
        tooltip: {
          trigger: 'axis',
          show: true
        }
      }
      chart.setOption(option)
    },

    render () {
      this.renderTaskPieStatus()
      this.renderTaskLine()
      this.renderTaskPieNode()
      this.renderDurationLine()
    },

    update () {
      this.$store.dispatch('spider/getSpiderStats')
        .then(() => {
          this.render()
        })
    },

    getPercentStr (value) {
      if (value === undefined) return 'NA'
      return (value * 100).toFixed(2) + '%'
    },

    getDecimal (value) {
      if (value === undefined) return 'NA'
      return value.toFixed(2)
    }
  },
  computed: {
    ...mapState('spider', [
      'overviewStats',
      'statusStats',
      'nodeStats',
      'dailyStats'
    ])
  },
  mounted () {
  }
}
</script>

<style scoped>
  .metric-list {
    display: flex;
  }

  .metric-list .metric-card {
    flex-basis: 25%;
  }

  .chart-wrapper {
    margin-top: 20px;
  }

  .chart {
    width: 100%;
    height: 240px;
  }

  .table {
    height: 240px;
  }

  h4 {
    display: inline-block;
    margin: 0
  }
</style>
