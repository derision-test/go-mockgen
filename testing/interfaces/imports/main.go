package imports

import (
	"github.com/efritz/go-mockgen/testing/interfaces/localtypes"
	. "github.com/efritz/go-mockgen/testing/interfaces/localtypes"
	localtyp "github.com/efritz/go-mockgen/testing/interfaces/localtypes"
)

type Imports interface {
	GetX() localtypes.X
	GetY() localtypes.Y
	GetZ() localtypes.Z
}

type LocalImports interface {
	GetX() X
	GetY() Y
	GetZ() Z
}

type RenamedImports interface {
	GetX() localtyp.X
	GetY() localtyp.Y
	GetZ() localtyp.Z
}
