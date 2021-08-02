export const EMPTY_OBJECT_ID = '000000000000000000000000';

export const isZeroObjectId = (id: string): boolean => {
  return !id || id === EMPTY_OBJECT_ID;
};
