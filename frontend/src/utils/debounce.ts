import {plainClone} from '@/utils/object';

interface DebounceOptions {
  delay?: number;
}

const defaultDebounceOptions: DebounceOptions = {
  delay: 500,
};

const getDefaultDebounceOptions = (): DebounceOptions => {
  return plainClone(defaultDebounceOptions);
};

const normalizeDebounceOptions = (options?: DebounceOptions): DebounceOptions => {
  if (!options) options = getDefaultDebounceOptions();
  if (!options.delay) options.delay = defaultDebounceOptions.delay;
  return options;
};

export const debounce = <T = any>(fn: Function, options?: DebounceOptions): Function => {
  let handle: number | null = null;
  return () => {
    if (handle) clearTimeout(handle);
    const {delay} = normalizeDebounceOptions(options);
    handle = setTimeout(fn, delay);
  };
};
