<template>
  <div class="task-detail-tab-data">
    <ListLayout
        :action-functions="actionFunctions"
        :pagination="tablePagination"
        :table-columns="tableColumns"
        :table-data="tableData"
        :table-total="tableTotal"
        class="result-list"
        no-actions
    >
      <template #extra>
      </template>
    </ListLayout>
  </div>
</template>
<script lang="ts">
import {computed, defineComponent} from 'vue';
import ListLayout from '@/layouts/ListLayout.vue';
import {useStore} from 'vuex';
import useTaskDetail from '@/views/task/detail/taskDetail';

export default defineComponent({
  name: 'TaskDetailTabOverview',
  components: {
    ListLayout,
  },
  setup() {
    // store
    const ns = 'task';
    const store = useStore();
    const {
      task: state,
    } = store.state as RootStoreState;

    // id
    const {
      activeId,
    } = useTaskDetail();

    // data
    const tableData = computed<TableData<Result>>(() => state.resultTableData);

    // total
    const tableTotal = computed<number>(() => state.resultTableTotal);

    // pagination
    const tablePagination = computed<TablePagination>(() => state.resultTablePagination);

    // columns
    const tableColumns = computed<TableColumns<Result>>(() => {
      const fields = store.getters[`${ns}/resultFields`] as ResultField[];
      return fields.map(f => {
        const {key} = f;
        return {
          key,
          label: key,
        };
      }) as TableColumns<Result>;
    });

    // action functions
    const actionFunctions = {
      setPagination: (pagination) => store.commit(`${ns}/setResultTablePagination`, pagination),
      getList: async () => {
        return store.dispatch(`${ns}/getResultData`, activeId.value);
      },
      getAll: async () => {
        console.warn('getAll is not implemented');
      },
      deleteList: (ids: string[]) => {
        console.warn('deleteList is not implemented');
      },
      deleteByIdConfirm: (row: BaseModel) => {
        console.warn('deleteByIdConfirm is not implemented');
      },
    } as ListLayoutActionFunctions;

    return {
      actionFunctions,
      tableData,
      tableTotal,
      tablePagination,
      tableColumns,
    };
  },
});
</script>
<style lang="scss" scoped>
.task-detail-tab-overview {
  margin: 20px;
}
</style>
