package users

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/sungora/app.v1/core"
)

// Модель
type User struct {
	ID         uint64     ``
	InvoiceID  uint64     `gorm:"default:null"`
	Names      string     `gorm:"default:null"`
	Cnt        int        ``
	Cnt8       int8       ``
	Cnt16      int16      ``
	Cnt32      int32      ``
	Cnt64      int64      ``
	Price32    float32    ``
	Price64    float64    ``
	Status     string     `gorm:"default:null"`
	StatusN    string     `gorm:"default:'Активный'"`
	IsCheck    bool       `gorm:"default:1;not null"`
	SampleSet  string     `gorm:"type:set('value 1','value 2','value 3');default:null"`
	SampleEnum string     `gorm:"type:enum('value 1','value 2','value 3');default:null"`
	CreatedAt  time.Time  ``
	UpdatedAt  time.Time  ``
	DeletedAt  *time.Time ``
}

// NewService создание модели
func NewUser(ID uint64) *User {
	self := new(User)
	if ID > 0 {
		core.DB.Find(self, ID)
	}
	return self
}

// TableName определение таблицы источника обьектов
func (self *User) TableName() string {
	return "users"
}

// BeforeCreate функция - хук вызовется перед вставкой - созданием записи
func (self *User) BeforeCreate(scope *gorm.Scope) error {
	return nil
}

// AfterCreate функция - хук вызовется после вставки - создания записи
func (self *User) AfterCreate(scope *gorm.Scope) error {
	return nil
}

// BeforeUpdate функция - хук вызовется перед обновлением записи
func (self *User) BeforeUpdate(scope *gorm.Scope) error {
	return nil
}

// AfterUpdate  функция - хук вызовется после обновления записи
func (self *User) AfterUpdate(scope *gorm.Scope) error {
	return nil
}

func (user *User) BeforeSave(scope *gorm.Scope) error {
	return nil
}

func (user *User) AfterSave(scope *gorm.Scope) error {
	return nil
}

// Load загрузка модели
func (self *User) Load(isCreate bool) error {
	if isCreate {
		return core.DB.Where(*self).FirstOrCreate(self).Error
	} else {
		return core.DB.Where(*self).FirstOrInit(self).Error
	}
}

// Save сохранение модели
func (self *User) Save() error {
	if self.ID > 0 {
		return core.DB.Save(self).Error
	} else {
		return core.DB.Create(self).Error
	}
}

// Delete удаление модели
func (self *User) Delete() error {
	err := core.DB.Delete(self).Error
	if err == nil {
		self.ID = 0
	}
	return err
}

// Sql получение хранилища пользовательских запросов
func (self *User) Sql() *query {
	return sql
}

func sampleOther() {
	// sample custom orm query
	var users []*User
	var count int
	err := core.DB.
		Select("id, name").
		Table("users").
		Joins("...", "...").
		Where("...", "...").
		Group("...").
		Having("...", "...").
		Order("id ASC").
		Limit(5).
		Find(&users).
		Count(&count).Error
	fmt.Println(users, count, err)
	// sample slice one column
	var names []string
	core.DB.Model(&User{}).Pluck("names", &names)
	//
}
