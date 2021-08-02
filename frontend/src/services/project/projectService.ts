import {Store} from 'vuex';
import {getDefaultService} from '@/utils/service';

const useProjectService = (store: Store<RootStoreState>): Services<Project> => {
  const ns = 'project';

  return {
    ...getDefaultService<Project>(ns, store),
  };
};

export default useProjectService;
