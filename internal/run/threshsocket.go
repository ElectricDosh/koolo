package run

import (
	"github.com/hectorgimenez/d2go/pkg/data"
	"github.com/hectorgimenez/d2go/pkg/data/area"
	"github.com/hectorgimenez/d2go/pkg/data/npc"
	"github.com/hectorgimenez/koolo/internal/action"
	"github.com/hectorgimenez/koolo/internal/config"
)

type Threshsocket struct {
	baseRun
}

func (m Threshsocket) Name() string {
	return string(config.ThreshsocketRun)
}

func (m Threshsocket) BuildActions() []action.Action {
	actions := []action.Action{
		m.builder.WayPoint(area.CrystallinePassage), // Moving to starting point (Durance of Hate Level 2)
		m.builder.MoveToArea(area.ArreatPlateau),
		a.char.KillMonsterSequence(func(d game.Data) (data.UnitID, bool) {
			if m, found := d.Monsters.FindOne(npc.BloodBringer, data.MonsterTypeSuperUnique); found {
				return m.UnitID, true
			}
      

			return 0, false
		}, nil),
			a.builder.ItemPickup(false, 35),
		)
	}

	return
}
