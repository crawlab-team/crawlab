<template>
  <codemirror
    class="file-content"
    :options="options"
    v-model="fileContent"
    @change="onChange"
  />
</template>

<script>
import { codemirror } from 'vue-codemirror'

require('codemirror/mode/python/python.js')
require('codemirror/mode/javascript/javascript.js')
require('codemirror/mode/go/go.js')
require('codemirror/mode/shell/shell.js')
require('codemirror/addon/fold/foldcode.js')
require('codemirror/addon/fold/foldgutter.js')
require('codemirror/addon/fold/brace-fold.js')
require('codemirror/addon/fold/xml-fold.js')
require('codemirror/addon/fold/indent-fold.js')
require('codemirror/addon/fold/markdown-fold.js')
require('codemirror/addon/fold/comment-fold.js')

export default {
  name: 'FileDetail',
  components: { codemirror },
  data () {
    return {
      internalFileContent: '',
      options: {
        theme: 'darcula',
        // mode: 'javascript',
        mode: { name: 'javascript', json: true },
        lineNumbers: true,
        line: true,
        matchBrackets: true
      }
    }
  },
  computed: {
    fileContent: {
      get () {
        return this.$store.state.file.fileContent
      },
      set (value) {
        return this.$store.commit('file/SET_FILE_CONTENT', value)
      }
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
    height: 480px;
  }

  .file-content >>> .CodeMirror {
    min-height: 100%;
  }
</style>
