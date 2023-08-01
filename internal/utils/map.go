package utils

func AddMap[KeyType comparable, ValueType any](m1 map[KeyType]ValueType, m2 map[KeyType]ValueType) {
	for key, value := range m2 {
		m1[key] = value
	}
}
