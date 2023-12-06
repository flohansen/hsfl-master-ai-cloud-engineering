package config

type Target struct {
}

type Configuration struct {
	Users     int // Number of users to simulate or workers to run concurrently.
	Requests  int // Number of requests to send per user.
	RateLimit int // Number of requests to send per second.
	Duration  int // Duration maximum of the test in seconds.
	RampUp    int // Duration of the ramp up in seconds.
	Targets   []Target
}
