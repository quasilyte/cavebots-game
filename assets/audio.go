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
)
