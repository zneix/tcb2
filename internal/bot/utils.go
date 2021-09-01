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

// SubEventDescriptions slice with all currently supported events, indexed by the values of SubEventType
// allowing to cast the index of a string in the slice to a corresponding SubEventType value
var SubEventDescriptions = []string{
	"when the game changes",          // SubEventTypeGame
	"when the title changes",         // SubEventTypeTitle
	"when the streamer goes live",    // SubEventTypeLive
	"when the streamer goes offline", // SubEventTypeOffline
}
