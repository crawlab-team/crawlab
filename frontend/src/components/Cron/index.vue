<style lang="scss" scoped>
  #change-crontab {
    .language {
      position: absolute;
      right: 25px;
      z-index: 1;
    }

    .el-tabs {
      box-shadow: none;
    }

    .tabBody {
      .el-row {
        margin: 10px 0;

        .long {
          .el-select {
            width: 350px;
          }
        }

        .el-input-number {
          width: 110px;
        }
      }
    }

    .bottom {
      width: 100%;
      text-align: center;
      margin-top: 5px;
      position: relative;

      .value {
        font-size: 18px;
        vertical-align: middle;
      }
    }
  }
</style>
<template>
  <div id="change-crontab">
    <!--        <el-button class="language" type="text" @click="i18n=(i18n==='en'?'cn':'en')">{{i18n}}</el-button>-->
    <el-tabs type="border-card">
      <el-tab-pane>
        <span slot="label"><i class="el-icon-date"></i> {{text.Minutes.name}}</span>
        <div class="tabBody">
          <el-row>
            <el-radio v-model="minute.cronEvery" label="1">{{text.Minutes.every}}</el-radio>
          </el-row>
          <el-row>
            <el-radio v-model="minute.cronEvery" label="2">{{text.Minutes.interval[0]}}
              <el-input-number size="small" v-model="minute.incrementIncrement" :min="0" :max="59"></el-input-number>
              {{text.Minutes.interval[1]}}
              <el-input-number size="small" v-model="minute.incrementStart" :min="0" :max="59"></el-input-number>
              {{text.Minutes.interval[2]||''}}
            </el-radio>
          </el-row>
          <el-row>
            <el-radio class="long" v-model="minute.cronEvery" label="3">{{text.Minutes.specific}}
              <el-select size="small" multiple v-model="minute.specificSpecific">
                <el-option v-for="val in 60" :key="val" :value="(val-1).toString()" :label="val-1"></el-option>
              </el-select>
            </el-radio>
          </el-row>
          <el-row>
            <el-radio v-model="minute.cronEvery" label="4">{{text.Minutes.cycle[0]}}
              <el-input-number size="small" v-model="minute.rangeStart" :min="0" :max="59"></el-input-number>
              {{text.Minutes.cycle[1]}}
              <el-input-number size="small" v-model="minute.rangeEnd" :min="0" :max="59"></el-input-number>
              {{text.Minutes.cycle[2]}}
            </el-radio>
          </el-row>
        </div>
      </el-tab-pane>
      <el-tab-pane>
        <span slot="label"><i class="el-icon-date"></i> {{text.Hours.name}}</span>
        <div class="tabBody">
          <el-row>
            <el-radio v-model="hour.cronEvery" label="1">{{text.Hours.every}}</el-radio>
          </el-row>
          <el-row>
            <el-radio v-model="hour.cronEvery" label="2">{{text.Hours.interval[0]}}
              <el-input-number size="small" v-model="hour.incrementIncrement" :min="0" :max="23"></el-input-number>
              {{text.Hours.interval[1]}}
              <el-input-number size="small" v-model="hour.incrementStart" :min="0" :max="23"></el-input-number>
              {{text.Hours.interval[2]}}
            </el-radio>
          </el-row>
          <el-row>
            <el-radio class="long" v-model="hour.cronEvery" label="3">{{text.Hours.specific}}
              <el-select size="small" multiple v-model="hour.specificSpecific">
                <el-option v-for="val in 24" :key="val" :value="(val-1).toString()" :label="val-1"></el-option>
              </el-select>
            </el-radio>
          </el-row>
          <el-row>
            <el-radio v-model="hour.cronEvery" label="4">{{text.Hours.cycle[0]}}
              <el-input-number size="small" v-model="hour.rangeStart" :min="0" :max="23"></el-input-number>
              {{text.Hours.cycle[1]}}
              <el-input-number size="small" v-model="hour.rangeEnd" :min="0" :max="23"></el-input-number>
              {{text.Hours.cycle[2]}}
            </el-radio>
          </el-row>
        </div>
      </el-tab-pane>
      <el-tab-pane>
        <span slot="label"><i class="el-icon-date"></i> {{text.Day.name}}</span>
        <div class="tabBody">
          <el-row>
            <el-radio v-model="day.cronEvery" label="1">{{text.Day.every}}</el-radio>
          </el-row>
          <el-row>
            <el-radio v-model="day.cronEvery" label="2">{{text.Day.intervalDay[0]}}
              <el-input-number size="small" v-model="day.incrementIncrement" :min="1" :max="31"></el-input-number>
              {{text.Day.intervalDay[1]}}
              <el-input-number size="small" v-model="day.incrementStart" :min="1" :max="31"></el-input-number>
              {{text.Day.intervalDay[2]}}
            </el-radio>
          </el-row>
          <el-row>
            <el-radio class="long" v-model="day.cronEvery" label="3">{{text.Day.specificDay}}
              <el-select size="small" multiple v-model="day.specificSpecific">
                <el-option v-for="val in 31" :key="val" :value="val.toString()" :label="val"></el-option>
              </el-select>
            </el-radio>
          </el-row>
          <el-row>
            <el-radio v-model="day.cronEvery" label="4">{{text.Day.cycle[0]}}
              <el-input-number size="small" v-model="day.rangeStart" :min="1" :max="31"></el-input-number>
              {{text.Day.cycle[1]}}
              <el-input-number size="small" v-model="day.rangeEnd" :min="1" :max="31"></el-input-number>
            </el-radio>
          </el-row>
        </div>
      </el-tab-pane>
      <el-tab-pane>
        <span slot="label"><i class="el-icon-date"></i> {{text.Month.name}}</span>
        <div class="tabBody">
          <el-row>
            <el-radio v-model="month.cronEvery" label="1">{{text.Month.every}}</el-radio>
          </el-row>
          <el-row>
            <el-radio v-model="month.cronEvery" label="2">{{text.Month.interval[0]}}
              <el-input-number size="small" v-model="month.incrementIncrement" :min="0" :max="12"></el-input-number>
              {{text.Month.interval[1]}}
              <el-input-number size="small" v-model="month.incrementStart" :min="0" :max="12"></el-input-number>
            </el-radio>
          </el-row>
          <el-row>
            <el-radio class="long" v-model="month.cronEvery" label="3">{{text.Month.specific}}
              <el-select size="small" multiple v-model="month.specificSpecific">
                <el-option v-for="val in 12" :key="val" :label="val" :value="val.toString()"></el-option>
              </el-select>
            </el-radio>
          </el-row>
          <el-row>
            <el-radio v-model="month.cronEvery" label="4">{{text.Month.cycle[0]}}
              <el-input-number size="small" v-model="month.rangeStart" :min="1" :max="12"></el-input-number>
              {{text.Month.cycle[1]}}
              <el-input-number size="small" v-model="month.rangeEnd" :min="1" :max="12"></el-input-number>
              {{text.Month.cycle[2]}}
            </el-radio>
          </el-row>
        </div>
      </el-tab-pane>
      <el-tab-pane>
        <span slot="label"><i class="el-icon-date"></i> {{text.Week.name}}</span>
        <div class="tabBody">
          <el-row>
            <el-radio v-model="week.cronEvery" label="1">{{text.Week.every}}</el-radio>
          </el-row>
          <el-row>
            <el-radio class="long" v-model="week.cronEvery" label="3">{{text.Week.specific}}
              <el-select size="small" multiple v-model="week.specificSpecific">
                <el-option v-for="i in 7" :key="i" :label="text.Week.list[i - 1]" :value="i.toString()"></el-option>
              </el-select>
            </el-radio>
          </el-row>
          <el-row>
            <el-radio v-model="week.cronEvery" label="4">{{text.Week.cycle[0]}}
              <el-input-number size="small" v-model="week.rangeStart" :min="1" :max="7"></el-input-number>
              {{text.Week.cycle[1]}}
              <el-input-number size="small" v-model="week.rangeEnd" :min="1" :max="7"></el-input-number>
            </el-radio>
          </el-row>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>
