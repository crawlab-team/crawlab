<template>
  <div class="spider-detail-tab-tasks">
    <TaskList no-actions/>
  </div>
</template>
<script lang="ts">
import {computed, defineComponent, onBeforeMount, onBeforeUnmount} from 'vue';
import {useRoute} from 'vue-router';
import {useStore} from 'vuex';
import {FILTER_OP_EQUAL} from '@/constants/filter';
import TaskList from '@/views/task/list/TaskList.vue';

export default defineComponent({
  name: 'SpiderDetailTabTasks',
  components: {
    TaskList,
  },
  setup() {
    // route
    const route = useRoute();

    // store
    const store = useStore();

    // id
    const id = computed<string>(() => route.params.id as string);

    onBeforeMount(() => {
      // set filter
      store.commit(`task/setTableListFilter`, [{
        key: 'spider_id',
        op: FILTER_OP_EQUAL,
        value: id.value,
      }]);
    });

    onBeforeUnmount(() => {
      store.commit(`task/resetTableListFilter`);
      store.commit(`task/resetTableData`);
    });

    return {};
  },
});
</script>

<style scoped>
.spider-detail-tab-tasks >>> .el-table {
  border: none;
}
</style>
