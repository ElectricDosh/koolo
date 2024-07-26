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

func (t Threshsocket) BuildActions() []action.Action {
	actions := []action.Action{
		t.builder.WayPoint(area.CrystallinePassage), // Moving to starting point (Crystalline Passage)
		t.builder.MoveToArea(area.ArreatPlateau),
		t.char.KillMonsterSequence(func(d game.Data) (data.UnitID, bool) {
			if m, found := d.Monsters.FindOne(npc.BloodBringer, data.MonsterTypeSuperUnique); found {
				return m.UnitID, true
			}
      

			return 0, false
		}, nil),
		t.builder.ItemPickup(false, 35)
		
	}

	return
}
