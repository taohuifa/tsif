package reflect

import (
	"fmt"
	"reflect"
)

// 对象表
type ObjTable struct {
	Name   string   // 表名
	Cloum  []string // key名
	Values []interface{}
}

//对象转成表
func Obj2Table(obj interface{}) *ObjTable {
	typeObj := reflect.TypeOf(obj)
	// typeElem := typeObj.Elem() 如果typeObj是指针, 用这个获取类
	// fmt.Println("typeObj", typeObj, typeElem)
	//表名
	tb := new(ObjTable)
	tb.Name = typeObj.Name()
	//列
	fildNum := typeObj.NumField()
	tb.Cloum = make([]string, fildNum)
	tb.Values = make([]interface{}, fildNum)
	for i := 0; i < fildNum; i++ {
		//列名
		cloumName := typeObj.Field(i).Name
		tb.Cloum[i] = cloumName
		//列值
		val := reflect.ValueOf(obj).FieldByName(cloumName)
		tb.Values[i] = val
	}
	return tb
}

//对象转成表
func Obj2Map(obj interface{}) *map[string]interface{} {
	objType := reflect.TypeOf(obj)
	// fmt.Println("typeObj", objType)

	// 遍历解析
	out_map := make(map[string]interface{})
	for i := 0; i < objType.NumField(); i++ {
		//变量数据
		field := objType.Field(i)
		cloumName := field.Name
		//列值
		refVal := reflect.ValueOf(obj).FieldByName(cloumName) //获取对应值的反射对象
		value := refVal.Interface()                           // 获取对应数据值
		//装入map
		out_map[cloumName] = value
		// fmt.Println(i, field, cloumName, value, reflect.TypeOf(value))
	}
	return &out_map
}

// 表转入对象
func Map2Obj(v_map map[string]interface{}, obj interface{}, objElem reflect.Value) bool {
	objType := objElem.Type() // 获取对应类型
	fmt.Println("map2obj type: ", objType, objElem)

	//objElem.Field(0).SetInt(32)
	//objElem.Field(0).Set(reflect.ValueOf(3))

	// 遍历
	for key, value := range v_map {
		field := objElem.FieldByName(key)
		// fmt.Println(key, value, reflect.TypeOf(value), " field=", field.Type(), field.CanSet())
		field.Set(reflect.ValueOf(value))
	}

	return true
}
