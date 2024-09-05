package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"xmserver/com"
	"xmserver/dao"
	"xmserver/util"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	configDir := ""
	if len(os.Args) > 1 {
		configDir = os.Args[1]
	}
	com.InitConfig(configDir)
	dao.InitDao("xms.db")
}

func main() {
	r := gin.Default()
	r.Use(TLSHandler())
	r.Use(ErrorHandler())
	r.POST("/getInitDoneTime", AuthMiddleware(), getInitDoneTime)
	r.POST("/setInitStartTime", AuthMiddleware(), setInitStartTime)
	r.POST("/setInitDoneTime", AuthMiddleware(), setInitDoneTime)
	r.POST("/upfile", AuthMiddleware(), upfile)
	r.POST("/uprecord", AuthMiddleware(), uprecord)
	r.POST("/uprecord_", AuthMiddleware(), uprecord_)
	r.POST("/drecord", AuthMiddleware(), delRecord)
	r.POST("/dfile", AuthMiddleware(), delFile)
	r.POST("/downfile", AuthMiddleware(), downfile)
	// r.POST("/downrecord", AuthMiddleware(), downrecord)
	// // 设置静态文件路由
	// r.Static("/uploads", "./uploads")
	// // 启用静态文件托管
	// r.StaticFile("/", "client.html")
	// 启动服务器
	fmt.Println(com.Cfg.XMS.SSL.Cert, com.Cfg.XMS.SSL.Key)
	r.RunTLS(fmt.Sprintf(":%d", com.Cfg.XMS.Port), com.Cfg.XMS.SSL.Cert, com.Cfg.XMS.SSL.Key)
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				er, ok := r.(*com.Err)
				if ok {
					log.Println("handler err", er.St, er.Msg)
					com.RespErr4(c, er)
				} else {
					com.RespErr(c)
				}
			}
		}()
		c.Next() // 处理请求
	}
}

// TLSHandler 返回一个用于HTTPS的中间件处理器
func TLSHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.TLS == nil /*|| len(c.Request.TLS.PeerCertificates) == 0*/ {
			com.RespErr2(c, "请使用HTTPS访问")
			c.Abort()
			return
		}
		c.Next()
	}
}

// AuthMiddleware 是一个验证用户是否登录的中间件处理器
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cid, e := strconv.Atoi(c.PostForm("cid"))
		user := c.PostForm("user")
		pwd := c.PostForm("pwd")
		// fmt.Println(cid, user, pwd)
		if e != nil {
			com.RespErr1(c, com.ST_Auth_Err)
			c.Abort()
			return
		}
		u := findUser(cid, user, pwd)
		if u == nil {
			com.RespErr1(c, com.ST_Auth_Err)
			c.Abort()
			return
		}
		c.Next()
	}
}
func getInitDoneTime(c *gin.Context) {
	user := c.PostForm("user")
	s := dao.GetEnv(user + "_" + com.K_Init_Done_Time)
	com.RespOK(c, "t", util.ToInt64(s))
}
func upfile(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		panic(com.Err1("无法解析表单数据"))
	}
	path := c.PostForm("path")
	user := c.PostForm("user")

	files := form.File["files"]
	if len(files) == 0 {
		panic(com.Err1("没有选择要上传的文件"))
	}
	pt := fmt.Sprintf("%s/%s/%s", com.Cfg.XMS.DataDir, user, path)
	err = os.MkdirAll(pt, os.ModePerm)
	if err != nil {
		panic(com.Err1("创建目录失败", pt))
	}

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			panic(com.Err1("无法打开上传的文件: %v", err))
		}
		defer src.Close()

		dst, err := os.Create(fmt.Sprintf("%s/%s", pt, file.Filename))
		if err != nil {
			panic(com.Err1("无法创建目标文件: %v", err))
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			panic(com.Err1("文件复制失败: %v", err))
		}
	}
	com.RespOK(c)
}
func delRecord(c *gin.Context) {
	user := c.PostForm("user")
	table := c.PostForm("table")
	k := c.PostForm("k")
	v := c.PostForm("v")
	dao.DelRecord(user, table, k, v)
	com.RespOK(c)
}
func delFile(c *gin.Context) {
	user := c.PostForm("user")
	file := c.PostForm("file")
	e := os.Remove(com.Cfg.XMS.DataDir + "/" + user + "/" + file)
	if e != nil {
		panic(com.Err1("del file err", e))
	}
	com.RespOK(c)
}
func uprecord(c *gin.Context) {
	user := c.PostForm("user")
	table := c.PostForm("table")
	k := c.PostForm("k")
	v := c.PostForm("v")
	data := c.PostForm("data")
	dao.UpsertRecord(user, table, k, v, data)
	com.RespOK(c)
}
func uprecord_(c *gin.Context) {
	user := c.PostForm("user")
	sql := c.PostForm("sql")
	dao.ExeSql(user, sql)
	com.RespOK(c)
}
func setInitStartTime(c *gin.Context) {
	cur := util.CurMills()
	user := c.PostForm("user")
	dao.InsertOrUpdateEnv(user+"_"+com.K_Init_Start_Time, strconv.Itoa(int(cur)))
	com.RespOK(c, "t", cur)
}
func setInitDoneTime(c *gin.Context) {
	cur := util.CurMills()
	user := c.PostForm("user")
	dao.InsertOrUpdateEnv(user+"_"+com.K_Init_Done_Time, strconv.Itoa(int(cur)))
	com.RespOK(c, "t", cur)
}
func downfile(c *gin.Context) {
	// 设置想要下载的文件的路径和文件名
	filePath := "path/to/your/file.txt"
	// 设置客户端保存的文件名
	filename := "downloaded-file.txt"

	// 提供文件供下载，这会自动设置响应的 Content-Type 和 Content-Disposition
	// Content-Disposition header 使浏览器会以下载文件的方式处理响应的内容
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	// 设置响应的 Content-Type
	// 你可以根据你的文件类型来设定具体的值，例如 "application/octet-stream" 或 "application/pdf" 等
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	// 发送文件
	c.File(filePath)
}

func findUser(cid int, user, pwd string) *com.User {
	for _, u := range com.Users {
		if u.CID == cid && u.User == user && u.Pwd == pwd {
			return u
		}
	}
	return nil
}
