package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	APPSTATE_TO_SHUTDOWN   = -1 // 开始停服
	APPSTATE_SHUTDOWN      = 0  // 停止
	APPSTATE_RUN           = 1  // 运行中
	APPSTATE_STARTING      = 2  // 启动中
	APPSTATE_WAIT_TO_CLOSE = 3  // 准备关闭

	APPCMD_SEG         = "\r\n" // cmd分隔符
	APPCMD_FILE_SUFFIX = ".buf" // cmdfile 后缀

	APPCMD_START = "start" // 启动指令
	APPCMD_STOP  = "stop"  // 结束指令
)

// App上下文
type Context struct {
	Name           string // appname
	UpdateInterval uint   // 间隔时间

	state        int   // 运行状态
	shutdownTime int64 // 关闭时间

	// app函数
	InitFunc    func(context *Context, params ...interface{}) bool                               // 初始化函数
	DestroyFunc func(context *Context)                                                           // 销毁函数
	UpdateFunc  func(context *Context, count int, dt int, nowtime time.Time, prevtime time.Time) // 更新函数
	CmdFunc     func(context *Context, cmd string, args []string)                                // 指令函数
}

// getCmd
func getCmd(pid string) *[]string {
	cmdFile := pid + APPCMD_FILE_SUFFIX
	// 读取文件
	cmdBody, err_rf := ioutil.ReadFile(cmdFile)
	if err_rf != nil {
		return nil // 读取失败, 不存在文件
	}
	// 读取后删除文件
	os.Remove(cmdFile)
	// 存在文件
	// fmt.Println("cmd body:", cmdFile, string(cmdBody))

	// 解析换行符
	cmdStr := string(cmdBody)
	cmdArray := strings.Split(cmdStr, APPCMD_SEG)
	// for _, cmd := range cmdArray {
	// 	if cmd == "" {
	// 		continue
	// 	}
	// 	fmt.Printf("cmd: \"%s\".\n", cmd)
	// }
	return &cmdArray
}

// sendCmd
func sendCmd(pid string, cmdStr string) error {
	cmdFile := pid + APPCMD_FILE_SUFFIX
	// pid文件存在, 转为发送指令
	file, err_o := os.OpenFile(cmdFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err_o != nil {
		return errors.New(fmt.Sprintf("open cmdfile fail! cmdfile=%s.\n", cmdFile))
	}
	defer file.Close()
	// write cmdfile
	_, err_w := file.Write([]byte(cmdStr + APPCMD_SEG))
	if err_w != nil {
		return errors.New(fmt.Sprintf("send cmd fail! file=%s cmd=\"%s\" err=%s\n", cmdFile, cmdStr, err_w.Error()))
	}
	// fmt.Printf("send cmd success! file=%s cmd=\"%s\".", cmdFile, cmd)
	return nil
}

// 处理Cmd
func disposeCmd(this *Context, cmds *[]string) {
	// 处理指令
	for _, cmd := range *cmds {
		if cmd == "" {
			continue // 空指令
		}
		// 拆分指令
		cmdArgs := strings.Split(cmd, " ")
		fmt.Printf("cmd=%d %s.\n", len(cmdArgs), cmdArgs)
		if cmdArgs == nil || len(cmdArgs) <= 0 {
			continue // 拆分失败
		}
		// 特殊指令处理
		if cmdArgs[0] == APPCMD_STOP {
			this.state = APPSTATE_TO_SHUTDOWN
			break
		}
		// 指令处理
		if this.CmdFunc != nil {
			this.CmdFunc(this, cmdArgs[0], cmdArgs[1:])
		}
	}
}

// 运行app
func run(this *Context, pid string, params ...interface{}) error {
	// fmt.Printf("run %s\n", pid)
	var result bool
	// init
	this.state = APPSTATE_STARTING
	// Log.Info("app init: " + fmt.Sprint(params))
	if this.InitFunc != nil {
		result = this.InitFunc(this, params)
		if !result {
			return errors.New("init error")
		}
	}
	// Log.Info("app start")

	// update
	this.state = APPSTATE_RUN
	prevtime := time.Now()
	count := 0
	for this.IsRunning() {
		UpdateInterval := time.Duration(math.Max(float64(this.UpdateInterval), 1000)) // 至少1s
		time.Sleep(UpdateInterval * time.Millisecond)
		//wait shutdown
		nowtime := time.Now()
		if this.shutdownTime > 0 {
			dtime := this.shutdownTime - nowtime.UnixNano()
			// Log.Infof("app wait stop", dtime)
			if dtime <= 0 {
				this.state = APPSTATE_TO_SHUTDOWN
				break
			}
		}

		// get cmd
		cmds := getCmd(pid)
		if cmds != nil {
			disposeCmd(this, cmds)
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
	this.state = APPSTATE_TO_SHUTDOWN
	// Log.Info("app destorying")
	if this.DestroyFunc != nil {
		this.DestroyFunc(this)
	}
	// Log.Info("app destory")
	this.state = APPSTATE_SHUTDOWN
	return nil
}

// 运行app
func (this *Context) Run(pidFile string, cmd string, params ...interface{}) error {
	// check params
	if pidFile == "" || cmd == "" {
		return errors.New(fmt.Sprintf("param fail! pidFile=\"%s\", cmd=\"%s\".\n", pidFile, cmd))
	}
	// 判断pid文件是否存在
	_, err_s := os.Stat(pidFile)
	if err_s == nil || os.IsExist(err_s) {
		// 存在pid文件, 认为程序运行中
		if cmd == APPCMD_START {
			return errors.New("proc is running...")
		}
		// 读取pid文件
		pid, err_rf := ioutil.ReadFile(pidFile)
		if err_rf != nil {
			return errors.New(fmt.Sprintf("read pid fail! pidfile=%s.\n", pidFile))
		}
		// send cmd
		return sendCmd(string(pid), cmd)
	}

	// 没有pid文件, 没有运行, 检测是否是启动指令
	if cmd != APPCMD_START {
		return errors.New(fmt.Sprintf("proc not running, can't action \"%s\".", cmd))
	}
	// 生成pid文件
	pid := strconv.Itoa(os.Getpid())
	ioutil.WriteFile(pidFile, []byte(pid), 0644)

	// 运行程序
	err := run(this, pid, params)

	// 清除文件
	cmdFile := pid + APPCMD_FILE_SUFFIX
	os.Remove(cmdFile) // 删除cmd文件
	os.Remove(pidFile) // 删除pid文件
	return err
}

func (this *Context) IsRunning() bool {
	return this.state > 0
}

func (this *Context) State() int {
	return this.state
}
func (this *Context) Stop(waittime int) {
	this.shutdownTime = time.Now().UnixNano() + int64(waittime)*1e6
}
