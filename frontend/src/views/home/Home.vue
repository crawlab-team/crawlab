<template>
  <div class="app-container">
    <el-row>
      <ul class="metric-list">
        <li class="metric-item" v-for="m in metrics" @click="onClickMetric(m)" :key="m.name">
          <el-card class="metric-card" shadow="hover">
            <el-col :span="6" class="icon-col">
              <font-awesome-icon :icon="m.icon" :color="m.color"/>
              <!--<i :class="m.icon" :style="{color:m.color}"></i>-->
            </el-col>
            <el-col :span="18" class="text-col">
              <el-row>
                <label class="label">{{$t(m.label)}}</label>
              </el-row>
              <el-row>
                <div class="value">{{overviewStats[m.name]}}</div>
              </el-row>
            </el-col>
          </el-card>
        </li>
      </ul>
    </el-row>
    <el-row>
      <el-card shadow="hover">
        <h4 class="title">{{$t('Daily New Tasks')}}</h4>
        <div id="echarts-daily-tasks" class="echarts-box"></div>
      </el-card>
    </el-row>
  </div>
</template>

<script>
import request from '../../api/request'
import echarts from 'echarts'

export default {
  name: 'Home',
  data () {
    return {
      echarts: {},
      overviewStats: {},
      dailyTasks: [],
      metrics: [
        { name: 'task_count', label: 'Total Tasks', icon: ['fa', 'play'], color: '#f56c6c', path: 'tasks' },
        { name: 'spider_count', label: 'Spiders', icon: ['fa', 'bug'], color: '#67c23a', path: 'spiders' },
        { name: 'active_node_count', label: 'Active Nodes', icon: ['fa', 'server'], color: '#409EFF', path: 'nodes' },
        { name: 'schedule_count', label: 'Schedules', icon: ['fa', 'clock'], color: '#409EFF', path: 'schedules' }
      ]
    }
  },
  methods: {
    initEchartsDailyTasks () {
      const option = {
        xAxis: {
          type: 'category',
          data: this.dailyTasks.map(d => d.date)
        },
        yAxis: {
          type: 'value'
        },
        series: [{
          data: this.dailyTasks.map(d => d.task_count),
          type: 'line',
          areaStyle: {},
          smooth: true
        }],
        tooltip: {
          trigger: 'axis',
          show: true
        }
      }
      this.echarts.dailyTasks = echarts.init(this.$el.querySelector('#echarts-daily-tasks'))
      this.echarts.dailyTasks.setOption(option)
    },
    onClickMetric (m) {
      this.$router.push(`/${m.path}`)
    }
  },
  created () {
    request.get('/stats/home')
      .then(response => {
        // overview stats
        this.overviewStats = response.data.data.overview

        // daily tasks
        this.dailyTasks = response.data.data.daily
        this.initEchartsDailyTasks()
      })
  },
  mounted () {
    // this.$ba.trackPageview('/')
  }
}
</script>

<style scoped lang="scss">
  .metric-list {
    margin-top: 0;
    padding-left: 0;
    list-style: none;
    display: flex;
    font-size: 16px;

    .metric-item:last-child .metric-card {
      margin-right: 0;
    }

    .metric-item {
      flex-basis: 25%;

      .metric-card:hover {
      }

      .metric-card {
        margin-right: 30px;
        cursor: pointer;

        .icon-col {
          text-align: right;

          i {
            margin-bottom: 15px;
            font-size: 56px;
          }
        }

        .text-col {
          padding-left: 20px;
          height: 76px;
          text-align: center;

          .label {
            cursor: pointer;
            font-size: 16px;
            display: block;
            height: 24px;
            color: grey;
            font-weight: 900;
          }

          .value {
            font-size: 24px;
            display: block;
            height: 32px;
          }
        }
      }
    }
  }

  .title {
    padding: 0;
    margin: 0;
  }

  #echarts-daily-tasks {
    height: 360px;
    width: 100%;
  }

  .el-card {
    /*border: 1px solid lightgrey;*/
  }

  .svg-inline--fa {
    width: 100%;
    height: 100%;
  }
</style>
