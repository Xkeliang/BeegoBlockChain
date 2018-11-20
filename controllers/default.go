package controllers

import (
	"github.com/astaxie/beego"
	"beegoBlockChain/models"
	//"fmt"
	"github.com/beego/orm"
	"fmt"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	//区块上链，验证区块
	o := orm.NewOrm()
	var preBlocks  []models.Block

	n,err := o.QueryTable("block").All(&preBlocks)
	if n == 0 {
		models.CreateFirstBlock()
		_,err = o.QueryTable("block").All(&preBlocks)
		goto tab
	}else {
	var preBlock models.Block
	if err != nil {

		beego.Info("读取错误")
		return
	}else {
		preBlock = preBlocks[len(preBlocks)-1]
	}

	isTrue :=models.IsValid(preBlock,preBlocks[(len(preBlocks)-2)])
	if !isTrue {
		beego.Info("验证失败")
		return
	}
	}
tab:
	c.Data["preBlocks"]=preBlocks
	c.TplName = "index.html"
}
func (c *MainController) Post() {
	c.TplName = "index.html"
}
func (c *MainController) ShowReg () {
	//c.TplName="registered.html"
	c.TplName="registered.html"
}
func (c *MainController) HandleReg () {
	var preKey string
	var pubKey string
	preKey = c.GetString("preKey")
	pubKey = models.RegisterAddress(preKey)
	c.Data["pubKey"] = pubKey
	o := orm.NewOrm()
	var user models.User
	user.UTXO=0
	user.Address= pubKey
	_,err := o.Insert(&user)
	if err != nil {
		beego.Info("注册失败")
		return
	}
	c.TplName = "registered.html"


}
func (c *MainController)ShowMining()  {
	c.TplName="mining.html"
}
func (c *MainController)HandleMining()  {
	Address := c.GetString("address")
	UTXO :=models.Mining(Address)
	if UTXO ==-1 {
		c.Ctx.WriteString("有人已经成功挖出该区块")
	}
	fmt.Println("挖出BTC=",UTXO)
	c.Data["UTXO"]=UTXO
	c.Data["address"]=Address


	c.TplName="mining.html"
}