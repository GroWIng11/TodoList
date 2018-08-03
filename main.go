package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"net/http"
	"fmt"
)

type Item struct{
	Id string `json:"id"`
	Content string `json:"content"`
}

var db *sql.DB

func main() {
	db = OpenSql()
	router := gin.Default()
	router.POST("/todo", posting)
	router.GET("/todo/:id", getting)
	router.DELETE("/todo/:id", deleting)
	router.PUT("/todo", putting)
	router.GET("/todo", gettingall)
	router.Run(":8088")
}

func posting(c *gin.Context) {
	var item Item
	c.BindJSON(&item)
	if (item.Id == "" || item.Content == ""){
		c.String(http.StatusBadRequest, "Creation failed")
		return
	}

	_, err := db.Exec("insert into tab_items values(" + "'" + item.Id + "', '" + item.Content + "')")
	if err != nil {
		c.String(http.StatusBadRequest, "Creation failed")
	} else {
		c.String(http.StatusOK, "Created")
	}
}

func getting(c *gin.Context) {
	id := c.Param("id")
	var content string
	err := db.QueryRow("select content from tab_items where id = " + id).Scan(&content)
	if err != nil {
		c.String(http.StatusBadRequest, "Item not found")
	} else {
		c.JSON(http.StatusOK, gin.H{
			"Item" + id : content,
		})
	}
}

func deleting(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("delete from tab_items where id = " + id)
	if err != nil {
		c.String(http.StatusBadRequest, "deleting failed")
	} else {
		c.String(http.StatusOK, "deleted")
	}
}

func putting(c *gin.Context) {
	var item Item
	c.BindJSON(&item)
	if (item.Id == "" || item.Content == ""){
		c.String(http.StatusBadRequest, "Creation failed")
		return
	}

	_, err := db.Exec("update tab_items set content=" + "'" + item.Content + "'where id = " + item.Id)
	if err != nil {
		c.String(http.StatusBadRequest, "updating failed")
	} else {
		c.String(http.StatusOK, "updated")
	}
}

func gettingall(c *gin.Context) {
	rows, _ := db.Query("select * from tab_items")
	mp := make(map[string]string)
	flag := false
	for rows.Next() {
		var id string
		var content string
		err := rows.Scan(&id, &content)
		if (err != nil) {
			continue
		}
		mp["Item" + id] = content
		flag = true
	}
	if flag {
		c.JSON(http.StatusOK, mp)
	} else {
		c.String(http.StatusBadRequest, "no itmes")
	}

}

func OpenSql() *sql.DB {
	fmt.Println("yes")
	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/?charset=utf8")
	checkErr(err)
	_, err1 := db.Exec("create database if not exists db_items")
	checkErr(err1)
	_, err2 := db.Exec("use db_items")
	checkErr(err2)
	_, err3 := db.Exec("create table if not exists tab_items(id varchar(20) not null primary key, " +
		"content varchar(100))")
	checkErr(err3)
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
