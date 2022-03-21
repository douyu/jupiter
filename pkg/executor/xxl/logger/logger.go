package logger

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/douyu/jupiter/pkg/executor/xxl/constants"
)

type LogIDKey string

var (
	__DefaultLogPath = "/home/www/logs/applogs/job/jobhandler/"
	DefaultLogIDKey  = LogIDKey("xxl-log-id")
)

func Info(logId int64, log string) {
	log += "\n"
	if err := writeLog(logId, log); err != nil {
		fmt.Println("xxl job Write Log Error:", err)
	}
}

func InfoWithContext(ctx context.Context, log string) {
	logId, ok := ctx.Value(DefaultLogIDKey).(int64)
	if !ok || logId == 0 {
		fmt.Println(log)
		return
	}
	log += "\n"
	if err := writeLog(logId, log); err != nil {
		fmt.Println("xxl job Write Log Error:", err)
	}
}

func GetLogPath(nowTime time.Time) string {
	return __DefaultLogPath + nowTime.Format(constants.DateFormat)
}

func InitLogPath(logPath string) error {
	__DefaultLogPath = logPath
	_, err := os.Stat(GetLogPath(time.Now()))
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(__DefaultLogPath, os.ModePerm)
	}
	return err
}

func writeLog(logId int64, log string) error {
	logPath := GetLogPath(time.Now())
	logFile := fmt.Sprintf("%d.log", logId)
	if strings.Trim(logFile, " ") != "" {
		fileFullPath := logPath + "/" + logFile
		file, err := os.OpenFile(fileFullPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil && os.IsNotExist(err) {
			err = os.MkdirAll(logPath, os.ModePerm)
			if err == nil {
				file, err = os.OpenFile(fileFullPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					return err
				}
			}
		}

		if file != nil {
			defer file.Close()
			res, err := file.Write([]byte(log))
			if err != nil {
				return err
			}
			if res <= 0 {
				return errors.New("write log failed")
			}
		}
	}
	return nil
}

func ReadLog(logDateTim, logId int64, fromLineNum int32) (line int32, content string) {
	nowtime := time.Unix(logDateTim/1000, 0)
	fileName := fmt.Sprintf("%s/%d.log", GetLogPath(nowtime), logId)
	file, err := os.Open(fileName)
	totalLines := int32(1)
	var buffer bytes.Buffer
	if err == nil {
		defer file.Close()
		rd := bufio.NewReader(file)
		for {
			line, err := rd.ReadString('\n')
			if totalLines >= fromLineNum {
				buffer.WriteString(line)
			}
			totalLines++
			if err != nil || io.EOF == err {
				break
			}
		}
	}
	return totalLines, buffer.String()
}
