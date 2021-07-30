package utils

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/UncleBig/goCache"
	"github.com/bitly/go-simplejson"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	Indent    = strings.Repeat(" ", 4)
	C         *goCache.Cache
	RedBold   = color.New(color.FgRed, color.Bold)
	BlueBold  = color.New(color.FgBlue, color.Bold)
	WhiteBold = color.New(color.FgWhite, color.Bold)
	Cyan      = color.New(color.FgCyan)
	Green     = color.New(color.FgGreen)
	Yellow    = color.New(color.FgYellow)
)

// init 公共初始化
func init() {
	C = goCache.New(10*time.Minute, 30*time.Second)
}

// Handler 错误处理
func Handler(err error) {
	if err != nil {
		_, _ = RedBold.Printf("\n%sError: %s\n\n", Indent, err)
		os.Exit(101)
	}
}

// ParseArgs 解析 cli 参数
func ParseArgs(menu map[string]string) (args map[string]interface{}) {

	p := os.Args
	args = make(map[string]interface{})

	for i := 1; i < len(p); i++ {
		if strings.HasPrefix(p[i], "--") {
			key := string([]byte(p[i])[2:])
			if strings.Contains(p[i], "=") {
				item := strings.Split(key, "=")
				args[item[0]] = item[1]
			} else {
				if len(p) < i+2 {
					Handler(fmt.Errorf("os.Args index out of range"))
				}
				args[key] = p[i+1]
				i += 1
			}
		} else if strings.HasPrefix(p[i], "-") {
			key := string([]byte(p[i])[1:])
			args[key] = true
		}
	}

	if _, ok := menu["h"]; !ok {
		menu["-h"] = "List argv help."
	}

	if _, ok := args["h"]; ok {
		_, _ = BlueBold.Printf("\n%sUsage: fortress [args...].\n\n", Indent)
		for key, value := range menu {
			_, _ = Green.Printf("%16s", key)
			_, _ = WhiteBold.Printf("%s%s\n", Indent, value)
		}
		fmt.Println()
		os.Exit(0)
	}

	return
}

// Menu 打印菜单
func Menu(list []string, multiple bool, def interface{}) (number []int, err error) {

	var input string
	if def != nil && def != "" {
		input = def.(string)
	} else {
		_, _ = BlueBold.Printf("\n%sAll of operation list.\n\n", Indent)

		for index, title := range list {
			index := fmt.Sprintf("%02d", index+1)
			fmt.Printf("\t%s  %s\n", Cyan.SprintFunc()("["+index+"]"), Green.SprintFunc()(title))
		}

		var mul string
		if multiple {
			mul = "(multiple split by comma)"
		}

		_, _ = BlueBold.Printf("\n%sSelect index %s>: ", Indent, mul)

		_, _ = fmt.Scan(&input)
	}

	inputArr := strings.Split(input, ",")
	if !multiple && len(inputArr) > 1 {
		err = fmt.Errorf("only one index allowed")
		return
	}

	max := len(list)
	for _, item := range inputArr {
		var num int
		num, err = strconv.Atoi(item)
		if err != nil || num < 1 || num > max {
			err = fmt.Errorf("index must be number than between 1 and %d, given -> `%s`", max, item)
			return
		}
		number = append(number, num)
	}

	return
}

// Md5 哈希值
func Md5(str string) string {
	hash := md5.New()
	_, _ = io.WriteString(hash, str)

	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Cache 缓存处理
func Cache(key string, fn func(key string) interface{}) interface{} {

	keyHash := Md5(key)
	val, found := C.Get(keyHash)

	if !found {
		val = fn(key)
		C.Set(keyHash, val, goCache.DefaultExpiration)
	}

	return val
}

// ReadJson 读取 JSON 格式的配置文件
func ReadJson(file string) (json *simplejson.Json, err error) {

	jsonData := Cache(file, func(key string) interface{} {

		// read file
		var data []byte
		data, err = ioutil.ReadFile(key)

		// convert to json
		json, err = simplejson.NewJson(data)

		return json
	})
	
	json = jsonData.(*simplejson.Json)
	return
}

// ExecShell 同步执行命令并返回标准输出
func ExecShell(command string) (string, error) {

	cmd := exec.Command("bash", "-c", command)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	return out.String(), err
}

// ExecShellFine 同步执行命令并返回标准输出(逐行实时返回)
func ExecShellFine(command string) (err error) {

	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	_ = cmd.Start()

	reader := bufio.NewReader(stdout)

	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Printf("%s%s", Indent, line)
	}

	err = cmd.Wait()

	return
}

// TplHandler 从 map 中处理模板字符串中的变量
func TplHandlerFromMap(tpl string, args map[string]interface{}, prefix string) string {

	for key, item := range args {
		from := fmt.Sprintf("{%s.%s}", prefix, key)
		tpl = strings.Replace(tpl, from, fmt.Sprint(item), -1)
	}

	return tpl
}

// TplHandler 从 struct 处理模板字符串中的变量
func TplHandlerFromStructure(tpl string, args interface{}, prefix string) string {

	typed := reflect.TypeOf(args)
	value := reflect.ValueOf(args)

	for k := 0; k < typed.NumField(); k++ {
		key := strings.ToLower(typed.Field(k).Name)
		item := value.Field(k).Interface()

		from := fmt.Sprintf("{%s.%s}", prefix, key)
		tpl = strings.Replace(tpl, from, fmt.Sprint(item), -1)
	}

	return tpl
}
