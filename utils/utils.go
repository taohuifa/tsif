package utils

func ArrayAdd(array []interface{}, adds []interface{}) []interface{} {
	asize := len(adds)
	for i := 0; i < asize; i++ {
		array = append(array, adds[i])
	}
	return array
}
