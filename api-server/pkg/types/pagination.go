package types

import "kubeall.io/api-server/pkg/infra/constants"

type Pagination struct {
	PageSize uint `json:"pageSize"`
	Page     uint `json:"page"`
}

type Query struct {
	Pagination Pagination          `json:"pagination"`
	SortBy     string              `json:"sortBy"`
	SortOrder  constants.SortOrder `json:"sortOrder"`
	Filters    map[string]string   `json:"filters"`
}

type PageResult struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalItems int   `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
	Items      []any `json:"items"`
}
