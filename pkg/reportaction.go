package pkg

const (
	ReportActionUnknown = 0
	ReportActionBan     = 1
	ReportActionTimeout = 2
	ReportActionDismiss = 3
	ReportActionUndo    = 4
)

func GetReportActionByName(actionName string) uint8 {
	switch actionName {
	case "ban":
		return ReportActionBan
	case "timeout":
		return ReportActionTimeout
	case "dismiss":
		return ReportActionDismiss
	case "undo":
		return ReportActionUndo
	default:
		return ReportActionUnknown
	}
}

func GetReportActionName(action uint8) string {
	switch action {
	case ReportActionBan:
		return "ban"
	case ReportActionTimeout:
		return "timeout"
	case ReportActionDismiss:
		return "dismiss"
	case ReportActionUndo:
		return "undo"
	default:
		return "unknown"
	}
}
