<template>
  <span class="atom-material-icon" v-html="html"/>
</template>
<script lang="ts">
import {computed, defineComponent} from 'vue';
import {getFileIconFromName, getFolderIconFromName} from 'atom-material-icons';

export default defineComponent({
  name: 'AtomMaterialIcon',
  props: {
    name: {
      type: String,
      required: true,
    },
    isDir: {
      type: Boolean,
      required: false,
    }
  },
  setup(props: AtomMaterialIconProps) {
    const html = computed<string>(() => {
      const {name, isDir} = props;
      const icon = isDir ? getFolderIconFromName(name) : getFileIconFromName(name);
      return icon.default;
    });

    return {
      html,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
