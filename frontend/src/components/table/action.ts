import {inject, Ref, ref, SetupContext} from 'vue';
import {Table} from 'element-plus/lib/el-table/src/table.type';
import {ElMessageBox} from 'element-plus';

const useAction = (props: TableProps, ctx: SetupContext, table: Ref<Table | undefined>, actionFunctions: ListLayoutActionFunctions) => {
  const {emit} = ctx;

  // store context
  const storeContext = inject<ListStoreContext<BaseModel>>('store-context');
  const ns = storeContext?.namespace;
  const store = storeContext?.store;

  // table selection
  const selection = ref<TableData>([]);
  const onSelectionChange = (value: TableData) => {
    selection.value = value;
    emit('selection-change', value);
  };

  // action functions
  const {
    getList,
    deleteList,
  } = actionFunctions;

  const onAdd = () => {
    emit('add');
  };

  const onEdit = async () => {
    emit('edit', selection.value);
    if (storeContext) {
      store?.commit(`${ns}/showDialog`, 'edit');
      store?.commit(`${ns}/setIsSelectiveForm`, true);
      store?.commit(`${ns}/setFormList`, selection.value);
    }
  };

  const onDelete = async () => {
    const res = await ElMessageBox.confirm('Are you sure to delete selected items?', 'Batch Delete', {
      type: 'warning',
      confirmButtonText: 'Delete',
      confirmButtonClass: 'el-button--danger',
    });
    if (!res) return;
    const ids = selection.value.map(d => d._id as string);
    await deleteList(ids);
    table.value?.store?.clearSelection();
    await getList();
    emit('delete', selection.value);
  };

  const onExport = () => {
    emit('export');
  };

  return {
    // public variables and methods
    selection,
    onSelectionChange,
    onAdd,
    onEdit,
    onDelete,
    onExport,
  };
};

export default useAction;
