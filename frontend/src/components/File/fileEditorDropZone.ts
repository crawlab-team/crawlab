import {useDropzone} from 'vue3-dropzone';

const useFileEditorDropZone = () => {
  const onDrop = (acceptedFiles: InputFile[], rejectReasons: FileRejectReason[], event: Event) => {
    console.log(acceptedFiles);
  };

  const {
    getRootProps,
    getInputProps,
    open,
  } = useDropzone({
    onDrop,
  });

  return {
    getRootProps,
    getInputProps,
    open,
  };
};

export default useFileEditorDropZone;
