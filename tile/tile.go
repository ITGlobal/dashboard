package tile

// Size is a size of a tile
type Size int

const (
	// Size1x defines smallest 1x1 tile
	Size1x Size = iota
	// Size2x defines 2x1 (horizontal) tile
	Size2x
	// Size4x defines 2x2 tile
	Size4x
)

// State is a state of a tile
type State int

const (
	// StateDefault is a default "Gray" tile state
	StateDefault State = iota
	// StateSuccess defines "Green" tile state
	StateSuccess
	// StateIndeterminate defines "Cyan" tile state
	StateIndeterminate
	// StateWarning defines "Yellow" tile state
	StateWarning
	// StateError defines "Red" tile state
	StateError
)

// Type is a type of a tile
type Type int

const (
	// TypeText is a tile with a text only
	TypeText Type = iota
	// TypeTextStatus is a tile with header and text
	TypeTextStatus
	// TypeTextStatus2 is a tile with header and large text
	TypeTextStatus2
	// TypeTextStatusProgress is a tile with header, text and progress bar with label
	TypeTextStatusProgress
)

// Manager defines methods to manage tiles
type Manager interface {
	BeginUpdate(provider Provider) Updater
}

// ID is a tile identifier
type ID string

// Updater defines methods to add, update or remove tiles
type Updater interface {
	GetTiles() []Tile
	AddOrUpdateTile(id ID) Tile
	RemoveTile(id ID)
	EndUpdate()
}

// Tile defines methods to modify properties of one specific tile
type Tile interface {
	ID() ID

	SetType(value Type)
	SetSize(size Size)
	SetState(state State)

	Clear()
	SetTitleText(value string)
	SetDescriptionText(value string)
	SetStatusValue(value int)
	SetNoStatusValue()
}
