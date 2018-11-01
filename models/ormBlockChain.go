package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"crypto/sha256"
	"encoding/hex"
	//"github.com/beego"
)

type Block struct {
	Id int  `orm:"pk;auto"`
	Index int64  //区块编号
	Timestamp int64   //时间戳
	PreBlockHash string
	Hash string   //当前hash  实际中没有当前hash，通过下一个区块计算得出

	Data string   //区块内容
}
type BlockChain struct {
	Blocks []*Block
}

func init() {
	orm.RegisterDataBase("default","mysql","root:root@/chain?charset=utf8")
	orm.RegisterModel(new(Block),new(User))
	orm.RunSyncdb("default",false,true)
}
//创建新区块
func CreatNewBlock(pre Block,Data string) Block {
	var newBlock Block
	newBlock.Index= pre.Index+1
	newBlock.Timestamp=time.Now().Unix()
	newBlock.PreBlockHash= pre.Hash
	newBlock.Data=Data
	newBlock.Hash=newBlock.calHash()
	return newBlock
}
//计算Hash
func (b Block)calHash() string  {
	blockData := string(b.Index)+string(b.Timestamp)+b.PreBlockHash+b.Data
	//哈希计算  结果
	value := sha256.Sum256([]byte(blockData))
	//返回字节流转字符串
	calHashValue := hex.EncodeToString(value[:])
	return calHashValue
}
//头区块的特殊性，data记录特殊意义的事件
//只调用一次，也可以直接初始化一个头区块
/*func CreateFirstBlock() Block  {
	o := orm.NewOrm()
	newBlock := Block{}
	preBlock := Block{}
	preBlock.Index = -1
	preBlock.Hash=""
	newBlock =CreatNewBlock(preBlock,"目前BTC的市值1089.77亿,20181031")
	//写入数据库
	_,err := o.Insert(&newBlock)
	if err != nil {
		beego.Info("创世区块没写入")
		return newBlock
	}
	return  newBlock
}*/
//数据校验
//其他节点数据上链，存入数据库中之前校验
func IsValid(newBlock,oldBlock Block) bool {
	if newBlock.Index-1 !=oldBlock.Index{
		return false
	}
	if newBlock.PreBlockHash != oldBlock.Hash{
		return false
	}
	if newBlock.Hash != newBlock.calHash(){
		return false
	}
	return true
}
//数据上链
/*func (bc *Blockchian)SendData(data string)  {
	newBlock := GenerateNewBlock(*bc.Blocks[len(bc.Blocks)-1],data)
	bc.ApendBlock(&newBlock)  //调用区块上链
}
//区块上链
func (bc *Blockchian)ApendBlock(newBlock *Block)  {
	if len(bc.Blocks) == 0 {
		bc.Blocks = append(bc.Blocks,newBlock)
		return
	}
	if isValid(*newBlock,*bc.Blocks[len(bc.Blocks)-1])==true {
		bc.Blocks = append(bc.Blocks,newBlock)
	} else {
		log.Fatal("invalid block")
	}

}
*/