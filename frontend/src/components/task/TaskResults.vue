<template>
  <Tag
      :key="data"
      :icon="data.icon"
      :label="data.label"
      :size="size"
      :spinning="data.spinning"
      :type="data.type"
      class="task-status"
      @click="$emit('click')"
  >
    <template #tooltip>
      <div v-html="data.tooltip"/>
    </template>
  </Tag>
</template>

<script lang="ts">
import {computed, defineComponent, PropType} from 'vue';
import Tag from '@/components/tag/Tag.vue';
import {isCancellable} from '@/utils/task';
import {TASK_STATUS_PENDING} from '@/constants/task';

export default defineComponent({
  name: 'TaskResults',
  components: {
    Tag,
  },
  props: {
    results: {
      type: Number,
      required: false,
    },
    status: {
      type: String as PropType<TaskStatus>,
      required: false,
    },
    size: {
      type: String as PropType<BasicSize>,
      required: false,
      default: 'mini',
    },
  },
  setup(props: TaskResultsProps, {emit}) {
    const data = computed<TagData>(() => {
      const {results, status} = props;
      if (isCancellable(status)) {
        if (status === TASK_STATUS_PENDING) {
          return {
            label: results?.toFixed(0),
            tooltip: `Results: ${results}`,
            type: 'primary',
            icon: ['fa', 'hourglass-start'],
            spinning: true,
          };
        } else {
          return {
            label: results?.toFixed(0),
            tooltip: `Results: ${results}`,
            type: 'warning',
            icon: ['fa', 'spinner'],
            spinning: true,
          };
        }
      } else {
        if (results === 0) {
          return {
            label: results?.toFixed(0),
            tooltip: `No results`,
            type: 'danger',
            icon: ['fa', 'exclamation'],
          };
        } else {
          return {
            label: results?.toFixed(0),
            tooltip: `Results: ${results}`,
            type: 'success',
            icon: ['fa', 'table'],
          };
        }
      }
    });

    return {
      data,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
