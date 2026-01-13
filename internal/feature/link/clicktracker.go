package link

// ClickTracker records a redirect (click) event
type ClickTracker interface {
	TrackClick(event ClickEvent)
}
