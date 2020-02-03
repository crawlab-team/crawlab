<template>
  <div class="general-table-view">
    <el-table
      :data="filteredData"
      :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
      border>
      <template v-for="col in columns">
        <el-table-column :key="col" :label="col" :property="col" min-width="120">
          <template slot-scope="scope">
            <el-popover trigger="hover" :content="getString(scope.row[col])" popper-class="cell-popover">
              <div slot="reference" class="wrapper">
                {{getString(scope.row[col])}}
              </div>
            </el-popover>
          </template>
        </el-table-column>
      </template>
    </el-table>
    <div class="pagination">
      <el-pagination
        @current-change="onPageChange"
        @size-change="onPageChange"
        :current-page.sync="pageNum"
        :page-sizes="[10, 20, 50, 100]"
        :page-size.sync="pageSize"
        layout="sizes, prev, pager, next"
        :total="total">
      </el-pagination>
    </div>
  </div>
</template>

<script>
export default {
  name: 'GeneralTableView',
  data () {
    return {}
  },
  props: {
    pageNum: {
      type: Number,
      default: 1
    },
    pageSize: {
      type: Number,
      default: 10
    },
    total: {
      type: Number,
      default: 0
    },
    columns: {
      type: Array,
      default () {
        return []
      }
    },
    data: {
      type: Array,
      default () {
        return []
      }
    }
  },
  computed: {
    filteredData () {
      return this.data
    }
  },
  methods: {
    onPageChange () {
      this.$emit('page-change', { pageNum: this.pageNum, pageSize: this.pageSize })
    },
    getString (value) {
      if (value === undefined) return ''
      const str = JSON.stringify(value)
      if (str.match(/^"(.*)"$/)) return str.match(/^"(.*)"$/)[1]
      return str
    }
  }
}
</script>

<style scoped>
  .general-table-view >>> .cell .wrapper:hover {
    text-decoration: underline;
  }

  .general-table-view >>> .cell .wrapper {
    font-size: 12px;
    height: 24px;
    line-height: 24px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>

<style>
  .cell-popover {
    max-width: 480px;
  }
</style>
