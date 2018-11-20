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
type MiningPerson struct {
	user *User
	message chan bool
}

var MiningPersons []MiningPerson
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
	p :=orm.NewOrm()
	var preBlocks  []Block
	var preBlock Block
	n,err := o.QueryTable("block").All(&preBlocks)
	if err != nil {
		beego.Info("读取错误")
		return  -1
	}else {
		preBlock = preBlocks[len(preBlocks)-1]
	}
	if n == 1 {

	}
	preBlock = preBlocks[n-1]
	var date = PackTransaction()
	newBlock :=CreatNewBlock(preBlock,date)
	var miner User
	miner.Address = pubKey

	err = p.Read(&miner,"Address")
	m := miner.UTXO
	if err != nil {
		beego.Info("查询失败")
		return 0
	}else {
		//挖矿UTXO
		//time.Sleep(2*time.Second)  //理论上是计算nounce
		pow :=NewProofOfWork(&newBlock)
		nonce,hash := pow.Run(&miner)
		if nonce == -1 {
			return -1
		}
		newBlock.Nonce = nonce
		newBlock.Hash = hex.EncodeToString(hash[:])
		fmt.Println(newBlock.Hash)
		_,err = o.Insert(&newBlock)
		beego.Info("insert",err)

	}
		miner.UTXO =miner.UTXO+12.5      //12.5理论上动态减半
		//更新user
		_,err = p.Update(&miner,"UTXO")
		if err != nil {
			beego.Info("更新失败")
			return 0
		}
		fmt.Println(miner.UTXO)
		/*-----------------------------------------------------*/


		Reset()

		/*-----------------------------------------------------*/


		return miner.UTXO-m
	}