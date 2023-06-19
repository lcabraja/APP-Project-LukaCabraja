package data

type NamedResource interface {
	GetName() string
	GetType() ResourceType
}
