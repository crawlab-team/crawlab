interface MetricProps {
  title?: string;
  value?: number | string;
  icon?: Icon;
  color?: string;
  clickable?: boolean;
}

interface MetricMeta {
  name: string;
  key: string;
  value: number | string;
  icon: Icon;
  color?: string | Function;
  path?: string;
}
