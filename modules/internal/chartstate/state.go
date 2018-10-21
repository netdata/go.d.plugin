package chartstate

import "fmt"

type Trigger string

const (
	Renew    Trigger = "renew"
	Rebuild  Trigger = "rebuild"
	Remove   Trigger = "remove"
	Obsolete Trigger = "obsolete"
)

type State string

const (
	Initial        State = ""
	New            State = "new"
	Created        State = "created"
	Rebuilt        State = "rebuilt"
	Recovered      State = "recovered"
	MarkedObsolete State = "marked obsolete"
	Obsoleted      State = "obsolete"
	MarkedRemove   State = "marked remove"
	MarkedDelete   State = "marked delete"
)

func (s State) String() string {
	switch s {
	default:
		return string(s)
	case Initial:
		return "initial"
	}
}

func (s State) Dispatch(trigger Trigger) State {
	switch s {
	case Initial:
		switch trigger {
		case Renew, Rebuild, Obsolete, Remove:
			return Initial
		}
	case New:
		switch trigger {
		case Renew:
			return New
		case Rebuild:
			return Rebuilt
		case Obsolete:
			return MarkedObsolete
		case Remove:
			return MarkedRemove
		}
	case Created:
		switch trigger {
		case Renew:
			return New
		case Rebuild:
			return Rebuilt
		case Obsolete:
			return MarkedObsolete
		case Remove:
			return MarkedRemove
		}
	case Rebuilt:
		switch trigger {
		case Renew:
			return Rebuilt
		case Rebuild:
			return Rebuilt
		case Obsolete:
			return MarkedObsolete
		case Remove:
			return MarkedRemove
		}
	case Recovered:
		switch trigger {
		case Renew:
			return Recovered
		case Rebuild:
			return Recovered
		case Obsolete:
			return Obsoleted
		case Remove:
			return MarkedDelete
		}
	case MarkedObsolete:
		switch trigger {
		case Obsolete:
			return MarkedObsolete
		case Remove:
			return MarkedRemove
		}
	case Obsoleted:
		switch trigger {
		case Renew:
			return Recovered
		case Rebuild:
			return Recovered
		case Obsolete:
			return Obsoleted
		case Remove:
			return MarkedDelete
		}
	case MarkedRemove, MarkedDelete:
	}

	panic(fmt.Sprintf("%s: wrong trigger (%s)", s, trigger))
}
