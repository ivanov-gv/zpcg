package utils

func AddMap[KeyType comparable, ValueType any](m1 map[KeyType]ValueType, m2 map[KeyType]ValueType) {
	for key, value := range m2 {
		m1[key] = value
	}
}

func RevertMap[KeyType comparable, ValueType comparable](_map map[KeyType]ValueType) map[ValueType]KeyType {
	result := make(map[ValueType]KeyType, len(_map))
	for key, value := range _map {
		result[value] = key
	}
	return result
}
