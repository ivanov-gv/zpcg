package utils

import (
	"fmt"
	"sync"
)

func AddMap[KeyType comparable, ValueType any](m1 map[KeyType]ValueType, m2 map[KeyType]ValueType) {
	for key, value := range m2 {
		m1[key] = value
	}
}

func SyncMapToMap[KeyType comparable, ValueType any](syncMap *sync.Map, size int) (map[KeyType]ValueType, error) {
	result := make(map[KeyType]ValueType, size)
	var (
		keyConversionValid   = true
		valueConversionValid = true
		syncMapKeyExample    any
		syncMapValueExample  any
	)
	syncMap.Range(func(key, value any) bool {
		syncMapKeyExample, syncMapValueExample = key, value
		_key, ok := key.(KeyType)
		if !ok {
			keyConversionValid = false
			return false
		}
		_value, ok := value.(ValueType)
		if !ok {
			valueConversionValid = false
			return false
		}
		result[_key] = _value
		return true
	})
	if !keyConversionValid || !valueConversionValid {
		return nil, fmt.Errorf("failed to convert sync.map(%T, %T) to a map[%T]%T",
			syncMapKeyExample, syncMapValueExample, *new(KeyType), *new(ValueType))
	}
	return result, nil
}
