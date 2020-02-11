<template>
  <div class="log-item">
    <div class="line-no">{{index}}</div>
    <div class="line-content">
      <span v-if="isAnsi" v-html="dataHtml"></span>
      <span v-else="" v-html="dataHtml"></span>
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
  computed: {
    dataHtml () {
      if (!this.searchString) return this.data
      return this.data.replace(new RegExp(`(${this.searchString})`, 'gi'), '<mark>$1</mark>')
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
</style>
