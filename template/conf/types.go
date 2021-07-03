package conf

type DatabaseConfig struct {
	Enable       bool
	UserPassword string `required:"true"`
	DB           string `required:"true"`
	Write        struct {
		HostPort string `required:"true"`
	}
	Read struct {
		HostPort string `required:"true"`
	}
	Conn struct {
		MaxLifeTime int `required:"true"`
		MaxIdle     int `required:"true"`
		MaxOpen     int `required:"true"`
	}
}

type ServerConfig struct {
	Listen    string `required:"true"`
	Env       string `required:"true"`
	AppName   string `required:"true"`
	ProjectID string `default:"zh-fs"`
}