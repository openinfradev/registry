package taco

// Phase is taco registry build log phase
type Phase struct {
	Status   string
	StartSeq int
}

// PhasePulling is pulling phase
var PhasePulling = &Phase{
	Status:   "pulling",
	StartSeq: 0,
}

// PhaseBuilding is building phase
var PhaseBuilding = &Phase{
	Status:   "building",
	StartSeq: 100,
}

// PhasePushing is pushing phase
var PhasePushing = &Phase{
	Status:   "pushing",
	StartSeq: 10000,
}

// PhaseComplete is complete phase
var PhaseComplete = &Phase{
	Status:   "complete",
	StartSeq: 20000,
}

// PhaseError is error phase
var PhaseError = &Phase{
	Status:   "error",
	StartSeq: 99999,
}
