package conf

import (
	"bufio"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
)

// 解析配置
func ConfigParse(path string) (*map[string]string, error) {
	// 打开文件指定目录，返回一个文件f和错误信息
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 创建数据
	confMap := make(map[string]string)
	prefix := ""

	// 创建一个输出流向该文件的缓冲流*Reader
	r := bufio.NewReader(f)
	for {
		// 读取，返回[]byte 单行切片给b
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		//去除单行属性两端的空格
		s := strings.TrimSpace(string(b))
		//fmt.Println(s)

		// 过滤空和注释
		if len(s) <= 0 || s[0] == '#' || s[0] == ';' {
			continue
		}

		// 判断前置变量[v]
		regCom, _ := regexp.Compile(`\[[a-z]+\]`)
		// fmt.Println(s, regCom.MatchString(s))
		if regCom.MatchString(s) {
			// 解析正常, 按照前缀处理
			reg := regexp.MustCompile(`[a-z]+`)
			regStrs := reg.FindAllString(s, -1)
			if len(regStrs) <= 0 {
				return nil, errors.New("str parse fail! " + s)
			}
			prefix = regStrs[0]
			// fmt.Println(s, regStrs, prefix)
			continue
		}

		// 判断等号=在该行的位置
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		// 取得等号左边的key值，判断是否为空
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		// 取得等号右边的value值，判断是否为空
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		// 这样就成功吧配置文件里的属性key=value对，成功载入到内存中c对象里
		if prefix != "" {
			key = prefix + "." + key
		}
		confMap[key] = value
	}
	return &confMap, nil
}
