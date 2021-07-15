<template>
  <el-tooltip :content="tooltip">
    <el-tag :type="type" class="task-mode" size="mini">
      <font-awesome-icon :icon="icon" class="icon"/>
      <span>{{ label }}</span>
    </el-tag>
  </el-tooltip>
</template>

<script lang="ts">
import {computed, defineComponent} from 'vue';
import {
  TASK_MODE_ALL,
  TASK_MODE_RANDOM,
  TASK_MODE_SELECTED_NODE_TAGS,
  TASK_MODE_SELECTED_NODES
} from '@/constants/task';

export default defineComponent({
  name: 'TaskMode',
  props: {
    mode: {
      type: String,
      required: false,
    }
  },
  setup(props: TaskModeProps, {emit}) {
    const type = computed<string>(() => {
      const {mode} = props;
      switch (mode) {
        case TASK_MODE_RANDOM:
          return 'warning';
        case TASK_MODE_ALL:
          return 'success';
        case TASK_MODE_SELECTED_NODES:
          return 'primary';
        case TASK_MODE_SELECTED_NODE_TAGS:
          return 'primary';
        default:
          return 'info';
      }
    });

    const label = computed<string>(() => {
      const {mode} = props;
      switch (mode) {
        case TASK_MODE_RANDOM:
          return 'Random';
        case TASK_MODE_ALL:
          return 'All Nodes';
        case TASK_MODE_SELECTED_NODES:
          return 'Nodes';
        case TASK_MODE_SELECTED_NODE_TAGS:
          return 'Tags';
        default:
          return 'Unknown';
      }
    });

    const icon = computed<string[]>(() => {
      const {mode} = props;
      switch (mode) {
        case TASK_MODE_RANDOM:
          return ['fa', 'random'];
        case TASK_MODE_ALL:
          return ['fa', 'sitemap'];
        case TASK_MODE_SELECTED_NODES:
          return ['fa', 'network-wired'];
        case TASK_MODE_SELECTED_NODE_TAGS:
          return ['fa', 'tags'];
        default:
          return ['fa', 'question'];
      }
    });

    const tooltip = computed<string>(() => {
      const {mode} = props;
      switch (mode) {
        case TASK_MODE_RANDOM:
          return 'Run on a random node';
        case TASK_MODE_ALL:
          return 'Run on all nodes';
        case TASK_MODE_SELECTED_NODES:
          return 'Run on selected nodes';
        case TASK_MODE_SELECTED_NODE_TAGS:
          return 'Run on nodes with selected tags';
        default:
          return 'Unknown task mode';
      }
    });

    return {
      type,
      label,
      icon,
      tooltip,
    };
  },
});
</script>

<style lang="scss" scoped>
.task-mode {
  width: 80px;
  cursor: default;

  .icon {
    width: 10px;
    margin-right: 5px;
  }
}
</style>
