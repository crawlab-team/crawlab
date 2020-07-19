<template>
  <div v-loading="loading" class="spider-stats">
    <!--overall stats-->
    <el-row>
      <div class="metric-list">
        <metric-card
          label="30-Day Tasks"
          icon="fa fa-play"
          :value="overviewStats.task_count"
          type="danger"
        />
        <metric-card
          label="30-Day Results"
          icon="fa fa-table"
          :value="overviewStats.result_count"
          type="primary"
        />
        <metric-card
          label="Success Rate"
          icon="fa fa-check"
          :value="getPercentStr(overviewStats.success_rate)"
          type="success"
        />
        <metric-card
          label="Avg Duration (sec)"
          icon="fa fa-hourglass"
          :value="getDecimal(overviewStats.avg_runtime_duration)"
          type="warning"
        />
      </div>
    </el-row>
    <!--./overall stats-->

    <el-row>
      <el-col :span="24">
        <el-card class="chart-wrapper">
          <h4>{{ $t('Daily Tasks') }}</h4>
          <div id="task-line" class="chart" />
        </el-card>
      </el-col>
    </el-row>

    <el-row>
      <el-col :span="24">
        <el-card class="chart-wrapper">
          <h4>{{ $t('Daily Avg Duration (sec)') }}</h4>
          <div id="duration-line" class="chart" />
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
    data() {
      return {
        loading: false
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

    mounted() {
    },

    methods: {
      renderTaskLine() {
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
            data: this.dailyStats.map(d => d.task_count),
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

      renderDurationLine() {
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
            data: this.dailyStats.map(d => d.avg_runtime_duration),
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

      render() {
        this.renderTaskLine()
        this.renderDurationLine()
      },

      update() {
        this.loading = true
        this.$store.dispatch('spider/getSpiderStats')
          .then(() => {
            this.render()
          })
          .catch(() => {
            this.$message.error(this.$t('An error happened when fetching the data'))
          })
          .finally(() => {
            this.loading = false
          })
      },

      getPercentStr(value) {
        if (value === undefined) return 'NA'
        return (value * 100).toFixed(2) + '%'
      },

      getDecimal(value) {
        if (value === undefined) return 'NA'
        return value.toFixed(2)
      }
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
