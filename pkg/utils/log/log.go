package log

type Config struct {
	// rotate options
	LogFile    string `mapstructure:"log-file"`
	MaxSize    int	   `mapstructure:"max-size"`
	MaxBackups int		`mapstructure:"max-backups"`
	MaxAge     int		`mapstructure:"max-age"`
	Compress   bool		 `mapstructure:"compress"`

	// logger options
	LogLevel        string  `mapstructure:"log-level"`
	JsonEncode      bool	`mapstructure:"json-encode"`
	StacktraceLevel string	 `mapstructure:"stacktrace-level"`
}
