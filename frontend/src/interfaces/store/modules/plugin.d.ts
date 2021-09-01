type PluginStoreModule = BaseModule<PluginStoreState, PluginStoreGetters, PluginStoreMutations, PluginStoreActions>;

type PluginStoreState = BaseStoreState<CPlugin>;

type PluginStoreGetters = BaseStoreGetters<CPlugin>;

type PluginStoreMutations = BaseStoreMutations<CPlugin>;

type PluginStoreActions = BaseStoreActions<CPlugin>;
