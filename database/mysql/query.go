// Библиотека работы с БД Mysql.
//
// Непосредственная работа с БД. Выполнение запросов.
// Работа с БД Mysql в парадигме ORM.
// Отложенное выполнение запросов к БД.
// Генерация идентификаторов объектов к БД.
package mysql

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
	"gopkg.in/kshamiev/sungora.v1/database/mysql/confmysql"
	"gopkg.in/kshamiev/sungora.v1/lg"
)

// Стек соединений с БД
var conn = make(map[string]*Db)

// служебная переменная для реализации блокировок
var mutexNew sync.Mutex

var primaryKeys = map[string]bool{
	"id": true,
	"Id": true,
	"ID": true,
}

// Структура по работе с БД
type Db struct {
	Connect mysql.Conn // Конннект к БД
	free    bool       // Статус блокировки (использования)
	time    time.Time  // Дата и время последнего использования коннекта
}

// Конструктор соединений с БД
//    - *db Объект по работе с БД для выполнения запросов
//    - error Ошибка операции
func NewDb() (obj *Db) {
	// блокировка
	mutexNew.Lock()
	defer mutexNew.Unlock()
	// поиск старого (свободного) коннекта
	//fmt.Println("SEARCH CONNECT")
	for len(conn) >= int(cfg.CntConn) && 0 < len(conn) {
		for i := range conn {
			if true == conn[i].free {
				conn[i].free = false
				conn[i].time = time.Now().In(timeLocation)
				return conn[i]
			}
		}
	}
	//fmt.Println("NEW CONNECT")
	obj = new(Db)
	if cfg.Type == `tcp` {
		obj.Connect = mysql.New("tcp", "", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), cfg.Login, cfg.Password, cfg.Name)
	} else {
		obj.Connect = mysql.New("unix", "", cfg.Socket, cfg.Login, cfg.Password, cfg.Name)
	}
	if err := obj.Connect.Connect(); err != nil {
		lg.Fatal(111, err)
		return nil
	}
	obj.time = time.Now().In(timeLocation)
	obj.free = false
	conn[obj.time.String()] = obj
	return obj
}

// Free Освобождение коннекта для других процессов (запросов)
func (self *Db) Free() {
	self.time = time.Now().In(timeLocation)
	self.free = true
}

// Загрузка БД в память
//    + object interface{} объект для заполнения табличными данными
//    Свойства object могут быть хешами либо срезами структур, одноименные загружаемым в них таблицам (пример Users []*Users)
//    Свойства структур соответсвуют полям или столбцам в таблицах.
//    Все загружаемые свойства должны быть публичны
//    Возможно использование тегов - алиасов в объявлении структур:
//    `db:"-"` - игнорировать и не загружать
//    `db:"cross"` - загружать как кросс таблицу без Id
//    `db:"Login" - загружать под указанным имененм в запросе к БД
//    - error Ошибка загрузки
func (self *Db) SelectData(object interface{}) (err error) {
	objValue := reflect.ValueOf(object)
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}
	num := objValue.NumField()
	for i := 0; i < num; i++ {
		field := objValue.Field(i)
		// определяем имя источника
		var source = objValue.Type().Field(i).Name
		var db = objValue.Type().Field(i).Tag.Get(`db`)
		if db == `-` {
			continue
		}
		if false == field.CanSet() {
			return lg.Error(181, source)
		}
		var sql string
		if db == `cross` {
			sql = "SELECT * FROM " + source
		} else {
			sql = "SELECT * FROM " + source + " ORDER BY Id ASC"
		}
		// читаем данные
		//if err = self.SelectSlice(field.Addr().Interface(), sql); err != nil {
		//	return
		//}

		switch field.Type().Kind() {
		case reflect.Slice:
			if err = self.SelectSlice(field.Addr().Interface(), sql); err != nil {
				return
			}
		case reflect.Map:
			if err = self.SelectMap(field.Addr().Interface(), sql); err != nil {
				return
			}
		default:
			return lg.Error(812, source)
		}

		//switch field.Type().Kind() {
		//case reflect.Slice:
		//	slc := reflect.MakeSlice(field.Type(), 0, 0)
		//	slc, err = self.loadArrayReflect(slc, sql)
		//	if err != nil {
		//		return
		//	}
		//	field.Set(slc)
		//case reflect.Map:
		//	mp := reflect.MakeMap(field.Type())
		//	mp, err = self.loadArrayReflect(mp, sql)
		//	if err != nil {
		//		return err
		//	}
		//	field.Set(mp)
		//default:
		//	return lg.Error(158, source)
		//}
	}
	return err
}

