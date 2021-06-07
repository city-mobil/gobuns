package barber

type StatHost struct {
	ServerID   int
	FailsCount int
}

type Stats struct {
	Hosts []*StatHost
}
