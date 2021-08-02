import {onBeforeUnmount, onMounted} from 'vue';

export const setupAutoUpdate = (fn: Function, interval?: number, handle?: number) => {
  if (!interval) interval = 5000;
  onMounted(() => {
    handle = setInterval(fn, interval);
  });
  onBeforeUnmount(() => {
    clearInterval(handle);
  });
};
