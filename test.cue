{
	id: int
	uid: string
	title?: string
	description?: string
	style: *"light" | "dark"
	timezone?: *"browser" | "utc"
	editable: bool | *true
	graphTooltip: int >= 0 <= 2 | *0
    refresh: string
	schemaVersion: int | *25
	version: string
}