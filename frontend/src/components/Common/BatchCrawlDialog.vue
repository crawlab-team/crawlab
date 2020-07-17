<template>
  <el-dialog
    :visible="visible"
    width="1200px"
    :before-close="beforeClose"
  >
    <el-table
      :data="batchCrawlList"
    >
      <el-table-column
        :label="$t('Spider')"
        width="150px"
      >
        <template slot-scope="scope">
          <el-select
            v-model="scope.row.spider_id"
            size="mini"
            filterable
            :placeholder="$t('Spider')"
            @change="onSpiderChange(scope.row, $event)"
          >
<!--            <el-option :label="$t('Same Above')" value="same-above" :placeholder="$t('Spider')"/>-->
            <el-option
              v-for="op in allSpiderList"
              :key="op._id"
              :label="`${op.display_name} (${op.name})`"
              :value="op._id"
            />
          </el-select>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('Run Type')"
        width="150px"
      >
        <template slot-scope="scope">
          <el-select v-model="scope.row.run_type" size="mini">
            <el-option value="all-nodes" :label="$t('All Nodes')"/>
            <el-option value="selected-nodes" :label="$t('Selected Nodes')"/>
            <el-option value="random" :label="$t('Random')"/>
          </el-select>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('Nodes')"
        width="250px"
      >
        <template v-if="scope.row.run_type === 'selected-nodes'" slot-scope="scope">
          <el-select
            v-model="scope.row.node_ids"
            size="mini"
            multiple
            :placeholder="$t('Nodes')"
          >
            <el-option
              v-for="n in activeNodeList"
              :key="n._id"
              :label="n.name"
              :value="n._id"
            />
          </el-select>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('Scrapy Spider')"
        width="150px"
      >
        <template v-if="getSpiderById(scope.row.spider_id).is_scrapy" slot-scope="scope">
          <el-select
            v-model="scope.row.scrapy_spider_name"
            size="mini"
            :placeholder="$t('Scrapy Spider')"
            :disabled="!scope.row.scrapy_spider_name"
          >
            <el-option
              v-for="(n, index) in getScrapySpiderNames(scope.row.spider_id)"
              :key="index"
              :label="n"
              :value="n"
            />
          </el-select>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('Scrapy Log Level')"
        width="120px"
      >
        <template v-if="getSpiderById(scope.row.spider_id).is_scrapy" slot-scope="scope">
          <el-select
            v-model="scope.row.scrapy_log_level"
            :placeholder="$t('Scrapy Log Level')"
            size="mini"
          >
            <el-option value="INFO" label="INFO" />
            <el-option value="DEBUG" label="DEBUG" />
            <el-option value="WARN" label="WARN" />
            <el-option value="ERROR" label="ERROR" />
          </el-select>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('Parameters')"
        min-width="150px"
      >
        <template slot-scope="scope">
          <el-input v-model="scope.param" size="mini" :placeholder="$t('Parameters')"/>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('Action')"
        fixed="right"
        width="150px"
      >
        <template slot-scope="scope">
          <el-button icon="el-icon-plus" size="mini" type="primary" @click="onAdd(scope.$index)"/>
          <el-button icon="el-icon-delete" size="mini" type="danger" @click="onRemove(scope.$index)"/>
        </template>
      </el-table-column>
    </el-table>
    <template slot="footer">
      <el-button type="plain" size="small" @click="$emit('close')">{{ $t('Cancel') }}</el-button>
      <el-button type="plain" size="small" @click="reset">
        {{ $t('Reset') }}
      </el-button>
      <el-button type="primary" size="small" :disabled="isConfirmDisabled" @click="onConfirm">
        {{ $t('Confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script>
  import {
    mapState
  } from 'vuex'

  export default {
    name: 'BatchCrawlDialog',
    props: {
      visible: {
        type: Boolean,
        default: false
      }
    },
    data() {
      return {
        scrapySpidersNamesDict: {}
      }
    },
    computed: {
      ...mapState('task', [
        'batchCrawlList'
      ]),
      ...mapState('spider', [
        'allSpiderList'
      ]),
      ...mapState('node', [
        'nodeList'
      ]),
      isBatchCrawlDialogVisible: {
        get() {
          return this.$store.state.task.isBatchCrawlDialogVisible
        },
        set(value) {
          this.$store.commit('task/SET_IS_BATCH_CRAWL_DIALOG_VISIBLE', value)
        }
      },
      activeNodeList() {
        return this.nodeList.filter(n => n.status === 'online')
      },
      validCrawlList() {
        return this.batchCrawlList.filter(d => !!d.spider_id)
      },
      isConfirmDisabled() {
        for (let i = 0; i < this.validCrawlList.length; i++) {
          const row = this.validCrawlList[i]
          const spider = this.getSpiderById(row.spider_id)
          if (!spider) {
            return true
          }
          if (spider.is_scrapy && !row.scrapy_spider_name) {
            return true
          }
        }
        return false
      },
      scrapySpidersIds() {
        return Array.from(new Set(this.validCrawlList.filter(d => {
          const spider = this.getSpiderById(d.spider_id)
          return spider && spider.is_scrapy
        }).map(d => d.spider_id)))
      }
    },
    watch: {
      visible() {
        if (this.visible) {
          this.fetchAllScrapySpiderNames()
        }
      }
    },
    methods: {
      beforeClose() {
        this.$emit('close')
      },
      reset() {
        this.$store.commit('task/SET_BATCH_CRAWL_LIST', [])
        for (let i = 0; i < 10; i++) {
          this.batchCrawlList.push({
            spider_id: '',
            run_type: 'random',
            param: '',
            scrapy_log_level: 'INFO'
          })
        }
      },
      getSpiderById(id) {
        return this.allSpiderList.filter(d => d._id === id)[0] || {}
      },
      async onSpiderChange(row, id) {
        const spider = this.getSpiderById(id)
        if (!spider) return
        if (spider.is_scrapy) {
          await this.fetchScrapySpiderNames(id)
          if (this.scrapySpidersNamesDict[id] && this.scrapySpidersNamesDict[id].length > 0) {
            this.$set(row, 'scrapy_spider_name', this.scrapySpidersNamesDict[id][0])
          }
        }
      },
      getScrapySpiderNames(id) {
        if (!this.scrapySpidersNamesDict[id]) return []
        return this.scrapySpidersNamesDict[id]
      },
      async onConfirm() {
        await this.$request.put('/tasks/batch', this.validCrawlList.map(d => {
          const spider = this.getSpiderById(d.spider_id)
          // Scrapy爬虫特殊处理
          if (spider.type === 'customized' && spider.is_scrapy) {
            d.param = `${this.scrapySpidersNamesDict[d.spider_id] ? this.scrapySpidersNamesDict[d.spider_id][0] : ''} --loglevel=${d.scrapy_log_level} ${d.param || ''}`
          }
          return d
        }))
        this.reset()
        this.$emit('close')
      },
      async fetchScrapySpiderNames(id) {
        if (!this.scrapySpidersNamesDict[id]) {
          const res = await this.$request.get(`/spiders/${id}/scrapy/spiders`)
          this.$set(this.scrapySpidersNamesDict, id, res.data.data)
        }
      },
      async fetchAllScrapySpiderNames() {
        await Promise.all(this.scrapySpidersIds.map(async id => {
          await this.fetchScrapySpiderNames(id)
        }))
        this.validCrawlList.filter(d => {
          const spider = this.getSpiderById(d.spider_id)
          return spider && spider.is_scrapy
        }).forEach(row => {
          const id = row.spider_id
          if (this.scrapySpidersNamesDict[id] && this.scrapySpidersNamesDict[id].length > 0) {
            this.$set(row, 'scrapy_spider_name', this.scrapySpidersNamesDict[id][0])
          }
        })
      },
      onAdd(rowIndex) {
        this.batchCrawlList.splice(rowIndex + 1, 0, {
          spider_id: '',
          run_type: 'random',
          param: '',
          scrapy_log_level: 'INFO'
        })
      },
      onRemove(rowIndex) {
        this.batchCrawlList.splice(rowIndex, 1)
      }
    }
  }
</script>

<style scoped>
  .el-table .el-button {
    padding: 7px;
  }
</style>
