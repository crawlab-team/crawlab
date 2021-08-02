interface TableHeaderProps {
  column: TableColumn;
  index?: number;
}

interface TableHeaderActionStatusMap {
  filter: TableHeaderActionStatus;
  sort: TableHeaderActionStatus;
}

interface TableHeaderActionStatus {
  active: boolean;
  focused: boolean;
}
