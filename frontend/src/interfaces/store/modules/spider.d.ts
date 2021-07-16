type SpiderStoreModule = BaseModule<SpiderStoreState, SpiderStoreGetters, SpiderStoreMutations, SpiderStoreActions>;

interface SpiderStoreState extends BaseStoreState<Spider> {
  fileNavItems: FileNavItem[];
  activeNavItem?: FileNavItem;
  fileContent: string;
}

type SpiderStoreGetters = BaseStoreGetters<SpiderStoreState>;

interface SpiderStoreMutations extends BaseStoreMutations<Spider> {
  setFileNavItems: StoreMutation<BaseStoreState<Spider>, FileNavItem[]>;
  setActiveFileNavItem: StoreMutation<BaseStoreState<Spider>, FileNavItem>;
  resetActiveFileNavItem: StoreMutation<BaseStoreState<Spider>>;
  setFileContent: StoreMutation<BaseStoreState<Spider>, string>;
  resetFileContent: StoreMutation<BaseStoreState<Spider>>;
}

interface SpiderStoreActions extends BaseStoreActions<Spider> {
  runById: StoreAction<BaseStoreState, { id: string; options: SpiderRunOptions }>;
  listDir: StoreAction<BaseStoreState, FileRequestPayload>;
  getFile: StoreAction<BaseStoreState, FileRequestPayload>;
  getFileInfo: StoreAction<BaseStoreState, FileRequestPayload>;
  saveFile: StoreAction<BaseStoreState, FileRequestPayload>;
  saveFileBinary: StoreAction<BaseStoreState, FileRequestPayload>;
  saveDir: StoreAction<BaseStoreState, FileRequestPayload>;
  renameFile: StoreAction<BaseStoreState, FileRequestPayload>;
  deleteFile: StoreAction<BaseStoreState, FileRequestPayload>;
  copyFile: StoreAction<BaseStoreState, FileRequestPayload>;
}
