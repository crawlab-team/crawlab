import {Store} from 'vuex';
import {getDefaultService} from '@/utils/service';

const useTokenService = (store: Store<RootStoreState>): Services<Token> => {
  const ns = 'task';

  return {
    ...getDefaultService<Token>(ns, store),
  };
};

export default useTokenService;
