package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// 转化类型
func ChangeTypeByString(value string, kind reflect.Kind) interface{} {
	if kind == reflect.String {
		return value
	} else if kind == reflect.Bool {
		v, err := strconv.ParseBool(value)
		if err != nil {
			// panic(err)
			return false
		}
		return v
	} else if kind == reflect.Int {
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			// panic(err)
			return int(0)
		}
		return int(v)
	} else if kind == reflect.Int8 {
		v, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			// panic(err)
			return int8(0)
		}
		return int8(v)
	} else if kind == reflect.Int16 {
		v, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			// panic(err)
			return int16(0)
		}
		return int16(v)
	} else if kind == reflect.Int32 {
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			// panic(err)
			return int32(0)
		}
		return int32(v)
	} else if kind == reflect.Int64 {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			// panic(err)
			return int64(0)
		}
		return int64(v)
	} else if kind == reflect.Uint {
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			// panic(err)
			return uint(0)
		}
		return uint(v)
	} else if kind == reflect.Uint8 {
		v, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			// panic(err)
			return uint8(0)
		}
		return uint8(v)
	} else if kind == reflect.Uint16 {
		v, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			// panic(err)
			return uint16(0)
		}
		return uint16(v)
	} else if kind == reflect.Uint32 {
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			// panic(err)
			return uint32(0)
		}
		return uint32(v)
	} else if kind == reflect.Uint64 {
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			// panic(err)
			return uint64(0)
		}
		return v
	} else if kind == reflect.Float32 {
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			// panic(err)
			return float32(0)
		}
		return float32(v)
	} else if kind == reflect.Float64 {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			// panic(err)
			return float64(0)
		}
		return v
	}
	// return nil
	panic(errors.New("can't change string to type " + fmt.Sprint(kind) + " by \"" + value + "\""))
}

// 类型转化
func ChangeType(value interface{}, kind reflect.Kind) interface{} {
	vstr := fmt.Sprint(value)
	return ChangeTypeByString(vstr, kind)
}

// 提取参数
func GetByMap(confMap *map[string]string, key string, defaultValue interface{}) interface{} {
	value, ok := (*confMap)[key]
	if !ok {
		return defaultValue
	}

	// 根据类型转化
	valType := reflect.TypeOf(defaultValue)
	// fmt.Println(value, valType)
	return ChangeType(value, valType.Kind())
}
