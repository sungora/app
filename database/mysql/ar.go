package mysql

import (
	"strconv"
	"strings"

	"gopkg.in/kshamiev/sungora.v1/database/face"
)

type ar struct {
	property []string
	from     string
	where    []string
	group    string
	having   string
	order    string
	limit    string
}

func NewAr() face.ArFace {
	var self = new(ar)
	return self

}

func (self *ar) Select(property string) face.ArFace {
	self.property = append(self.property, property)
	return self
}

func (self *ar) From(from string) face.ArFace {
	self.from = from
	return self
}

func (self *ar) Where(where string) face.ArFace {
	self.where = append(self.where, where)
	return self
}

func (self *ar) Group(group string) face.ArFace {
	self.group = group
	return self
}

func (self *ar) Having(having string) face.ArFace {
	self.having = having
	return self
}

func (self *ar) Order(order string) face.ArFace {
	self.order = order
	return self
}

func (self *ar) Limit(start, step int) face.ArFace {
	self.limit = strconv.Itoa(start) + `, ` + strconv.Itoa(step)
	return self
}

func (self *ar) Get() (query string) {
	query += "SELECT\n\t" + strings.Join(self.property, `, `) + "\n"
	query += "FROM " + self.from + "\n"
	if len(self.where) > 0 {
		query += "WHERE 1\n\t" + strings.Join(self.where, "\n\t") + "\n"
	}
	if self.group != `` {
		query += `GROUP BY ` + self.group + "\n"
	}
	if self.having != `` {
		query += `HAVING ` + self.having + "\n"
	}
	if self.order != `` {
		query += `ORDER BY ` + self.order + "\n"
	}
	if self.limit != `` {
		query += `LIMIT ` + self.limit + "\n"
	}
	query = strings.TrimSpace(query)
	return
}
