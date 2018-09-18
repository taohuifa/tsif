package app

import (
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
	RunState int // 运行状态

	shutdownTime    uint64 // 关闭时间
	updateIinterval uint64 // 间隔时间

	InitFunc    func(params ...interface{}) bool                               // 初始化函数
	DestroyFunc func()                                                         // 销毁函数
	UpdateFunc  func(count int, dt int, nowtime time.Time, prevtime time.Time) // 更新函数
}

func (this *AppContext) Start(params ...interface{}) (bool, string) {
	var result bool
	// init
	this.RunState = APPSTATE_STARTING
	Log.Info("app init", false)
	if this.InitFunc != nil {
		result = this.InitFunc(params)
		if !result {
			return result, "init error."
		}
	}
	Log.Info("app start", false)

	// update
	this.RunState = APPSTATE_RUN
	prevtime := time.Now()
	count := 0
	for this.IsRunning() {
		count++
		nowtime := time.Now()
		dt := (nowtime.Nanosecond() - prevtime.Nanosecond())
		if this.UpdateFunc != nil {
			this.UpdateFunc(count, dt, nowtime, prevtime)
		}
	}

	// destory
	this.RunState = APPSTATE_TO_SHUTDOWN
	Log.Info("app destorying", false)
	if this.DestroyFunc != nil {
		this.DestroyFunc()
	}
	Log.Info("app destory", false)
	this.RunState = APPSTATE_SHUTDOWN
	return true, "success"
}

func (this *AppContext) IsRunning() bool {
	return this.RunState > 0
}

func (this *AppContext) Stop(waittime int) {

}
