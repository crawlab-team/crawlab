import {createStore, Store} from 'vuex';
import login from '@/store/modules/login';
import layout from '@/store/modules/layout';
import node from '@/store/modules/node';
import project from '@/store/modules/project';
import spider from '@/store/modules/spider';
import task from '@/store/modules/task';
import file from '@/store/modules/file';
import tag from '@/store/modules/tag';
import dataCollection from '@/store/modules/dataCollection';
import schedule from '@/store/modules/schedule';
import user from '@/store/modules/user';
import token from '@/store/modules/token';

export default createStore<RootStoreState>({
  modules: {
    login,
    layout,
    node,
    project,
    spider,
    task,
    file,
    tag,
    dataCollection,
    schedule,
    user,
    token,
  },
}) as Store<RootStoreState>;
