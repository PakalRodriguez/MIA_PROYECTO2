package Estructuras

type BloqueArchivo struct {
	B_content [64]byte
}

func NuevoBloqueArchivo() BloqueArchivo {
	var bloqueArchivo BloqueArchivo
	bloqueArchivo.B_content = [64]byte{}
	return bloqueArchivo
}
