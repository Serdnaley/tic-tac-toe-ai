package map_builder

type MapStatus int

const (
	MapStatusNotStarted MapStatus = iota
	MapStatusBuilding
	MapStatusReady
)

func (ms MapStatus) String() string {
	switch ms {
	case MapStatusNotStarted:
		return "NOT_STARTED"
	case MapStatusBuilding:
		return "BUILDING"
	case MapStatusReady:
		return "READY"
	default:
		return "UNKNOWN"
	}
}
