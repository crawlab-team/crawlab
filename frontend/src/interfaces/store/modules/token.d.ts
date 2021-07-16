type TokenStoreModule = BaseModule<TokenStoreState, TokenStoreGetters, TokenStoreMutations, TokenStoreActions>;

type TokenStoreState = BaseStoreState<Token>;

type TokenStoreGetters = BaseStoreGetters<Token>;

type TokenStoreMutations = BaseStoreMutations<Token>;

type TokenStoreActions = BaseStoreActions<Token>;
