package log

import (
	"io"
	"log"
	"os"
	"sync"
)

// 定义我们自己的输出类型
var (
	errorLog = log.New(os.Stdout, "\033[31m[Error ]\033[0m", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[32m[Info ]\033[0m", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	//这里设置一个锁，用来设置不同的输出登记
	mutex sync.Mutex
)

// 定义自己的输出方法
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

// 下面设置输出的登记
const (
	ErrorLevel = iota
	InfoLevel
	Disable
)

// 设置当前需要输出的登记
func SetLevel(level int) {

	//先将所有的输出都定向到io.Discard
	for _, logger := range loggers {
		logger.SetOutput(io.Discard)
	}

	//通过level来决定使用那个等级的输出
	switch level {
	case ErrorLevel:
		errorLog.SetOutput(os.Stdout)
	case InfoLevel:
		infoLog.SetOutput(os.Stdout)
	case Disable:
	default:

	}
}
