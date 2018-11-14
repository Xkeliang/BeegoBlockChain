package models

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/beego/orm"
	"github.com/beego"
	"fmt"

)

//BTC地址（公钥）与UTXO

type User struct {
	Id int	`orm:"pk;auto"`
	Address  string
	UTXO float64
}
//随机私钥,生成地址
func RegisterAddress(preKey string) string  {
	var pubKey string
	value := sha256.Sum256([]byte(preKey))
	pubKey = hex.EncodeToString(value[:])
	return pubKey

}


//挖矿产生BTC
func Mining(pubKey string) float64 {

	o := orm.NewOrm()
	var preBlocks  []Block
	var preBlock Block
	n,err := o.QueryTable("block").All(&preBlocks)
	if n == 1 {

	}
	preBlock = preBlocks[n-1]
	var date = ""
	newBlock :=CreatNewBlock(preBlock,date)
	var miner User
	miner.Address = pubKey

	err = o.Read(&miner,"Address")
	m := miner.UTXO
	vt :=false
	if err != nil {
		beego.Info("查询失败")
		return 0
	}else {
		//挖矿UTXO
		//time.Sleep(2*time.Second)  //理论上是计算nounce
		pow :=NewProofOfWork(&newBlock)
		nonce,hash := pow.Run()
		newBlock.Nonce = nonce
		newBlock.Hash = string(hash)
		}
		if vt==false {
			beego.Info("挖矿失败")
		}
		miner.UTXO =miner.UTXO+12.5      //12.5理论上动态减半
		//更新user
		_,err = o.Update(&miner,"UTXO")
		if err != nil {
			beego.Info("更新失败")
			return 0
		}
		fmt.Println(miner.UTXO)
		return miner.UTXO-m
	}