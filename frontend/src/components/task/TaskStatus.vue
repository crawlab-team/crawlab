<template>
  <Tag
      class="task-status"
      :key="data"
      :icon="data.icon"
      :label="data.label"
      :spinning="data.spinning"
      :type="data.type"
      :size="size"
      @click="$emit('click')"
  >
    <template #tooltip>
      <div v-html="data.tooltip"/>
    </template>
  </Tag>
</template>

<script lang="ts">
import {computed, defineComponent, PropType} from 'vue';
import {
  TASK_STATUS_ABNORMAL,
  TASK_STATUS_CANCELLED,
  TASK_STATUS_ERROR,
  TASK_STATUS_FINISHED,
  TASK_STATUS_PENDING,
  TASK_STATUS_RUNNING
} from '@/constants/task';
import Tag from '@/components/tag/Tag.vue';
import colors from '@/styles/color.scss';

export default defineComponent({
  name: 'TaskStatus',
  components: {
    Tag,
  },
  props: {
    status: {
      type: String as PropType<TaskStatus>,
      required: false,
    },
    size: {
      type: String as PropType<BasicSize>,
      required: false,
      default: 'mini',
    },
    error: {
      type: String,
      required: false,
    }
  },
  emits: ['click'],
  setup(props: TaskStatusProps, {emit}) {
    const data = computed<TagData>(() => {
      const {status, error} = props;
      switch (status) {
        case TASK_STATUS_PENDING:
          return {
            label: 'Pending',
            tooltip: 'Task is pending in the queue',
            type: 'primary',
            icon: ['fa', 'hourglass-start'],
            spinning: true,
          };
        case TASK_STATUS_RUNNING:
          return {
            label: 'Running',
            tooltip: 'Task is currently running',
            type: 'warning',
            icon: ['fa', 'spinner'],
            spinning: true,
          };
        case TASK_STATUS_FINISHED:
          return {
            label: 'Finished',
            tooltip: 'Task finished successfully',
            type: 'success',
            icon: ['fa', 'check'],
          };
        case TASK_STATUS_ERROR:
          return {
            label: 'Error',
            tooltip: `Task ended with an error:<br><span style="color: ${colors.red}">${error}</span>`,
            type: 'danger',
            icon: ['fa', 'times'],
          };
        case TASK_STATUS_CANCELLED:
          return {
            label: 'Cancelled',
            tooltip: 'Task has been cancelled',
            type: 'info',
            icon: ['fa', 'ban'],
          };
        case TASK_STATUS_ABNORMAL:
          return {
            label: 'Cancelled',
            tooltip: 'Task ended abnormally',
            type: 'info',
            icon: ['fa', 'exclamation'],
          };
        default:
          return {
            label: 'Unknown',
            tooltip: 'Unknown task status',
            type: 'info',
            icon: ['fa', 'question'],
          };
      }
    });

    return {
      data,
    };
  },
});
</script>

<style lang="scss" scoped>
.task-status {
}
</style>
