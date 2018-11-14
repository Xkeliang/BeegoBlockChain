package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"crypto/sha256"
	"github.com/beego"
	"bytes"
	"math/big"
	"math"
	"fmt"
	"encoding/binary"
	"os"
	"encoding/hex"
)

type Block struct {
	//Id int
	Index int64  `orm:"pk;auto"`//区块编号
	Version int64
	Timestamp int64   //时间戳
	PreBlockHash string
	Hash string  //当前hash  实际中没有当前hash，通过下一个区块计算得出
	MerkelRoot string
	Bit int64
	Nonce int64

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
//设置hash
func (block *Block) calHash() string {
	temp := [][]byte{
		IntToByte(block.Version),
		[]byte(block.PreBlockHash),
		[]byte(block.MerkelRoot),
		IntToByte(block.Timestamp),
		IntToByte(block.Bit),
		IntToByte(block.Nonce),
		[]byte(block.Data),
	}
	data := bytes.Join(temp,[]byte{})
	hash :=sha256.Sum256(data)
	beego.Info(hash[:])
	block.Hash=hex.EncodeToString(hash[:])

	return block.Hash
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

//头区块的特殊性，data记录特殊意义的事件
//只调用一次，也可以直接初始化一个头区块
func CreateFirstBlock() Block  {
	o := orm.NewOrm()
	newBlock := Block{}
	preBlock := Block{}
	preBlock.Index = -1
	preBlock.Hash=""
	data := "目前BTC的市值1089.77亿,20181031"
	newBlock =CreatNewBlock(preBlock,data)
	//写入数据库
	beego.Info(newBlock.Hash)
	_,err := o.Insert(&newBlock)
	if err != nil {
		beego.Info("创世区块没写入",err)
		return newBlock
	}
	return  newBlock
}
//数据校验
//其他节点数据上链，存入数据库中之前校验
func IsValid(newBlock,oldBlock Block) bool {
	if newBlock.Index-1 !=oldBlock.Index{
		return false
	}
	if string(newBlock.PreBlockHash) != string(oldBlock.Hash){
		return false
	}
	if string(newBlock.Hash) != string(newBlock.calHash()){
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
const targetBits  =16

type ProofOfWork struct {
	block *Block
	target *big.Int
}

func (pow *ProofOfWork)PrepareData(nonce int64) []byte  {
	block := pow.block
	temp := [][]byte{
		IntToByte(block.Version),
		[]byte(block.PreBlockHash),
		[]byte(block.MerkelRoot),
		IntToByte(block.Timestamp),
		IntToByte(targetBits),
		IntToByte(nonce),
		[]byte(block.Data),
	}
	data :=bytes.Join(temp,[]byte{})
	return data
}

func (pow *ProofOfWork)Run()(int64,[]byte){
	var hash [32]byte
	var nonce int64=0
	var hashInt big.Int
	for nonce < math.MaxInt64{
		data := pow.PrepareData(nonce)

		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("found hash:%x,nonce:%d\n",hash,nonce)
			break
		}else {
			nonce++
		}
	}
	return nonce,hash[:]
}

func NewProofOfWork(block *Block)*ProofOfWork  {
	target := big.NewInt(1)
	target.Lsh(target,uint(256-targetBits))
	pow := ProofOfWork{target:target,block:block}
	return &pow
}
//int转byte
func IntToByte(num int64)[]byte  {
	var buffer bytes.Buffer
	err := binary.Write(&buffer,binary.BigEndian,num)
	CheckErr("IntToByte",err)
	return  buffer.Bytes()
}

func CheckErr(pos string,err error)  {
	if err != nil {
		fmt.Println("pos error=",pos,err)
		os.Exit(1)
	}
}
