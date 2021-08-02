interface TabProps {
  tab?: Tab;
  icon?: Icon;
  showTitle: boolean;
  showClose: boolean;
}

interface Tab {
  id?: number;
  path: string;
  dragging?: boolean;
  isAction?: boolean;
}
