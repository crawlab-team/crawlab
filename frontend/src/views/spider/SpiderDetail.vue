<template>
  <div class="app-container">
    <!--tour-->
    <v-tour
      name="spider-detail"
      :steps="tourSteps"
      :callbacks="tourCallbacks"
      :options="$utils.tour.getOptions(true)"
    />
    <!--./tour-->

    <!--selector-->
    <div class="selector">
      <label class="label">{{$t('Spider')}}: </label>
      <el-select id="spider-select" v-model="spiderForm._id" @change="onSpiderChange">
        <el-option v-for="op in spiderList" :key="op._id" :value="op._id" :label="op.name"></el-option>
      </el-select>
    </div>

    <!--tabs-->
    <el-tabs v-model="activeTabName" @tab-click="onTabClick" type="border-card">
      <el-tab-pane :label="$t('Overview')" name="overview">
        <spider-overview/>
      </el-tab-pane>
      <el-tab-pane v-if="isGit" :label="$t('Git')" name="git-settings">
        <git-settings/>
      </el-tab-pane>
      <el-tab-pane v-if="isScrapy" :label="$t('Scrapy Settings')" name="scrapy-settings">
        <spider-scrapy
          @click-spider="onClickScrapySpider"
          @click-pipeline="onClickScrapyPipeline"
        />
      </el-tab-pane>
      <el-tab-pane v-if="isConfigurable" :label="$t('Config')" name="config">
        <config-list ref="config" @convert="onConvert"/>
      </el-tab-pane>
      <el-tab-pane :label="$t('Files')" name="files">
        <file-list
          ref="file-list"
        />
      </el-tab-pane>
      <el-tab-pane :label="$t('Environment')" name="environment">
        <environment-list/>
      </el-tab-pane>
      <el-tab-pane :label="$t('Analytics')" name="analytics">
        <spider-stats ref="spider-stats"/>
      </el-tab-pane>
      <el-tab-pane :label="$t('Schedules')" name="schedules">
        <spider-schedules/>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import FileList from '../../components/File/FileList'
import SpiderOverview from '../../components/Overview/SpiderOverview'
import EnvironmentList from '../../components/Environment/EnvironmentList'
import SpiderStats from '../../components/Stats/SpiderStats'
import ConfigList from '../../components/Config/ConfigList'
import SpiderSchedules from './SpiderSchedules'
import SpiderScrapy from '../../components/Scrapy/SpiderScrapy'
import GitSettings from '../../components/Settings/GitSettings'

