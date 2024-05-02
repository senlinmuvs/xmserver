package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"log"
	"math/rand"
	"net"
	"net/smtp"
	"os"
	"runtime"
	"time"

	"github.com/nfnt/resize"
)

var CHARSE = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-+")
var CHARSE1 = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
var CHARSE_NUM = []rune("0123456789")

func IfThen(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	} else {
		return falseVal
	}
}
func Rand(n int) int {
	cur := time.Now().UnixNano()
	rand.Seed(cur)
	return rand.Intn(n)
}

func Rand1(n int, seed int64) int {
	rand.Seed(seed)
	return rand.Intn(n)
}

func RandStr(n int) string {
	b := make([]rune, n)
	len := len(CHARSE)
	cur := time.Now().UnixNano()
	for i := range b {
		b[i] = CHARSE[Rand1(len, cur+int64(i))]
	}
	return string(b)
}

func RandStr1(n int) string {
	b := make([]rune, n)
	cur := time.Now().UnixNano()
	len := len(CHARSE1)
	for i := range b {
		b[i] = CHARSE1[Rand1(len, cur+int64(1))]
	}
	return string(b)
}

func RandNumberStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = CHARSE_NUM[Rand(len(CHARSE_NUM))]
	}
	return string(b)
}

func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ToJson(d map[string]interface{}) string {
	data, err := json.Marshal(d)
	if err != nil {
		fmt.Println("json.marshal failed, err:", err)
		return ""
	}
	return string(data)
}

func AppendJson(args ...interface{}) string {
	if args == nil {
		return ""
	}
	m := map[string]interface{}{}
	for i := 0; i < len(args); i += 2 {
		m[args[i].(string)] = args[i+1]
	}
	return ToJson(m)
}

func AppendJsonMap(args ...interface{}) (m map[string]any) {
	if args == nil {
		return
	}
	m = map[string]any{}
	for i := 0; i < len(args); i += 2 {
		m[args[i].(string)] = args[i+1]
	}
	return
}

func ObjToJsonStr(o interface{}) string {
	data, err := json.Marshal(o)
	if err != nil {
		log.Println("json.marshal failed, err:", err)
		return ""
	}
	return string(data)
}

func ObjToJsonStyle(o interface{}) (string, error) {
	data, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	err = json.Indent(&out, data, ",", "  ")
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func Stack() string {
	var buf [2 << 10]byte
	return string(buf[:runtime.Stack(buf[:], true)])
}

func Humanmills(mills int64) int64 {
	s := time.Unix(mills/1000, 0).Format("20060102150405")
	return int64(ToInt(s))
}

func SendToMail(user, nick, password, host string, port int, to, subject, body string) (err error) {
	header := make(map[string]string)
	header["From"] = nick + "<" + user + ">"
	header["To"] = to
	header["Subject"] = subject
	header["Content-Type"] = "text/html; charset=UTF-8"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	auth := smtp.PlainAuth("", user, password, host)
	err = SendMailUsingTLS(fmt.Sprintf("%s:%d", host, port), auth, user, []string{to}, []byte(message))
	return err
}

// return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// 参考net/smtp的func SendMail()
// 使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
// len(to)>1时,to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func ExistsInStrArr(arr []string, item string) bool {
	for _, it := range arr {
		if it == item {
			return true
		}
	}
	return false
}
func ExistsInArr(arr []interface{}, item interface{}) bool {
	for _, it := range arr {
		if it == item {
			return true
		}
	}
	return false
}
func ExistsInIntArr(arr []int, item int) bool {
	for _, it := range arr {
		if it == item {
			return true
		}
	}
	return false
}
func GetErrStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func CurMills() int64 {
	t := time.Now()
	return int64(time.Nanosecond) * t.UnixNano() / int64(time.Millisecond)
}

func CurTime() int64 {
	t := time.Now()
	ts := t.Format("20060102150405")
	return ToInt64(ts)
}
func CurTimeStr() string {
	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")
	return ts
}

func Mkdir(d string) (err error) {
	if _, err = os.Stat(d); os.IsNotExist(err) {
		err = os.Mkdir(d, os.ModePerm)
	}
	return
}

func ResizeImage(f, newf string, w uint) error {
	// open "test.jpg"
	file, err := os.Open(f)
	if err != nil {
		print("ResizeImage err1", err)
		return err
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		print("ResizeImage err2", err)
		return err
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(w, 0, img, resize.Lanczos3)

	out, err := os.Create(newf)
	if err != nil {
		print("ResizeImage err3", err)
		return err
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
	return nil
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func MinusArr(arr1, arr2 []string) (c []string) {
	for _, a1 := range arr1 {
		found := false
		for _, a2 := range arr2 {
			if a1 == a2 {
				found = true
			}
		}
		if !found {
			c = append(c, a1)
		}
	}
	return
}

func StrArrToAnyArr(arr []string) (arr2 []interface{}) {
	for _, a := range arr {
		arr2 = append(arr2, a)
	}
	return
}

func ArrEq(arr1, arr2 []string) bool {
	for i, ar := range arr1 {
		if ar != arr2[i] {
			return false
		}
	}
	return true
}
func ArrEqIg(arr1, arr2 []string, igEle ...string) bool {
	for i, ar := range arr1 {
		isIg := false
		for _, ie := range igEle {
			if ar == ie {
				isIg = true
				break
			}
		}
		if isIg {
			continue
		}
		if ar != arr2[i] {
			return false
		}
	}
	return true
}
func ByteArrEq(arr1, arr2 []byte) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	for i, e := range arr1 {
		if e != arr2[i] {
			return false
		}
	}
	return true
}
