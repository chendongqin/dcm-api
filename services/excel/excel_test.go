package excel

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SaveAs(t *testing.T) {
	excel := New()
	excel.SetHeader([]string{"第一行", "第二行", "第三行"})
	rows := make(Rows, 0)
	rows = append(rows, []Value{"111", "111", "111"})
	rows = append(rows, []Value{"111", "222", "222"})
	excel.SetRows(rows)
	err := excel.SaveAs("here.xls")
	assert.Error(t, err)
}

func Test_TransferCellName(t *testing.T) {
	for i := 1; i <= 30; i++ {
		str := TransferCellName(i, 1)
		fmt.Println(str)
	}
}
