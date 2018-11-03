package controllers

import "github.com/beego"

type MiningController struct {
	beego.Controller
}


func (c *MiningController)Get()  {
	c.TplName="mining.html"
}
func (c *MiningController)Post()  {
	//Address := c.GetString("address")
	c.TplName="mining.html"
}