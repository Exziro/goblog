package logger

import "log"

//查错函数
func LogError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
