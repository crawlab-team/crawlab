interface CNode extends BaseModel {
  name?: string;
  ip?: string;
  mac?: string;
  hostname?: string;
  description?: string;
  key?: string;
  is_master?: boolean;
  status?: string;
  enabled?: boolean;
  active?: boolean;
  available_runners?: number;
  max_runners?: number;
}
