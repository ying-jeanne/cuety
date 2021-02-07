{
	id: int
	uid: string
	title?: string
	description?: string
	style: *"light" | "dark"
	timezone?: *"browser" | "utc"
	editable: bool | *true
	graphTooltip: >= 0 <= 2 | *0 // to make the API works correctly here int needs to be removed, ticket opened https://github.com/cuelang/cue/issues/742
    refresh: string
	schemaVersion: int | *25
	version: string
	// Time range for dashboard, e.g. last 6 hours, last 7 days, etc
	time?: {
		from: string | *"now-6h"
		to:   string | *"now"
	}
	timepicker?: {
		// Whether timepicker is collapsed or not.
		collapse: bool | *false
		// Whether timepicker is enabled or not.
		enable: bool | *true
		// Whether timepicker is visible or not.
		hidden: bool | *false
		// Selectable intervals for auto-refresh.
		refresh_intervals: [...string] | *["5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"]
	},
	templating?: list: [...{}]
	// Annotations.
	annotations?: list: [...{
		builtIn: int | *0
		// Datasource to use for annotation.
		datasource: string
		// Whether annotation is enabled.
		enable?: bool | *true
		// Whether to hide annotation.
		hide?: bool | *false
		// Annotation icon color.
		iconColor?: string
		// Name of annotation.
		name?: string
		// Query for annotation data.
		rawQuery: string
		showIn:   int | *0
	}] | *[]
}