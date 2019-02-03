package workflow

// главная конфигурация
type configMain struct {
	Workflow configFile
}

// конфигурация поджгружаемая из файла конфигурации
type configFile struct {
	LimitCh   int // Лимит канала задач
	LimitPool int // Лимит пула (количество воркеров)
}

type manager struct {
	Name      string
	IsExecute bool
	Minute    string
	Hour      string
	Day       string
	Month     string
	Week      string
}
