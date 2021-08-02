const useIcon = () => {
  // implementation
  const isFaIcon = (icon: Icon) => {
    if (Array.isArray(icon)) {
      return icon.length > 0 && icon[0].substr(0, 2) === 'fa';
    } else {
      return icon?.substr(0, 2) === 'fa';
    }
  };

  const getFontSize = (size: IconSize) => {
    switch (size) {
      case 'large':
        return '24px';
      case 'normal':
        return '16px';
      case 'small':
        return '12px';
      case 'mini':
        return '10px';
      default:
        return size || '16px';
    }
  };

  return {
    // public variables and methods
    isFaIcon,
    getFontSize,
  };
};

export default useIcon;
