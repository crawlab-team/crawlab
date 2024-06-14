package generic

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

type ListSort struct {
	Key       string
	Direction SortDirection
}
