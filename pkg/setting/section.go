package setting

type DatabaseSettingS struct {
	DBType       string
	UserName     string
	Password     string
	Host         string
	Port         string
	DBName       string
	Charset      string
	SSLMode      string
	MaxIdleConns int
	MaxOpenConns int
}

type AppSettingS struct {
	LogSavePath string
	LogFileName string
	LogFileExt  string
}

type RPCSettingS struct {
	URL          string
	Mode         string
	WebSocketURL string
}

type SystemCodeHashS struct {
	Dao   string
	Sudts []string
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	return nil
}
