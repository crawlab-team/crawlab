<template>
  <ListLayout
      :action-functions="actionFunctions"
      :no-actions="noActions"
      :pagination="tablePagination"
      :table-columns="tableColumns"
      :table-data="tableData"
      :table-total="tableTotal"
      class="result-list"
  >
    <template #extra>
    </template>
  </ListLayout>
</template>

<script lang="ts">
import {computed, defineComponent, watch} from 'vue';
import ListLayout from '@/layouts/ListLayout.vue';
import {useStore} from 'vuex';

export default defineComponent({
  name: 'ResultList',
  components: {
    ListLayout,
  },
  props: {
    id: {
      type: String,
      required: true,
    },
    noActions: {
      type: Boolean,
      default: false,
    },
  },
  setup(props: ResultListProps) {
    // store
    const ns = 'dataCollection';
    const store = useStore();
    const {
      dataCollection: state,
    } = store.state as RootStoreState;

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
        const {id} = props;
        if (!id) return;
        return store.dispatch(`${ns}/getResultData`, {
          id,
          params: tablePagination.value,
        });
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

    const {
      getList,
    } = actionFunctions;

    watch(() => props.id, getList);

    watch(() => tablePagination.value, getList);

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

</style>