<script>
import Language from './language/index'

export default {
  name: 'VueCronLinux',
  props: ['data', 'i18n'],
  data () {
    return {
      second: {
        cronEvery: '',
        incrementStart: '3',
        incrementIncrement: '5',
        rangeStart: '',
        rangeEnd: '',
        specificSpecific: [0]
      },
      minute: {
        cronEvery: '',
        incrementStart: '3',
        incrementIncrement: '5',
        rangeStart: '',
        rangeEnd: '',
        specificSpecific: ['0']
      },
      hour: {
        cronEvery: '',
        incrementStart: '3',
        incrementIncrement: '5',
        rangeStart: '',
        rangeEnd: '',
        specificSpecific: ['0']
      },
      day: {
        cronEvery: '',
        incrementStart: '1',
        incrementIncrement: '1',
        rangeStart: '',
        rangeEnd: '',
        specificSpecific: ['1'],
        cronLastSpecificDomDay: 1,
        cronDaysBeforeEomMinus: '',
        cronDaysNearestWeekday: ''
      },
      month: {
        cronEvery: '',
        incrementStart: '3',
        incrementIncrement: '5',
        rangeStart: '',
        rangeEnd: '',
        specificSpecific: ['1']
      },
      week: {
        cronEvery: '',
        incrementStart: '1',
        incrementIncrement: '1',
        specificSpecific: ['1'],
        cronNthDayDay: 1,
        cronNthDayNth: '1',
        rangeStart: '',
        rangeEnd: ''
      },
      output: {
        second: '',
        minute: '',
        hour: '',
        day: '',
        month: '',
        Week: '',
        year: ''
      }
    }
  },
  watch: {
    data () {
      this.updateCronFromData()
    },
    cron () {
      this.$emit('change', this.cron)
    }
  },
  computed: {
    text () {
      return Language[this.i18n || 'cn']
    },
    minutesText () {
      let minutes = ''
      let cronEvery = this.minute.cronEvery
      switch (cronEvery.toString()) {
        case '1':
          minutes = '*'
          break
        case '2':
          minutes = this.minute.incrementStart + '/' + this.minute.incrementIncrement
          break
        case '3':
          this.minute.specificSpecific.map(val => {
            minutes += val + ','
          })
          minutes = minutes.slice(0, -1)
          break
        case '4':
          minutes = this.minute.rangeStart + '-' + this.minute.rangeEnd
          break
      }
      return minutes
    },
    hoursText () {
      let hours = ''
      let cronEvery = this.hour.cronEvery
      switch (cronEvery.toString()) {
        case '1':
          hours = '*'
          break
        case '2':
          hours = this.hour.incrementStart + '/' + this.hour.incrementIncrement
          break
        case '3':
          this.hour.specificSpecific.map(val => {
            hours += val + ','
          })
          hours = hours.slice(0, -1)
          break
        case '4':
          hours = this.hour.rangeStart + '-' + this.hour.rangeEnd
          break
      }
      return hours
    },
    daysText () {
      let days = ''
      let cronEvery = this.day.cronEvery
      switch (cronEvery.toString()) {
        case '1':
          days = '*'
          break
        case '2':
          days = this.day.incrementStart + '/' + this.day.incrementIncrement
          break
        case '3':
          this.day.specificSpecific.map(val => {
            days += val + ','
          })
          days = days.slice(0, -1)
          break
        case '4':
          days = this.day.rangeStart + '-' + this.day.rangeEnd
          break
      }
      return days
    },
    monthsText () {
      let months = ''
      let cronEvery = this.month.cronEvery
      switch (cronEvery.toString()) {
        case '1':
          months = '*'
          break
        case '2':
          months = this.month.incrementStart + '/' + this.month.incrementIncrement
          break
        case '3':
          this.month.specificSpecific.map(val => {
            months += val + ','
          })
          months = months.slice(0, -1)
          break
        case '4':
          months = this.month.rangeStart + '-' + this.month.rangeEnd
          break
      }
      return months
    },
    weeksText () {
      let weeks = ''
      let cronEvery = this.week.cronEvery
      switch (cronEvery.toString()) {
        case '1':
          weeks = '*'
          break
        case '3':
          this.week.specificSpecific.map(val => {
            weeks += val + ','
          })
          weeks = weeks.slice(0, -1)
          break
        case '4':
          weeks = this.week.rangeStart + '-' + this.week.rangeEnd
          break
      }
      return weeks
    },
    cron () {
      return [this.minutesText, this.hoursText, this.daysText, this.monthsText, this.weeksText]
        .filter(v => !!v)
        .join(' ')
    }
  },
  methods: {
    getValue () {
      return this.cron
    },
    change () {
      this.$emit('change', this.cron)
      this.close()
    },
    close () {
      this.$emit('close')
    },
    updateCronItem (key, value) {
      if (value === undefined) {
        this[key].cronEvery = '0'
        return
      }
      if (value.match(/^\*$/)) {
        this[key].cronEvery = '1'
      } else if (value.match(/\//)) {
        this[key].cronEvery = '2'
        this[key].incrementStart = value.split('/')[0]
        this[key].incrementIncrement = value.split('/')[1]
      } else if (value.match(/,|^\d+$/)) {
        this[key].cronEvery = '3'
        this[key].specificSpecific = value.split(',')
      } else if (value.match(/-/)) {
        this[key].cronEvery = '4'
        this[key].rangeStart = value.split('-')[0]
        this[key].rangeEnd = value.split('-')[1]
      } else {
        this[key].cronEvery = '0'
      }
    },
    updateCronFromData () {
      const arr = this.data.split(' ')
      const minute = arr[0]
      const hour = arr[1]
      const day = arr[2]
      const month = arr[3]
      const week = arr[4]

      this.updateCronItem('minute', minute)
      this.updateCronItem('hour', hour)
      this.updateCronItem('day', day)
      this.updateCronItem('month', month)
      this.updateCronItem('week', week)
    }
  },
  mounted () {
    this.updateCronFromData()
  }
}</script>
