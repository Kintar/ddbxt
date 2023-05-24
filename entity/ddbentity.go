package entity

// KeyValue represents a dynamoDB key value, either partition or sort
type KeyValue[T any] struct {
	Name  string
	Value T
}

type PartitionKey[T any] KeyValue[T]

type SortKey[T any] KeyValue[T]

type SimpleKey[T any] PartitionKey[T]

type CompositeKey[T, S any] struct {
	PartitionKey[T]
	SortKey[S]
}

type DdbEntity[T SimpleKey[any] | CompositeKey[any, any]] interface {
	Key() T
	TableName() string
}
