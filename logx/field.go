package logx

import "go.uber.org/zap"

type Field = zap.Field

var (
	Bool     = zap.Bool
	Bools    = zap.Bools
	Int      = zap.Int
	Ints     = zap.Ints
	Int8     = zap.Int8
	Int8s    = zap.Int8s
	Int16    = zap.Int16
	Int16s   = zap.Int16s
	Int32    = zap.Int32
	Int32s   = zap.Int32s
	Int64    = zap.Int64
	Int64s   = zap.Int64s
	Uint     = zap.Uint
	Uints    = zap.Uints
	Uint8    = zap.Uint8
	Uint8s   = zap.Uint8s
	Uint16   = zap.Uint16
	Uint16s  = zap.Uint16s
	Uint32   = zap.Uint32
	Uint32s  = zap.Uint32s
	Uint64   = zap.Uint64
	Uint64s  = zap.Uint64s
	String   = zap.String
	Strings  = zap.Strings
	Float32  = zap.Float32
	Float32s = zap.Float32s
	Float64  = zap.Float64
	Float64s = zap.Float64s

	Any        = zap.Any
	Time       = zap.Time
	NamedError = zap.NamedError
	Duration   = zap.Duration
)
