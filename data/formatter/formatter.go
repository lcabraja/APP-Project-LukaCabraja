package formatter

type DataFormatter interface {
	Format(value interface{}) (interface{}, error)
}
