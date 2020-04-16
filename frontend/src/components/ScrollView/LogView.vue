<template>
  <div class="log-view-container">
    <div class="filter-wrapper">
      <div class="left">
        <el-switch
          v-model="isLogAutoScroll"
          :inactive-text="$t('Auto-Scroll')"
          style="margin-right: 10px"
        >
        </el-switch>
        <el-input
          v-model="logKeyword"
          size="small"
          suffix-icon="el-icon-search"
          :placeholder="$t('Search Log')"
          style="width: 240px; margin-right: 10px"
        />
        <el-button
          size="small"
          type="primary"
          icon="el-icon-search"
          @click="onSearchLog"
        >
          {{$t('Search Log')}}
        </el-button>
      </div>
      <div class="right">
        <el-pagination
          size="small"
          :total="taskLogTotal"
          :current-page.sync="taskLogPage"
          :page-sizes="[1000, 2000, 5000, 10000]"
          :page-size.sync="taskLogPageSize"
          :pager-count="3"
          layout="sizes, prev, pager, next"
        />
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
            <span class="line-content">
              {{item.msg}}
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
      'taskForm',
      'taskLogTotal',
      'logKeyword'
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
    logKeyword: {
      get () {
        return this.$store.state.task.logKeyword
      },
      set (value) {
        this.$store.commit('task/SET_LOG_KEYWORD', value)
      }
    },
    taskLogPage: {
      get () {
        return this.$store.state.task.taskLogPage
      },
      set (value) {
        this.$store.commit('task/SET_TASK_LOG_PAGE', value)
      }
    },
    taskLogPageSize: {
      get () {
        return this.$store.state.task.taskLogPageSize
      },
      set (value) {
        this.$store.commit('task/SET_TASK_LOG_PAGE_SIZE', value)
      }
    },
    isLogAutoScroll: {
      get () {
        return this.$store.state.task.isLogAutoScroll
      },
      set (value) {
        this.$store.commit('task/SET_IS_LOG_AUTO_SCROLL', value)
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
    taskLogPage () {
      this.$emit('search')
      this.$st.sendEv('任务详情', '日志', '改变页数')
    },
    taskLogPageSize () {
      this.$emit('search')
      this.$st.sendEv('任务详情', '日志', '改变日志每页条数')
    },
    isLogAutoScroll () {
      if (this.isLogAutoScroll) {
        this.$store.dispatch('task/getTaskLog', {
          id: this.$route.params.id,
          keyword: this.logKeyword
        }).then(() => {
          this.toBottom()
        })
        this.$st.sendEv('任务详情', '日志', '点击自动滚动')
      } else {
        this.$st.sendEv('任务详情', '日志', '取消自动滚动')
      }
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
          logItem,
          data: isAnsi ? convert.toHtml(logItem.data) : logItem.data,
          searchString: this.logKeyword,
          active: logItem.active,
          isAnsi
        }
      }
    },
    onToBottom () {
    },
    onScroll () {
    },
    toBottom () {
      this.$el.querySelector('.log-view').scrollTo({ top: 99999999 })
    },
    onAutoScroll () {

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
      this.isLogAutoScroll = false
      const handle = setInterval(() => {
        this.isLogAutoScroll = false
      }, 10)
      setTimeout(() => {
        clearInterval(handle)
      }, 500)
    },
    onSearchLog () {
      this.$emit('search')
      this.$st.sendEv('任务详情', '日志', '搜索日志')
    }
  },
  mounted () {
    this.currentLogIndex = 0
    this.handle = setInterval(() => {
      if (this.isLogAutoScroll) {
        this.toBottom()
      }
    }, 200)
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
    overflow: auto;
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

  .right {
    display: flex;
    align-items: center;
  }

  .right .el-pagination {
    margin-right: 10px;
  }
</style>
