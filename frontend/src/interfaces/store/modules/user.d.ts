type UserStoreModule = BaseModule<UserStoreState, UserStoreGetters, UserStoreMutations, UserStoreActions>;

type UserStoreState = BaseStoreState<User>;

type UserStoreGetters = BaseStoreGetters<User>;

type UserStoreMutations = BaseStoreMutations<User>;

interface UserStoreActions extends BaseStoreActions<User> {
  changePassword: StoreAction<UserStoreState, { id: string; password: string }>;
}
