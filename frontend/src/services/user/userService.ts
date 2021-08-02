import {Store} from 'vuex';
import {getDefaultService} from '@/utils/service';

const useUserService = (store: Store<RootStoreState>): Services<User> => {
  const ns = 'user';

  return {
    ...getDefaultService<User>(ns, store),
  };
};

export default useUserService;
