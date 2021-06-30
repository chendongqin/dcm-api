package utils

import (
	"strconv"
	"strings"
)

/**
 * 版本号比较
 * 左边比右边大返回1
 * 左边比右边小返回-1
 * 左边等于右边返回0
 */
func CompareSimple(version1, version2 string) int {
	maxLen, left, right := 0, 0, 0

	v1 := strings.Split(version1, ".")
	v2 := strings.Split(version2, ".")
	len1 := len(v1)
	len2 := len(v2)

	if len1 > len2 {
		maxLen = len1
	} else {
		maxLen = len2
	}

	for i := 0; i < maxLen; i++ {
		if i < len1 && i < len2 {
			if v1[i] == v2[i] {
				continue
			}
		}

		left = 0
		if i < len1 {
			left, _ = strconv.Atoi(v1[i])
		}

		right = 0
		if i < len2 {
			right, _ = strconv.Atoi(v2[i])
		}

		if left < right {
			return -1
		} else if left > right {
			return 1
		}
	}
	return 0
}

func Compare(version1, version2, operator string) bool {
	compare := CompareSimple(version1, version2)

	switch {
	case operator == ">" || operator == "gt":
		return compare > 0
	case operator == ">=" || operator == "ge":
		return compare >= 0
	case operator == "<=" || operator == "le":
		return compare <= 0
	case operator == "==" || operator == "=" || operator == "eq":
		return compare == 0
	case operator == "<>" || operator == "!=" || operator == "ne":
		return compare != 0
	case operator == "" || operator == "<" || operator == "lt":
		return compare < 0
	}

	return false
}
