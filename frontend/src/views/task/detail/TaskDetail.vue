<template>
  <DetailLayout no-sidebar store-namespace="task">
    <template #actions>
      <TaskDetailActionsCommon/>
      <TaskDetailActionsLogs v-if="activeTabName === 'logs'"/>
    </template>
  </DetailLayout>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import useTaskDetail from '@/views/task/detail/taskDetail';
import DetailLayout from '@/layouts/DetailLayout.vue';
import TaskDetailActionsCommon from '@/views/task/detail/actions/TaskDetailActionsCommon.vue';
import TaskDetailActionsLogs from '@/views/task/detail/actions/TaskDetailActionsLogs.vue';

export default defineComponent({
  name: 'TaskDetail',
  components: {
    TaskDetailActionsLogs,
    DetailLayout,
    TaskDetailActionsCommon,
  },
  setup() {
    return {
      ...useTaskDetail(),
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
