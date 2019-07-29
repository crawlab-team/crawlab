<template>
  <virtual-list
    class="log-view"
    :size="6"
    :remain="100"
    :item="item"
    :itemcount="logData.length"
    :itemprops="getItemProps"
  >
  </virtual-list>
</template>

<script>
import LogItem from './LogItem'
import VirtualList from 'vue-virtual-scroll-list'

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
      item: LogItem
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
    }
  },
  methods: {
    getItemProps (index) {
      const logItem = this.logData[index]
      return {
        // <item/> will render with itemProps.
        // https://vuejs.org/v2/guide/render-function.html#createElement-Arguments
        props: {
          index: logItem.index,
          data: logItem.data
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
    margin-top: 0!important;
    min-height: 100%;
    overflow-y: scroll;
    list-style: none;
    color: #A9B7C6;
    background: #2B2B2B;
  }

</style>
