interface TableHeaderDialogProps {
  visible: boolean;
  column: TableColumn;
  actionStatusMap: TableHeaderActionStatusMap;
  filter?: TableHeaderDialogFilterData;
  sort?: SortData;
}

interface TableHeaderDialogValue {
  sort?: SortData;
  filter?: TableHeaderDialogFilterData;
}
