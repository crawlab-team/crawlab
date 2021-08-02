import {Dayjs} from 'dayjs';

declare global {
  interface StatsResult extends Result {
    date?: string;
  }

  interface DateRange {
    start: Dayjs;
    end: Dayjs;
  }
}
