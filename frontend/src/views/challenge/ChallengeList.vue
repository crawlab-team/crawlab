<template>
  <div class="app-container">
    <ul class="challenge-list">
      <li
        v-for="(c, $index) in challenges"
        :key="$index"
        class="challenge-item"
      >
        <el-card>
          <div class="title" :title="lang === 'zh' ? c.title_cn : c.title_en">
            {{lang === 'zh' ? c.title_cn : c.title_en}}
          </div>
          <div class="rating block">
            <span class="label">{{$t('Difficulty')}}: </span>
            <el-rate
              v-model="c.difficulty"
              disabled
            >
            </el-rate>
          </div>
          <div class="achieved block">
            <span class="label">{{$t('Status')}}: </span>
            <div class="content">
              <div v-if="c.achieved" class="status is-achieved">
                <i class="fa fa-check-square-o"></i>
                <span>{{$t('Achieved')}}</span>
              </div>
              <div v-else class="status is-not-achieved">
                <i class="fa fa-square-o"></i>
                <span>{{$t('Not Achieved')}}</span>
              </div>
            </div>
          </div>
          <div class="description">
            {{lang === 'zh' ? c.description_cn : c.description_en}}
          </div>
          <div class="actions">
            <el-button
              v-if="c.achieved"
              size="mini"
              type="success"
              icon="el-icon-check"
              disabled
            >
              {{$t('Achieved')}}
            </el-button>
            <el-button
              v-else
              size="mini"
              type="primary"
              icon="el-icon-s-flag"
            >
              {{$t('Start Challenge')}}
            </el-button>
          </div>
        </el-card>
      </li>
    </ul>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
export default {
  name: 'ChallengeList',
  data () {
    return {
      challenges: []
    }
  },
  computed: {
    ...mapState('lang', [
      'lang'
    ])
  },
  methods: {
    async getData () {
      await this.$request.post('/challenges-check')
      const res = await this.$request.get('/challenges')
      this.challenges = res.data.data || []
    }
  },
  async created () {
    await this.getData()
  }
}
</script>

<style scoped>
  .challenge-list {
    list-style: none;
    display: flex;
    flex-wrap: wrap;
  }

  .challenge-list .challenge-item {
    flex-basis: 240px;
    width: 240px;
    margin: 10px;
  }

  .challenge-list .challenge-item .title {
    padding-bottom: 10px;
    border-bottom: 1px solid #e9e9eb;
    height: 30px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .challenge-list .challenge-item .el-card {
    height: 275px;
  }

  .challenge-list .challenge-item .block {
    margin-top: 10px;
    margin-bottom: 10px;
  }

  .challenge-list .challenge-item .rating {
  }

  .challenge-list .challenge-item .rating .el-rate {
    display: inline-block;
  }

  .challenge-list .challenge-item .label {
    display: inline-flex;
    align-items: center;
    font-size: 12px;
    line-height: 21px;
    height: 21px;
    margin-right: 5px;
    text-align: right;
  }

  .challenge-list .challenge-item .content {
    display: inline-flex;
    align-items: center;
    font-size: 12px;
    line-height: 21px;
    height: 21px;
    font-weight: bolder;
  }

  .challenge-list .challenge-item .block.achieved {
    display: flex;
    align-items: center;
  }

  .challenge-list .challenge-item .achieved .content .status {
    margin-top: 0;
    display: flex;
    align-items: center;
  }

  .challenge-list .challenge-item .achieved .content .status.is-achieved {
    color: #67c23a;
  }

  .challenge-list .challenge-item .achieved .content .status.is-not-achieved {
    color: #E6A23C;
  }

  .challenge-list .challenge-item .achieved .content .status i {
    margin: 0 3px;
    font-size: 18px;
  }

  .challenge-list .challenge-item .description {
    box-sizing: border-box;
    font-size: 12px;
    padding-top: 10px;
    padding-bottom: 10px;
    line-height: 20px;
    height: 100px;
    border-top: 1px solid #e9e9eb;
    border-bottom: 1px solid #e9e9eb;
    overflow: auto;
  }

  .challenge-list .challenge-item .actions {
    text-align: right;
    padding-top: 10px;
  }

</style>
