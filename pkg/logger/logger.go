package logger

import "log"

//当存在错误时，记录日志
func LogError(err error) {
	if err != nil {
		log.Println(err)
	}
}
