package com

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"xmserver/util"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Config struct {
	XMS struct {
		Port    int
		DataDir string
	}
}

type User struct {
	CID  int
	User string
	Pwd  string
}
type Sync struct {
	ID    int64  `json:"id"`
	Table string `json:"table"`
	RowID int64  `json:"row_id"`
	T     int64  `json:"t"`
}

const (
	ST_OK       = 0
	ST_Err      = -1
	ST_Auth_Err = 1

	//
	K_Init_Start_Time = "init_start_time"
	K_Init_Done_Time  = "init_done_time"
)

var (
	Cfg   Config
	Users = []*User{}
)
var (
	ST_Msg_Map = map[int]string{
		ST_Auth_Err: "账户验证失败",
	}
)

func InitConfig() {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}
	err = yaml.Unmarshal(data, &Cfg)
	if err != nil {
		log.Fatalf("无法解析配置文件: %v", err)
	}

	initUserConfig()
}
func initUserConfig() {
	file, err := os.Open("user")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			fmt.Println("Invalid line format:", line)
			continue
		}
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("Error parsing ID:", parts[0])
			continue
		}
		u := &User{id, parts[1], parts[2]}
		Users = append(Users, u)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from file:", err)
	}
}

func RespOK(c *gin.Context, args ...any) {
	if args == nil {
		args = []interface{}{}
	}
	args = append(args, "st", ST_OK)
	js := util.AppendJsonMap(args...)
	c.JSON(http.StatusOK, js)
}

func RespErr(c *gin.Context, args ...any) {
	if args == nil {
		args = []interface{}{}
	}
	args = append(args, "st", ST_Err)
	js := util.AppendJsonMap(args...)
	c.JSON(http.StatusOK, js)
}

func RespErr1(c *gin.Context, st int, args ...any) {
	if args == nil {
		args = []interface{}{}
	}
	msg, ok := ST_Msg_Map[st]
	if !ok {
		msg = ""
	}
	args = append(args, "st", st, "msg", msg)
	js := util.AppendJsonMap(args...)
	c.JSON(http.StatusOK, js)
}

func RespErr2(c *gin.Context, msg string, args ...any) {
	if args == nil {
		args = []interface{}{}
	}
	args = append(args, "st", ST_Err, "msg", msg)
	js := util.AppendJsonMap(args...)
	c.JSON(http.StatusOK, js)
}

func RespErr3(c *gin.Context, st int, msg string, args ...any) {
	if args == nil {
		args = []interface{}{}
	}
	args = append(args, "st", st, "msg", msg)
	js := util.AppendJsonMap(args...)
	c.JSON(http.StatusOK, js)
}

func RespErr4(c *gin.Context, er *Err, args ...any) {
	if args == nil {
		args = []interface{}{}
	}
	args = append(args, "st", er.St, "msg", er.Msg)
	js := util.AppendJsonMap(args...)
	c.JSON(http.StatusOK, js)
}

type Err struct {
	St  int
	Msg string
}

func Err1(texts ...any) *Err {
	return Err2(ST_Err, texts...)
}
func Err2(code int, texts ...any) *Err {
	return &Err{code, util.Join(texts, " ")}
}
