interface TagProps {
  label?: string;
  tooltip?: string;
  type?: BasicType;
  color?: string;
  backgroundColor?: string;
  borderColor?: string;
  icon?: Icon;
  spinning?: boolean;
  width?: string;
  effect?: BasicEffect;
  clickable?: boolean;
  closable?: boolean;
  disabled?: boolean;
  tag?: Tag;
}

type TagData = TagProps;
