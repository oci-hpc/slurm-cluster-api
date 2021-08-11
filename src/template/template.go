package template

type SlurmTemplate struct {
	Id          int
	Name        string
	Body        string
	Keys        []TemplateKey
	Description string
	IsPublished bool
	Version     int
	OriginalId  int
}

type TemplateKey struct {
	Id          int
	Key         string
	Type        string
	Description string
	TemplateId  int
	Picklist    []string
}
