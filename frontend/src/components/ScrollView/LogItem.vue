<template>
  <div class="log-item" :style="style">
    <div class="line-no">{{index}}</div>
    <div class="line-content">
      <span v-if="isLogEnd" style="color: #E6A23C" class="loading-text">
        {{$t('Updating log...')}}
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
    }
  },
  data () {
    const token = ' :,.'
    return {
      errorRegex: new RegExp(`(?:[${token}]|^)((?:error|exception|traceback)s?)(?:[${token}]|$)`, 'gi')
      // errorRegex: new RegExp('(error|exception|traceback)', 'gi')
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
    display: table;
  }

  .log-item:first-child .line-no {
    padding-top: 10px;
  }

  .log-item .line-no {
    display: table-cell;
    color: #A9B7C6;
    background: #313335;
    padding-right: 10px;
    text-align: right;
    flex-basis: 40px;
    width: 70px;
  }

  .log-item .line-content {
    padding-left: 10px;
    display: table-cell;
    /*display: inline-block;*/
    word-break: break-word;
    flex-basis: calc(100% - 50px);
  }

  .loading-text {
    animation: blink;
  }

  @keyframes blink {
    0% {
      opacity: 0;
    }

    100% {
      opacity: 100%;
    }
  }
</style>
