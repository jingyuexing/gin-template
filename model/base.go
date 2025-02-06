package model

type ModelOption[T any] func(*T)

// ModelOption 通用创建
func NewModel[T any](opts ...ModelOption[T]) *T {
	model := new(T)
	for _, opt := range opts {
		opt(model)
	}
	return model
}

// WithString 通用字符串字段设置
func WithString[T any](field *string, value string) ModelOption[T] {
	return func(m *T) {
		*field = value
	}
}

// WithInt64 通用int64字段设置
func WithInt64[T any](field *int64, value int64) ModelOption[T] {
	return func(m *T) {
		*field = value
	}
}