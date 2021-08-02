import {Store} from 'vuex';
import {getDefaultService} from '@/utils/service';

const useScheduleService = (store: Store<RootStoreState>): Services<Schedule> => {
  const ns = 'schedule';

  return {
    ...getDefaultService<Schedule>(ns, store),
  };
};

export default useScheduleService;
