package Estructuras

type Journaling struct {
	contenido [15]ContenidoJournaling
}

type ContenidoJournaling struct {
	Operacion [8]byte
	Path      [80]byte
	Contenido [80]byte
	Time      [16]byte
}

func NuevoJournaling() Journaling {
	var journaling Journaling
	journaling.contenido = [15]ContenidoJournaling{}
	return journaling
}
