import {Placement} from '@popperjs/core';

declare global {
  interface FileEditorNavTabsShowMoreContextMenuProps {
    tabs: FileNavItem[];
    visible?: boolean;
    placement?: Placement;
  }
}
