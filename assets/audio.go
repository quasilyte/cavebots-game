package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
)

func registerAudioResources(ctx *ge.Context) {
	audioResources := map[resource.AudioID]resource.AudioInfo{
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

		AudioBowShot1: {Path: "audio/bow_shot1.wav"},
		AudioBowShot2: {Path: "audio/bow_shot2.wav"},
		AudioBowShot3: {Path: "audio/bow_shot3.wav"},

		AudioWarriorHit1: {Path: "audio/warrior_hit1.wav", Volume: -0.15},
		AudioWarriorHit2: {Path: "audio/warrior_hit2.wav", Volume: -0.15},
		AudioWarriorHit3: {Path: "audio/warrior_hit3.wav", Volume: -0.15},
		AudioWarriorHit4: {Path: "audio/warrior_hit4.wav", Volume: -0.15},

		AudioGatlingShot: {Path: "audio/gatling_shot.wav", Volume: -0.1},
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
	default:
		return 1
	}
}

const (
	AudioNone resource.AudioID = iota

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

	AudioBowShot1
	AudioBowShot2
	AudioBowShot3

	AudioWarriorHit1
	AudioWarriorHit2
	AudioWarriorHit3
	AudioWarriorHit4

	AudioGatlingShot
)
