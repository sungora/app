package lg

// конфигурация
type Config struct {
	ServiceName string `yaml:"ServiceName"`
	Info        bool   `yaml:"Info"`
	Notice      bool   `yaml:"Notice"`
	Warning     bool   `yaml:"Warning"`
	Error       bool   `yaml:"Error"`
	Critical    bool   `yaml:"Critical"`
	Fatal       bool   `yaml:"Fatal"`
	Traces      bool   `yaml:"Traces"`
	OutStd      bool   `yaml:"OutStd"`
	OutFile     string `yaml:"OutFile"` // лог файл
	OutHttp     string `yaml:"OutHttp"` // url куда отправляются логи
}

type msg struct {
	Datetime   string
	Level      string
	LineNumber int
	Action     string
	Service    string
	Login      string
	Message    string
	Traces     []trace
}

type trace struct {
	FuncName   string // Название функции
	FileName   string // Имя исходного файла
	LineNumber int    // Номер строки внутри функции
}
