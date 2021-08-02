import dayjs, {Dayjs} from 'dayjs';

export const DEFAULT_DATE_FORMAT = 'YYYY-MM-DD';

export const spanDateRange = (start: Dayjs | string, end: Dayjs | string, data: StatsResult[], dateKey?: string): StatsResult[] => {
  // date key
  const key = dateKey || 'date';

  // format
  const format = DEFAULT_DATE_FORMAT;

  // cache data
  const cache = new Map<string, StatsResult>();
  data.forEach(d => cache.set(d[key], d));

  // results
  const results = [] as StatsResult[];

  // iterate
  for (let date = dayjs(start, format); date.format(format) <= dayjs(end, format).format(format); date = date.add(1, 'day')) {
    let item = cache.get(date.format(format));
    if (!item) item = {};
    item[key] = date.format(format);
    results.push(item);
  }

  return results;
};
