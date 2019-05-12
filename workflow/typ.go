package workflow

// компонент
type Component struct {
	p             *pool
	cronTaskRun   []Task
	cronControlCH chan struct{}
}

// конфигурация
type Config struct {
	LimitCh   int `yaml:"Limitch"`   // Лимит канала задач
	LimitPool int `yaml:"Limitpool"` // Лимит пула (количество воркеров)
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
