type ProjectStoreModule = BaseModule<ProjectStoreState, ProjectStoreGetters, ProjectStoreMutations, ProjectStoreActions>;

interface ProjectStoreState extends BaseStoreState<Project> {
  allProjectSelectOptions: SelectOption[];
  allProjectTags: string[];
}

type ProjectStoreGetters = BaseStoreGetters<Project>;

interface ProjectStoreMutations extends BaseStoreMutations<Project> {
  setAllProjectSelectOptions: StoreMutation<ProjectStoreState, SelectOption[]>;
  setAllProjectTags: StoreMutation<ProjectStoreState, string[]>;
}

interface ProjectStoreActions extends BaseStoreActions<Project> {
  getAllProjectSelectOptions: StoreAction<ProjectStoreState>;
  getAllProjectTags: StoreAction<ProjectStoreState>;
}
