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
	//send 确认地址
	address := c.GetString("formAddress")
	BTC,err := c.GetFloat("BTC",64)
	if err!=nil {
		beego.Info("BTC格式有误",err)
		return
	}
	//var user  []models.User
	var user1 = models.User{}
	user1.Address = address

	o := orm.NewOrm()

	err =o.Read(&user1,"Address")
	if err != nil {
		beego.Info("当前用户不从在")

		return

	}
	if BTC < user1.UTXO {
		beego.Info("BTC 不足")
	}
	c.Data["UTXO"]=user1.UTXO
	user1.UTXO=user1.UTXO-BTC

	//更新数据
	_,err =o.Update(&user1)
	if err != nil {
		beego.Info("更新失败",err)
		return
	}
	////找到最后一个区块
	data := c.GetString("matter")
	var preBlocks []models.Block
	preBlock := models.Block{}
	newBlock := models.Block{}
	_,err =o.QueryTable("block").All(&preBlocks)

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