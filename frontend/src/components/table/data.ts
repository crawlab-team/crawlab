import {computed, SetupContext} from 'vue';

const useData = (props: TableProps, ctx: SetupContext) => {
  const tableData = computed(() => {
    const {data} = props;
    return data;
  });

  return {
    tableData,
  };
};

export default useData;
