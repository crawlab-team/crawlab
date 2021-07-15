import {
  getDefaultStoreActions,
  getDefaultStoreGetters,
  getDefaultStoreMutations,
  getDefaultStoreState
} from '@/utils/store';
import useRequest from '@/services/request';
import {TAB_NAME_DATA, TAB_NAME_LOGS, TAB_NAME_OVERVIEW} from '@/constants/tab';
import {Editor} from 'codemirror';
import {getFieldsFromData} from '@/utils/result';
import {getDefaultPagination} from '@/utils/pagination';

const {
  put,
  getList,
} = useRequest();

const state = {
  ...getDefaultStoreState<Task>('task'),
  tabs: [
    {id: TAB_NAME_OVERVIEW, title: 'Overview'},
    {id: TAB_NAME_LOGS, title: 'Logs'},
    {id: TAB_NAME_DATA, title: 'Data'},
  ],
  logContent: '',
  logPagination: {
    page: 1,
    size: 1000,
  },
  logTotal: 0,
  logAutoUpdate: false,
  logCodeMirrorEditor: undefined,
  resultTableData: [],
  resultTablePagination: getDefaultPagination(),
  resultTableTotal: 0,
} as TaskStoreState;

const getters = {
  ...getDefaultStoreGetters<Task>(),
  resultFields: (state: TaskStoreState) => {
    return getFieldsFromData(state.resultTableData);
  },
} as TaskStoreGetters;

const mutations = {
  ...getDefaultStoreMutations<Task>(),
  setLogContent: (state: TaskStoreState, content: string) => {
    state.logContent = content;
  },
  resetLogContent: (state: TaskStoreState) => {
    state.logContent = '';
  },
  setLogPagination: (state: TaskStoreState, pagination: TablePagination) => {
    state.logPagination = pagination;
  },
  resetLogPagination: (state: TaskStoreState) => {
    state.logPagination = {page: 1, size: 1000};
  },
  setLogTotal: (state: TaskStoreState, total: number) => {
    state.logTotal = total;
  },
  resetLogTotal: (state: TaskStoreState) => {
    state.logTotal = 0;
  },
  enableLogAutoUpdate: (state: TaskStoreState) => {
    state.logAutoUpdate = true;
  },
  disableLogAutoUpdate: (state: TaskStoreState) => {
    state.logAutoUpdate = false;
  },
  setLogCodeMirrorEditor: (state: TaskStoreState, cm: Editor) => {
    state.logCodeMirrorEditor = cm;
  },
  setResultTableData: (state: TaskStoreState, data: Result[]) => {
    state.resultTableData = data;
  },
  resetResultTableData: (state: TaskStoreState) => {
    state.resultTableData = [];
  },
  setResultTablePagination: (state: TaskStoreState, pagination: TablePagination) => {
    state.resultTablePagination = pagination;
  },
  resetResultTablePagination: (state: TaskStoreState) => {
    state.resultTablePagination = getDefaultPagination();
  },
  setResultTableTotal: (state: TaskStoreState, total: number) => {
    state.resultTableTotal = total;
  },
  resetResultTableTotal: (state: TaskStoreState) => {
    state.resultTableTotal = 0;
  },
} as TaskStoreMutations;

const actions = {
  ...getDefaultStoreActions<Task>('/tasks'),
  getList: async ({state, commit}: StoreActionContext<TaskStoreState>) => {
    const payload = {
      ...state.tablePagination,
      conditions: JSON.stringify(state.tableListFilter),
      sort: JSON.stringify(state.tableListSort),
      stats: true,
    };
    const res = await getList(`/tasks`, payload);
    commit('setTableData', {data: res.data || [], total: res.total});
    return res;
  },
  create: async ({state, commit}: StoreActionContext<TaskStoreState>, form: Task) => {
    return await put(`/tasks/run`, form);
  },
  getLogs: async ({state, commit}: StoreActionContext<TaskStoreState>, id: string) => {
    const {page, size} = state.logPagination;
    const res = await getList(`/tasks/${id}/logs`, {page, size});
    commit('setLogContent', res.data?.join('\n'));
    commit('setLogTotal', res.total);
    return res;
  },
  getResultData: async ({state, commit}: StoreActionContext<TaskStoreState>, id: string) => {
    const {page, size} = state.resultTablePagination;
    const res = await getList(`/tasks/${id}/data`, {page, size});
    commit('setResultTableData', res.data || []);
    commit('setResultTableTotal', res.total);
    return res;
  },
} as TaskStoreActions;

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions,
} as TaskStoreModule;
