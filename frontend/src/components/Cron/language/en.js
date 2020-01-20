export default {
  Seconds: {
    name: 'Seconds',
    every: 'Every second',
    interval: ['Every', 'second(s) starting at second'],
    specific: 'Specific second (choose one or many)',
    cycle: ['Every second between second', 'and second']
  },
  Minutes: {
    name: 'Minutes',
    every: 'Every minute',
    interval: ['Every', 'minute(s) starting at minute'],
    specific: 'Specific minute (choose one or many)',
    cycle: ['Every minute between minute', 'and minute']
  },
  Hours: {
    name: 'Hours',
    every: 'Every hour',
    interval: ['Every', 'hour(s) starting at hour'],
    specific: 'Specific hour (choose one or many)',
    cycle: ['Every hour between hour', 'and hour']
  },
  Day: {
    name: 'Day',
    every: 'Every day',
    intervalWeek: ['Every', 'day(s) starting on'],
    intervalDay: ['Every', 'day(s) starting at the', 'of the month'],
    specificWeek: 'Specific day of week (choose one or many)',
    specificDay: 'Specific day of month (choose one or many)',
    lastDay: 'On the last day of the month',
    lastWeekday: 'On the last weekday of the month',
    lastWeek: ['On the last', ' of the month'],
    beforeEndMonth: ['day(s) before the end of the month'],
    nearestWeekday: ['Nearest weekday (Monday to Friday) to the', 'of the month'],
    someWeekday: ['On the', 'of the month'],
    cycle: ['From', 'to']
  },
  Week: {
    name: 'Week',
    every: 'Every day',
    specific: 'Specific weekday (choose on or many)',
    list: ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'],
    cycle: ['From', 'to']
  },
  // Week:['Sunday','Monday','Tuesday','Wednesday','Thursday','Friday','Saturday'],
  Month: {
    name: 'Month',
    every: 'Every month',
    interval: ['Every', 'month(s) starting in'],
    specific: 'Specific month (choose one or many)',
    cycle: ['Every month between', 'and']
  },
  Year: {
    name: 'Year',
    every: 'Any year',
    interval: ['Every', 'year(s) starting in'],
    specific: 'Specific year (choose one or many)',
    cycle: ['Every year between', 'and']
  },
  Save: 'Save',
  Close: 'Close'
}
