package rbac

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"go-zentao-task/pkg/db"
	"log"
)

var Enforcer *casbin.Enforcer

func Setup() {
	adapter, err := gormadapter.NewAdapterByDB(db.Orm)
	if err != nil {
		log.Fatalln(err)
	}

	Enforcer, err = casbin.NewEnforcer("docs/casbin/rbac_model.conf", adapter)
	if err != nil {
		log.Fatalln(err)
	}
}
