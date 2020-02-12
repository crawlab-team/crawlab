<template>
  <div class="log-view-wrapper">
    <div class="filter-wrapper">
      <el-button
        size="small"
        type="primary"
        icon="el-icon-download"
        style="margin-right: 10px"
        :disabled="isToBottom"
        @click="onAutoScroll"
      >
        {{$t('Auto-Scroll')}}
      </el-button>
      <el-input
        v-model="searchString"
        size="small"
        suffix-icon="el-icon-search"
        :placeholder="$t('Search Log')"
        style="width: 240px"
      />
    </div>
    <div class="log-view-wrapper" ref="log-view-wrapper">
      <virtual-list
        class="log-view"
        ref="log-view"
        :size="6"
        :remain="100"
        :item="item"
        :itemcount="filteredLogData.length"
        :itemprops="getItemProps"
        :tobottom="onToBottom"
        :onscroll="onScroll"
      />
    </div>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import VirtualList from 'vue-virtual-scroll-list'
import Convert from 'ansi-to-html'
import hasAnsi from 'has-ansi'

import LogItem from './LogItem'

const convert = new Convert()
export default {
  name: 'LogView',
  components: {
    VirtualList
  },
  props: {
    data: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      item: LogItem,
      searchString: '',
      isToBottom: false,
      isScrolling: false,
      isScrolling2nd: false
    }
  },
  computed: {
    ...mapState('task', [
      'taskForm'
    ]),
    logData () {
      const data = this.data.split('\n')
        .map((d, i) => {
          return {
            index: i + 1,
            data: d
          }
        })
      if (this.taskForm && this.taskForm.status === 'running') {
        data.push({
          index: data.length + 1,
          data: '###LOG_END###'
        })
        data.push({
          index: data.length + 2,
          data: ''
        })
      }
      return data
    },
    filteredLogData () {
      return this.logData.filter(d => {
        if (!this.searchString) return true
        return !!d.data.toLowerCase().match(this.searchString.toLowerCase())
      })
    }
  },
  watch: {
    searchString () {
      this.$st.sendEv('任务详情', '日志', '搜索日志')
    }
  },
  methods: {
    getItemProps (index) {
      const logItem = this.filteredLogData[index]
      const isAnsi = hasAnsi(logItem.data)
      return {
        // <item/> will render with itemProps.
        // https://vuejs.org/v2/guide/render-function.html#createElement-Arguments
        props: {
          index: logItem.index,
          data: isAnsi ? convert.toHtml(logItem.data) : logItem.data,
          searchString: this.searchString,
          isAnsi
        }
      }
    },
    onToBottom () {
      if (this.isScrolling) return
      this.isToBottom = true
    },
    onScroll () {
      if (this.isScrolling2nd) {
        this.isToBottom = false
      }
      this.isScrolling = true
      setTimeout(() => {
        this.isScrolling2nd = true
        setTimeout(() => {
          this.isScrolling2nd = false
        }, 50)
      }, 50)
      setTimeout(() => {
        this.isScrolling = false
      }, 100)
    },
    toBottom () {
      this.$el.querySelector('.log-view').scrollTo({ top: 99999999 })
      setTimeout(() => {
        this.isToBottom = true
      }, 50)
    },
    onAutoScroll () {
      this.toBottom()
    }
  },
  mounted () {
    this.handle = setInterval(() => {
      if (this.isToBottom) {
        this.toBottom()
      }
    }, 100)
  },
  destroyed () {
    clearInterval(this.handle)
  }
}
</script>

<style scoped>
  .log-view {
    margin-top: 0 !important;
    overflow-y: scroll !important;
    list-style: none;
    color: #A9B7C6;
    background: #2B2B2B;
  }

  .filter-wrapper {
    margin-bottom: 10px;
  }
</style>
