<template>
  <el-tree
    :data="docData"
    ref="documentation-tree"
    node-key="path"
  >
    <span class="custom-tree-node" :class="data.active ? 'active' : ''" slot-scope="{ node, data }">
      <template v-if="data.level === 1 && data.children && data.children.length">
        <span>{{node.label}}</span>
      </template>
      <template v-else>
        <span>
          <a :href="data.url" target="_blank" style="display: block">
            {{node.label}}
          </a>
        </span>
      </template>
    </span>
  </el-tree>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'Documentation',
  data () {
    return {
      data: []
    }
  },
  computed: {
    ...mapState('doc', [
      'docData'
    ]),
    pathLv1 () {
      if (this.$route.path === '/') return '/'
      const m = this.$route.path.match(/(^\/\w+)/)
      return m[1]
    }
  },
  watch: {
    pathLv1 () {
      this.update()
    }
  },
  methods: {
    update () {
      // expand related documentation list
      setTimeout(() => {
        this.docData.forEach(d => {
          const node = this.$refs['documentation-tree'].getNode(d)
          let isActive = false
          for (let i = 0; i < this.$utils.doc.docs.length; i++) {
            const pattern = this.$utils.doc.docs[i]
            if (d.url.match(pattern)) {
              isActive = true
              break
            }
          }
          node.expanded = isActive
          this.$set(d, 'active', isActive)
        })
      }, 100)
    },
    async getDocumentationData () {
      // fetch api data
      await this.$store.dispatch('doc/getDocData')
    }
  },
  async created () {
  },
  mounted () {
    this.update()
  }
}
</script>
<style scoped>
  .el-tree >>> .custom-tree-node.active {
    color: #409eff;
  }
</style>
