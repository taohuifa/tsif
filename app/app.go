package app

import (
	"errors"
	"fmt"
	"math"
	"time"

	Log "github.com/tsif/component/log"
)

const (
	APPSTATE_TO_SHUTDOWN   = -1 // 开始停服
	APPSTATE_SHUTDOWN      = 0  // 停止
	APPSTATE_RUN           = 1  // 运行中
	APPSTATE_STARTING      = 2  // 启动中
	APPSTATE_WAIT_TO_CLOSE = 3  // 准备关闭

)

// App上下文
type AppContext struct {
	Name            string // appname
	RunState        int    // 运行状态
	UpdateIinterval uint   // 间隔时间

	shutdownTime int64 // 关闭时间

	InitFunc    func(context *AppContext, params ...interface{}) bool                               // 初始化函数
	DestroyFunc func(context *AppContext)                                                           // 销毁函数
	UpdateFunc  func(context *AppContext, count int, dt int, nowtime time.Time, prevtime time.Time) // 更新函数
}

func (this *AppContext) Start(params ...interface{}) error {
	var result bool
	// init
	this.RunState = APPSTATE_STARTING
	Log.Info("app init: " + fmt.Sprint(params))
	if this.InitFunc != nil {
		result = this.InitFunc(this, params)
		if !result {
			return errors.New("init error")
		}
	}
	Log.Info("app start")

	// update
	this.RunState = APPSTATE_RUN
	prevtime := time.Now()
	count := 0
	for this.IsRunning() {
		updateIinterval := time.Duration(math.Max(float64(this.UpdateIinterval), 1000))
		time.Sleep(updateIinterval * time.Millisecond)
		//wait shutdown
		nowtime := time.Now()
		if this.shutdownTime > 0 {
			dtime := this.shutdownTime - nowtime.UnixNano()
			// Log.Infof("app wait stop", dtime)
			if dtime <= 0 {
				this.RunState = APPSTATE_TO_SHUTDOWN
				break
			}
		}
		// counter
		count++
		dt := int((nowtime.UnixNano() - prevtime.UnixNano()) / 1e6)
		prevtime = nowtime
		// Log.Infof("app update", dt, count)
		// update
		if this.UpdateFunc != nil {
			this.UpdateFunc(this, count, dt, nowtime, prevtime)
		}
	}

	// destory
	this.RunState = APPSTATE_TO_SHUTDOWN
	Log.Info("app destorying")
	if this.DestroyFunc != nil {
		this.DestroyFunc(this)
	}
	Log.Info("app destory")
	this.RunState = APPSTATE_SHUTDOWN
	return nil
}

func (this *AppContext) IsRunning() bool {
	return this.RunState > 0
}

func (this *AppContext) Stop(waittime int) {
	this.shutdownTime = time.Now().UnixNano() + int64(waittime)*1e6
}
