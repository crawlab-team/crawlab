interface Schedule {
  _id?: string;
  name?: string;
  description?: string;
  spider_id?: string;
  spider_name?: string;
  cron?: string;
  cmd?: string;
  param?: string;
  mode?: TaskMode;
  node_ids?: string[];
  node_tags?: string[];
  enabled?: boolean;
}
