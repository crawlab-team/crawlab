interface FileUploadProps {
  mode?: string;
  getInputProps?: Function;
  open?: Function;
}

interface FileUploadModeOption {
  label: string;
  value: string;
}

interface FileUploadInfo {
  dirName?: string;
  fileCount?: number;
  filePaths?: string[];
}
