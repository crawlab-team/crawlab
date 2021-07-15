interface TagInputProps {
  modelValue: Tag[];
  disabled: boolean;
}

interface TagInputOption extends Tag {
  isEdit?: boolean;
}
