package lib

import (
	"fmt"
	"time"
)

/**
*	时间格式化,由于golang库的时间格式化难以满足我的需求,因此搞出个简单的替换方式
*	layout：
*		YYYY/YY/Y 表示年,会格式化为4/2位数
*		MM/M 表示月,会格式化为两位数,MM个位数会带0
*		DD/D 表示日
*		HH/H 24进制时间格式
*		mm/m 分
*		ss/s 表示秒
 */
func TimeFormat(t time.Time, layout string) string {
	fmtStr := ""
	var fmtArr []interface{}
	for i := 0; i < len(layout); {
		switch layout[i] {
		case 'Y':
			if layout[i:i+4] == "YYYY" {
				fmtStr = fmtStr + "%04d"
				fmtArr = append(fmtArr, t.Year())
				i = i + 4
				continue
			}
			if layout[i:i+2] == "YY" {
				fmtStr = fmtStr + "%02d"
				fmtArr = append(fmtArr, t.Year()%100)
				i = i + 2
				continue
			}
			fmtStr = fmtStr + "%d"
			fmtArr = append(fmtArr, t.Year())
			i = i + 1
		case 'M':
			if layout[i:i+2] == "MM" {
				fmtStr = fmtStr + "%02d"
				fmtArr = append(fmtArr, t.Month())
				i = i + 2
				continue
			}
			fmtStr = fmtStr + "%d"
			fmtArr = append(fmtArr, t.Month())
			i = i + 1
		case 'D':
			if layout[i:i+2] == "DD" {
				fmtStr = fmtStr + "%02d"
				fmtArr = append(fmtArr, t.Day())
				i = i + 2
				continue
			}
			fmtStr = fmtStr + "%d"
			fmtArr = append(fmtArr, t.Day())
			i = i + 1
		case 'H':
			if layout[i:i+2] == "HH" {
				fmtStr = fmtStr + "%02d"
				fmtArr = append(fmtArr, t.Hour())
				i = i + 2
				continue
			}
			fmtStr = fmtStr + "%d"
			fmtArr = append(fmtArr, t.Hour())
			i = i + 1
		case 'm':
			if layout[i:i+2] == "mm" {
				fmtStr = fmtStr + "%02d"
				fmtArr = append(fmtArr, t.Minute())
				i = i + 2
				continue
			}
			fmtStr = fmtStr + "%d"
			fmtArr = append(fmtArr, t.Minute())
			i = i + 1
		case 's':
			if layout[i:i+2] == "ss" {
				fmtStr = fmtStr + "%02d"
				fmtArr = append(fmtArr, t.Second())
				i = i + 2
				continue
			}
			fmtStr = fmtStr + "%d"
			fmtArr = append(fmtArr, t.Second())
			i = i + 1
		default:
			fmtStr = fmtStr + layout[i:i+1]
			i = i + 1
		}
	}
	return fmt.Sprintf(fmtStr, fmtArr...)
}
