package kilt

type Kilt struct {
	definition LanguageInterface
}

func NewKilt(impl LanguageInterface) *Kilt {
	k := new(Kilt)
	k.definition = impl
	return k
}

func (k *Kilt) Build(info *TargetInfo) (*Build, error) {
	return k.definition.Build(info)
}

func (k *Kilt) Runtime(info *TargetInfo) (*Runtime, error) {
	return k.definition.Runtime(info)
}

func (k *Kilt) Task() (*Task, error) {
	return k.definition.Task()
}
