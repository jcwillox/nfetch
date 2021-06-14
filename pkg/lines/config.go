package lines

import "github.com/spf13/cast"

type LineConfig struct {
	config map[interface{}]interface{}
}

func GetLineConfig(entry interface{}) (string, LineConfig) {
	if s, ok := entry.(string); ok {
		return s, LineConfig{}
	}
	if config, ok := entry.(map[interface{}]interface{}); ok {
		for k, v := range config {
			if config, ok := v.(map[interface{}]interface{}); ok {
				return k.(string), LineConfig{config}
			}
			return k.(string), LineConfig{}
		}
	}
	return "unknown", LineConfig{}
}

func (c LineConfig) Get(key string) interface{} {
	return c.config[key]
}

func (c LineConfig) GetString(key string) string {
	return cast.ToString(c.Get(key))
}

func (c LineConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(c.Get(key))
}
