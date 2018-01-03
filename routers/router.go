package routers

import (
	"github.com/edwardhey/flow/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
