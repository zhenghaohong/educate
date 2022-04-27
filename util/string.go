package util
import "strconv"

func StringToInt(e string) (int, error) {
	return strconv.Atoi(e)
}


func StringToFloat(e string) (float64, error) {
	return strconv.ParseFloat(e, 64)
}



