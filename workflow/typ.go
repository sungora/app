package workflow

// главная конфигурация
type configFile struct {
	Workflow configTyp
}

// конфигурация поджгружаемая из файла конфигурации
type configTyp struct {
	LimitCh   int // Лимит канала задач
	LimitPool int // Лимит пула (количество воркеров)
}

// Task Задача
type Task interface {
	Manager() Manager
	Execute()
}

type Manager struct {
	Name      string
	IsExecute bool
	Minute    string
	Hour      string
	Day       string
	Month     string
	Week      string
}
