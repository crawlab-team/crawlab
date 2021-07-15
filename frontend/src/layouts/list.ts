import {computed, onBeforeUnmount, provide, readonly, watch} from 'vue';
import {Store} from 'vuex';
import {ElMessage, ElMessageBox} from 'element-plus';
import {FILTER_OP_CONTAINS, FILTER_OP_IN, FILTER_OP_NOT_SET} from '@/constants/filter';

const getFilterConditions = (column: TableColumn, filter: TableHeaderDialogFilterData) => {
  // allow filter search/items
  const {allowFilterSearch, allowFilterItems} = column;

  // conditions
  const conditions = [] as FilterConditionData[];

  // filter conditions
  if (filter.conditions) {
    filter.conditions
      .filter(d => d.op !== FILTER_OP_NOT_SET)
      .forEach(d => {
        conditions.push(d);
      });
  }

  if (allowFilterItems) {
    // allow filter items (only relevant to items)
    if (filter.items && filter.items.length > 0) {
      conditions.push({
        op: FILTER_OP_IN,
        value: filter.items,
      });
    }
  } else if (allowFilterSearch) {
    // not allow filter items and allow filter search (only relevant to search string)
    if (filter.searchString) {
      conditions.push({
        op: FILTER_OP_CONTAINS,
        value: filter.searchString,
      });
    }
  }

  return conditions;
};

const useList = <T = any>(ns: ListStoreNamespace, store: Store<RootStoreState>, opts?: UseListOptions<T>): ListLayoutComponentData => {
  // store state
  const state = store.state[ns] as BaseStoreState;

  // table
  const tableData = computed<TableData<T>>(() => state.tableData as TableData<T>);
  const tableTotal = computed<number>(() => state.tableTotal);
  const tablePagination = computed<TablePagination>(() => state.tablePagination);

  // action functions
  const actionFunctions = readonly<ListLayoutActionFunctions>({
    setPagination: (pagination: TablePagination) => store.commit(`${ns}/setTablePagination`, pagination),
    getList: () => store.dispatch(`${ns}/getList`),
    getAll: () => store.dispatch(`${ns}/getAllList`),
    deleteList: (ids: string[]) => store.dispatch(`${ns}/deleteList`, ids),
    deleteByIdConfirm: async (row: BaseModel) => {
      await ElMessageBox.confirm('Are you sure to delete?', 'Delete', {
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      });
      await store.dispatch(`${ns}/deleteById`, row._id);
      await ElMessage.success('Deleted successfully');
      await store.dispatch(`${ns}/getList`);
    },
    onHeaderChange: async (column, sort, filter) => {
      const {key} = column;

      // filter
      if (!filter) {
        // no filter
        store.commit(`${ns}/resetTableListFilterByKey`, key);
      } else {
        // has filter
        const conditions = getFilterConditions(column, filter);
        store.commit(`${ns}/setTableListFilterByKey`, {key, conditions});
      }

      // sort
      if (!sort) {
        // no sort
        store.commit(`${ns}/resetTableListSortByKey`, key);
      } else {
        // has sort
        store.commit(`${ns}/setTableListSortByKey`, {key, sort});
      }

      // get list
      await store.dispatch(`${ns}/getList`);
    },
  });

  // active dialog key
  const activeDialogKey = computed<DialogKey | undefined>(() => state.activeDialogKey);

  // get list when pagination changes
  watch(() => tablePagination.value, actionFunctions.getList);

  // reset form when active dialog key is changed
  watch(() => state.activeDialogKey, () => {
    if (!state.activeDialogKey) {
      store.commit(`${ns}/resetForm`);
      store.commit(`${ns}/resetFormList`);
    }
  });

  // store context
  provide<ListStoreContext<T>>('store-context', {
    namespace: ns,
    store,
    state,
  });

  onBeforeUnmount(() => {
    store.commit(`${ns}/resetTableData`);
    store.commit(`${ns}/resetTablePagination`);
    store.commit(`${ns}/resetTableListFilter`);
  });

  return {
    ...opts,
    tableData,
    tableTotal,
    tablePagination,
    actionFunctions,
    activeDialogKey,
  };
};

export default useList;
