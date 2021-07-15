<template>
  <div class="home">
    <el-row class="row-overview-metrics">
      <el-col
          v-for="(m, i) in metrics"
          :key="i"
          :span="24 / Math.min(metrics.length, 4)"
      >
        <Metric
            :icon="m.icon"
            :title="m.name"
            :value="m.value"
            :clickable="!!m.path"
            :color="getColor(m)"
            @click="onMetricClick(m)"
        />
      </el-col>
    </el-row>
    <el-row class="row-line-chart">
      <LineChart
          :config="dailyConfig"
          is-time-series
          label-key="date"
      />
    </el-row>
    <el-row class="row-pie-chart">
      <el-col :span="8">
        <PieChart
            :config="tasksByStatusConfig"
            label-key="status"
            value-key="tasks"
        />
      </el-col>
      <el-col :span="8">
        <PieChart
            :config="tasksByNodeConfig"
            label-key="node_name"
            value-key="tasks"
        />
      </el-col>
      <el-col :span="8">
        <PieChart
            :config="tasksBySpiderConfig"
            label-key="spider_name"
            value-key="tasks"
        />
      </el-col>
    </el-row>
  </div>
</template>

<script lang="ts">
import {defineComponent, onMounted, ref} from 'vue';
import useRequest from '@/services/request';
import LineChart from '@/components/chart/LineChart.vue';
import dayjs from 'dayjs';
import {spanDateRange} from '@/utils/stats';
import Metric from '@/components/chart/Metric.vue';
import variables from '@/styles/variables.scss';
import {useRouter} from 'vue-router';
import PieChart from '@/components/chart/PieChart.vue';
import {
  TASK_STATUS_CANCELLED,
  TASK_STATUS_ERROR,
  TASK_STATUS_FINISHED,
  TASK_STATUS_PENDING,
  TASK_STATUS_RUNNING
} from '@/constants/task';

const {
  get,
} = useRequest();

export default defineComponent({
  name: 'Home',
  components: {PieChart, Metric, LineChart},
  setup() {
    const router = useRouter();

    const dateRange = ref<DateRange>({
      start: dayjs().subtract(30, 'day'),
      end: dayjs(),
    });

    const metrics = ref<MetricMeta[]>([
      {
        name: 'Active Nodes',
        icon: ['fa', 'server'],
        key: 'nodes',
        value: 0,
        path: '/nodes',
      },
      {
        name: 'Projects',
        icon: ['fa', 'project-diagram'],
        key: 'projects',
        value: 0,
        path: '/projects'
      },
      {
        name: 'Spiders',
        icon: ['fa', 'spider'],
        key: 'spiders',
        value: 0,
        path: '/spiders',
      },
      {
        name: 'Schedules',
        icon: ['fa', 'clock'],
        key: 'schedules',
        value: 0,
        path: '/schedules',
      },
      {
        name: 'Total Tasks',
        icon: ['fa', 'tasks'],
        key: 'tasks',
        value: 0,
        path: '/tasks',
      },
      {
        name: 'Error Tasks',
        icon: ['fa', 'exclamation'],
        key: 'error_tasks',
        value: 0,
        path: '/tasks',
        color: (m: MetricMeta) => m.value > 0 ? variables.dangerColor : variables.successColor,
      },
      {
        name: 'Total Results',
        icon: ['fa', 'table'],
        key: 'results',
        value: 0,
        color: (m: MetricMeta) => m.value > 0 ? variables.successColor : variables.infoMediumColor,
      },
      {
        name: 'Users',
        icon: ['fa', 'users'],
        key: 'users',
        value: 0,
        path: '/users',
      },
    ]);

    const dailyConfig = ref<EChartsConfig>({
      dataMetas: [
        {
          key: 'tasks',
          name: 'Tasks',
          yAxisIndex: 0,
        },
        {
          key: 'results',
          name: 'Results',
          yAxisIndex: 1,
        },
      ],
      data: [],
      option: {
        title: {
          text: 'Daily Stats',
        },
        yAxis: [
          {name: 'Tasks', position: 'left'},
          {name: 'Results', position: 'right'},
        ],
        color: [
          variables.primaryColor,
          variables.successColor,
        ],
      },
    });

    const tasksByStatusConfig = ref<EChartsConfig>({
      data: [],
      option: {
        title: {
          text: 'Tasks by Status',
        },
      },
      itemStyleColorFunc: ({data}: any) => {
        const {name} = data;
        switch (name) {
          case TASK_STATUS_PENDING:
            return variables.primaryColor;
          case TASK_STATUS_RUNNING:
            return variables.warningColor;
          case TASK_STATUS_FINISHED:
            return variables.successColor;
          case TASK_STATUS_ERROR:
            return variables.dangerColor;
          case TASK_STATUS_CANCELLED:
            return variables.infoMediumColor;
          default:
            return 'red';
        }
      },
    });

    const tasksByNodeConfig = ref<EChartsConfig>({
      data: [],
      option: {
        title: {
          text: 'Tasks by Node',
        },
      },
    });

    const tasksBySpiderConfig = ref<EChartsConfig>({
      data: [],
      option: {
        title: {
          text: 'Tasks by Spider',
        },
      },
    });

    const getOverview = async () => {
      // TODO: filter by date range?
      // const {start, end} = dateRange.value;
      const res = await get(`/stats/overview`);
      metrics.value.forEach(m => {
        m.value = res.data[m.key];
      });
    };

    const getDaily = async () => {
      // TODO: filter by date range?
      const {start, end} = dateRange.value;
      const res = await get(`/stats/daily`);
      dailyConfig.value.data = spanDateRange(start, end, res.data, 'date');
    };

    const getTasks = async () => {
      // TODO: filter by date range?
      const {start, end} = dateRange.value;
      const res = await get(`/stats/tasks`);
      tasksByStatusConfig.value.data = res.data.by_status;
      tasksByNodeConfig.value.data = res.data.by_node;
      tasksBySpiderConfig.value.data = res.data.by_spider;
    };

    const getData = async () => Promise.all([
      getOverview(),
      getDaily(),
      getTasks(),
    ]);

    const onMetricClick = (m: MetricMeta) => {
      if (m.path) {
        router.push(m.path);
      }
    };

    const defaultColorFunc = (value: string | number) => {
      if (typeof value === 'number') {
        // number
        if (value === 0) {
          return variables.infoMediumColor;
        } else {
          return variables.primaryColor;
        }
      } else {
        // string
        const v = Number(value);
        if (isNaN(v) || v == 0) {
          return variables.infoMediumColor;
        } else {
          return variables.primaryColor;
        }
      }
    };

    const getColor = (m: MetricMeta) => {
      if (!m.color) {
        return defaultColorFunc(m.value);
      } else if (typeof m.color === 'function') {
        return m.color(m);
      } else {
        return m.color;
      }
    };

    onMounted(async () => {
      await getData();
    });

    return {
      metrics,
      dailyConfig,
      tasksByStatusConfig,
      tasksByNodeConfig,
      tasksBySpiderConfig,
      onMetricClick,
      getColor,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables.scss";

.home {
  background: white;
  min-height: calc(100vh - #{$headerHeight} - #{$tabsViewHeight});
  padding: 20px;

  .row-overview-metrics {
    display: flex;
    flex-wrap: wrap;
    margin-bottom: 20px;
  }

  .row-line-chart {
    height: 400px;
  }

  .row-pie-chart {
    height: 400px;
  }
}
</style>
