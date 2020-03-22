<template>
  <div class="app-container">
    <el-row>
      <ul class="metric-list">
        <li class="metric-item" v-for="m in metrics" @click="onClickMetric(m)" :key="m.name">
          <div class="metric-icon" :class="m.color">
            <!--            <font-awesome-icon :icon="m.icon"/>-->
            <i :class="m.icon"></i>
          </div>
          <div class="metric-content" :class="m.color">
            <div class="metric-wrapper">
              <div class="metric-number">
                {{overviewStats[m.name]}}
              </div>
              <div class="metric-name">
                {{$t(m.label)}}
              </div>
            </div>
          </div>
          <!--          <el-card class="metric-card" shadow="hover">-->
          <!--            <el-col :span="6" class="icon-col">-->
          <!--              <font-awesome-icon :icon="m.icon" :color="m.color"/>-->
          <!--            </el-col>-->
          <!--            <el-col :span="18" class="text-col">-->
          <!--              <el-row>-->
          <!--                <label class="label">{{$t(m.label)}}</label>-->
          <!--              </el-row>-->
          <!--              <el-row>-->
          <!--                <div class="value">{{overviewStats[m.name]}}</div>-->
          <!--              </el-row>-->
          <!--            </el-col>-->
          <!--          </el-card>-->
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
        { name: 'task_count', label: 'Total Tasks', icon: 'fa fa-check', color: 'blue', path: 'tasks' },
        { name: 'spider_count', label: 'Spiders', icon: 'fa fa-bug', color: 'green', path: 'spiders' },
        { name: 'active_node_count', label: 'Active Nodes', icon: 'fa fa-server', color: 'red', path: 'nodes' },
        { name: 'schedule_count', label: 'Schedules', icon: 'fa fa-clock-o', color: 'orange', path: 'schedules' },
        { name: 'project_count', label: 'Projects', icon: 'fa fa-code-fork', color: 'grey', path: 'projects' }
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

    .metric-item:hover {
      transform: scale(1.05);
      transition: transform 0.5s ease;
    }

    .metric-item {
      flex-basis: 20%;
      height: 64px;
      display: flex;
      color: white;
      cursor: pointer;

      .metric-icon {
        display: inline-flex;
        width: 64px;
        align-items: center;
        justify-content: center;
        border-top-left-radius: 5px;
        border-bottom-left-radius: 5px;
        font-size: 24px;

        svg {
          width: 24px;
        }
      }

      .metric-content {
        display: flex;
        width: calc(100% - 80px);
        align-items: center;
        opacity: 0.85;
        font-size: 14px;
        padding-left: 15px;
        border-top-right-radius: 5px;
        border-bottom-right-radius: 5px;

        .metric-number {
          font-weight: bolder;
          margin-bottom: 5px;
        }
      }

      .metric-icon.blue,
      .metric-content.blue {
        background: #409eff;
      }

      .metric-icon.green,
      .metric-content.green {
        background: #67c23a;
      }

      .metric-icon.red,
      .metric-content.red {
        background: #f56c6c;
      }

      .metric-icon.orange,
      .metric-content.orange {
        background: #E6A23C;
      }

      .metric-icon.grey,
      .metric-content.grey {
        background: #97a8be;
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
