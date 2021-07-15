interface FileEditorNavMenuProps {
  items: FileNavItem[];
}

interface FileEditorNavMenuClickStatus {
  clicked: boolean;
  item?: FileNavItem;
}

interface FileEditorNavMenuCache<T = any> {
  [key: string]: T;
}
