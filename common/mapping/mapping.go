package mapping

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// KeyFunc defines a function type for generating map keys from struct fields
type KeyFunc func(field reflect.StructField) string

// omit specify field in target struct
//
//	 type People struct {
//	    Name    string
//	    Age     int
//	    Address string
//	}
//
//	p := &People{
//	    Name:    "william",
//	    Age:     20,
//	    Address: "BC",
//	}
//
//	Omit(p,[]string{"Name"})
func Omit[T any](target T, fields ...string) map[string]any {
	mapping := map[string]bool{}
	for _, value := range fields {
		mapping[SnakeCase(value)] = true
		mapping[value] = true
	}
	result := StructFilter(target, func(field string, val reflect.Value) bool {
		_, ok := mapping[field]
		return !ok
	}, nil)
	return result
}

// pick specify field in target struct
//
//	 type People struct {
//	    Name    string
//	    Age     int
//	    Address str
//	}
//
//	p := &People{
//	    Name:    "william",
//	    Age:     20,
//	    Address: "BC",
//	}
//
//	Pick(p,[]string{"Name"})
func Pick[T any](target T, fields ...string) map[string]any {
	mapping := map[string]bool{}
	for _, value := range fields {
		mapping[value] = true
		mapping[SnakeCase(value)] = true
	}
	result := StructFilter(target, func(field string, val reflect.Value) bool {
		_, ok := mapping[field]
		return ok
	}, nil)
	return result
}

func jsonKeyFunc(field reflect.StructField) string {
	if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
		return jsonTag
	}
	return ""
}

// SnakeCase 将 CamelCase 转换为 snake_case，并处理缩写词
func SnakeCase(s string) string {
	var result []rune
	n := len(s)
	allUpper := true

	// 检查字符串是否全是大写
	for _, r := range s {
		if unicode.IsLower(r) {
			allUpper = false
			break
		}
	}

	// 如果字符串是全大写的，直接转换为小写并返回
	if allUpper {
		return strings.ToLower(s)
	}

	for i := 0; i < n; i++ {
		r := rune(s[i])

		// 如果是大写字母
		if unicode.IsUpper(r) {
			// 如果当前大写字母是开头，或者之前的字符也是大写，则不加下划线
			if i > 0 && !unicode.IsUpper(rune(s[i-1])) {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// GetFieldName 根据字段的 JSON 标签或字段名返回正确的字段名称
func GetFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return SnakeCase(field.Name)
	}
	jsonName := strings.Split(jsonTag, ",")[0] // 获取 JSON 标签的第一个部分
	if jsonName != "" && jsonName != "-" {
		return jsonName
	}
	return SnakeCase(field.Name)
}

// GetFieldNames 获取结构体字段的名称
func GetFieldNames(i interface{}) []string {
	var result []string
	t := reflect.TypeOf(i)

	// 如果是指针类型，获取指针指向的元素类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 处理结构体类型
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldType := field.Type

			// 处理嵌入字段（embed）
			if field.Anonymous {
				embeddedFields := GetFieldNames(reflect.New(fieldType).Elem().Interface())
				result = append(result, embeddedFields...)
				continue
			}

			// 获取字段名称
			fieldName := GetFieldName(field)
			result = append(result, fieldName)
		}
	}
	return result
}

func structFilterRecursive(
	value reflect.Value,
	t reflect.Type,
	callback func(field string, val reflect.Value) bool,
	keyFunc KeyFunc,
	result map[string]any,
) {
	for i := 0; i < value.NumField(); i++ {
		field := t.Field(i)
		fieldValue := value.Field(i)

		key := keyFunc(field)

		if field.Anonymous {
			// If the field is an embedded struct, recurse into it
			structFilterRecursive(fieldValue, fieldValue.Type(), callback, keyFunc, result)
		} else if callback(key, fieldValue) && key != "" {
			result[key] = fieldValue.Interface()
		}
	}
}
func StructFilter[T any](
	target T,
	callback func(field string, val reflect.Value) bool,
	key KeyFunc,
) map[string]any {
	reflectValue := reflect.ValueOf(target)
	reflectType := reflect.TypeOf(target)
	if reflectValue.Kind() == reflect.Pointer {
		reflectValue = reflectValue.Elem()
		reflectType = reflectType.Elem()
	}
	GenKey := key
	if key == nil {
		GenKey = GetFieldName
	}
	result := make(map[string]any)
	structFilterRecursive(reflectValue, reflectType, callback, GenKey, result)
	return result
}

// BindMapToStruct 将一个 map 中的值绑定到结构体的字段上
func BindMapToStruct(
	data map[string]interface{},
	target interface{},
	callback func(reflect.StructField) string,
) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr && v.Kind() != reflect.Struct && target == nil {
		return errors.New("target must be a pointer to a struct")
	}
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		tag := callback(fieldType)
		if tag == "" {
			continue
		}
		value, ok := data[tag]
		if !ok {
			continue
		}

		if !field.CanSet() {
			continue
		}

		if err := setFieldValue(field, value); err != nil {
			return fmt.Errorf("error setting field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// setFieldValue 设置结构体字段的值
func setFieldValue(field reflect.Value, value interface{}) error {
	switch field.Kind() {
	case reflect.String:
		strValue, ok := value.(string)
		if !ok {
			return errors.New("value is not a string")
		}
		field.SetString(strValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, ok := value.(int64)
		if !ok {
			return errors.New("value is not an int")
		}
		field.SetInt(intValue)
	case reflect.Float32, reflect.Float64:
		floatValue, ok := value.(float64)
		if !ok {
			return errors.New("value is not a float")
		}
		field.SetFloat(floatValue)
	case reflect.Bool:
		booleanValue, ok := value.(bool)
		if !ok {
			return errors.New("value is not a boolean")
		}
		field.SetBool(booleanValue)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}

// GroupBy 用于对数据进行分组，支持结构体和map
func GroupBy[T any](data T, keyField string) map[any][]T {
	grouped := make(map[any][]T)

	// 通过反射获取传入的slice
	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Slice {
		panic("data should be a slice")
	}

	// 遍历slice中的每个元素
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)

		// 判断当前元素是结构体还是map
		var keyVal reflect.Value
		if item.Kind() == reflect.Struct {
			// 结构体类型，获取字段的值
			keyVal = item.FieldByName(keyField)
			if !keyVal.IsValid() {
				panic(fmt.Sprintf("field %s does not exist", keyField))
			}
		} else if item.Kind() == reflect.Map {
			// map类型，获取key对应的值
			keyVal = item.MapIndex(reflect.ValueOf(keyField))
			if !keyVal.IsValid() {
				panic(fmt.Sprintf("key %s does not exist", keyField))
			}
		} else {
			panic("element should be a struct or map")
		}

		// 使用interface{}作为key
		groupKey := keyVal.Interface()
		grouped[groupKey] = append(grouped[groupKey], item.Interface().(T))
	}

	return grouped
}
