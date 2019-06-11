<template>
  <div class="log-item">
    <div class="line-no">{{index}}</div>
    <div class="line-content">{{data}}</div>
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
    }
  }
}
</script>

<style scoped>
  .log-item {
    display: flex;
  }

  .log-item .line-no {
    margin-right: 10px;
    text-align: right;
    flex-basis: 40px;
  }

  .log-item .line-content {
    display: inline-block;
    flex-basis: calc(100% - 50px);
  }
</style>
