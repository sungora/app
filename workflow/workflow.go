package workflow

import (
	"strconv"
	"strings"
)

// TaskAdd добавляем задачу в пул
func TaskAdd(task Task) {
	component.p.tasks <- task
}

// TaskAddCron добавляем фоновую задачу в пул
func TaskAddCron(task Task) {
	component.cronTaskRun = append(component.cronTaskRun, task)
}

func checkRuntime(val int, mask string) bool {
	var number int
	var sl []string
	//  any valid value or exact match
	number, _ = strconv.Atoi(mask)
	if "*" == mask || val == number {
		return true
	}
	//  range
	sl = strings.Split(mask, "-")
	if 1 < len(sl) {
		number1, _ := strconv.Atoi(sl[0])
		number2, _ := strconv.Atoi(sl[1])
		if number1 <= val && val <= number2 {
			return true
		}
		return false
	}
	//  fold
	sl = strings.Split(mask, "/")
	if 1 < len(sl) {
		number, _ = strconv.Atoi(sl[1])
		if 0 < val%number {
			return false
		}
		return true
	}
	//  list
	sl = strings.Split(mask, ",")
	if 1 < len(sl) {
		for _, v := range sl {
			number, _ = strconv.Atoi(v)
			if number == val {
				return true
			}
		}
		return false
	}
	//
	return false
}
