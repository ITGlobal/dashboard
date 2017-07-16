package app

import "github.com/itglobal/dashboard/tile"
import "time"

type tileObject struct {
	Type            tile.Type
	Size            tile.Size
	State           tile.State
	TitleText       string
	DescriptionText string
	StatusValue     int
	HasStatusValue  bool
	LastChangeTime  time.Time
	id              tile.ID
	container       *tileContainer
}

func newTileObject(container *tileContainer, id tile.ID) *tileObject {
	return &tileObject{id: id, container: container}
}

func (t *tileObject) ID() tile.ID {
	return t.id
}

func (t *tileObject) SetType(value tile.Type) {
	if t.Type != value {
		t.Type = value
		t.onChanged()
	}
}

func (t *tileObject) SetSize(value tile.Size) {
	if t.Size != value {
		t.Size = value
		t.onChanged()
	}
}

func (t *tileObject) SetState(value tile.State) {
	if t.State != value {
		t.State = value
		t.onChanged()
	}
}

func (t *tileObject) Clear() {
	t.SetTitleText("")
	t.SetDescriptionText("")
	t.SetNoStatusValue()
}

func (t *tileObject) SetTitleText(value string) {
	if t.TitleText != value {
		t.TitleText = value
		t.onChanged()
	}
}

func (t *tileObject) SetDescriptionText(value string) {
	if t.DescriptionText != value {
		t.DescriptionText = value
		t.onChanged()
	}
}

func (t *tileObject) SetStatusValue(value int) {
	if t.StatusValue != value || !t.HasStatusValue {
		t.StatusValue = value
		t.HasStatusValue = true
		t.onChanged()
	}
}

func (t *tileObject) SetNoStatusValue() {
	if t.HasStatusValue {
		t.StatusValue = 0
		t.HasStatusValue = false
		t.onChanged()
	}
}

func (t *tileObject) onChanged() {
	t.LastChangeTime = time.Now().UTC()
	t.container.Invalidate()
}

func (t *tileObject) getJSON() tileJSON {
	var obj tileJSON
	obj.ID = string(t.id)
	obj.LastChangeTime = t.LastChangeTime
	obj.Source = t.container.provider.Type()

	switch t.Type {
	case tile.TypeText:
		obj.Type = TypeText
	case tile.TypeTextStatus:
		obj.Type = TypeTextStatus
	case tile.TypeTextStatus2:
		obj.Type = TypeTextStatus2
	case tile.TypeTextStatusProgress:
		obj.Type = TypeTextStatusProgress
	default:
		obj.Type = TypeText
	}

	switch t.Size {
	case tile.Size1x:
		obj.Size = Size1x
	case tile.Size2x:
		obj.Size = Size2x
	case tile.Size4x:
		obj.Size = Size4x
	default:
		obj.Size = Size1x
	}

	switch t.State {
	case tile.StateDefault:
		obj.State = StateDefault
	case tile.StateSuccess:
		obj.State = StateSuccess
	case tile.StateIndeterminate:
		obj.State = StateIndeterminate
	case tile.StateWarning:
		obj.State = StateWarning
	case tile.StateError:
		obj.State = StateError
	default:
		obj.State = StateDefault
	}

	obj.TitleText = t.TitleText
	obj.DescriptionText = t.DescriptionText

	if t.HasStatusValue {
		obj.StatusValue = &t.StatusValue
	} else {
		obj.StatusValue = nil
	}

	return obj
}
