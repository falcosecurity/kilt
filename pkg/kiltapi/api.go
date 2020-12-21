package kiltapi

import (
	"github.com/falcosecurity/kilt/pkg/hocon"
	"github.com/falcosecurity/kilt/pkg/kilt"
)

func NewKiltFromHocon(definition string) *kilt.Kilt {
	impl := hocon.NewKiltHocon(definition)
	return kilt.NewKilt(impl)
}

func NewKiltFromHoconWithConfig(definition string, config kilt.RecipeConfig) *kilt.Kilt {
	impl := hocon.NewKiltHoconWithConfig(definition, config)
	return kilt.NewKilt(impl)
}