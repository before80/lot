package opdf

import "fmt"

func CellHeight(fontSize float64) float64 {
	return fontSize * 0.3528 * 1.2
}

func LineHeight(fontSize float64) float64 {
	return fontSize*0.3528*1.2 + 2
}

//func SliceToStr(s []any) (str string) {
//	for _, v := range s {
//		str += fmt.Sprintf("%v", v)
//	}
//	return str
//}

func SliceToStr[T any](s []T, delim string) string {
	var str string
	for i, v := range s {
		if i != len(s)-1 {
			str += fmt.Sprintf("%v%s ", v, delim)
		} else {
			str += fmt.Sprintf("%v", v)
		}
	}
	return str
}
