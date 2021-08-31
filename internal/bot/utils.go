package bot

// ParseSubEventType tries to convert value to a ChannelEvent event
func ParseSubEventType(value string) (valid bool, event SubEventType) {
	switch value {
	case "game":
		return true, SubEventTypeGame
	case "title":
		return true, SubEventTypeTitle
	case "live":
		return true, SubEventTypeLive
	case "offline":
		return true, SubEventTypeOffline
	case "partnered":
		return true, SubEventTypePartnered
	default:
		return false, SubEventTypeInvalid
	}
}
