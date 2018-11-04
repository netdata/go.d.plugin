package modules

//
//import "fmt"
//
//type trigger string
//
//const (
//	renewTrigger    trigger = "renew"
//	rebuildTrigger  trigger = "rebuild"
//	removeTrigger   trigger = "remove"
//	obsoleteTrigger trigger = "obsolete"
//)
//
//type state string
//
//const (
//	initialState        state = ""
//	newState            state = "new"
//	createdState        state = "created"
//	rebuiltState        state = "rebuilt"
//	recoveredState      state = "recovered"
//	markedObsoleteState state = "marked obsolete"
//	obsoletedState      state = "obsolete"
//	markedRemoveState   state = "marked remove"
//	markedDeleteState   state = "marked delete"
//)
//
//func (s state) String() string {
//	switch s {
//	default:
//		return string(s)
//	case initialState:
//		return "initial"
//	}
//}
//
//func (s state) dispatch(trigger trigger) state {
//	switch s {
//	case initialState:
//		switch trigger {
//		case renewTrigger, rebuildTrigger, obsoleteTrigger, removeTrigger:
//			return initialState
//		}
//	case newState:
//		switch trigger {
//		case renewTrigger:
//			return newState
//		case rebuildTrigger:
//			return rebuiltState
//		case obsoleteTrigger:
//			return markedObsoleteState
//		case removeTrigger:
//			return markedRemoveState
//		}
//	case createdState:
//		switch trigger {
//		case renewTrigger:
//			return newState
//		case rebuildTrigger:
//			return rebuiltState
//		case obsoleteTrigger:
//			return markedObsoleteState
//		case removeTrigger:
//			return markedRemoveState
//		}
//	case rebuiltState:
//		switch trigger {
//		case renewTrigger:
//			return rebuiltState
//		case rebuildTrigger:
//			return rebuiltState
//		case obsoleteTrigger:
//			return markedObsoleteState
//		case removeTrigger:
//			return markedRemoveState
//		}
//	case recoveredState:
//		switch trigger {
//		case renewTrigger:
//			return recoveredState
//		case rebuildTrigger:
//			return recoveredState
//		case obsoleteTrigger:
//			return obsoletedState
//		case removeTrigger:
//			return markedDeleteState
//		}
//	case markedObsoleteState:
//		switch trigger {
//		case obsoleteTrigger:
//			return markedObsoleteState
//		case removeTrigger:
//			return markedRemoveState
//		}
//	case obsoletedState:
//		switch trigger {
//		case renewTrigger:
//			return recoveredState
//		case rebuildTrigger:
//			return recoveredState
//		case obsoleteTrigger:
//			return obsoletedState
//		case removeTrigger:
//			return markedDeleteState
//		}
//	case markedRemoveState, markedDeleteState:
//	}
//
//	panic(fmt.Sprintf("%s: wrong trigger (%s)", s, trigger))
//}