// Загрузка объектов (записей) из БД в хеш.
//    + ObjectMap interface{} Хеш для заполнения табличными данными (пример map[uint64]*Users)
//    Передается по ссылке.
//    Свойства структуры соответсвуют полям или столбцам в таблицах.
//    Все загружаемые свойства должны быть публичны
//    Возможно использование тегов - алиасов в объявлении структуры:
//    `db:"-"` - игнорировать и не загружать свойство
//    `db:"Login" - загружать под указанным имененм в запросе к БД
//    + sql string Запрос к БД ("Select * FROM TableName WHERE Comment LIKE ? OR Id = ?")
//    + params ...interface{} Обциональное количество параметров запроса любого скалярного типа ('%кошка%', 25)
//    - error Ошибка загрузки
func (self *Db) SelectMap(ObjectMap interface{}, sql string, params ...interface{}) (err error) {
	// рефлеския объекта
	objType := reflect.TypeOf(ObjectMap)
	if objType.Kind() != reflect.Ptr {
		return lg.Error(381, objType.String(), sql)
	}
	if objType.Elem().Kind() != reflect.Map {
		return lg.Error(813, objType.String(), sql)
	}
	var objValue = reflect.MakeMap(objType.Elem())
	var fieldMap = make(map[string]string)
	var obj = reflect.New(objType.Elem().Elem().Elem())
	obj = obj.Elem()
	num := obj.NumField()
	for i := 0; i < num; i++ {
		field := obj.Field(i)
		fieldName := obj.Type().Field(i).Name
		fieldTag := obj.Type().Field(i).Tag.Get(`db`)
		if fieldTag == `-` {
			continue
		}
		if false == field.IsValid() || false == field.CanSet() {
			return lg.Error(804, fieldName, sql)
		}
		if fieldTag == `` {
			fieldMap[fieldName] = fieldName
		} else {
			fieldMap[fieldTag] = fieldName
		}
	}
	// запрос
	if strings.LastIndex(sql, " ") == -1 {
		if sql, err = confmysql.GetQuery(sql); err != nil {
			return
		}
	}
	//sql, params = sqlParse(sql, params)
	var res mysql.Result
	var rows []mysql.Row
	var row mysql.Row
	var stm mysql.Stmt
	defer func() {
		if stm != nil {
			stm.Delete()
		}
	}()
	stm, err = self.Connect.Prepare(sql)
	if err != nil {
		return lg.Error(801, sql, err)
	}
	if len(params) > 0 {
		stm.Bind(params...)
	}
	rows, res, err = stm.Exec()
	if err != nil {
		return lg.Error(802, sql, err)
	}
	// пустой результат
	if 0 == len(rows) {
		return
	}
	// соответствие структуры типа и структуры запроса
	var fieldRes = make(map[string]string)
	for _, field := range res.Fields() {
		if _, ok := fieldMap[field.Name]; ok == false {
			lg.Warning(183, field.Name, sql)
			continue
		}
		fieldRes[field.Name] = fieldMap[field.Name]
	}
	// наполнение результата
	for _, row = range rows {
		objectId := uint64(0)
		var obj = reflect.New(objType.Elem().Elem().Elem())
		obj = obj.Elem()
		for fAlias, fName := range fieldRes {
			// logs.Dumper(fAlias, res.Map(fAlias))
			prop := obj.FieldByName(fName)
			// заносим полученные значения согласно типам свойств
			switch prop.Type().String() {
			case "bool":
				prop.SetBool(row.Bool(res.Map(fAlias)))
			case "int8", "int64", "int32", "int16":
				prop.SetInt(row.Int64(res.Map(fAlias)))
			case "uint8", "uint64", "uint32", "uint16":
				// prop.SetUint(row.Uint64(res.Map(fAlias)))
				val := row.Uint64(res.Map(fAlias))
				prop.SetUint(val)
				if _, ok := primaryKeys[fAlias]; ok {
					objectId = val
				} else if _, ok := primaryKeys[fName]; ok {
					objectId = val
				}
			case "float32", "float64":
				prop.SetFloat(row.Float(res.Map(fAlias)))
			case "string":
				prop.SetString(row.Str(res.Map(fAlias)))
			case "[]uint8":
				prop.SetBytes(row.Bin(res.Map(fAlias)))
			case "[]string":
				val := row.Str(res.Map(fAlias))
				slc := reflect.MakeSlice(prop.Type(), 0, 0)
				if "" != val {
					list := strings.Split(val, ",")
					for i := range list {
						slc = reflect.Append(slc, reflect.ValueOf(list[i]))
					}
				}
				prop.Set(slc)
			case "time.Time":
				prop.Set(reflect.ValueOf(row.Time(res.Map(fAlias), timeLocation)))
			case "time.Duration":
				prop.Set(reflect.ValueOf(row.Duration(res.Map(fAlias))))
			}
		}
		if objectId == 0 {
			objectId++
		}
		objValue.SetMapIndex(reflect.ValueOf(objectId), obj.Addr())
	}
	reflect.ValueOf(ObjectMap).Elem().Set(objValue)
	return
}

