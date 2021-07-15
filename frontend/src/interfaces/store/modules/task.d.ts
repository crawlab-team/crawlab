import {Editor} from 'codemirror';

declare global {
  type TaskStoreModule = BaseModule<TaskStoreState, TaskStoreGetters, TaskStoreMutations, TaskStoreActions>;

  interface TaskStoreState extends BaseStoreState<Task> {
    logContent: string;
    logPagination: TablePagination;
    logTotal: number;
    logAutoUpdate: boolean;
    logCodeMirrorEditor?: Editor;
    resultTableData: TableData<Result>;
    resultTablePagination: TablePagination;
    resultTableTotal: number;
  }

  interface TaskStoreGetters extends BaseStoreGetters<TaskStoreState> {
    resultFields: StoreGetter<TaskStoreState, ResultField[]>;
  }

  interface TaskStoreMutations extends BaseStoreMutations<Task> {
    setLogContent: StoreMutation<TaskStoreState, string>;
    resetLogContent: StoreMutation<TaskStoreState>;
    setLogPagination: StoreMutation<TaskStoreState, TablePagination>;
    resetLogPagination: StoreMutation<TaskStoreState>;
    setLogTotal: StoreMutation<TaskStoreState, number>;
    resetLogTotal: StoreMutation<TaskStoreState>;
    enableLogAutoUpdate: StoreMutation<TaskStoreState>;
    disableLogAutoUpdate: StoreMutation<TaskStoreState>;
    setLogCodeMirrorEditor: StoreMutation<TaskStoreState, Editor>;
    setResultTableData: StoreMutation<TaskStoreState, TableData<Result>>;
    resetResultTableData: StoreMutation<TaskStoreState>;
    setResultTablePagination: StoreMutation<TaskStoreState, TablePagination>;
    resetResultTablePagination: StoreMutation<TaskStoreState>;
    setResultTableTotal: StoreMutation<TaskStoreState, number>;
    resetResultTableTotal: StoreMutation<TaskStoreState>;
  }

  interface TaskStoreActions extends BaseStoreActions<Task> {
    getLogs: StoreAction<TaskStoreState, string>;
    getResultData: StoreAction<TaskStoreState, string>;
  }
}
