package confx

import (
	"github.com/spf13/viper"
)

type FileConfigure struct {
	v    *viper.Viper
	next Configure
}

func NewFileConfigure(fname string, next Configure) (Configure, error) {
	v := viper.New()
	v.SetConfigFile(fname)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return &FileConfigure{v: v, next: next}, nil
}

func (fc *FileConfigure) Exist(key string) bool {
	if fc.v.IsSet(key) {
		return true
	}
	if fc.next == nil {
		return false
	}
	return fc.next.Exist(key)
}

func (fc *FileConfigure) GetString(key string, def string) (string, error) {
	val := fc.v.GetString(key)
	if val == "" && !fc.Exist(key) {
		return def, nil
	}
	return val, nil
}

func (fc *FileConfigure) GetBool(key string, def bool) (bool, error) {
	val := fc.v.GetBool(key)
	if val == false && !fc.Exist(key) {
		return def, nil
	}
	return val, nil
}

func (fc *FileConfigure) GetInt(key string, def int) (int, error) {
	val := fc.v.GetInt(key)
	if val == 0 && !fc.Exist(key) {
		return def, nil
	}
	return val, nil
}

func (fc *FileConfigure) GetInt32(key string, def int32) (int32, error) {
	val := fc.v.GetInt32(key)
	if val == 0 && !fc.Exist(key) {
		return def, nil
	}
	return val, nil
}

func (fc *FileConfigure) GetInt64(key string, def int64) (int64, error) {
	val := fc.v.GetInt64(key)
	if val == 0 && !fc.Exist(key) {
		return def, nil
	}
	return val, nil
}

func (fc *FileConfigure) GetFloat32(key string, def float32) (float32, error) {
	val := fc.v.GetFloat64(key)
	if val == 0 && !fc.Exist(key) {
		return def, nil
	}
	return float32(val), nil
}

func (fc *FileConfigure) GetFloat64(key string, def float64) (float64, error) {
	val := fc.v.GetFloat64(key)
	if val == 0 && !fc.Exist(key) {
		return def, nil
	}
	return val, nil
}

func (fc *FileConfigure) SetString(key string, val string) error {
	fc.v.Set(key, val)
	return nil
}

func (fc *FileConfigure) SetBool(key string, val bool) error {
	fc.v.Set(key, val)
	return nil
}

func (fc *FileConfigure) SetInt(key string, val int) error {
	fc.v.Set(key, val)
	return nil
}

func (fc *FileConfigure) SetInt32(key string, val int32) error {
	fc.v.Set(key, val)
	return nil
}

func (fc *FileConfigure) SetInt64(key string, val int64) error {
	fc.v.Set(key, val)
	return nil
}

func (fc *FileConfigure) SetFloat32(key string, val float32) error {
	fc.v.Set(key, val)
	return nil
}

func (fc *FileConfigure) SetFloat64(key string, val float64) error {
	fc.v.Set(key, val)
	return nil
}
