package utils

func Intersection[KeyType comparable, ValueType comparable](m1, m2 map[KeyType]ValueType) map[KeyType]ValueType {
	var lowerLenMap, higherLenMap map[KeyType]ValueType
	if len(m1) < len(m2) {
		lowerLenMap = m1
		higherLenMap = m2
	} else {
		lowerLenMap = m2
		higherLenMap = m1
	}
	// iterate over the map with lower num of elems
	intersection := make(map[KeyType]ValueType, len(lowerLenMap))
	for key, value := range lowerLenMap {
		if _, found := higherLenMap[key]; found {
			intersection[key] = value
		}
	}
	return intersection
}
