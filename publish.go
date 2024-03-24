package oidcc

type Publish int

const (
	NoPublish Publish = iota
	SummaryPublish
	EverythingPublish
)

func (p Publish) String() string {
	switch p {
	case SummaryPublish:
		return "summary"
	case EverythingPublish:
		return "everything"
	default:
		return ""
	}
}
