import useDetail from '@/layouts/detail';

const useUserDetail = () => {
  return {
    ...useDetail('user'),
  };
};

export default useUserDetail;
