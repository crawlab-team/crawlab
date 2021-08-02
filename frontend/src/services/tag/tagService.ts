import {Store} from 'vuex';
import {getDefaultService} from '@/utils/service';

const useTagService = (store: Store<RootStoreState>): Services<Tag> => {
  const ns = 'tag';

  return {
    ...getDefaultService<Tag>(ns, store),
  };
};

export default useTagService;
