package streams

// RedisStreams represents the different types of streams that can be used to send messages between services
type RedisStreams int

// Const representing possible stream types
const (
	Unknown RedisStreams = iota
	JobControl
	ExperimentControl
	ConfigControl
)

func (r RedisStreams) String() string {
	streamTypes := [...]string{
		"Unknown",
		"job_control",
		"experiment_control",
		"config_control",
	}
	if int(r) < 0 || int(r) >= len(streamTypes) {
		return "Unknown"
	}
	return streamTypes[r]
}
