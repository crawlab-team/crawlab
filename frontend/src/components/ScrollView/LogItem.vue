<template>
  <div class="log-item" :style="style" :class="`log-item-${source.index} ${source.active ? 'active' : ''}`">
    <div class="line-no">{{ source.index }}</div>
    <div class="line-content">
      <span v-if="isLogEnd" style="color: #E6A23C">
        <span class="loading-text">{{ $t('Updating log...') }}</span>
        <i class="el-icon-loading" />
      </span>
      <span v-else-if="source.isAnsi" v-html="dataHtml" />
      <span v-else v-html="dataHtml" />
    </div>
  </div>
</template>

<script>
  import {
    mapGetters
  } from 'vuex'

  export default {
    name: 'LogItem',
    props: {
      source: {
        type: Object,
        default() {
          return {}
        }
      }
    },
    data() {
      return {
      }
    },
    computed: {
      ...mapGetters('user', [
        'userInfo'
      ]),
      errorRegex() {
        if (!this.userInfo.setting.error_regex_pattern) {
          return this.$utils.log.errorRegex
        }
        console.log(this.userInfo.setting.error_regex_pattern)
        return new RegExp(this.userInfo.setting.error_regex_pattern, 'i')
      },
      dataHtml() {
        let html = this.source.data.replace(this.errorRegex, ' <span style="font-weight: bolder; text-decoration: underline">$1</span> ')
        if (!this.source.searchString) return html
        html = html.replace(new RegExp(`(${this.source.searchString})`, 'gi'), '<mark>$1</mark>')
        return html
      },
      style() {
        let color = ''
        if (this.source.data.match(this.errorRegex)) {
          color = '#F56C6C'
        }
        return {
          color
        }
      },
      isLogEnd() {
        return this.source.data === '###LOG_END###'
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
