package drawing

import (
	"github.com/Hrubon/pivot/model"
)

type Drawer interface {
	Draw(rlist *model.RouteList) error
}
