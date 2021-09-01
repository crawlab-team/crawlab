interface CPlugin extends BaseModel {
  name?: string;
  description?: string;
  type?: string;
  proto?: string;
  active?: boolean;
  endpoint?: string;
  cmd?: string;
  ui_components?: PluginUIComponent[];
  ui_sidebar_navs?: MenuItem[];
}

interface PluginUIComponent {
  name?: string;
  title?: string;
  src?: string;
  type?: string;
  path?: string;
  parent_paths?: string[];
}
