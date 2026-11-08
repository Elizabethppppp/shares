package filtres

type ListFilter struct {
	Offset   int
	Limit    int
	Title    *string
	IsActive *bool
	SortBy   string
	SortDir  string
}
