package main

import (
	"zzzzzzzzz/models"

	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sungora/app/core"
	"github.com/sungora/app/lg"
)

func main() {
	var err error
	core.DB, err = gorm.Open("mysql", "root:root@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		lg.Dumper(core.DB, err.Error())
	}
	defer core.DB.Close()
	//
	u := models.NewUser(0)

	core.DB.AutoMigrate(u)

	u.InvoiceIDSet(546)
	u.NamSet("Вася Пупкин")
	u.SampleJsonSet(`{"ServiceID": 24,"ClientID": 24,"InvoiceID": 24,"BillingCode": "NL","LocationCode": "NL"}`)
	u.Save()

	u.InvoiceID = nil
	u.Nam = nil
	u.Save()

	// u.Save()
	// if err != nil {
	// 	lg.Dumper(u, err.Error())
	// } else {
	// 	lg.Dumper(u)
	// }

	//
	// var namesq []string
	// core.DB.Model(&users.Users{}).Pluck("names", &namesq)
	//
	// lg.Dumper(namesq)

	// err = s.Delete()
	// lg.Dumper(s, err)

	// db.CreateTable(&user{})

	// create
	// u := new(User)
	// u.Name = "Вася 1"
	// u.Age = 2365
	// u.Birthday = time.Now().Add(time.Hour * 5)
	// u.Description = "Описание 1"
	// lg.Dumper(u)
	// err = db.Create(u).Error
	// //db.Table("users_copy").Create(u)
	// lg.Dumper(u, err)

	// QUERY

	// ORM
	// var sample []*User
	// err = db.Find(&sample, "name = ? AND ID >= ?", "Вася", 3).Error
	// lg.Dumper(sample, err)

	// var sample []*User
	// var count int32
	// err = db.Select("id, name").Table("users_copy").Where("ID > ?", 1056).Order("ID ASC").Limit(2).Find(&sample).Error
	// err = db.Select("id, name").Table("users_copy").Where("ID > ?", 10).Order("ID ASC").Find(&sample).Count(&count).Error
	// err = db.Select("id, name").Where("ID > ?", 10).Order("ID ASC").Find(&sample).Count(&count).Error
	// lg.Dumper(sample, count, err)

	// получение только количества записей
	// var count1 int32
	// err = db.Model(&User{}).Count(&count1).Error
	// lg.Dumper(count1, err)
	// err = db.Table("users_copy").Count(&count1).Error
	// lg.Dumper(count1, err)

	// если нужно получить данные из одной колонки в срез
	// var names []string
	// db.Model(&User{}).Pluck("name", &names)
	// lg.Dumper(names)

	// если нужно создать объект с учетом наличия его в БД по условиям
	// var user User
	// db.Where(User{Name: "Вася"}).FirstOrInit(&user) // не создает просто ищет  ( см. Attrs(User{Age: 30}) Assign(User{Age: 20}) )
	// db.Where(User{Name: "Муся"}).FirstOrCreate(&user) // создает если нет ( см. Attrs(User{Age: 30}) Assign(User{Age: 20}) )
	// lg.Dumper(user)

	// Raw SQL 1
	// type Result struct {
	// 	Name string
	// 	Age  int
	// }
	// var result []Result
	// db.Select("name, age").Table("sample").Where("name = ?", "Вася").Scan(&result)
	// lg.Dumper(result, err)

	// Raw SQL 2
	// var sample []*User
	// db.Raw("SELECT c.name, c.age FROM sample as c WHERE c.name = ?", "Вася").Scan(&sample)
	// lg.Dumper(sample)

	// транзакции
	// db.Set("gorm:query_option", "FOR UPDATE").First(&user, 10)

}

type Sample struct {
	ID          uint64    `gorm:"primary_key"`
	Name        string    `gorm:"default:NULL"`                      // если значение может быть NULL
	Age         int64     ``                                         // если ничего - для форматирвоания
	Birthday    time.Time `gorm:"default:NULL"`                      //
	Cnt         int32     `gorm:"default:20"`                        // если значение не может быть NULL
	Description string    `gorm:"column:opisanie;default:'popcorn'"` // если поле в бд имеет другое имя и не null
	Data        []byte    `gorm:"-"`                                 // исключает использование с БД
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
