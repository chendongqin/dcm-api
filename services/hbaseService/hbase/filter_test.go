package hbase

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestNewFilters(t *testing.T) {
	filters := NewFilters().
		Add(NewSingleColumnValueFilter("r", "dc", CompareOperatorLess, BinaryComparator(1000), true, true)).
		And(&PrefixFilter{Prefix: "V1"}).ToByteArray()
	assert.Equal(t, "SingleColumnValueFilter('r', 'dc', <, 'binary:1000', true, true) AND PrefixFilter('V1')", string(filters))
}
