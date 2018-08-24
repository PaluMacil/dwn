package gamelibrary

type Providers struct {
	Houses HouseProvider
	Games  GameProvider
}

type HouseProvider interface {
}

type GameProvider interface {
}