export default {
  name: 'SpiderDetail',
  components: {
    GitSettings,
    SpiderScrapy,
    SpiderSchedules,
    ConfigList,
    SpiderStats,
    EnvironmentList,
    FileList,
    SpiderOverview
  },
  watch: {
    configListTs () {
      this.onConvert()
    }
  },
  data () {
    return {
      activeTabName: 'overview',
      tourSteps: [
        // top bar
        {
          target: '.el-tabs__nav.is-top',
          content: this.$t('You can switch to each section of the spider detail.')
        },
        {
          target: '#spider-select',
          content: this.$t('You can switch to different spider using this selector.')
        },
        // overview
        {
          target: '.task-table-view',
          content: this.$t('You can view latest tasks for this spider and click each row to view task detail.'),
          params: {
            placement: 'right'
          }
        },
        {
          target: '.spider-form',
          content: this.$t('You can edit the detail info for this spider.'),
          params: {
            placement: 'left'
          }
        },
        {
          target: '.button-container',
          content: this.$t('Here you can action on the spider, including running a task, uploading a zip file and save the spider info.'),
          params: {
            placement: 'top'
          }
        },
        // file
        {
          target: '.tree',
          content: this.$t('File navigation panel.<br><br>You can right click on <br>each item to create or delete<br> a file/directory.')
        },
        {
          target: '.add-btn',
          content: this.$t('Click to add a file or directory<br> on the root directory.')
        },
        {
          target: '.main-content',
          content: this.$t('You can edit, save, rename<br> and delete the selected file <br>in this box.'),
          params: {
            placement: 'left'
          }
        },
        // environment
        {
          target: '.environment-list',
          content: this.$t('Here you can add environment variables that will be passed to the spider program when running a task.')
        },
        // schedules
        {
          target: '.schedule-list',
          content: this.$t('You can add, edit and delete schedules (cron jobs) for the spider.')
        }
      ],
      tourCallbacks: {
        onStop: () => {
          this.$utils.tour.finishTour('spider-detail')
        },
        onPreviousStep: (currentStep) => {
          if (currentStep === 5) {
            this.activeTabName = 'overview'
          } else if (currentStep === 8) {
            this.activeTabName = 'files'
          } else if (currentStep === 9) {
            this.activeTabName = 'environment'
          }
          this.$utils.tour.prevStep('spider-detail', currentStep)
        },
        onNextStep: (currentStep) => {
          if (currentStep === 4) {
            this.activeTabName = 'files'
          } else if (currentStep === 7) {
            this.activeTabName = 'environment'
          } else if (currentStep === 8) {
            this.activeTabName = 'schedules'
          }
          this.$utils.tour.nextStep('spider-detail', currentStep)
        }
      },
      redirectType: ''
    }
  },
  computed: {
    ...mapState('spider', [
      'spiderList',
      'spiderForm',
      'configListTs'
    ]),
    ...mapState('file', [
      'currentPath'
    ]),
    ...mapState('deploy', [
      'deployList'
    ]),
    isCustomized () {
      return this.spiderForm.type === 'customized'
    },
    isConfigurable () {
      return this.spiderForm.type === 'configurable'
    },
    isScrapy () {
      return this.isCustomized && this.spiderForm.is_scrapy
    },
    isGit () {
      return this.spiderForm.is_git
    }
  },
  methods: {
    async onTabClick (tab) {
      if (this.activeTabName === 'analytics') {
        setTimeout(() => {
          this.$refs['spider-stats'].update()
        }, 0)
      } else if (this.activeTabName === 'config') {
        setTimeout(() => {
          this.$refs['config'].update()
        }, 0)

        if (!this.$utils.tour.isFinishedTour('spider-detail-config')) {
          setTimeout(() => {
            this.$utils.tour.startTour(this, 'spider-detail-config')
          }, 100)
        }
      } else if (this.activeTabName === 'scrapy-settings') {
        await this.getScrapyData()
      } else if (this.activeTabName === 'files') {
        await this.$store.dispatch('spider/getFileTree')
        if (this.currentPath) {
          await this.$store.dispatch('file/getFileContent', { path: this.currentPath })
        }
      }
      this.$st.sendEv('爬虫详情', '切换标签', tab.name)
    },
    onSpiderChange (id) {
      this.$router.push(`/spiders/${id}`)
      this.$st.sendEv('爬虫详情', '切换爬虫')
    },
    async getScrapyData () {
      await Promise.all([
        this.$store.dispatch('spider/getSpiderScrapySpiders', this.$route.params.id),
        this.$store.dispatch('spider/getSpiderScrapyItems', this.$route.params.id),
        this.$store.dispatch('spider/getSpiderScrapySettings', this.$route.params.id),
        this.$store.dispatch('spider/getSpiderScrapyPipelines', this.$route.params.id)
      ])
    },
    async onClickScrapySpider (filepath) {
      this.activeTabName = 'files'
      await this.$store.dispatch('spider/getFileTree')
      this.$refs['file-list'].clickSpider(filepath)
    },
    async onClickScrapyPipeline () {
      this.activeTabName = 'files'
      await this.$store.dispatch('spider/getFileTree')
      this.$refs['file-list'].clickPipeline()
    },
    onConvert () {
      this.activeTabName = 'overview'
    }
  },
  async created () {
    // get spider basic info
    await this.$store.dispatch('spider/getSpiderData', this.$route.params.id)

    // get spider file info
    await this.$store.dispatch('file/getFileList', this.spiderForm.src)

    // get spider tasks
    await this.$store.dispatch('spider/getTaskList', this.$route.params.id)

    // get spider list
    await this.$store.dispatch('spider/getSpiderList', { owner_type: 'all' })
  },
  mounted () {
    if (!this.$utils.tour.isFinishedTour('spider-detail')) {
      this.$utils.tour.startTour(this, 'spider-detail')
    }
  }
}
</script>

<style scoped>
  .selector {
    display: flex;
    align-items: center;
    position: absolute;
    right: 48px;
    /*float: right;*/
    z-index: 999;
    margin-top: 5px;
  }

  .selector .el-select {
    height: 30px;
    line-height: 30px;
    padding-left: 10px;
    width: 180px;
    border-radius: 0;
  }

  .selector .el-select >>> .el-input__icon,
  .selector .el-select >>> .el-input__inner {
    border-radius: 0;
    height: 30px;
    line-height: 30px;
  }

  .label {
    text-align: right;
    width: 80px;
    color: #909399;
    font-weight: 100;
  }
</style>
