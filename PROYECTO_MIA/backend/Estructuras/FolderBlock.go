package Estructuras

type BloqueCarpeta struct {
	B_content [4]Contenido
}

type Contenido struct {
	B_name  [12]byte
	B_inodo int32
}

func NuevoBloqueCarpeta() BloqueCarpeta {
	var bloqueCarpeta BloqueCarpeta
	for i := 0; i < len(bloqueCarpeta.B_content); i++ {
		bloqueCarpeta.B_content[i] = NuevoContenido()
	}
	return bloqueCarpeta
}

func NuevoContenido() Contenido {
	var contenido Contenido
	contenido.B_name = [12]byte{}
	contenido.B_inodo = -1
	return contenido
}
