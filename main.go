package main

import (
	_ "beegoBlockChain/routers"
	"github.com/astaxie/beego"
	_"beegoBlockChain/models"
)

func main() {
	beego.Run()
	//models.CreateFirstBlock()

}

