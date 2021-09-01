interface ListRequestParams {
  page?: number;
  size?: number;
  conditions?: FilterConditionData[] | string;
  all?: boolean | string | number;
}

interface BatchRequestPayload {
  ids: string[];
}

interface BatchRequestPayloadWithData<T = any> extends BatchRequestPayload {
  data: T;
  fields: string[];
}

type BatchRequestPayloadWithJsonStringData = BatchRequestPayloadWithData<string>;

interface FileRequestPayload {
  id?: string;
  path?: string;
  new_path?: string;
  data?: string;
  file?: File;
}

interface Response {
  status: string;
  message: string;
  error?: string;
}

type HttpResponse = Response;

interface ResponseWithData<T = any> extends Response {
  data?: T;
}

interface ResponseWithListData<T = any> extends ResponseWithData {
  data?: T[];
  total: number;
}
