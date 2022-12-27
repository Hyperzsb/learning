package question

import (
	"fmt"
	"strconv"
)

type Question struct {
	Operand  [2]int
	Operator byte
	Result   int
}

func (q *Question) Parse(raw []string) error {
	if len(raw) != 2 {
		return fmt.Errorf("invalid field number")
	}

	question, ans := raw[0], raw[1]
	for i := 0; i < len(question); i++ {
		switch {
		case question[i] >= '0' && question[i] <= '9':
			continue
		case question[i] == '+' || question[i] == '-':
			q.Operator = question[i]
			if val, err := strconv.Atoi(string(question[:i])); err == nil {
				q.Operand[0] = val
			} else {
				return err
			}
			if val, err := strconv.Atoi(string(question[i+1:])); err == nil {
				q.Operand[1] = val
			} else {
				return err
			}
			break
		default:
			return fmt.Errorf("invalid character: %s", string(question[i]))
		}
	}
	if val, err := strconv.Atoi(ans); err == nil {
		q.Result = val
	} else {
		return err
	}

	return nil
}
