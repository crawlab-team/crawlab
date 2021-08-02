interface IconProps {
  icon?: Icon;
  spinning?: boolean;
  size: IconSize;
}

type Icon = string | string[];

type IconSize = 'mini' | 'small' | 'normal' | 'large' | string;
