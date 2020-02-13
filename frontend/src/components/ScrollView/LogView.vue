<template>
  <div class="log-view-container">
    <div class="filter-wrapper">
      <div class="left">
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
          style="width: 240px; margin-right: 10px"
        />
      </div>
      <div class="right">
        <el-badge
          v-if="errorLogData.length > 0"
          :value="errorLogData.length"
        >
          <el-button
            type="danger"
            size="small"
            icon="el-icon-warning-outline"
            @click="toggleErrors"
          >
            {{$t('Error Count')}}
          </el-button>
        </el-badge>
      </div>
    </div>
    <div class="content">
      <div
        class="log-view-wrapper"
        :class="isErrorsCollapsed ? 'errors-collapsed' : ''"
      >
        <virtual-list
          class="log-view"
          ref="log-view"
          :start="currentLogIndex - 1"
          :offset="0"
          :size="18"
          :remain="remainSize"
          :item="item"
          :itemcount="filteredLogData.length"
          :itemprops="getItemProps"
          :tobottom="onToBottom"
          :onscroll="onScroll"
        />
      </div>
      <div
        v-show="!isErrorsCollapsed && !isErrorCollapsing"
        class="errors-wrapper"
        :class="isErrorsCollapsed ? 'collapsed' : ''"
      >
        <ul class="error-list">
          <li
            v-for="item in errorLogData"
            :key="item.index"
            class="error-item"
            :class="currentLogIndex === item.index ? 'active' : ''"
            @click="onClickError(item)"
          >
            <span class="line-no">
              {{item.index}}
            </span>
            <span class="line-content">
              {{item.data}}
            </span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
import {
  mapState,
  mapGetters
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
      isScrolling2nd: false,
      errorRegex: this.$utils.log.errorRegex,
      currentOffset: 0,
      isErrorsCollapsed: true,
      isErrorCollapsing: false
    }
  },
  computed: {
    ...mapState('task', [
      'taskForm'
    ]),
    ...mapGetters('task', [
      'logData',
      'errorLogData'
    ]),
    currentLogIndex: {
      get () {
        return this.$store.state.task.currentLogIndex
      },
      set (value) {
        this.$store.commit('task/SET_CURRENT_LOG_INDEX', value)
      }
    },
    filteredLogData () {
      return this.logData.filter(d => {
        if (!this.searchString) return true
        return !!d.data.toLowerCase().match(this.searchString.toLowerCase())
      })
    },
    remainSize () {
      const height = document.querySelector('body').clientHeight
      return (height - 240) / 18
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
          active: logItem.active,
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
    },
    toggleErrors () {
      this.isErrorsCollapsed = !this.isErrorsCollapsed
      this.isErrorCollapsing = true
      setTimeout(() => {
        this.isErrorCollapsing = false
      }, 300)
    },
    onClickError (item) {
      this.currentLogIndex = item.index
      this.isToBottom = false
      const handle = setInterval(() => {
        this.isToBottom = false
      }, 10)
      setTimeout(() => {
        clearInterval(handle)
      }, 500)
    }
  },
  mounted () {
    this.handle = setInterval(() => {
      if (this.isToBottom) {
        this.toBottom()
      }
    }, 500)
  },
  destroyed () {
    clearInterval(this.handle)
  }
}
</script>

<style scoped>
  .filter-wrapper {
    display: flex;
    justify-content: space-between;
    margin-bottom: 10px;
  }

  .content {
    display: block;
  }

  .log-view-wrapper {
    float: left;
    flex-basis: calc(100% - 240px);
    width: calc(100% - 300px);
    transition: width 0.3s;
  }

  .log-view-wrapper.errors-collapsed {
    flex-basis: 100%;
    width: 100%;
  }

  .log-view {
    margin-top: 0 !important;
    overflow-y: scroll !important;
    list-style: none;
    color: #A9B7C6;
    background: #2B2B2B;
    border: none;
  }

  .errors-wrapper {
    float: left;
    display: inline-block;
    margin: 0;
    padding: 0;
    flex-basis: 240px;
    width: 300px;
    transition: opacity 0.3s;
    border-top: 1px solid #DCDFE6;
    border-right: 1px solid #DCDFE6;
    border-bottom: 1px solid #DCDFE6;
    height: calc(100vh - 240px);
    font-size: 16px;
  }

  .errors-wrapper.collapsed {
    width: 0;
  }

  .errors-wrapper .error-list {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .errors-wrapper .error-list .error-item {
    white-space: nowrap;
    text-overflow: ellipsis;
    overflow: hidden;
    /*height: 18px;*/
    border-bottom: 1px solid white;
    padding: 5px 0;
    background: #F56C6C;
    color: white;
    cursor: pointer;
  }

  .errors-wrapper .error-list .error-item.active {
    background: #E6A23C;
    font-weight: bolder;
    text-decoration: underline;
  }

  .errors-wrapper .error-list .error-item:hover {
    font-weight: bolder;
    text-decoration: underline;
  }

  .errors-wrapper .error-list .error-item .line-no {
    display: inline-block;
    text-align: right;
    width: 70px;
  }

  .errors-wrapper .error-list .error-item .line-content {
    display: inline;
    width: calc(100% - 70px);
    padding-left: 10px;
  }
</style>
