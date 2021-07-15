<template>
  <div :style="style" class="line-chart">
    <div ref="elRef" class="echarts-element"></div>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, PropType, ref, watch} from 'vue';
import {init} from 'echarts';

export default defineComponent({
  name: 'LineChart',
  props: {
    config: {
      type: Object as PropType<EChartsConfig>,
      required: true,
    },
    width: {
      type: String,
      default: '100%',
    },
    height: {
      type: String,
      default: '100%',
    },
    theme: {
      type: String,
      default: 'light',
    },
    labelKey: {
      type: String,
    },
    valueKey: {
      type: String,
    },
    isTimeSeries: {
      type: Boolean,
      default: true,
    },
  },
  setup(props: LineChartProps, {emit}) {
    const style = computed<Partial<CSSStyleDeclaration>>(() => {
      const {width, height} = props;
      return {
        width,
        height,
      };
    });

    const elRef = ref<HTMLDivElement>();
    const chart = ref<ECharts>();

    const isMultiSeries = computed<boolean>(() => {
      const {config} = props;
      const {dataMetas} = config;
      return dataMetas ? dataMetas.length > 1 : false;
    });

    const getSeriesData = (data: StatsResult[], key?: string) => {
      const {valueKey, labelKey, isTimeSeries} = props;
      const _valueKey = !key ? valueKey : key;

      if (_valueKey) {
        if (isTimeSeries) {
          // time series
          return data.map(d => [d[labelKey || 'date'], d[_valueKey] || 0]);
        } else {
          // not time series
          return data.map(d => d[_valueKey] || 0);
        }
      } else {
        // default
        return data;
      }
    };

    const getSeries = (): EChartSeries[] => {
      const {config} = props;
      const {data, dataMetas} = config;

      if (!isMultiSeries.value) {
        // single series
        return [{
          type: 'line',
          data: getSeriesData(data),
        }];
      } else {
        // multiple series
        const series = [] as EChartSeries[];
        if (!dataMetas) return series;
        dataMetas.forEach(({key, name, yAxisIndex}) => {
          series.push({
            name,
            yAxisIndex,
            type: 'line',
            data: getSeriesData(data, key),
          });
        });
        return series;
      }
    };

    const render = () => {
      const {config, theme, isTimeSeries} = props;
      const {option} = config;

      // dom
      const el = elRef.value;
      if (!el) return;

      // xAxis
      if (!option.xAxis) {
        option.xAxis = {};
        if (isTimeSeries) {
          option.xAxis.type = 'time';
        }
      }

      // yAxis
      if (!option.yAxis) {
        option.yAxis = {};
      }

      // series
      option.series = getSeries();

      // tooltip
      if (!option.tooltip) {
        option.tooltip = {
          // trigger: 'axis',
          // position: ['50%', '50%'],
          // axisPointer: {
          //   type: 'cross',
          // },
        };
      }

      // legend
      option.legend = {};

      // render
      if (!chart.value) {
        // init
        chart.value = init(el, theme);
      }
      (chart.value as ECharts).setOption(option);
    };

    watch(() => props.config.data, render);
    watch(() => props.theme, render);

    onMounted(() => {
      render();
    });

    return {
      style,
      elRef,
      render,
    };
  },
});
</script>

<style lang="scss" scoped>
.line-chart {
  .echarts-element {
    width: 100%;
    height: 100%;
  }
}
</style>
