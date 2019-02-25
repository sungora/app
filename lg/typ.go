package lg

// главная конфигурация
type configFile struct {
	ServiceName string
	Lg          configTyp
}

// конфигурация поджгружаемая из файла конфигурации
type configTyp struct {
	Info     bool
	Notice   bool
	Warning  bool
	Error    bool
	Critical bool
	Fatal    bool
	Traces   bool
	OutStd   bool
	OutFile  bool
	OutHttp  string // url куда отправляются логи
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
