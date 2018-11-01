package controllers

import (
	"github.com/astaxie/beego"
	"github.com/beego/orm"
	"beegoBlockChain/models"

)
func init() {
	orm.RegisterDataBase("default","mysql","root:root@/chain?charset=utf8")
	orm.RegisterModel(new(models.Block),new(models.User))
	orm.RunSyncdb("default",false,true)
}
type SendController struct {
	beego.Controller
}
func (c *SendController)Post()  {
	//beego.Info("-----1111111")
	data := c.GetString("matter")
	//beego.Info("-----1111111")
	o := orm.NewOrm()
	////找到最后一个区块
	//beego.Info("---------------")
	var preBlocks []models.Block
	preBlock := models.Block{}
	newBlock := models.Block{}
	_,err :=o.QueryTable("block").All(&preBlocks)

	if err != nil {
		beego.Info("读取错误")
		return
	}else {
		preBlock = preBlocks[len(preBlocks)-1]
	}
	newBlock =models.CreatNewBlock(preBlock,data)
	_,err = o.Insert(&newBlock)
	c.Redirect("/",302)
}
func (c *SendController)Get()  {
	c.TplName="send.html"
}