package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	REDIS_TIMEOUT = 10 * time.Second // 连接超时
)

// redis服务
type RedisClient struct {
	Url      string // redis地址, 例:127.0.0.1:6379
	Password string // redis 密钥
	Prefix   string // redis操作key前缀
}

// 统一的连接处理
func connect(this *RedisClient) (redis.Conn, error) {
	return redis.Dial("tcp", this.Url, redis.DialPassword(this.Password),
		redis.DialConnectTimeout(REDIS_TIMEOUT),
		redis.DialReadTimeout(REDIS_TIMEOUT), redis.DialWriteTimeout(REDIS_TIMEOUT))
}

// 调用方法
func (this *RedisClient) do(cmd string, args ...interface{}) (reply interface{}, err error) {
	// 访问连接
	c, err_c := connect(this)
	if err_c != nil {
		return nil, err_c
	}
	defer c.Close()
	// 执行do
	return c.Do(cmd, args...)
}

// 设置数据
// 额外参数: livetime(int) 存活时间(s), 0为永久
func (this *RedisClient) Set(key string, value interface{}, extparams ...interface{}) error {
	// redis set args
	vstr := fmt.Sprint(value)
	args := []interface{}{this.Prefix + key, vstr}
	// 识别参数
	pnum := len(extparams)
	if pnum > 0 {
		// 存活时间
		livetime := extparams[0].(int)
		if livetime > 0 {
			args = append(args, "EX", fmt.Sprint(livetime))
		}
	}

	// 执行写入
	//fmt.Printf("redis set %s=%v\n", key, value)
	_, err_s := this.do("SET", args...)
	return err_s
}

// 提取数据
func (this *RedisClient) get(key string) (interface{}, error) {
	// 执行读取
	// Log.Debugf("redis get %s=%v", key, reply)
	return this.do("GET", this.Prefix+key)
}

// 提取数据
func (this *RedisClient) Get(key string) (string, error) {
	return redis.String(this.get(key))
}

// 设置数据
func (this *RedisClient) SetObj(key string, value interface{}, extparams ...interface{}) error {
	b, err_j := json.Marshal(value)
	if err_j != nil {
		return err_j
	}
	// set to redis
	return this.Set(key, string(b), extparams...)
}

// 提取数据
func (this *RedisClient) GetObj(key string, v interface{}) error {
	// 提取数据
	reply, err_g := this.get(key)
	if err_g != nil {
		return err_g
	}
	// 提取byte
	b, err_b := redis.Bytes(reply, err_g)
	if err_b != nil {
		return err_b
	}

	// 解析
	err_j := json.Unmarshal(b, v)
	if err_j != nil {
		return err_j
	}
	return nil
}

// 列表头添加
func (this *RedisClient) Lpull(key string, value interface{}) error {
	_, err_g := this.do("lpush", this.Prefix+key, fmt.Sprint(value))
	return err_g
}

// 列表尾添加
func (this *RedisClient) Rpull(key string, value interface{}) error {
	_, err_g := this.do("rpush", this.Prefix+key, fmt.Sprint(value))
	return err_g
}

// 列表长度
func (this *RedisClient) Llen(key string) (int, error) {
	r, err_g := redis.Int(this.do("llen", this.Prefix+key))
	return r, err_g
}

// 解析出string列表
func listString(reply interface{}, err error) ([]string, error) {
	if err != nil {
		return nil, err
	} else if reply == nil {
		return nil, errors.New("no reply")
	}
	// 获取数量
	list := reply.([]interface{})
	llen := len(list)
	// 遍历提取列表
	rlist := make([]string, 0, llen)
	for k, v := range list {
		istr, err_s := redis.String(v, nil)
		if err_s != nil {
			return nil, errors.New(fmt.Sprintf("list[%d] %s", k, err_s.Error()))
		}
		rlist = append(rlist, istr)
	}
	// Log.Info("r", r)
	return rlist, nil
}

// 局部获取
func (this *RedisClient) Lrange(key string, start int, end int) ([]string, error) {
	return listString(this.do("lrange", this.Prefix+key, start, end))
}

// 设置列表中某个值
func (this *RedisClient) Lset(key string, index int, value interface{}) error {
	_, err_g := this.do("lset", this.Prefix+key, index, fmt.Sprint(value))
	return err_g
}

// 返回并删除名称为key的list中的首元素
func (this *RedisClient) Lpop(key string) (string, error) {
	r, err_g := redis.String(this.do("lpop", this.Prefix+key))
	return r, err_g
}

// 返回并删除名称为key的list中的尾元素
func (this *RedisClient) Rpop(key string) (string, error) {
	r, err_g := redis.String(this.do("rpop", this.Prefix+key))
	return r, err_g
}

// 返回并删除名称为srckey的list的尾元素，并将该元素添加到名称为dstkey的list的头部
func (this *RedisClient) Rpoplpush(skey string, dkey string) (string, error) {
	r, err_g := redis.String(this.do("rpoplpush", this.Prefix+skey, this.Prefix+dkey))
	return r, err_g
}
