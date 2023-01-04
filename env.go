package unixodbc

import "C"
import (
	"github.com/ninthclowd/unixodbc/internal/api"
)

type PoolOption int

const (
	PoolOptionDefault PoolOption = iota
	PoolOptionOff
	PoolOptionPerDriver
	PoolOptionPerEnv
)

func (p PoolOption) toOption() uint64 {
	switch p {
	case PoolOptionOff:
		return api.SQL_CP_OFF
	case PoolOptionPerDriver:
		return api.SQL_CP_ONE_PER_DRIVER
	case PoolOptionPerEnv:
		return api.SQL_CP_ONE_PER_HENV
	default:
		return api.SQL_CP_DEFAULT
	}
}

type PoolMatch int

const (
	PoolMatchDefault PoolMatch = iota
	PoolMatchStrict
	PoolMatchRelaxed
)

func (m PoolMatch) toOption() uint64 {
	switch m {
	case PoolMatchStrict:
		return api.SQL_CP_STRICT_MATCH
	case PoolMatchRelaxed:
		return api.SQL_CP_RELAXED_MATCH
	default:
		return api.SQL_CP_MATCH_DEFAULT
	}
}

type EnvConfig struct {
	PoolOption PoolOption
	PoolMatch  PoolMatch
}

var _env *api.Environment

func env(config *EnvConfig) (*api.Environment, error) {
	if _env == nil {
		hnd, err := api.NewEnvironment()
		if err != nil {
			return nil, err
		}

		envOptions := map[api.SQLINTEGER]uint64{
			api.SQL_ATTR_ODBC_VERSION: api.SQL_OV_ODBC3_80,
		}
		if config != nil {
			envOptions[api.SQL_ATTR_CONNECTION_POOLING] = config.PoolOption.toOption()
			envOptions[api.SQL_ATTR_CP_MATCH] = config.PoolMatch.toOption()
		}

		for key, value := range envOptions {
			if err = hnd.SetAttrInt(key, int32(value)); err != nil {
				return nil, err
			}
		}
		_env = hnd
	}
	return _env, nil
}
