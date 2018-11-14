package sources;

import (
	"github.com/Hrubon/pivot/model"
)

type Source interface {
	GetRoutes() (*model.RouteList, error)
}
