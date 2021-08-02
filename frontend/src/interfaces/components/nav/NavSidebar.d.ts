interface NavSidebar {
  scroll: (id: string) => void;
}

interface NavSidebarProps {
  items: NavItem[];
  activeKey?: string;
  collapsed?: boolean;
  showActions?: boolean;
}
