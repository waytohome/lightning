package ginx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// BindAndValidate 绑定并校验属性，要求传入指针类型对象
func BindAndValidate(c *gin.Context, ptr interface{}) error {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		return New(CodeParamValidateFailed, fmt.Sprintf("require ptr kind object, %s kind not support", v.Kind().String()))
	}
	if err := parse(c, ptr); err != nil {
		return New(CodeParamValidateFailed, err.Error())
	}
	if err := binding.Validator.ValidateStruct(ptr); err != nil {
		return New(CodeParamValidateFailed, err.Error())
	}
	return nil
}

func parse(c *gin.Context, ptr interface{}) error {
	m := make(map[string][]string)
	for _, p := range c.Params {
		m[p.Key] = append(m[p.Key], p.Value)
	}
	values := c.Request.URL.Query()
	for key, items := range values {
		m[key] = append(m[key], items...)
	}
	if len(m) > 0 {
		// parse path param
		if err := setVal(ptr, m, "uri"); err != nil {
			return err
		}
		// parse query param
		if err := setVal(ptr, m, "query"); err != nil {
			return err
		}
	}
	defer func() {
		_ = c.Request.Body.Close()
	}()
	// parse json
	all, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	if len(all) <= 0 {
		return nil
	}
	if err := json.Unmarshal(all, ptr); err != nil {
		return err
	}
	return nil
}

func setVal(ptr interface{}, m map[string][]string, tag string) error {
	elem := reflect.ValueOf(ptr).Elem()
	for i := 0; i < elem.NumField(); i++ {
		if !elem.Field(i).CanSet() {
			continue
		}

		field := elem.Type().Field(i)
		tagVal := field.Tag.Get(tag)
		if tagVal == "" {
			continue
		}

		var setVal interface{}
		if field.Type.String() == "int" {
			parseInt, err := strconv.ParseInt(m[tagVal][0], 10, 64)
			if err != nil {
				return err
			}
			setVal = int(parseInt)
		} else if field.Type.String() == "int8" {
			parseInt, err := strconv.ParseInt(m[tagVal][0], 10, 8)
			if err != nil {
				return err
			}
			setVal = int8(parseInt)
		} else if field.Type.String() == "int16" {
			parseInt, err := strconv.ParseInt(m[tagVal][0], 10, 16)
			if err != nil {
				return err
			}
			setVal = int16(parseInt)
		} else if field.Type.String() == "int32" {
			parseInt, err := strconv.ParseInt(m[tagVal][0], 10, 32)
			if err != nil {
				return err
			}
			setVal = int32(parseInt)
		} else if field.Type.String() == "int64" {
			parseInt, err := strconv.ParseInt(m[tagVal][0], 10, 64)
			if err != nil {
				return err
			}
			setVal = parseInt
		} else if field.Type.String() == "float32" {
			parseInt, err := strconv.ParseFloat(m[tagVal][0], 32)
			if err != nil {
				return err
			}
			setVal = float32(parseInt)
		} else if field.Type.String() == "float64" {
			parseInt, err := strconv.ParseFloat(m[tagVal][0], 64)
			if err != nil {
				return err
			}
			setVal = parseInt
		} else if field.Type.String() == "bool" {
			parseInt, err := strconv.ParseBool(m[tagVal][0])
			if err != nil {
				return err
			}
			setVal = parseInt
		} else if field.Type.String() == "uint" {
			parseInt, err := strconv.ParseUint(m[tagVal][0], 10, 64)
			if err != nil {
				return err
			}
			setVal = uint(parseInt)
		} else if field.Type.String() == "uint8" {
			parseInt, err := strconv.ParseUint(m[tagVal][0], 10, 8)
			if err != nil {
				return err
			}
			setVal = uint8(parseInt)
		} else if field.Type.String() == "uint16" {
			parseInt, err := strconv.ParseUint(m[tagVal][0], 10, 16)
			if err != nil {
				return err
			}
			setVal = uint16(parseInt)
		} else if field.Type.String() == "uint32" {
			parseInt, err := strconv.ParseUint(m[tagVal][0], 10, 32)
			if err != nil {
				return err
			}
			setVal = uint32(parseInt)
		} else if field.Type.String() == "uint64" {
			parseInt, err := strconv.ParseUint(m[tagVal][0], 10, 64)
			if err != nil {
				return err
			}
			setVal = parseInt
		} else if field.Type.String() == "string" {
			setVal = m[tagVal][0]
		} else {
			return errors.New(fmt.Sprintf("%s type not match", field.Type.String()))
		}
		elem.FieldByName(field.Name).Set(reflect.ValueOf(setVal))
	}
	return nil
}
