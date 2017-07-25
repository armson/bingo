package utils

import(
    "bytes"
	"sort"
	"strings"
)

type binSlice string
var Slice *binSlice

func (_ *binSlice) Join(slices []string, sep string) (s string) {
    if len(slices) < 1 { return ""}
    buf := bytes.Buffer{}
    for _, value := range slices {
        buf.WriteString(value)
        buf.WriteString(sep)
    }
    s = buf.String()
    s = s[0 : len(s)-len(sep)]
    return
}

func (_ *binSlice) In(needle string, slices []string) (bool) {
    if len(slices) < 1 || len(needle) < 1 { return false}
    for _, value := range slices {
        if needle == value {
            return true
        }
    }
    return false
}

func (_ *binSlice) JoinSQL(slices []string) (s string) {
    if len(slices) < 1 { return ""}
    buf := bytes.Buffer{}
    for _, value := range slices {
		buf.WriteString("'")
        buf.WriteString(value)
        buf.WriteString("',")
    }
	s = buf.String()
	s = s[0 : len(s)-1]
    return
}

// 注：只针对 []map[string]interface{} 类型，并 rowKey的值是int类型
func (_ *binSlice) SortStringInterfaceInt(slices []map[string]interface{} , rowKey string , order string) ([]map[string]interface{}) {
	sorted := []map[string]interface{}{}

	length := len(slices)
	if length < 1 { return sorted }
	if length < 2 { return slices }

	rows := []int{}
	maps := map[int][]map[string]interface{}{}

	for _ , s := range slices {
		key := s[rowKey].(int)
		rows = append(rows, key)
		if _, m := maps[key]; m {
			maps[key] = append(maps[key], s)
		} else {
			maps[key] = []map[string]interface{}{s}
		}
	}
	if strings.ToLower(order) == "desc" {
		sort.Sort(sort.Reverse(sort.IntSlice(rows)))
	} else {
		sort.Ints(rows)
	}

	preValue := 0
	pos := 0
	for k , value := range rows {
		if k != 0 {
			if value == preValue {
				pos = pos + 1
			} else {
				pos = 0
			}
			preValue = value
		}
		sorted = append(sorted, maps[value][pos])
	}
	return sorted
}

func (_ *binSlice) CutStringInterface(slices []map[string]interface{} , offset, limit int) ([]map[string]interface{}) {
	length := len(slices)
	end := offset + limit

	if end < length {
		return slices[offset:end]
	}
	return slices[offset:length]
}

func (_ *binSlice) ColumnStringString(arrs []map[string]string, rowKey string) []string {
	rows := []string{}
	if len(arrs) < 1 {
		return rows
	}
	for _, arr := range arrs {
		rows = append(rows, arr[rowKey])
	}
	return rows
}

func (_ *binSlice) ColumnStringInterface(arrs []map[string]interface{}, rowKey string) []string {
	rows := []string{}
	if len(arrs) < 1 {
		return rows
	}
	for _, arr := range arrs {
		rows = append(rows, arr[rowKey].(string))
	}
	return rows
}

func (_ *binSlice) CombineStringString(arrs []map[string]string, rowKey string) map[string]map[string]string {
	rows := map[string]map[string]string{}
	if len(arrs) < 1 {
		return rows
	}
	for _, arr := range arrs {
		rows[arr[rowKey]] = arr
	}
	return rows
}







