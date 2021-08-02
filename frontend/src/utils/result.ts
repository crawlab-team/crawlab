export const getFieldsFromData = (data: TableData<Result>) => {
  if (data.length === 0) {
    return [];
  }
  const item = data[0];
  if (typeof item !== 'object') return [];
  return Object.keys(item).map(key => {
    return {
      key,
    };
  });
};
