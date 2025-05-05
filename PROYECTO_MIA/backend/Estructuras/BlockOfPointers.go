package Estructuras

type BloqueApuntadores struct {
	B_pointers [16]int32
}

func NuevoBloqueApuntadores() BloqueApuntadores {
	var bloqueApuntadores BloqueApuntadores
	for i := 0; i < len(bloqueApuntadores.B_pointers); i++ {
		bloqueApuntadores.B_pointers[i] = -1
	}
	return bloqueApuntadores
}
