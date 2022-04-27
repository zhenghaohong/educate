package util
import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// 字符串转十六进制
func HexToString(arrInt []int) (weight, height string) {
	arr1 := IntArr(arrInt)
	bs, err := hex.DecodeString(arr1)
	if err != nil {
		return "", ""
	}
	fmt.Println(string(bs))

	str_arr := strings.Split(string(bs), "")
	fmt.Println(str_arr[1])
	weightStr := str_arr[3:7]
	heightStr := str_arr[8:11]
	fmt.Printf("weight: %+v  height: %+v \n", weightStr, heightStr)

	// weight = str_arr[3] + str_arr[4] + str_arr[5] + str_arr[6]
	weight = str_arr[3] + str_arr[4]+"."+ str_arr[5] 
	height = str_arr[8] + str_arr[9] + str_arr[10]
	fmt.Println(weight)
	fmt.Println(height)
	return weight, height
}




func DealHeightWeight(str string) (weight, height string) {
	bs, err := hex.DecodeString(str)
	if err != nil {
		return "", ""
	}
	str_arr := strings.Split(string(bs), "")
	weight = str_arr[3] + str_arr[4] + str_arr[5] + str_arr[6]
	height = str_arr[8] + str_arr[9] + str_arr[10]
	return weight, height
}

func BytetoH(b []byte) (H string) {
	H = fmt.Sprintf("%x", b)
	return
}

func Hextob(str string) []byte {
	slen := len(str)
	bHex := make([]byte, len(str)/2)
	ii := 0
	for i := 0; i < len(str); i = i + 2 {
		if slen != 1 {
			ss := string(str[i]) + string(str[i+1])
			bt, _ := strconv.ParseInt(ss, 16, 32)
			bHex[ii] = byte(bt)
			ii = ii + 1
			slen = slen - 2
		}
	}
	return bHex

}

func Tool_DecimalByteSlice2HexString(DecimalSlice []byte) string {
	var sa = make([]string, 0)
	for _, v := range DecimalSlice {
		sa = append(sa, fmt.Sprintf("%02X", v))
	}
	ss := strings.Join(sa, "")
	return ss
}

func IntArr(arrInts []int) string {
	str := ""
	for i := 0; i < len(arrInts); i++ {
		fmt.Printf("%.2x", arrInts[i])
		str = str + fmt.Sprintf("%.2x", arrInts[i])
	}
	return str

}





