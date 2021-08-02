interface BaseModel {
  _id?: string;
  tags?: Tag[];

  [field: string]: any;
}
