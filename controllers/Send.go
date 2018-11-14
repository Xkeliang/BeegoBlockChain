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
	address2 := c.GetString("toAddress")
	BTC,err := c.GetFloat("BTC",64)
	if err!=nil {
		beego.Info("BTC格式有误",err)
		return
	}
	//var user  []models.User
	var user1 = models.User{}
	var user2 = models.User{}
	user1.Address = address
	user2.Address=address2

	o := orm.NewOrm()

	err =o.Read(&user1,"Address")
	if err != nil {
		beego.Info("当前用户不存在")

		c.Ctx.WriteString("地址有误")
		return

	}
	err =o.Read(&user2,"Address")
	if err != nil {
		beego.Info("接受用户不存在")

		c.Ctx.WriteString("接受地址有误")
		return

	}
	if BTC > user1.UTXO {
		beego.Info("BTC 不足")
		c.Ctx.WriteString("BTC不足")
		return
	}
	c.Data["UTXO"]=user1.UTXO
	user1.UTXO=user1.UTXO-BTC
	user2.UTXO +=BTC

	//更新数据
	_,err =o.Update(&user1)
	if err != nil {
		beego.Info("更新失败",err)
		c.Ctx.WriteString("发送失败")
		return
	}
	_,err =o.Update(&user2)
	if err != nil {
		beego.Info("更新失败",err)
		c.Ctx.WriteString("发送失败")
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