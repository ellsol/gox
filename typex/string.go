package typex

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

func CompareStringList(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for k, v := range s1 {
		if s2[k] != v {
			return false
		}
	}
	return true
}

func MapStringList(list []string, fn func(string) string) []string {
	newList := make([]string, len(list))

	for i, v := range list {
		newList[i] = fn(v)
	}

	return newList
}

func MapStringListWithPos(list []string, fn func(int, string) string) []string {
	newList := make([]string, len(list))

	for i, v := range list {
		newList[i] = fn(i, v)
	}

	return newList
}

func FilterStringList(list []string, fn func(string) bool) []string {
	newList := make([]string, 0)

	for _, v := range list {
		if fn(v) {
			newList = append(newList, v)
		}
	}

	return newList
}

func UniqueStringList(list []string) []string {
	tempMap := make(map[string]int, 0)

	newList := make([]string, 0)

	for _, v := range list {
		tempMap[v] = 1
	}

	for k, _ := range tempMap {
		newList = append(newList, k)
	}

	return newList
}

func CommaSeparatedString(set []string) string {
	var buffer bytes.Buffer

	length := len(set)

	for k, v := range set {
		buffer.WriteString(v)

		if k < length-1 {
			buffer.WriteString(",")
		}
	}

	return buffer.String()
}

func StringListFromInterfaceList(val []interface{}) ([]string, error) {
	result := make([]string, len(val))

	for k, v := range val {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("couldn't cast interface of type %v to string", reflect.TypeOf(v))
		}
		result[k] = s
	}
	return result, nil
}

func StringListContains(value string, list []string) bool{
	for _,v := range list{
		if v == value {
			return true
		}
	}
	return false
}


func StringListContainsWord(value string, list []string) bool{
	for _,v := range list{
		if strings.Contains(v, value) {
			return true
		}
	}
	return false
}