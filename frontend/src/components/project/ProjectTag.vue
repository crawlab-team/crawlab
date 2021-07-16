<template>
  <el-tag :color="color" :type="type" class="project-tag" size="mini">
    <span>{{ label }}</span>
  </el-tag>
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue';

export default defineComponent({
  name: 'ProjectTag',
  props: {
    label: {
      type: String,
    },
    type: {
      type: String as PropType<BasicType>,
    },
    color: {
      type: String,
    }
  },
  setup(props: ProjectTagProps, {emit}) {
    return {};
  },
});
</script>

<style lang="scss" scoped>
.project-tag {
  margin-right: 10px;
}
</style>
