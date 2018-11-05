package users

func init() {
	sql.GetListFilter = `
	SELECT * FROM users WHERE updated_at < NOW() - INTERVAL 1 HOUR 
	`
}

type query struct {
	GetListFilter string
}

// хранилище пользовательских запросов
var sql = new(query)

// Sql получение хранилища пользовательских запросов
func Sql() *query {
	return sql
}