// Загрузка объектов (записей) из БД в срез.
//    + ObjectSlice interface{} Срез для заполнения табличными данными (пример []*Users)
//    Передается по ссылке.
//    Свойства структуры соответсвуют полям или столбцам в таблицах.
//    Все загружаемые свойства должны быть публичны
//    Возможно использование тегов - алиасов в объявлении структуры:
//    `db:"-"` - игнорировать и не загружать свойство
//    `db:"Login" - загружать под указанным имененм в запросе к БД
//    + sql string Запрос к БД ("Select * FROM TableName WHERE Comment LIKE ? OR Id = ?")
//    + params ...interface{} Обциональное количество параметров запроса любого скалярного типа ('%кошка%', 25)
//    - error Ошибка загрузки
func (self *Db) SelectSlice(ObjectList interface{}, sql string, params ...interface{}) (err error) {
	// рефлеския объекта
	objType := reflect.TypeOf(ObjectList)
	if objType.Kind() != reflect.Ptr {
		return lg.Fatal(157, objType.String(), sql)
	}
	if objType.Elem().Kind() != reflect.Slice {
		return lg.Fatal(169, objType.String(), sql)
	}
	var objValue = reflect.MakeSlice(objType.Elem(), 0, 0)
	var fieldMap = make(map[string]string)
	var obj = reflect.New(objType.Elem().Elem().Elem())
	obj = obj.Elem()
	num := obj.NumField()
	for i := 0; i < num; i++ {
		field := obj.Field(i)
		fieldName := obj.Type().Field(i).Name
		fieldTag := obj.Type().Field(i).Tag.Get(`db`)
		if fieldTag == `-` {
			continue
		}
		if false == field.IsValid() || false == field.CanSet() {
			return lg.Error(170, fieldName, sql)
		}
		if fieldTag == `` {
			fieldMap[fieldName] = fieldName
		} else {
			fieldMap[fieldTag] = fieldName
		}
	}
	// запрос
	if strings.LastIndex(sql, " ") == -1 {
		if sql, err = confmysql.GetQuery(sql); err != nil {
			return
		}
	}
	sql, params = sqlParse(sql, params)
	var res mysql.Result
	var rows []mysql.Row
	var row mysql.Row
	var stm mysql.Stmt
	defer func() {
		if stm != nil {
			stm.Delete()
		}
	}()
	stm, err = self.Connect.Prepare(sql)
	if err != nil {
		return lg.Error(110, sql, err)
	}
	if len(params) > 0 {
		stm.Bind(params...)
	}
	rows, res, err = stm.Exec()
	if err != nil {
		return lg.Error(109, sql, err)
	}
	// пустой результат
	if 0 == len(rows) {
		return
	}
	// соответствие структуры типа и структуры запроса
	var fieldRes = make(map[string]string)
	for _, field := range res.Fields() {
		if _, ok := fieldMap[field.Name]; ok == false {
			lg.Warning(171, field.Name, sql)
			continue
		}
		fieldRes[field.Name] = fieldMap[field.Name]
	}
	// наполнение результата
	for _, row = range rows {

		var obj = reflect.New(objType.Elem().Elem().Elem())
		obj = obj.Elem()
		for fAlias, fName := range fieldRes {
			// lg.Dumper(fAlias, res.Map(fAlias))
			prop := obj.FieldByName(fName)
			// заносим полученные значения согласно типам свойств
			switch prop.Type().String() {
			case "bool":
				prop.SetBool(row.Bool(res.Map(fAlias)))
			case "int8", "int64", "int32", "int16":
				prop.SetInt(row.Int64(res.Map(fAlias)))
			case "uint8", "uint64", "uint32", "uint16":
				prop.SetUint(row.Uint64(res.Map(fAlias)))
			case "float32", "float64":
				prop.SetFloat(row.Float(res.Map(fAlias)))
			case "string":
				prop.SetString(row.Str(res.Map(fAlias)))
			case "[]uint8":
				prop.SetBytes(row.Bin(res.Map(fAlias)))
			case "[]string":
				val := row.Str(res.Map(fAlias))
				slc := reflect.MakeSlice(prop.Type(), 0, 0)
				if "" != val {
					list := strings.Split(val, ",")
					for i := range list {
						slc = reflect.Append(slc, reflect.ValueOf(list[i]))
					}
				}
				prop.Set(slc)
			case "time.Time":
				prop.Set(reflect.ValueOf(row.Time(res.Map(fAlias), timeLocation)))
			case "time.Duration":
				prop.Set(reflect.ValueOf(row.Duration(res.Map(fAlias))))
			}
		}
		objValue = reflect.Append(objValue, obj.Addr())
	}
	reflect.ValueOf(ObjectList).Elem().Set(objValue)
	return

}

