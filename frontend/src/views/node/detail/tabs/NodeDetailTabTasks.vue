<template>
  <div class="node-detail-tab-tasks">
    <TaskList no-actions/>
  </div>
</template>
<script lang="ts">
import {computed, defineComponent, onBeforeMount, onBeforeUnmount} from 'vue';
import TaskList from '@/views/task/list/TaskList.vue';
import {useStore} from 'vuex';
import {useRoute} from 'vue-router';
import {FILTER_OP_EQUAL} from '@/constants/filter';

export default defineComponent({
  name: 'NodeDetailTabTasks',
  components: {
    TaskList
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
        key: 'node_id',
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
<style lang="scss" scoped>
.node-detail-tab-overview {
  margin: 20px;
}
</style>

<style scoped>
.node-detail-tab-tasks >>> .el-table {
  border: none;
}
</style>
