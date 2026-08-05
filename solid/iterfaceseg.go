package solid

type BEWorker interface {
	Code()
	AttendMeetings()
	DeployServers()
}

type UIWorker interface {
	DesignUI()
}

type BE struct {
}

func (b *BE) Code() {
	// Implementation of coding logic
}

func (b *BE) AttendMeetings() {
	// Implementation of attending meetings logic
}

func (b *BE) DeployServers() {
	// Implementation of deploying servers logic
}

type UI struct {
}

func (u *UI) DesignUI() {
	// Implementation of designing UI logic
}
