package utils

func GetSpiderCol(col string, name string) string {
	if col == "" {
		return "results_" + name
	}
	return col
}
