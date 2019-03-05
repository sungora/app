package workflow

// конфигурация поджгружаемая из файла конфигурации
type configFile struct {
	Workflow configTyp
}

// конфигурация workerPool
type configTyp struct {
	LimitCh   int // Лимит канала задач
	LimitPool int // Лимит пула (количество воркеров)
}

// задача
type Task interface {
	Manager() Manager // режим выполнения задачи
	Execute()         // тело задачи
}

// управление режимом работы фоновой задачи
type Manager struct {
	Name      string
	IsExecute bool
	Minute    string
	Hour      string
	Day       string
	Month     string
	Week      string
}
