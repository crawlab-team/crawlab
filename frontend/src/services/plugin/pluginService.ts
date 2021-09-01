import {Store} from 'vuex';
import {getDefaultService} from '@/utils/service';

type Plugin = CPlugin;

const usePluginService = (store: Store<RootStoreState>): Services<Plugin> => {
  const ns = 'plugin';

  return {
    ...getDefaultService<Plugin>(ns, store),
  };
};

export default usePluginService;
