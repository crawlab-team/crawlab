import {
  getDefaultStoreActions,
  getDefaultStoreGetters,
  getDefaultStoreMutations,
  getDefaultStoreState
} from '@/utils/store';
import useRequest from '@/services/request';
import {TAB_NAME_OVERVIEW, TAB_NAME_TASKS} from '@/constants/tab';

const {
  post,
} = useRequest();

const state = {
  ...getDefaultStoreState<Schedule>('schedule'),
  tabs: [
    {id: TAB_NAME_OVERVIEW, title: 'Overview'},
    {id: TAB_NAME_TASKS, title: 'Tasks'},
  ],
} as ScheduleStoreState;

const getters = {
  ...getDefaultStoreGetters<Schedule>(),
} as ScheduleStoreGetters;

const mutations = {
  ...getDefaultStoreMutations<Schedule>(),
} as ScheduleStoreMutations;

const actions = {
  ...getDefaultStoreActions<Schedule>('/schedules'),
  enable: async (ctx: StoreActionContext, id: string) => {
    return await post(`/schedules/${id}/enable`);
  },
  disable: async (ctx: StoreActionContext, id: string) => {
    return await post(`/schedules/${id}/disable`);
  },
} as ScheduleStoreActions;

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions,
} as ScheduleStoreModule;
