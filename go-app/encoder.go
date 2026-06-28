package main

import (
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const base = int(62)

func Encode(id int) string {
	if id == 0 {
		return "a"
	}
	result := []byte{}
	for true{
		if id <= 0 {
			break
		}
		rem := id % base
		id = id / base
		result = append(result , alphabet[int(rem)])
	}

	for i , j := 0 , len(result) - 1 ; i < j ; i , j = i+1 , j-1 {
		result[i] , result[j] = result[j] , result[i]
	}
	return string(result)
}

func Decode(short_code string) int {
	num := int(0)
	for _ , char := range short_code {
		idx := strings.IndexRune(alphabet , char)
		num = (num*base) + idx
	}
	return num
}