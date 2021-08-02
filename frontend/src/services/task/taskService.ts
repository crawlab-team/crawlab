import {Store} from 'vuex';
import {getDefaultService} from '@/utils/service';

const useTaskService = (store: Store<RootStoreState>): Services<Task> => {
  const ns = 'task';

  return {
    ...getDefaultService<Task>(ns, store),
  };
};

export default useTaskService;
