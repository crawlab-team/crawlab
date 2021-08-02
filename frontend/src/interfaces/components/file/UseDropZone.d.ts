import {FileWithPath} from 'file-selector';

declare global {
  type FileAccept = string | string[];
  type FileHandler = (evt: Event) => void;
  type FileErrorCode = 'file-invalid-type' | 'file-too-large' | 'file-too-small' | 'too-many-files' | string;
  type FileRejectionError = {
    code: FileErrorCode;
    message: string;
  } | null | boolean;
  type InputFile = (FileWithPath | DataTransferItem) & {
    path?: string;
    size?: number;
  };
  type FileRejectReason = {
    file: InputFile;
    errors: FileRejectionError[];
  };
}
