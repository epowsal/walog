// walog by suirosu exgaya epowsal wlb iwlb@outlook.com exgaya@gmail.com 20241230;

package walog

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var logf *os.File

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.LstdFlags)
	var loge error
	logf, loge = os.OpenFile(os.Args[0]+".log", os.O_CREATE|os.O_WRONLY, 0666)
	if loge == nil {
		logf.Seek(0, os.SEEK_END)
		log.SetOutput(logf)
	} else {
		Pa("log file", os.Args[0]+".log", "open error", loge)
	}
}

func ChangeLogPath(path string) error {
	var loge error
	if strings.LastIndexAny(path, "/\\") != -1 {
		os.MkdirAll(path[:strings.LastIndexAny(path, "/\\")], 0666)
	}
	if logf != nil {
		logf.Close()
	}
	logf, loge = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if loge == nil {
		logf.Seek(0, os.SEEK_END)
		log.SetOutput(logf)
	} else {
		Pa("log file", path, "open error", loge)
	}
	return nil
}

var LogStackDeep int = 8
var LogStackLineEnd string = "<<" //\n or <<;
var LogWrapSize int = 512

func Er(as ...any) error {
	return errors.New(string(append(sprint(as...), gostack()...)))
}

func Pa(as ...any) error {
	panic(string(append(sprint(as...), gostack()...)))
}

func Log(as ...any) {
	log.Println(string(append(sprint(as...), gostack()...)), time.Now().Format(time.RFC3339Nano))
}

func Pr(as ...any) {
	fmt.Println(string(append(sprint(as...), gostack()...)), time.Now().Format(time.RFC3339Nano))
}

func Ck(rl string, as ...any) {
	srl := sprint(as...)
	if len(srl) > 0 {
		srl = srl[:len(srl)-1]
	}
	if rl != string(srl) {
		Pa(append([]any{rl}, as...)...)
	}
}

// list to any list
func C(a ...any) []any {
	return a
}

func gostack() []byte {
	var stackbs []byte
	tbn := 0
	for i := 0; i <= LogStackDeep; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if ok {
			sn := file[strings.LastIndex(file[:strings.LastIndex(file, "/")], "/")+1:]
			if !(strings.HasSuffix(sn, "/testing.go") || strings.HasSuffix(sn, "/walog.go")) && strings.HasSuffix(sn, ".s") == false {
				fc := runtime.FuncForPC(pc)
				stackbs = append(stackbs, []byte(fmt.Sprintf(logStrRepeat('\t', tbn)+"[%d %s]%s"+LogStackLineEnd, line, fc.Name(), sn))...)
				tbn += 1
			}
		}
	}
	return stackbs
}

func sprint(as ...any) []byte {
	var err []byte
	prel := 0
	for _, a := range as {
		switch a.(type) {
		case string:
			err = append(err, []byte(fmt.Sprintf("%s^", a.(string)))...)
		case int:
			err = append(err, []byte(fmt.Sprintf("%d^", a.(int)))...)
		case uint:
			err = append(err, []byte(fmt.Sprintf("%d^", a.(uint)))...)
		case int8:
			err = append(err, []byte(fmt.Sprintf("%d^", a.(int8)))...)
		case uint8:
			err = append(err, []byte(fmt.Sprintf("%d^", a.(uint8)))...)
		case int16:
			err = append(err, []byte(fmt.Sprintf("%d^", a.(int16)))...)
		case uint16:
			err = append(err, []byte(fmt.Sprintf("%d^", a.(uint16)))...)
		case int32:
			err = append(err, []byte(fmt.Sprintf("%d^", a.(int32)))...)
		case uint32:
			err = append(err, []byte(fmt.Sprintf("%d^", a.(uint32)))...)
		case int64:
			err = append(err, []byte(fmt.Sprintf("%d^", a.(int64)))...)
		case uint64:
			err = append(err, []byte(fmt.Sprintf("%d^", a.(uint64)))...)
		case float32:
			err = append(err, []byte(f32s(a.(float32))+"^")...)
		case float64:
			err = append(err, []byte(f64s(a.(float64))+"^")...)
		case error:
			err = append(err, []byte("error{"+a.(error).Error()+"}^")...)
		case bool:
			err = append(err, []byte(fmt.Sprintf("%t^", a.(bool)))...)
		default:
			err = append(err, []byte(fmt.Sprintf("%#v^", a))...)
		}
		if len(err)-prel >= LogWrapSize {
			err = append(err, []byte("\n\t")...)
		}
		prel = len(err)
	}
	if len(err) > 0 {
		if err[len(err)-1] == '^' {
			err[len(err)-1] = '\t'
		} else if len(err) >= 2 {
			err = err[:len(err)-1]
			err[len(err)-1] = '\t'
		}
	}
	return err
}

func logStrRepeat(b byte, n int) string {
	if LogStackLineEnd == "<<" {
		return ""
	}
	bs := []byte{}
	for i := 0; i < n; i += 1 {
		bs = append(bs, b)
	}
	return string(bs)
}

func f32s(v float32) string {
	return strconv.FormatFloat(float64(v), 'g', 6, 32)
}
func f64s(v float64) string {
	return strconv.FormatFloat(v, 'g', 15, 64)
}
