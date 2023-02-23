package model

type Target interface {
	Hash() uint64
	Tags() Tags
	TUID() string
}

type TargetGroup interface {
	Targets() []Target
	Source() string
}

type Config struct {
	Tags  Tags
	Conf  string
	Stale bool
}
