export default {
  Seconds: {
    name: '秒',
    every: '每一秒钟',
    interval: ['每隔', '秒执行 从', '秒开始'],
    specific: '具体秒数(可多选)',
    cycle: ['周期从', '到', '秒']
  },
  Minutes: {
    name: '分',
    every: '每一分钟',
    interval: ['每隔', '分执行 从', '分开始'],
    specific: '具体分钟数(可多选)',
    cycle: ['周期从', '到', '分']
  },
  Hours: {
    name: '时',
    every: '每一小时',
    interval: ['每隔', '小时执行 从', '小时开始'],
    specific: '具体小时数(可多选)',
    cycle: ['周期从', '到', '小时']
  },
  Day: {
    name: '天',
    every: '每一天',
    intervalWeek: ['每隔', '周执行 从', '开始'],
    intervalDay: ['每隔', '天执行 从', '天开始'],
    specificWeek: '具体星期几(可多选)',
    specificDay: '具体天数(可多选)',
    lastDay: '在这个月的最后一天',
    lastWeekday: '在这个月的最后一个工作日',
    lastWeek: ['在这个月的最后一个'],
    beforeEndMonth: ['在本月底前', '天'],
    nearestWeekday: ['最近的工作日（周一至周五）至本月', '日'],
    someWeekday: ['在这个月的第', '个'],
    cycle: ['从', '到']
  },
  Week: {
    name: '周',
    every: '每天',
    specific: '具体天数(可多选)',
    list: ['一', '二', '三', '四', '五', '六', '天'].map(val => '星期' + val),
    cycle: ['从', '到']
  },
  Month: {
    name: '月',
    every: '每一月',
    interval: ['每隔', '月执行 从', '月开始'],
    specific: '具体月数(可多选)',
    cycle: ['从', '到', '月之间的每个月']
  },
  Year: {
    name: '年',
    every: '每一年',
    interval: ['每隔', '年执行 从', '年开始'],
    specific: '具体年份(可多选)',
    cycle: ['从', '到', '年之间的每一年']
  },
  Save: '保存',
  Close: '关闭'
}
