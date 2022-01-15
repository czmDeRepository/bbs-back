package common

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type DateTime struct {
	time.Time
}

const FormatDateTime = "2006-01-02 15:04:05"

func Now() *DateTime {
	dateTime := new(DateTime)
	dateTime.Time = time.Now()
	return dateTime
}

// MarshalJSON 序列化为JSON
func (t DateTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("\"\""), nil
	}
	stamp := fmt.Sprintf("\"%s\"", t.Format(FormatDateTime))
	return []byte(stamp), nil
}

// UnmarshalJSON 反序列化为JSON
func (t *DateTime) UnmarshalJSON(data []byte) error {
	var err error
	t.Time, err = time.Parse(FormatDateTime, string(data))
	return err
}

// String 重写String方法
func (t *DateTime) String() string {
	data, _ := json.Marshal(t)
	return string(data)
}

// FieldType 数据类型
func (t *DateTime) FieldType() int {
	return orm.TypeDateTimeField

}

// SetRaw 读取数据库值
func (t *DateTime) SetRaw(value interface{}) error {

	switch value.(type) {
	case time.Time:
		t.Time = value.(time.Time)
	case []uint8:
		t.UnmarshalJSON(value.([]uint8))
	}
	return nil
}

// RawValue 写入数据库
func (t *DateTime) RawValue() interface{} {
	str := t.Format(FormatDateTime)
	if str == "0001-01-01 00:00:00" {
		return nil
	}
	return str
}
