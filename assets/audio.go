package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
)

const (
	SoundGroupEffect uint = iota
	SoundGroupMusic
)

func VolumeMultiplier(level int) float64 {
	switch level {
	case 1:
		return 0.01
	case 2:
		return 0.15
	case 3:
		return 0.45
	case 4:
		return 0.8
	case 5:
		return 1.0
	default:
		return 0
	}
}

func registerAudioResources(ctx *ge.Context) {
	audioResources := map[resource.AudioID]resource.AudioInfo{
		AudioMusic1: {Path: "audio/otomata_track.ogg", Group: SoundGroupMusic},

		AudioUnitAck1: {Path: "audio/unit_ack1.wav"},
		AudioUnitAck2: {Path: "audio/unit_ack2.wav"},
		AudioUnitAck3: {Path: "audio/unit_ack3.wav"},
		AudioUnitAck4: {Path: "audio/unit_ack4.wav"},
		AudioUnitAck5: {Path: "audio/unit_ack5.wav"},
		AudioUnitAck6: {Path: "audio/unit_ack6.wav"},
		AudioUnitAck7: {Path: "audio/unit_ack7.wav"},
		AudioUnitAck8: {Path: "audio/unit_ack8.wav"},

		AudioPatrolLaser1: {Path: "audio/patrol_laser1.wav", Volume: -0.15},
		AudioPatrolLaser2: {Path: "audio/patrol_laser2.wav", Volume: -0.15},
		AudioPatrolLaser3: {Path: "audio/patrol_laser3.wav", Volume: -0.15},
		AudioPatrolLaser4: {Path: "audio/patrol_laser4.wav", Volume: -0.15},

		AudioVanguardShot1: {Path: "audio/vanguard_shot1.wav", Volume: -0.05},
		AudioVanguardShot2: {Path: "audio/vanguard_shot2.wav", Volume: -0.05},
		AudioVanguardShot3: {Path: "audio/vanguard_shot3.wav", Volume: -0.05},

		AudioLaserTurretShot1: {Path: "audio/lasertower1.wav", Volume: -0.2},
		AudioLaserTurretShot2: {Path: "audio/lasertower2.wav", Volume: -0.2},
		AudioLaserTurretShot3: {Path: "audio/lasertower3.wav", Volume: -0.2},

		AudioBowShot1: {Path: "audio/bow_shot1.wav"},
		AudioBowShot2: {Path: "audio/bow_shot2.wav"},
		AudioBowShot3: {Path: "audio/bow_shot3.wav"},

		AudioWarriorHit1: {Path: "audio/warrior_hit1.wav", Volume: -0.15},
		AudioWarriorHit2: {Path: "audio/warrior_hit2.wav", Volume: -0.15},
		AudioWarriorHit3: {Path: "audio/warrior_hit3.wav", Volume: -0.15},
		AudioWarriorHit4: {Path: "audio/warrior_hit4.wav", Volume: -0.15},

		AudioGatlingShot: {Path: "audio/gatling_shot.wav", Volume: -0.1},

		AudioExplosion1: {Path: "audio/explosion1.wav", Volume: -0.2},
		AudioExplosion2: {Path: "audio/explosion2.wav", Volume: -0.2},
		AudioExplosion3: {Path: "audio/explosion3.wav", Volume: -0.2},
		AudioExplosion4: {Path: "audio/explosion4.wav", Volume: -0.2},

		AudioResourceAdded:     {Path: "audio/resource_added.wav", Volume: -0.2},
		AudioUnitReady:         {Path: "audio/unit_ready.wav", Volume: -0.35},
		AudioDig:               {Path: "audio/dig.wav", Volume: -0.1},
		AudioRepair:            {Path: "audio/repair.wav", Volume: -0.5},
		AudioBuildingPlaced:    {Path: "audio/building_placed.wav", Volume: +0.5},
		AudioProductionStarted: {Path: "audio/production_started.wav", Volume: -0.35},
	}

	for id, res := range audioResources {
		ctx.Loader.AudioRegistry.Set(id, res)
		ctx.Loader.LoadAudio(id)
	}
}

func NumSamples(a resource.AudioID) int {
	switch a {
	case AudioUnitAck1:
		return 8
	case AudioPatrolLaser1:
		return 4
	case AudioBowShot1:
		return 3
	case AudioWarriorHit1:
		return 4
	case AudioExplosion1:
		return 4
	case AudioLaserTurretShot1:
		return 3
	default:
		return 1
	}
}

const (
	AudioNone resource.AudioID = iota

	AudioMusic1

	AudioUnitAck1
	AudioUnitAck2
	AudioUnitAck3
	AudioUnitAck4
	AudioUnitAck5
	AudioUnitAck6
	AudioUnitAck7
	AudioUnitAck8

	AudioPatrolLaser1
	AudioPatrolLaser2
	AudioPatrolLaser3
	AudioPatrolLaser4

	AudioVanguardShot1
	AudioVanguardShot2
	AudioVanguardShot3

	AudioLaserTurretShot1
	AudioLaserTurretShot2
	AudioLaserTurretShot3

	AudioBowShot1
	AudioBowShot2
	AudioBowShot3

	AudioWarriorHit1
	AudioWarriorHit2
	AudioWarriorHit3
	AudioWarriorHit4

	AudioGatlingShot

	AudioExplosion1
	AudioExplosion2
	AudioExplosion3
	AudioExplosion4

	AudioResourceAdded
	AudioUnitReady
	AudioDig
	AudioRepair
	AudioBuildingPlaced
	AudioProductionStarted
)
