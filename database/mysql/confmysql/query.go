// mysql запросы приложения к БД
package confmysql

import (
	"errors"
)

// Поиск и получение запроса по его индексу
//    + index string Индекс запроса (`ModuleName/FileName/NumberQuery` -> `base/users/0`)
//    - string Найденный запрос или пустая строка в случае ошибки
//    - error Ошибка поиска и получение запроса
func GetQuery(index string) (sql string, err error) {
	if sql, ok := queries[index]; ok {
		return sql, nil
	}
	return "", errors.New("Error retrieving request by index" + index)
}

func SetQuery(query map[string]string) {
	for k, v := range query {
		queries[k] = v
	}
}

var queries = map[string]string{
	`test`: "SELECT * FROM `Test`;",
	`test/one`: "SELECT * FROM Test WHERE ID = ?;",
}
