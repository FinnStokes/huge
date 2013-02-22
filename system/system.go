// Package system manages individual systems that operate periodically on entities.
package system

import (
	"time"

	"github.com/FinnStokes/huge/entity"
)

// System is an interface that must be satisfied by any system that is added to the system manager.
type System interface {
	Update(dt time.Duration, entities *entity.Manager)
	Draw(entities *entity.Manager)
}
