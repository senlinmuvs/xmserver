package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gopkg.in/yaml.v2"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func Tr(s string) string {
	return strings.Trim(s, " \n\t")
}

func SubUnicode(s string, i, j int) string {
	if i < 0 || j < 0 {
		return ""
	}
	srune := []rune(s)
	return string(srune[i:j])
}

func LenUnicode(s string) int {
	return len([]rune(s))
}

func LeftUnicode(s string, len int) string {
	if len <= 0 {
		return ""
	}
	return SubUnicode(s, 0, Min(len, LenUnicode(s)))
}

func RightUnicode(s string, i int) string {
	return SubUnicode(s, Min(i, LenUnicode(s)), LenUnicode(s))
}

func ToInt(s string) int {
	if s == "" {
		return 0
	}
	x, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Println("err: toInt() err", err)
		return 0
	}
	return int(x)
}
func ToIntDef(s string, def int) int {
	if s == "" {
		return def
	}
	x, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return int(x)
}
func ToInt64(s string) int64 {
	return ToInt64Def(s, 0)
}
func ToInt64Def(s string, def int64) int64 {
	if s == "" {
		return def
	}
	x, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Println("err: toInt64() for value", s)
		return def
	}
	return int64(x)
}
func ToIntArr(s string, sp string, except int) []int {
	arr := []int{}
	if s == "" {
		return arr
	}
	arr2 := strings.Split(s, sp)
	for _, id := range arr2 {
		if id != "" {
			x := ToInt(id)
			if x != except {
				arr = append(arr, x)
			}
		}
	}
	return arr
}

func ToStr(arr []string, sep string) string {
	str := ""
	for _, s := range arr {
		str += s + sep
	}
	if len(str) > 0 {
		str = LeftUnicode(str, len(str)-1)
	}
	return "[" + str + "]"
}

/*
 * 分割后，提取第几个值, 从0开始
 */
func ExtractVal(s string, i int, sep string) int {
	arr := strings.Split(s, sep)
	n := 0
	for _, e := range arr {
		if n == i {
			return ToInt(e)
		}
		n++
	}
	return 0
}

func ExtractArea(s string) (string, int, int) {
	if len(s) < 4 {
		return "", -1, -1
	}
	i0 := strings.Index(s, "(")
	if i0 == -1 {
		i0 = strings.Index(s, "[")
	}
	i1 := strings.Index(s, ")")
	if i1 == -1 {
		i1 = strings.Index(s, "]")
	}
	return SubUnicode(s, i0, i1+1), i0, i1
}

func Index(rstr []rune, sub string, start int) int {
	for i, r := range rstr {
		if i > start {
			if string(r) == sub {
				return i
			}
		}
	}
	return -1
}

func IndexUnicode(str, substr string) int {
	// 子串在字符串的字节位置
	result := strings.Index(str, substr)
	if result >= 0 {
		// 获得子串之前的字符串并转换成[]byte
		prefix := []byte(str)[0:result]
		// 将子串之前的字符串转换成[]rune
		rs := []rune(string(prefix))
		// 获得子串之前的字符串的长度，便是子串在字符串的字符位置
		result = len(rs)
	}
	return result
}

func LastIndexUnicode(str, substr string) int {
	// 子串在字符串的字节位置
	result := strings.LastIndex(str, substr)
	if result >= 0 {
		// 获得子串之前的字符串并转换成[]byte
		prefix := []byte(str)[0:result]
		// 将子串之前的字符串转换成[]rune
		rs := []rune(string(prefix))
		// 获得子串之前的字符串的长度，便是子串在字符串的字符位置
		result = len(rs)
	}
	return result
}

func Extract(s_, a, b string) (string, int, int) {
	s := []rune(s_)
	l := len(s)
	if l == 0 || l < len(a) || l < len(b) {
		return "", -1, -1
	}
	i0 := Index(s, a, 0)
	if i0 < 0 {
		return "", i0, -1
	}
	i1 := Index(s, b, i0)
	if i1 < 0 || i0 == i1 {
		return "", i0, i1
	}
	// fmt.Println(string(s), i0, i1)
	return SubUnicode(s_, i0+len(a), i1), i0, i1
}

func ExtractPart(s, flag1, flag2 string) string {
	i := IndexUnicode(s, flag1)
	j := IndexUnicode(s, flag2)
	return SubUnicode(s, i+len(flag1), j)
}

func NumStrInArea(s, area string) bool {
	s = strings.Trim(s, " ")
	if s == "" {
		return false
	}
	x, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return false
	}
	return InArea(int(x), area)
}

/*
 * 判断某个值在不在某个区间内
 */
