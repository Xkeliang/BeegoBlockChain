package controllers

import (
	"github.com/beego"
	"beegoBlockChain/models"
)

type RegController struct {
	beego.Controller
}

func (c *RegController) Get () {
	//c.TplName="registered.html"
	c.TplName="registered.html"
}
func (c *RegController) Post () {
	var preKey string
	var pubKey string
	preKey = c.GetString("preKey")
	pubKey = models.RegisterAddress(preKey)
	c.Data["pubKey"] = pubKey
	c.TplName = "registered.html"

}