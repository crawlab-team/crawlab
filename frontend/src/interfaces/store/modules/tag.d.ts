type TagStoreModule = BaseModule<TagStoreState, TagStoreGetters, TagStoreMutations, TagStoreActions>;

type TagStoreState = BaseStoreState<Tag>;

type TagStoreGetters = BaseStoreGetters<Tag>;

interface TagStoreMutations extends BaseStoreMutations<Tag> {
  setAllTagSelectOptions: StoreMutation<BaseStoreState<Tag>, SelectOption[]>;
  setAllTagTags: StoreMutation<BaseStoreState<Tag>, string[]>;
}

type TagStoreActions = BaseStoreActions<Tag>;
