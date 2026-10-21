package public

import (
	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

var DB *xorm.Engine
