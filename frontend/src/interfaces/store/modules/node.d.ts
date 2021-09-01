type NodeStoreModule = BaseModule<NodeStoreState, NodeStoreGetters, NodeStoreMutations, NodeStoreActions>;

type NodeStoreState = BaseStoreState<CNode>;

type NodeStoreGetters = BaseStoreGetters<CNode>;

interface NodeStoreMutations extends BaseStoreMutations<CNode> {
  setAllNodeSelectOptions: StoreMutation<BaseStoreState<CNode>, SelectOption[]>;
  setAllNodeTags: StoreMutation<BaseStoreState<CNode>, string[]>;
}

type NodeStoreActions = BaseStoreActions<CNode>;
