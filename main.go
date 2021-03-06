package main

import (
	_ "Heihei/routers"
	"github.com/astaxie/beego"
  "github.com/astaxie/beego/plugins/cors"
)

func main() {
  beego.InsertFilter("*", beego.BeforeRouter,cors.Allow(&cors.Options{
  AllowOrigins: []string{"*"},
  AllowMethods: []string{"*"},
  AllowHeaders: []string{"Origin"},
  ExposeHeaders: []string{"Content-Length"},
  AllowCredentials: true,
  }))

  beego.SetStaticPath("/", "static")
  beego.SetStaticPath("/index_files", "static/index_files")
  beego.SetStaticPath("/intl", "static/intl")

	beego.Run()
}
