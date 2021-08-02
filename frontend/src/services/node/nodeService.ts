import {Store} from 'vuex';
import {getDefaultService} from '@/utils/service';

type Node = CNode;

const useNodeService = (store: Store<RootStoreState>): Services<Node> => {
  const ns = 'node';

  return {
    ...getDefaultService<Node>(ns, store),
  };
};

export default useNodeService;
