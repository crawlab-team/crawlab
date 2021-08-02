import {Store} from 'vuex';
import {plainClone} from '@/utils/object';
import {computed, Ref} from 'vue';

const useFormTable = (ns: ListStoreNamespace, store: Store<RootStoreState>, services: Services<BaseModel>, data: FormComponentData<BaseModel>) => {
  const {
    form,
    formTableFieldRefsMap,
  } = data;

  // state
  const state = store.state[ns];

  // form list
  const formList = computed(() => state.formList);

  const getNewForm = () => {
    return {...form.value};
  };

  const onAdd = (index: number) => {
    formList.value.splice(index, 0, getNewForm());
  };

  const onClone = (index: number) => {
    const form = plainClone(formList.value[index]);
    formList.value.splice(index, 0, form);
  };

  const onDelete = (index: number) => {
    formList.value.splice(index, 1);
    for (const key of formTableFieldRefsMap.value.keys()) {
      const rowIndex = key[0];
      if (rowIndex === index) {
        formTableFieldRefsMap.value.delete(key);
      }
    }
  };

  const onFieldChange = (rowIndex: number, prop: string, value: any) => {
    if (rowIndex !== -1) {
      // one row change
      const item = formList.value[rowIndex] as BaseModel;
      item[prop] = value;
    } else {
      // all rows change
      for (let i = 0; i < formList.value.length; i++) {
        onFieldChange(i, prop, value);
      }
    }
  };

  const onFieldRegister = (rowIndex: number, prop: string, formRef: Ref) => {
    const key = [rowIndex, prop] as FormTableFieldRefsMapKey;
    formTableFieldRefsMap.value.set(key, formRef);
  };

  return {
    onAdd,
    onClone,
    onDelete,
    onFieldChange,
    onFieldRegister,
  };
};

export default useFormTable;
