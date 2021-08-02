import {
  getDefaultStoreActions,
  getDefaultStoreGetters,
  getDefaultStoreMutations,
  getDefaultStoreState
} from '@/utils/store';
import {TAB_NAME_OVERVIEW, TAB_NAME_TASKS} from '@/constants/tab';

const state = {
  ...getDefaultStoreState<CNode>('node'),
  tabs: [
    {id: TAB_NAME_OVERVIEW, title: 'Overview'},
    {id: TAB_NAME_TASKS, title: 'Tasks'},
  ],
  // TODO: dummy data
  allNodeSelectOptions: [
    {label: 'Master', value: 'master'},
    {label: 'Worker 1', value: 'worker-1'},
    {label: 'Worker 2', value: 'worker-2'},
    {label: 'Worker 3', value: 'worker-3'},
    {label: 'Worker 4', value: 'worker-4'},
    {label: 'Worker 5', value: 'worker-5'},
  ],
  allNodeTags: [
    '1c2g',
    '2c4g',
    '2c8g',
    '4c16g',
  ],
} as NodeStoreState;

const getters = {
  ...getDefaultStoreGetters<CNode>(),
} as NodeStoreGetters;

const mutations = {
  ...getDefaultStoreMutations<CNode>(),
} as NodeStoreMutations;

const actions = {
  ...getDefaultStoreActions<CNode>('/nodes'),
} as NodeStoreActions;

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions,
} as NodeStoreModule;