// Загрузка одного объекта (записи) из БД.
//    + Object interface{} Объект структуры переданный по ссылке (пример var Object = new(Users)).
//    Свойства структуры соответсвуют полям или столбцам в таблице.
//    Все загружаемые свойства должны быть публичны
//    Возможно использование тегов - алиасов в объявлении структуры:
//    `db:"-"` - игнорировать и не загружать свойство
//    `db:"Login" - загружать под указанным имененм в запросе к БД
//    + sql string Запрос к БД ("Select * FROM TableName WHERE Comment LIKE ? OR Id = ?")
//    + params ...interface{} Обциональное количество параметров запроса любого скалярного типа ('%кошка%', 25)
//    - error Ошибка загрузки
func (self *Db) Select(Object interface{}, sql string, params ...interface{}) (err error) {
	// рефлеския объекта
	var objValue = reflect.ValueOf(Object)
	if objValue.Kind() != reflect.Ptr {
		return lg.Fatal(176, objValue.Type().String(), sql)
	}
	if objValue.IsNil() == true {
		return lg.Fatal(180, objValue.Type().String(), sql)
	}
	objValue = objValue.Elem()
	var fieldMap = make(map[string]string)
	num := objValue.NumField()
	for i := 0; i < num; i++ {
		field := objValue.Field(i)
		fieldName := objValue.Type().Field(i).Name
		fieldTag := objValue.Type().Field(i).Tag.Get(`db`)
		if fieldTag == `-` {
			continue
		}
		if false == field.IsValid() || false == field.CanSet() {
			return lg.Error(177, fieldName, sql)
		}
		if fieldTag == `` {
			fieldMap[fieldName] = fieldName
		} else {
			fieldMap[fieldTag] = fieldName
		}
	}
	// запрос
	if strings.LastIndex(sql, " ") == -1 {
		if sql, err = confmysql.GetQuery(sql); err != nil {
			return
		}
	}
	sql, params = sqlParse(sql, params)
	var res mysql.Result
	var rows []mysql.Row
	var row mysql.Row
	var stm mysql.Stmt
	defer func() {
		if stm != nil {
			stm.Delete()
		}
	}()
	stm, err = self.Connect.Prepare(sql)
	if err != nil {
		return lg.Error(110, sql, err)
	}
	if len(params) > 0 {
		stm.Bind(params...)
	}
	rows, res, err = stm.Exec()
	if err != nil {
		return lg.Error(109, sql, err)
	}
	// пустой результат
	if 0 == len(rows) {
		return lg.Info(156, sql)
	}
	// соответствие структуры типа и структуры запроса
	var fieldRes = make(map[string]string)
	for _, field := range res.Fields() {
		if _, ok := fieldMap[field.Name]; ok == false {
			lg.Warning(178, field.Name, sql)
			continue
		}
		fieldRes[field.Name] = fieldMap[field.Name]
	}
	// наполнение результата
	row = rows[0]
	for fAlias, fName := range fieldRes {
		// lg.Dumper(fAlias, res.Map(fAlias))
		prop := objValue.FieldByName(fName)
		// заносим полученные значения согласно типам свойств
		switch prop.Type().String() {
		case "bool":
			prop.SetBool(row.Bool(res.Map(fAlias)))
		case "int8", "int64", "int32", "int16":
			prop.SetInt(row.Int64(res.Map(fAlias)))
		case "uint8", "uint64", "uint32", "uint16":
			prop.SetUint(row.Uint64(res.Map(fAlias)))
		case "float32", "float64":
			prop.SetFloat(row.Float(res.Map(fAlias)))
		case "string":
			prop.SetString(row.Str(res.Map(fAlias)))
		case "[]uint8":
			prop.SetBytes(row.Bin(res.Map(fAlias)))
		case "[]string":
			val := row.Str(res.Map(fAlias))
			slc := reflect.MakeSlice(prop.Type(), 0, 0)
			if "" != val {
				list := strings.Split(val, ",")
				for i := range list {
					slc = reflect.Append(slc, reflect.ValueOf(list[i]))
				}
			}
			prop.Set(slc)
		case "time.Time":
			prop.Set(reflect.ValueOf(row.Time(res.Map(fAlias), timeLocation)))
		case "time.Duration":
			prop.Set(reflect.ValueOf(row.Duration(res.Map(fAlias))))
		}
	}
	return
}

