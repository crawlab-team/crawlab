<template>
  <div class="log-view-wrapper">
    <div class="filter-wrapper">
      <el-input
        v-model="searchString"
        suffix-icon="el-icon-search"
        :placeholder="$t('Search Log')"
        style="width: 240px"
      />
    </div>
    <virtual-list
      class="log-view"
      :size="6"
      :remain="100"
      :item="item"
      :itemcount="filteredLogData.length"
      :itemprops="getItemProps"
    />
  </div>
</template>

<script>
import LogItem from './LogItem'
import VirtualList from 'vue-virtual-scroll-list'
import Convert from 'ansi-to-html'
import hasAnsi from 'has-ansi'

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
      searchString: ''
    }
  },
  computed: {
    logData () {
      return this.data.split('\n')
        .map((d, i) => {
          return {
            index: i + 1,
            data: d
          }
        })
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
    }
  },
  mounted () {
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
