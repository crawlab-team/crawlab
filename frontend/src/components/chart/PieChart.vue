<template>
  <div :style="style" class="pie-chart">
    <div v-if="isEmpty" class="empty-placeholder">
      No Data Available
    </div>
    <div ref="elRef" class="echarts-element"></div>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, PropType, ref, watch} from 'vue';
import {init} from 'echarts';

export default defineComponent({
  name: 'PieChart',
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
  },
  setup(props: PieChartProps, {emit}) {
    const style = computed<Partial<CSSStyleDeclaration>>(() => {
      const {width, height} = props;
      return {
        width,
        height,
      };
    });

    const elRef = ref<HTMLDivElement>();
    const chart = ref<ECharts>();

    const isEmpty = computed<boolean>(() => {
      const {config} = props;
      const {data} = config;
      if (!data) return true;
      return data.length === 0;

    });

    const getSeriesData = (data: StatsResult[], key?: string) => {
      const {valueKey, labelKey} = props;
      const _valueKey = !key ? valueKey : key;

      if (_valueKey) {
        return data.map(d => {
          return {
            name: d[labelKey || '_id'],
            value: d[_valueKey] || 0,
          };
        });
      } else {
        // default
        return data;
      }
    };

    const getSeries = (): EChartSeries[] => {
      const {config} = props;
      const {data, itemStyleColorFunc} = config;

      const seriesItem = {
        type: 'pie',
        data: getSeriesData(data || []),
        radius: ['40%', '70%'],
        alignTo: 'labelLine',
      } as EChartSeries;

      if (itemStyleColorFunc) {
        seriesItem.itemStyle = {color: itemStyleColorFunc};
      }

      return [seriesItem];
    };

    const render = () => {
      const {config, theme} = props;
      const {option} = config;

      // dom
      const el = elRef.value;
      if (!el) return;

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
      isEmpty,
      style,
      elRef,
      render,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables";

.pie-chart {
  position: relative;

  .empty-placeholder {
    position: absolute;
    top: 0;
    left: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    width: 100%;
  }

  .echarts-element {
    width: 100%;
    height: 100%;
  }
}
</style>
