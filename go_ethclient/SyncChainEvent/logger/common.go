package logger

// LogConf 配置信息
type LogConf struct {
	ServiceName string `toml:"service_name" mapstructure:"service_name" json:"service_name"`
	Mode        string `toml:"mode" json:"mode"`
	Path        string `toml:"path" json:"path"`
	Level       string `toml:"level" json:"level"`
	Compress    bool   `toml:"compress" json:"compress"`
	KeepDays    int    `toml:"keep_days" mapstructure:"keep_days" json:"keep_days"`
}
