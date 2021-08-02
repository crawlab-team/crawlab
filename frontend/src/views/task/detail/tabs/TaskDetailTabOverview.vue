<template>
  <div class="task-detail-tab-overview">
    <TaskForm readonly/>
  </div>
</template>
<script lang="ts">
import {defineComponent} from 'vue';
import TaskForm from '@/components/task/TaskForm.vue';

export default defineComponent({
  name: 'TaskDetailTabOverview',
  components: {
    TaskForm,
  },
  setup() {
    return {};
  },
});
</script>
<style lang="scss" scoped>
.task-detail-tab-overview {
  margin: 20px;
}
</style>
