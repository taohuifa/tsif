package main

import (
	"fmt"
	"reflect"

	"github.com/tsif/app"
	"github.com/tsif/utils"
	"github.com/tsif/utils/http"

	Log "github.com/tsif/component/log"
)

func test01() {
	fmt.Println("hello world")

	kinds := []reflect.Kind{
		reflect.Bool,
		reflect.Float32,
		reflect.Float64,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.String,
	}

	str := "01234567"
	for i := 0; i < len(kinds); i++ {
		kind := kinds[i]
		v := utils.ChangeType(str, kind)
		fmt.Println(str, kind, v, reflect.TypeOf(v))
	}

}

func test02() {
	// init app
	initFunc := func(context *app.AppContext, params ...interface{}) bool {
		context.Stop(5 * 1000)
		return true
	}
	destroyFunc := func(context *app.AppContext) {

	}
	appCtx := app.AppContext{Name: "app", InitFunc: initFunc, DestroyFunc: destroyFunc}
	// start
	err := appCtx.Start(1, 2, "3")
	if err != nil {
		Log.Info("app start fail! " + err.Error())
	}
}

func test03() {
	params := make(map[string]interface{})
	params["1"] = 1
	params["2"] = 2
	params["a"] = "松松散散"
	mstr := http.CreateParams(params)

	Log.Infof("params ", params, mstr)
}

func main() {
	// test01()
	// test02()
	test03()
}