func InArea(v int, area string) bool {
	if len(area) < 4 {
		return false
	}

	limit0 := area[0]
	limit1 := area[len(area)-1]
	if limit0 != '(' && limit0 != '[' {
		return false
	}
	if limit1 != ')' && limit1 != ']' {
		return false
	}
	area_ := SubUnicode(area, 1, len(area)-1)
	arr := strings.Split(area_, ",")
	if len(arr) < 2 {
		return false
	}
	a_ := arr[0]
	b_ := arr[1]
	if a_ == "" && b_ == "" {
		return false
	}

	if a_ != "" {
		a := ToInt(a_)
		if limit0 == '(' {
			if v <= a {
				return false
			}
		} else if limit0 == '[' {
			if v < a {
				return false
			}
		}
	}
	if b_ != "" {
		b := ToInt(b_)
		if limit1 == ')' {
			if v >= b {
				return false
			}
		} else if limit1 == ']' {
			if v > b {
				return false
			}
		}
	}
	return true
}

func Split(s, sep string) []string {
	arr := []string{}
	arr_ := strings.Split(s, sep)
	for _, a := range arr_ {
		a = strings.Trim(a, " ")
		if a != "" {
			arr = append(arr, a)
		}
	}
	return arr
}

func IsNumber(s string) bool {
	if s == "" {
		return false
	}
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func Uint8ArrToStr(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

func ObjToStr(v interface{}) (string, error) {
	switch r := v.(type) {
	case int64:
		return strconv.Itoa(int(r)), nil
	// case []byte:
	// 	return string(r), nil
	case string:
		return r, nil
	case []uint8:
		return Uint8ArrToStr(r), nil
	default:
		return "", fmt.Errorf("not found type %s", r)
	}
}

func IsMail(s string) bool {
	return len(s) > 3 && len(s) < 254 && emailRegex.MatchString(s)
}

func IsFile(s string) bool {
	return strings.Index(s, "file://") == 0
}

func IsHttp(s string) bool {
	return strings.Index(s, "http://") == 0 || strings.Index(s, "https://") == 0
}

func EndStart(s, endStr string) bool {
	i := strings.Index(s, endStr)
	return i == len(s)-len(endStr)
}

func MinusStr(s1, s2 string) (s string) {
	str := []rune{}
	s1_ := []rune(s1)
	s2_ := []rune(s2)
	for _, c1 := range s1_ {
		found := false
		for _, c2 := range s2_ {
			if c1 == c2 {
				found = true
				break
			}
		}
		if !found {
			str = append(str, c1)
		}
	}
	return string(str)
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func ConvertGBK2Str(gbkStr string) (s string) {
	//将GBK编码的字符串转换为utf-8编码
	s, e := simplifiedchinese.GBK.NewDecoder().String(gbkStr)
	if e != nil {
		return
	}
	//如果是[]byte格式的字符串，可以使用Bytes方法
	b, e := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(gbkStr))
	if e != nil {
		return
	}
	return string(b)
}

func ConvertStr2GBK(str string) (s string) {
	//将utf-8编码的字符串转换为GBK编码
	s, e := simplifiedchinese.GBK.NewEncoder().String(str)
	if e != nil {
		return
	}
	//如果是[]byte格式的字符串，可以使用Bytes方法
	b, e := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(str))
	if e != nil {
		return
	}
	return string(b)
}

func SplitLast(s, sep string) (arr []string) {
	i := strings.LastIndex(s, sep)
	if i < 0 {
		return []string{s, ""}
	}
	return []string{s[:i], s[i:]}
}

// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Str(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

func Str2Obj(s string, obj interface{}) {
	err := yaml.Unmarshal(([]byte)(s), obj)
	if err != nil {
		log.Fatalf("parse json err: %v", err)
	}
}

func Join(arr []interface{}, sep string) string {
	if arr == nil {
		return ""
	}
	s := ""
	l := len(arr)
	for i, e := range arr {
		if i == l-1 {
			s += Str(e)
		} else {
			s += Str(e) + sep
		}
	}
	return s
}

func Join2(arr []string, sep string) string {
	if arr == nil {
		return ""
	}
	s := ""
	l := len(arr)
	for i, e := range arr {
		if i == l-1 {
			s += Str(e)
		} else {
			s += Str(e) + sep
		}
	}
	return s
}

// go1.13.8 unicode11
func IsChinese(w string) bool {
	for _, a := range w {
		if !(a >= 0x3400 && a <= 0x4db5) && !(a >= 0x4e00 && a <= 0x9fef) {
			return false
		}
	}
	return true
}

// go1.13.8 unicode11
func IsEnglish(w string) bool {
	for _, a := range w {
		if !(a >= 0x0041 && a <= 0x005a) &&
			!(a >= 0x0061 && a <= 0x007A) {
			return false
		}
	}
	return true
}

func Compact(s string, cs ...string) string {
	for _, c := range cs {
		for {
			cc := c + c
			i := strings.Index(s, cc)
			if i < 0 {
				break
			}
			s = strings.ReplaceAll(s, cc, c)
		}
	}
	return s
}
