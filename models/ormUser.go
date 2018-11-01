package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"github.com/beego/orm"
	"github.com/beego"
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
	var miner User
	miner.Address = pubKey
	o := orm.NewOrm()
	err := o.Read(&miner,"Address")
	if err != nil {
		beego.Info("查询失败")
		return 0
	}else {
		//挖矿UTXO
		time.Sleep(2*time.Second)  //理论上是计算nounce
		miner.UTXO =+12.5      //12.5理论上动态减半
	}
	//更新user
	_,err = o.Update(&miner,"UTXO")
	if err != nil {
		beego.Info("更新失败")
		return 0
	}
	return miner.UTXO
}