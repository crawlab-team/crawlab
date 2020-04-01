<template>
  <div class="log-item" :style="style" :class="`log-item-${index} ${active ? 'active' : ''}`">
    <div class="line-no">{{index}}</div>
    <div class="line-content">
      <span v-if="isLogEnd" style="color: #E6A23C">
        <span class="loading-text">{{$t('Updating log...')}}</span>
        <i class="el-icon-loading"></i>
      </span>
      <span v-else-if="isAnsi" v-html="dataHtml"></span>
      <span v-else v-html="dataHtml"></span>
    </div>
  </div>
</template>

<script>
export default {
  name: 'LogItem',
  props: {
    index: {
      type: Number,
      default: 1
    },
    logItem: {
      type: Object,
      default () {
        return {}
      }
    },
    data: {
      type: String,
      default: ''
    },
    isAnsi: {
      type: Boolean,
      default: false
    },
    searchString: {
      type: String,
      default: ''
    },
    active: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      errorRegex: this.$utils.log.errorRegex
    }
  },
  computed: {
    dataHtml () {
      let html = this.data.replace(this.errorRegex, ' <span style="font-weight: bolder; text-decoration: underline">$1</span> ')
      if (!this.searchString) return html
      html = html.replace(new RegExp(`(${this.searchString})`, 'gi'), '<mark>$1</mark>')
      return html
    },
    style () {
      let color = ''
      if (this.data.match(this.errorRegex)) {
        color = '#F56C6C'
      }
      return {
        color
      }
    },
    isLogEnd () {
      return this.data === '###LOG_END###'
    }
  }
}
</script>

<style scoped>
  .log-item {
    display: block;
  }

  .log-item:hover {
    background: rgba(55, 57, 59, 0.5);
  }

  .log-item:first-child .line-no {
    padding-top: 10px;
    text-align: right;
  }

  .log-item .line-no {
    display: inline-block;
    width: 70px;
    color: #A9B7C6;
    background: #313335;
    padding-right: 10px;
    text-align: right;
  }

  .log-item.active .line-no {
    background: #E6A23C;
    color: white;
    font-weight: bolder;
  }

  .log-item .line-content {
    padding-left: 10px;
    display: inline-block;
    width: calc(100% - 70px);
    white-space: nowrap;
  }

  .loading-text {
    margin-right: 5px;
    animation: blink 2s ease-in infinite;
  }

  @keyframes blink {
    0% {
      opacity: 1;
    }

    50% {
      opacity: 0.3;
    }

    100% {
      opacity: 1;
    }
  }
</style>
