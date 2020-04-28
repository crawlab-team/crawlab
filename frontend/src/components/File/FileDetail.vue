<template>
  <div class="file-detail">
    <codemirror
      class="file-content"
      :options="options"
      v-model="fileContent"
    />
  </div>
</template>

<script>
import {
  mapState,
  mapGetters
} from 'vuex'
import { codemirror } from 'vue-codemirror-lite'

import 'codemirror/lib/codemirror.js'

// language
import 'codemirror/mode/python/python.js'
import 'codemirror/mode/javascript/javascript.js'
import 'codemirror/mode/go/go.js'
import 'codemirror/mode/shell/shell.js'
import 'codemirror/mode/markdown/markdown.js'
import 'codemirror/mode/php/php.js'
import 'codemirror/mode/yaml/yaml.js'

export default {
  name: 'FileDetail',
  components: { codemirror },
  data () {
    return {
      internalFileContent: ''
    }
  },
  computed: {
    ...mapState('spider', [
      'spiderForm'
    ]),
    ...mapGetters('user', [
      'userInfo'
    ]),
    fileContent: {
      get () {
        return this.$store.state.file.fileContent
      },
      set (value) {
        return this.$store.commit('file/SET_FILE_CONTENT', value)
      }
    },
    options () {
      return {
        mode: this.language,
        theme: 'darcula',
        styleActiveLine: true,
        smartIndent: true,
        indentUnit: 4,
        lineNumbers: true,
        line: true,
        matchBrackets: true,
        readOnly: this.isDisabled ? 'nocursor' : false
      }
    },
    language () {
      const fileName = this.$store.state.file.currentPath
      if (!fileName) return ''
      if (fileName.match(/\.js$/)) {
        return 'text/javascript'
      } else if (fileName.match(/\.py$/)) {
        return 'text/x-python'
      } else if (fileName.match(/\.go$/)) {
        return 'text/x-go'
      } else if (fileName.match(/\.sh$/)) {
        return 'text/x-shell'
      } else if (fileName.match(/\.php$/)) {
        return 'text/x-php'
      } else if (fileName.match(/\.md$/)) {
        return 'text/x-markdown'
      } else if (fileName.match('Spiderfile')) {
        return 'text/x-yaml'
      } else {
        return 'text'
      }
    },
    isDisabled () {
      return this.spiderForm.is_public && this.spiderForm.username !== this.userInfo.username && this.userInfo.role !== 'admin'
    }
  },
  created () {
    this.internalFileContent = this.fileContent
  }
}
</script>

<style scoped>
  .file-content {
    border: 1px solid #eaecef;
    height: calc(100vh - 256px);
  }

  .file-content >>> .CodeMirror {
    min-height: 100%;
  }
</style>
