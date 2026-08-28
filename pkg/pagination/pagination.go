package pagination

import (
	"strconv"
)

const DefaultLimit = 20
const MaxLimit = 100

type Params struct {
	Cursor string
	Limit  int
	Sort   string
	Order  string
}

func FromQuery(cursor string, limitStr string, sort, order string) Params {
	limit := DefaultLimit
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		if l > MaxLimit {
			limit = MaxLimit
		} else {
			limit = l
		}
	}
	if sort == "" {
		sort = "created_at"
	}
	if order == "" {
		order = "DESC"
	}
	return Params{
		Cursor: cursor,
		Limit:  limit,
		Sort:   sort,
		Order:  order,
	}
}
