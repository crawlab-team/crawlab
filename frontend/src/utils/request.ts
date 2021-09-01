export const getRequestBaseUrl = (): string => {
  return process.env.VUE_APP_API_BASE_URL || 'http://localhost:8000';
};
