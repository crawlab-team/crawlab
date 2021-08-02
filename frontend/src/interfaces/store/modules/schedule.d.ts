type ScheduleStoreModule = BaseModule<ScheduleStoreState, ScheduleStoreGetters, ScheduleStoreMutations, ScheduleStoreActions>;

type ScheduleStoreState = BaseStoreState<Schedule>;

type ScheduleStoreGetters = BaseStoreGetters<Schedule>;

type ScheduleStoreMutations = BaseStoreMutations<Schedule>;

interface ScheduleStoreActions extends BaseStoreActions<Schedule> {
  enable: StoreAction<ScheduleStoreState, string>;
  disable: StoreAction<ScheduleStoreState, string>;
}
