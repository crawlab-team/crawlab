interface TableHeaderDialogFilterProps {
  column?: TableColumn;
  searchString?: string;
  conditions?: FilterConditionData[];
}

interface TableHeaderDialogFilterData {
  searchString?: string;
  conditions?: FilterConditionData[];
  items?: string[];
}
