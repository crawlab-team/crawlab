interface SpiderServices extends Services<Spider> {
  listDir: (id: string, path: string) => Promise<ResponseWithData<FileNavItem[]>>;
  listRootDir: (id: string) => Promise<ResponseWithData<FileNavItem[]>>;
  getFile: (id: string, path: string) => Promise<ResponseWithData<string>>;
  getFileInfo: (id: string, path: string) => Promise<ResponseWithData<FileNavItem>>;
  saveFile: (id: string, path: string, data: string) => Promise<Response>;
  saveFileBinary: (id: string, path: string, file: File) => Promise<Response>;
  saveDir: (id: string, path: string) => Promise<Response>;
  renameFile: (id: string, path: string, new_path: string) => Promise<Response>;
  deleteFile: (id: string, path: string) => Promise<Response>;
  copyFile: (id: string, path: string, new_path: string) => Promise<Response>;
}
