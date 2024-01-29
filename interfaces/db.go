package interfaces

type IDatabase interface {
	CreateDocument()
	GetDocument()
	UpdateDocument()
	DeleteDocument()
}
