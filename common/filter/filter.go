package filter

import "template/common/mapping"

func FilterResult[T any](target []T) []map[string]any {
	result := make([]map[string]any, 0)
	result = AdvanceFilter(target, func(item T) map[string]any {
		return mapping.Omit(item, "ID", "CreatedAt", "UpdatedAt", "DeletedAt")
	})
	return result
}

// 泛型函数 ArrayFilter，接收一个目标数组 target 和回调函数 cb
// 泛型函数 ArrayFilter，接收一个目标数组 target 和回调函数 cb
func ArrayFilter[T any](target []T, cb func(T) bool) []T {
	var result []T // 存储过滤后的结果

	for _, item := range target {
		// 调用回调函数 cb，对每个元素进行判断
		if cb(item) {
			result = append(result, item) // 如果条件成立，添加到结果中
		}
	}

	return result
}

func Filter[T any](target []T, cb func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range target {
		if cb(item) {
			result = append(result, item)
		}
	}
	return result
}

func GORMFieldsFilter[T any](target T) map[string]any {
	return mapping.Omit(target, "ID", "CreatedAt", "UpdatedAt", "DeletedAt")
}

func AdvanceFilter[T any](target []T, call func(item T) map[string]any) []map[string]any {
	result := make([]map[string]any, 0)
	for _, item := range target {
		filted := call(item)
		if filted == nil {
			continue
		}
		result = append(result, filted)
	}
	return result
}
