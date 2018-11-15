package terminal

import (
	"github.com/Hrubon/pivot/model"
)

type Terminal interface {
	Draw(rlist *model.RouteList) error
}
