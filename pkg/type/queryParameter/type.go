package queryParameter

import (
	"architecture_go_2/pkg/type/pagination"
	"architecture_go_2/pkg/type/sort"
)

type QueryParameter struct {
	Sorts      sort.Sorts
	Pagination pagination.Pagination
	/*Тут можно добавить фильтр*/
}
