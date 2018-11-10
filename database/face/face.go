package face

type ArFace interface {
	Select(property string) ArFace
	From(from string) ArFace
	Where(where string) ArFace
	Group(group string) ArFace
	Having(having string) ArFace
	Order(order string) ArFace
	Limit(start, step int) ArFace
	Get() (query string)
}

// Интерфейс к БД
type DbFace interface {
	Select(typ interface{}, sql string, params ...interface{}) (err error)
	SelectMap(typMap interface{}, sql string, params ...interface{}) (err error)
	SelectSlice(typSlice interface{}, sql string, params ...interface{}) (err error)
	SelectData(typType interface{}) (err error)
	Insert(typ interface{}, source string, properties ...map[string]string) (insertId uint64, err error)
	Update(typ interface{}, source, key string, properties ...map[string]string) (affectedRow uint64, err error)
	Query(sql string, params ...interface{}) (err error)
	QueryByte(data []byte) (messages []string, err error)
	CallExec(typ interface{}, nameCall string, params ...interface{}) (err error)
	CallFunc(typ interface{}, nameCall string, params ...interface{}) (err error)
	Free()
}
