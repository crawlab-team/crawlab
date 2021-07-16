type ECharts = echarts.ECharts;
type EChartOption = echarts.EChartOption;
type EChartSeries = echarts.EChartOption.Series;
type EChartYAxis = echarts.EChartOption.YAxis;

interface EchartsDataMeta {
  // value key
  key: string;
  // name of series
  name: string;
  // x-axis index
  xAxisIndex?: number;
  // y-axis index
  yAxisIndex?: number;
}

interface EChartsConfig {
  dataMetas?: EchartsDataMeta[];
  data: StatsResult[];
  option: EChartOption;
  itemStyleColorFunc?: Function;
}
