import {defineAsyncComponent} from 'vue';
import {RouteRecordRaw} from 'vue-router';
import {getLoadModuleOptions, loadModule} from '@/utils/sfc';

const endpoint = 'plugins';

export default [
  {
    path: endpoint,
    // component: defineAsyncComponent(() => loadModule('/vue/HelloWorld.vue', getLoadModuleOptions())),
    component: defineAsyncComponent(() => loadModule('/vue/App.vue', getLoadModuleOptions())),
  },
] as Array<RouteRecordRaw>;
