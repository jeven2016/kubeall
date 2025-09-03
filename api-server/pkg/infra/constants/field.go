package constants

type Field string

const (
	FieldName              Field = "name"
	FieldCreationTimeStamp Field = "creationTimestamp"
	FieldNamespace         Field = "namespace"
	FieldStatus            Field = "status"
)

type SortOrder string

var (
	Asc  SortOrder = "asc"
	Desc SortOrder = "desc"
)

var SupportedSortFields = []Field{
	FieldName,
	FieldCreationTimeStamp,
	FieldNamespace,
	FieldStatus,
}

var SupportedSortOrders = []SortOrder{
	Asc,
	Desc,
}

var SupportedFilterFields = []Field{
	FieldName,
	FieldCreationTimeStamp,
	FieldNamespace,
	FieldStatus,
}
