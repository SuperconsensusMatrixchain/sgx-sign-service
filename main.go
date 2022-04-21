package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

/// 此应用放在sgx 运行
/// 需要做好数据冗余处理

func init() {
	// 初始化数据库dbsoucename
	dbsoucename := "xuperchain"
	db, err := InitDB(dbsoucename)
	if err != nil {
		panic("init db error")
	}
	GDB = db
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		// 使用统一返回结构
		c.JSON(http.StatusOK, NewResponse(http.StatusOK, SUCCESSMSG))
	})

	// 创建地址
	r.GET("/create", func(c *gin.Context) {
		// 创建地址
		addr, err := CreateXuperAccount()
		if err != nil {
			c.JSON(http.StatusOK, NewResponse(ErrCode, ERRORMSG))
			return
		}
		c.JSON(http.StatusOK, NewResponse(OKCode, SUCCESSMSG).WithData([]byte(addr)))
	})

	// 签名
	r.POST("/sign", func(c *gin.Context) {
		// 传入需要地址和数据
		paramters := struct {
			Address string `json:"address"`
			Msg     []byte `json:"msg"`
		}{}
		err := c.ShouldBindJSON(&paramters)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewResponse(ErrCode, PAMATERSERR))
			return
		}
		// 校验参数
		if paramters.Address == "" || paramters.Msg == nil {
			c.JSON(http.StatusBadRequest, NewResponse(ErrCode, PAMATERSERR))
			return
		}
		// 签名
		signserve := NewXuperchainAccount(paramters.Address)
		sign, err := signserve.Sign(paramters.Msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewResponse(500, "sign error"))
			return
		}
		c.JSON(http.StatusOK, NewResponse(OKCode, SUCCESSMSG).WithData(sign))
	})

	// 验证
	r.POST("/verify", func(c *gin.Context) {
		// sign  and  msg
		paramters := struct {
			Address string `json:"address"`
			Sign    []byte `json:"sign"`
			Msg     string `json:"msg"`
		}{}
		err := c.ShouldBindJSON(&paramters)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewResponse(ErrCode, PAMATERSERR))
			return
		}
		if paramters.Address == "" || paramters.Sign == nil || paramters.Msg == "" {
			c.JSON(http.StatusBadRequest, NewResponse(ErrCode, PAMATERSERR))
			return
		}
		// 校验
		signserve := NewXuperchainAccount(paramters.Address)
		result, err := signserve.verify(paramters.Sign, []byte(paramters.Msg))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, NewResponse(500, "verify error"))
			return
		}
		c.JSON(http.StatusOK, NewResponse(OKCode, "verify success").WithData([]byte(strconv.FormatBool(result))))
	})

	// 验证 address 是否已经存在
	r.POST("/is-exist", func(c *gin.Context) {
		paramters := struct {
			Address string `json:"address"`
		}{}
		err := c.ShouldBindJSON(&paramters)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewResponse(ErrCode, PAMATERSERR))
			return
		}
		if paramters.Address == "" {
			c.JSON(http.StatusBadRequest, NewResponse(ErrCode, PAMATERSERR))
			return
		}
		result := IsExist(paramters.Address)
		c.JSON(http.StatusOK, NewResponse(OKCode, "query address").WithData([]byte(strconv.FormatBool(result))))
	})

	//监听端口默认为8080
	r.Run(":8080")
}
