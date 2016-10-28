package test

type Identification struct {
	Position string
}

type ConfigData struct {
	Names   []string
	Default Identification
	Mode    Identification
	Local   Identification
	Partial struct {
		One, Two, Three string
		Sub             struct {
			One, Two, Three string
		}
	}
}
