import {
  getDefaultStoreActions,
  getDefaultStoreGetters,
  getDefaultStoreMutations,
  getDefaultStoreState
} from '@/utils/store';
import {TAB_NAME_OVERVIEW, TAB_NAME_SPIDERS} from '@/constants/tab';

const state = {
  ...getDefaultStoreState<Project>('project'),
  tabs: [
    {id: TAB_NAME_OVERVIEW, title: 'Overview'},
    {id: TAB_NAME_SPIDERS, title: 'Spiders'},
  ],
  // TODO: dummy data
  allProjectSelectOptions: [
    {label: 'Taobao', value: '000000000000000000000000'},
    {label: 'Tmall', value: '000000000000000000000001'},
    {label: 'JD', value: '000000000000000000000002'},
    {label: '163', value: '000000000000000000000003'},
    {label: 'Sina', value: '000000000000000000000004'},
    {label: '36kr', value: '000000000000000000000005'},
  ],
  allProjectTags: [
    'ecommerce',
    'news',
  ],
} as ProjectStoreState;

const getters = {
  ...getDefaultStoreGetters<Project>(),
} as ProjectStoreGetters;

const mutations = {
  ...getDefaultStoreMutations<Project>(),
  setAllProjectSelectOptions: (state: ProjectStoreState, options: SelectOption[]) => {
    state.allProjectSelectOptions = options;
  },
  setAllProjectTags: (state: ProjectStoreState, tags: string[]) => {
    state.allProjectTags = tags;
  },
} as ProjectStoreMutations;

const actions = {
  ...getDefaultStoreActions<Project>('/projects'),
  getAllProjectSelectOptions: async (state: ProjectStoreState) => {
    // TODO: implement
  },
  getAllProjectTags: async (state: ProjectStoreState) => {
    // TODO: implement
  },
} as ProjectStoreActions;

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions,
} as ProjectStoreModule;
