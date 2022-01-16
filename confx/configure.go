package confx

const (
	DefaultConfigFileName = "config.yaml"
)

type Configure interface {
	GetString(key string, def string) (string, error)
	GetBool(key string, def bool) (bool, error)
	GetInt(key string, def int) (int, error)
	GetInt32(key string, def int32) (int32, error)
	GetInt64(key string, def int64) (int64, error)
	GetFloat32(key string, def float32) (float32, error)
	GetFloat64(key string, def float64) (float64, error)

	SetString(key string, val string) error
	SetBool(key string, val bool) error
	SetInt(key string, val int) error
	SetInt32(key string, val int32) error
	SetInt64(key string, val int64) error
	SetFloat32(key string, val float32) error
	SetFloat64(key string, val float64) error

	Exist(key string) bool
}
