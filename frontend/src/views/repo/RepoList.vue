<template>
  <div class="app-container repo-list">
    <el-card>
      <div class="filter">
        <el-form inline>
          <el-form-item :label="$t('Search Keyword')">
            <el-input
              v-model="keyword"
              size="small"
              :placeholder="$t('Search Keyword')"
              @keyup.enter.native="getRepos"
            />
          </el-form-item>
          <el-form-item :label="$t('Sort')">
            <el-select
              v-model="sortKey"
              size="small"
            >
              <el-option :label="$t('Default Sort')" value=""/>
              <el-option :label="$t('Most Stars')" value="stars"/>
              <el-option :label="$t('Most Forks')" value="forks"/>
              <el-option :label="$t('Latest Pushed')" value="pushed_at"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="success" size="small" @click="getRepos">{{ $t('Search') }}</el-button>
          </el-form-item>
        </el-form>
      </div>
      <el-table
        ref="table"
        :data="repos"
        :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
        :row-class-name="getRowClassName"
        row-key="id"
        border
        v-loading="isLoading"
        @expand-change="onRowExpand"
      >
        <el-table-column type="expand">
          <template slot-scope="scope">
            <ul class="sub-dir-list">
              <li
                v-for="sub in getSubDirList(scope.row)"
                :key="sub.full_name"
                class="sub-dir-item"
              >
                <div class="sub-dir-title">
                  {{ sub.name }}
                </div>
                <div class="action">
                  <el-tooltip :content="$t('Download')" placement="top">
                    <el-button
                      type="primary"
                      icon="fa fa-download"
                      size="mini"
                      @click="onDownload(scope.row, sub.full_name, $event)"
                      v-loading="scope.row.loading"
                    />
                  </el-tooltip>
                </div>
              </li>
            </ul>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('Name')"
          prop="full_name"
          width="300px"
        />
        <el-table-column
          :label="$t('Description')"
          prop="description"
          min-width="500px"
        />
        <el-table-column
          label="Stars"
          prop="stars"
          width="80px"
          align="right"
        />
        <el-table-column
          label="Forks"
          prop="forks"
          width="80px"
          align="right"
        />
        <el-table-column
          :label="$t('Pushed At')"
          prop="pushed_at"
          width="150px"
        >
          <template slot-scope="scope">
            {{ getTime(scope.row.pushed_at) }}
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('Action')"
          width="120px"
        >
          <template slot-scope="scope">
            <el-tooltip :content="$t('Download')" placement="top">
              <el-button
                type="primary"
                icon="fa fa-download"
                size="mini"
                @click="onDownload(scope.row, null, $event)"
                v-loading="scope.row.loading"
              />
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          :current-page.sync="pageNum"
          :page-sizes="[10, 20, 50, 100]"
          :page-size.sync="pageSize"
          layout="sizes, prev, pager, next"
          :total="total"
        />
      </div>
    </el-card>
  </div>
</template>

<script>
  import dayjs from 'dayjs'

  export default {
    name: 'RepoList',
    data() {
      return {
        repos: [],
        total: 0,
        pageNum: 1,
        pageSize: 20,
        keyword: '',
        sortKey: '',
        isLoading: false,
        subDirCache: {}
      }
    },
    watch: {
      pageNum() {
        this.getRepos()
      },
      pageSize() {
        this.getRepos()
      },
      sortKey() {
        this.getRepos()
      }
    },
    async created() {
      await this.getRepos()
    },
    methods: {
      async getRepos() {
        this.isLoading = true
        try {
          const res = await this.$request.get('/repos', {
            page_num: this.pageNum,
            page_size: this.pageSize,
            keyword: this.keyword,
            sort_key: this.sortKey
          })
          this.repos = res.data.data
          this.total = res.data.total
        } finally {
          this.isLoading = false
        }
      },
      getTime(t) {
        return dayjs(t).format('YYYY-MM-DD HH:mm:ss')
      },
      onDownload(row, fullName, ev) {
        ev.stopPropagation()
        this.$confirm(this.$t('Are you sure to download this spider?'), this.$t('Notification'), {
          confirmButtonText: this.$t('Confirm'),
          cancelButtonText: this.$t('Cancel'),
          type: 'warning'
        }).then(async() => {
          this.$set(row, 'loading', true)
          try {
            await this.download(fullName || row.full_name)
            this.$message.success('Downloaded successfully')
            this.$st.sendEv('爬虫市场', '下载仓库')
          } finally {
            this.$set(row, 'loading', false)
          }
        })
      },
      async download(fullName) {
        this.$request.post('/repos/download', {
          full_name: fullName
        })
      },
      getRowClassName({ row }) {
        return row.is_sub_dir ? '' : 'non-expandable'
      },
      async onRowExpand(row, expandedRows) {
        if (!this.subDirCache[row.full_name]) {
          const res = await this.$request.get('/repos/sub-dir', {
            full_name: row.full_name
          })
          this.$set(this.subDirCache, row.full_name, res.data.data)
        }
        this.$st.sendEv('爬虫市场', '点击展开')
      },
      getSubDirList(row) {
        if (!this.subDirCache[row.full_name]) return []
        return this.subDirCache[row.full_name].map(n => {
          return {
            name: n,
            full_name: `${row.full_name}/${n}`
          }
        })
      }
    }
  }
</script>

<style scoped>
  .el-table .el-button {
    padding: 7px;
  }

  .el-table >>> .non-expandable .el-table__expand-icon {
    display: none;
  }

  .el-table .sub-dir-list {
    list-style: none;
    font-size: 12px;
    margin: 0;
    padding: 0;
  }

  .el-table .sub-dir-list .sub-dir-item {
    padding: 6px 0 6px 60px;
    border-bottom: 1px dashed #EBEEF5;
    display: flex;
    justify-content: space-between;
  }

  .el-table .sub-dir-list .sub-dir-item:last-child {
    border-bottom: none;
  }

  .el-table .sub-dir-list .sub-dir-item .sub-dir-title {
    line-height: 29px;
  }

  .el-table .sub-dir-list .sub-dir-item .action {
    width: 120px;
    padding: 0 12px;
  }
</style>
