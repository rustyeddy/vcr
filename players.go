package main

// Players are the playa's of the redeye network!
var (
	players map[string]interface{} = make(map[string]interface{}, 3)
)

func AddPlayer(name string) {
	players[name] = 1
}
