package routers

import (
	"beegoBlockChain/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	//beego.Router("/", &controllers.MainController{})
    beego.Router("/send",&controllers.SendController{})
	beego.Router("/reg",&controllers.MainController{},"Get:ShowReg;Post:HandlerReg")
}