// Сохранение объекта (записи) в БД.
//    + Object interface{} Объект структуры переданный по ссылке (пример var Object = new(Users)).
//    Свойства структуры соответсвуют полям или столбцам в таблице.
//    Все сохраняемые свойства должны быть публичны
//    Возможно использование тегов - алиасов в объявлении структуры:
//    `db:"-"` - игнорировать и не сохранять свойство
//    `db:"Login" - сохранять под указанным именем
//    + source string Имя таблицы в которую сохраняется объект (запись)
//    + properties ...map[string]string Хеш с обциональным количеством свойств которые нужно сохранить
//    Если хеш пустой сохраняются все доступные для сохранения свойства
//    - uint64 Идентификатор (сохраненной новой записи в БД) объекта
//    - error Ошибка сохранения
func (self *Db) Insert(Object interface{}, source string, properties ...map[string]string) (insertId uint64, err error) {
	// инициализация
	var stm mysql.Stmt
	defer func() {
		if stm != nil {
			stm.Delete()
		}
	}()
	var res mysql.Result
	var sql string
	var paramList []interface{}
	// запрос
	objValue := reflect.ValueOf(Object)
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}
	num := objValue.NumField()
	for i := 0; i < num; i++ {
		prop := objValue.Field(i)
		propName := objValue.Type().Field(i).Name
		propTag := objValue.Type().Field(i).Tag.Get(`db`)
		// пропускаем защищенные от изменения и не указанные свойства
		if `-` == propTag || false == prop.CanSet() {
			continue
		}
		//	пропускаем свойства которые не надо сохранять
		if 0 < len(properties) {
			if _, ok := properties[0][propName]; ok == false {
				continue
			}
		}
		// определение алиаса для БД
		if propTag != `` {
			propName = propTag
		}
		switch prop.Type().String() {
		case "bool":
			val := prop.Interface().(bool)
			if false == val {
				paramList = append(paramList, 0)
			} else {
				paramList = append(paramList, 1)
			}
			sql += "`" + propName + "` = ?, "
		case "int64":
			val := prop.Interface().(int64)
			if 0 == val {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "int32":
			val := prop.Interface().(int32)
			if 0 == val {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "int16":
			val := prop.Interface().(int16)
			if 0 == val {
				paramList = append(paramList, 0)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "int8":
			val := prop.Interface().(int8)
			if 0 == val {
				paramList = append(paramList, 0)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "uint64":
			val := prop.Interface().(uint64)
			if 0 == val {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "uint32":
			val := prop.Interface().(uint32)
			if 0 == val {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "uint16":
			val := prop.Interface().(uint16)
			if 0 == val {
				paramList = append(paramList, 0)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "uint8":
			val := prop.Interface().(uint8)
			if 0 == val {
				paramList = append(paramList, 0)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "float32":
			val := prop.Interface().(float32)
			if 0 == val {
				paramList = append(paramList, 0)
			} else {
				paramList = append(paramList, val)
			}
			sql += "`" + propName + "` = ?, "
		case "float64":
			val := prop.Interface().(float64)
			if 0 == val {
				paramList = append(paramList, 0)
			} else {
				paramList = append(paramList, val)
			}
			sql += "`" + propName + "` = ?, "
		case "string":
			val := prop.Interface().(string)
			if "" == val {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, val)
			}
			sql += "`" + propName + "` = ?, "
		case "[]string":
			val := strings.Join(prop.Interface().([]string), ",")
			if "" == val {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, val)
			}
			sql += "`" + propName + "` = ?, "
		case "[]uint8":
			val := prop.Interface().([]uint8)
			if 0 == len(val) {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, val)
			}
			sql += "`" + propName + "` = ?, "
		case "time.Time":
			val := prop.Interface().(time.Time)
			if val.IsZero() == true {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, val)
			}
			sql += "`" + propName + "` = ?, "
		case "time.Duration":
			val := prop.Interface().(time.Duration)
			if val <= 0 {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, val)
			}
			sql += "`" + propName + "` = ?, "
		}
	}
	sql = "INSERT " + source + " SET " + sql[:len(sql)-2]
	// запрос и параметры
	if stm, err = self.Connect.Prepare(sql); err != nil {
		return insertId, lg.Error(110, sql, err)
	}
	stm.Bind(paramList...)
	// выполнение запроса
	if _, res, err = stm.Exec(); err != nil {
		return insertId, lg.Error(109, sql, err)
	}
	insertId = res.InsertId()
	return
}

// Изменение объекта (записи) в БД.
//    + Object interface{} Объект структуры переданный по ссылке (пример var Object = new(Users)).
//    Свойства структуры соответсвуют полям или столбцам в таблице.
//    Все изменяемые свойства должны быть публичны
//    Возможно использование тегов - алиасов в объявлении структуры:
//    `db:"-"` - игнорировать и не изменять свойство
//    `db:"Login" - сохранять под указанным именем
//    + source string Имя таблицы в которой изменяется объект (запись)
//    + key string Ключевое поле по которому изменяется объект (запись)
//    + properties ...map[string]string Хеш с обциональным количеством свойств которые нужно изменить
//    Если хеш пустой изменяются все доступные для изменения свойства
//    - uint64 Количество измененых записей (всегдя либо 0 либо 1)
//    - error Ошибка изменения
func (self *Db) Update(Object interface{}, source, key string, properties ...map[string]string) (affectedRow uint64, err error) {
	// инициализация
	var stm mysql.Stmt
	defer func() {
		if stm != nil {
			stm.Delete()
		}
	}()
	var res mysql.Result
	var sql string
	var paramList []interface{}
	// запрос
	objValue := reflect.ValueOf(Object)
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}
	field := objValue.FieldByName(key)
	if field.IsValid() == false {
		return affectedRow, lg.Error(113, key, source)
	}
	num := objValue.NumField()
	for i := 0; i < num; i++ {
		prop := objValue.Field(i)
		propName := objValue.Type().Field(i).Name
		propTag := objValue.Type().Field(i).Tag.Get(`db`)
		// пропускаем защищенные от изменения и не указанные свойства
		if `-` == propTag || false == prop.CanSet() {
			continue
		}
		//	пропускаем свойства которые не надо сохранять
		if 0 < len(properties) {
			if _, ok := properties[0][propName]; ok == false {
				continue
			}
		}
		// определение алиаса для БД
		if propTag != `` {
			propName = propTag
		}
		switch prop.Type().String() {
		case "bool":
			val := prop.Interface().(bool)
			if false == val {
				paramList = append(paramList, 0)
			} else {
				paramList = append(paramList, 1)
			}
			sql += "`" + propName + "` = ?, "
		case "int64":
			val := prop.Interface().(int64)
			if 0 == val {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "int32":
			val := prop.Interface().(int32)
			if 0 == val {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "int16":
			val := prop.Interface().(int16)
			if 0 == val {
				paramList = append(paramList, 0)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "int8":
			val := prop.Interface().(int8)
			if 0 == val {
				paramList = append(paramList, 0)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "uint64":
			val := prop.Interface().(uint64)
			if 0 == val {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "uint32":
			val := prop.Interface().(uint32)
			if 0 == val {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "uint16":
			val := prop.Interface().(uint16)
			if 0 == val {
				paramList = append(paramList, 0)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "uint8":
			val := prop.Interface().(uint8)
			if 0 == val {
				paramList = append(paramList, 0)
			} else {
				paramList = append(paramList, prop.Interface())
			}
			sql += "`" + propName + "` = ?, "
		case "float32":
			val := prop.Interface().(float32)
			if 0 == val {
				paramList = append(paramList, 0)
			} else {
				paramList = append(paramList, val)
			}
			sql += "`" + propName + "` = ?, "
		case "float64":
			val := prop.Interface().(float64)
			if 0 == val {
				paramList = append(paramList, 0)
			} else {
				paramList = append(paramList, val)
			}
			sql += "`" + propName + "` = ?, "
		case "string":
			val := prop.Interface().(string)
			if "" == val {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, val)
			}
			sql += "`" + propName + "` = ?, "
		case "[]string":
			val := strings.Join(prop.Interface().([]string), ",")
			if "" == val {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, val)
			}
			sql += "`" + propName + "` = ?, "
		case "[]uint8":
			val := prop.Interface().([]uint8)
			if 0 == len(val) {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, val)
			}
			sql += "`" + propName + "` = ?, "
		case "time.Time":
			val := prop.Interface().(time.Time)
			if val.IsZero() == true {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, val)
			}
			sql += "`" + propName + "` = ?, "
		case "time.Duration":
			val := prop.Interface().(time.Duration)
			if val <= 0 {
				paramList = append(paramList, nil)
			} else {
				paramList = append(paramList, val)
			}
			sql += "`" + propName + "` = ?, "
		}
	}
	sql = "UPDATE " + source + " SET " + sql[:len(sql)-2] + " WHERE `" + key + "` = ?"
	paramList = append(paramList, field.Interface())
	// запрос и параметры
	if stm, err = self.Connect.Prepare(sql); err != nil {
		return affectedRow, lg.Error(110, sql, err)
	}
	stm.Bind(paramList...)
	// выполнение запроса
	if _, res, err = stm.Exec(); err != nil {
		return affectedRow, lg.Error(109, sql, err)
	}
	affectedRow = res.AffectedRows()
	return
}

// Пользовательский запрос на обновление в БД (как правило из файлов).
//    + data []byte Дамп запросов на изменения в БД (К примеру бекап БД)
//    + []string Сообщения о результате обновления
//    + error Ошибка обновления
func (self *Db) QueryByte(data []byte) (messages []string, err error) {
	var result mysql.Result
	_, result, err = self.Connect.Query(string(data))
	if err != nil {
		return nil, lg.Error(114)
	} else {
		if result != nil {
			for result.MoreResults() == true {
				if message := result.Message(); message != "" {
					messages = append(messages, message)
				}
				result, err = result.NextResult()
				if err != nil {
					err = nil
					break
				}
			}
		}
	}
	return
}

// Пользовательский запрос на обновление в БД.
//    + sql string Запрос на обновление БД (UPDATE Users SET Email = ?, Status = ?, Age = ? WHERE Id = ?)
//    + params ...interface{} Обциональное количество параметров запроса любого скалярного типа ("funtik@yandex.ru", true, 25, 4536)
//    - error Ошибка обновления
func (self *Db) Query(sql string, params ...interface{}) (err error) {
	if strings.LastIndex(sql, " ") == -1 {
		if sql, err = confmysql.GetQuery(sql); err != nil {
			return
		}
	}
	// инициализация
	var stm mysql.Stmt
	defer func() {
		if stm != nil {
			stm.Delete()
		}
	}()
	// запрос
	sql, params = sqlParse(sql, params)
	stm, err = self.Connect.Prepare(sql)
	if err != nil {
		return lg.Error(110, sql, err)
	}
	stm.Bind(params...)
	if _, _, err = stm.Exec(); err != nil {
		return lg.Error(109, sql, err)
	}
	return
}

// Выполнение функций
//    + Object interface{} Срез структур принимающий данные результата выполнения функции.
//    + nameCall string Имя функции
//    + params ...interface{} Обциональное количество параметров функции любого скалярного типа ("funtik@yandex.ru", true, 25)
func (self *Db) CallFunc(Object interface{}, nameCall string, params ...interface{}) (err error) {
	return
}

// Выполнение хранимых процедур
//    + Object interface{} Срез структур принимающий данные результата выполнения хранимой процедуры.
//    + nameCall string Имя хранимой процедуры
//    + params ...interface{} Обциональное количество параметров процедуры любого скалярного типа ("funtik@yandex.ru", true, 25)
func (self *Db) CallExec(Object interface{}, nameCall string, params ...interface{}) (err error) {
	return
}

// Разбор входных параметров
// Корректировка параметров для запросов с параметрами IN
// Использовать не рекомендуется.
// Рекомендуется подготавливать параметры запросов до их вызова.
//    + sql string Запрос к БД
//    + param []interface{} Срез параметров запроса
//    - string Запрос к БД со вставленными параметрами из срезов
//    - []interface{} Срез параметров запроса, за вычетом параметрами из срезов
func sqlParse(sql string, param []interface{}) (string, []interface{}) {
	// инициализируем данные для встапвки и находим их местоположения в запросе
	var str []string
	var strPos []int
	var params []interface{}
	for i := range param {
		var typ = reflect.TypeOf(param[i])
		switch typ.String() {
		case "[]string":
			d := "'" + strings.Join(param[i].([]string), "', '") + "'"
			str = append(str, d)
			strPos = append(strPos, i)
		case "[]int64":
			var s []string
			for _, elm := range param[i].([]int64) {
				s = append(s, fmt.Sprintf("%d", elm))
			}
			d := strings.Join(s, ", ")
			str = append(str, d)
			strPos = append(strPos, i)
		case "[]int32":
			var s []string
			for _, elm := range param[i].([]int32) {
				s = append(s, fmt.Sprintf("%d", elm))
			}
			d := strings.Join(s, ", ")
			str = append(str, d)
			strPos = append(strPos, i)
		case "[]int16":
			var s []string
			for _, elm := range param[i].([]int16) {
				s = append(s, fmt.Sprintf("%d", elm))
			}
			d := strings.Join(s, ", ")
			str = append(str, d)
			strPos = append(strPos, i)
		case "[]int8":
			var s []string
			for _, elm := range param[i].([]int8) {
				s = append(s, fmt.Sprintf("%d", elm))
			}
			d := strings.Join(s, ", ")
			str = append(str, d)
			strPos = append(strPos, i)
		case "[]uint64":
			var s []string
			for _, elm := range param[i].([]uint64) {
				s = append(s, fmt.Sprintf("%d", elm))
			}
			d := strings.Join(s, ", ")
			str = append(str, d)
			strPos = append(strPos, i)
		case "[]uint32":
			var s []string
			for _, elm := range param[i].([]uint32) {
				s = append(s, fmt.Sprintf("%d", elm))
			}
			d := strings.Join(s, ", ")
			str = append(str, d)
			strPos = append(strPos, i)
		case "[]uint16":
			var s []string
			for _, elm := range param[i].([]uint16) {
				s = append(s, fmt.Sprintf("%d", elm))
			}
			d := strings.Join(s, ", ")
			str = append(str, d)
			strPos = append(strPos, i)
		case "[]uint8":
			return sql, param
			var s []string
			for _, elm := range param[i].([]uint8) {
				s = append(s, fmt.Sprintf("%d", elm))
			}
			d := strings.Join(s, ", ")
			str = append(str, d)
			strPos = append(strPos, i)
		case "[]float64":
			var s []string
			for _, elm := range param[i].([]float64) {
				s = append(s, fmt.Sprintf("%f", elm))
			}
			d := strings.Join(s, ", ")
			str = append(str, d)
			strPos = append(strPos, i)
		case "[]float32":
			var s []string
			for _, elm := range param[i].([]float32) {
				s = append(s, fmt.Sprintf("%f", elm))
			}
			d := strings.Join(s, ", ")
			str = append(str, d)
			strPos = append(strPos, i)
		default:
			params = append(params, param[i])
		}
	}
	// вставляем данные срезов
	var data = strings.Split(sql, "?")
	var dataNew []string
	var pos int
	for i := range strPos {
		dataNew = append(dataNew, strings.Join(data[pos:strPos[i]+1], "?"))
		dataNew = append(dataNew, str[i])
		pos = strPos[i] + 1
	}
	dataNew = append(dataNew, strings.Join(data[pos:], "?"))
	return strings.Join(dataNew, ""), params
}

/*
     0  []*mysql.Field (len = 23) {
    36  .  ...
    37  .  3: *mysql.Field {
    38  .  .  Catalog: "def"
    39  .  .  Db: "Test"
    40  .  .  Table: "Test"
    41  .  .  OrgTable: "Test"
    42  .  .  Name: "LastName"
    43  .  .  OrgName: "LastName"
    44  .  .  DispLen: 64
    45  .  .  Flags: 16392
    46  .  .  Type: 8
    47  .  .  Scale: 0
    48  .  }
    49  .  ...
   277  }
*/
