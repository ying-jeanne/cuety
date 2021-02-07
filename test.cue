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
}