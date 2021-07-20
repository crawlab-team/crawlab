interface FileUploadProps {
  mode?: string;
  getInputProps?: Function;
  open?: Function;
}

interface FileUploadModeOption {
  label: string;
  value: string;
}

interface FileUploadDirInfo {
  dirName: string;
  fileCount: number;
}
