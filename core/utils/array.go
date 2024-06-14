package utils

import (
	"errors"
	"math/rand"
	"reflect"
	"time"
)

func StringArrayContains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func GetArrayItems(array interface{}) (res []interface{}, err error) {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			obj, ok := s.Index(i).Interface().(interface{})
			if !ok {
				return nil, errors.New("invalid type")
			}
			res = append(res, obj)
		}
	default:
		return nil, errors.New("invalid type")
	}
	return res, nil
}

func ShuffleArray(slice []interface{}) (err error) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}

	return nil
}
