import {ref} from 'vue';

export const getDefaultFormComponentData = <T>(defaultFn: DefaultFormFunc<T>) => {
  return {
    form: ref<T>(defaultFn()),
    formRef: ref(),
    formList: ref<T[]>([]),
    formTableFieldRefsMap: ref<FormTableFieldRefsMap>(new Map()),
  } as FormComponentData<T>;
};
