package Comandos

import (
	"PROYECTO_MIA/backend/Structs"
	"fmt"
	"strings"
)

// Mounted muestra todas las particiones montadas en memoria
func Mounted(entrada []string) string {
	var respuesta string

	// Verificar que no haya parámetros adicionales
	if len(entrada) > 1 {
		mensaje := "ERROR MOUNTED: Este comando no acepta parámetros"
		fmt.Println(mensaje)
		return mensaje
	}

	// Verificar si hay particiones montadas
	if len(Structs.Montadas) == 0 {
		mensaje := "INFO MOUNTED: No hay particiones montadas en el sistema"
		fmt.Println(mensaje)
		return mensaje
	}

	// Construir la respuesta con los IDs de todas las particiones montadas
	var ids []string
	for _, particion := range Structs.Montadas {
		ids = append(ids, particion.Id)
	}

	respuesta = "PARTICIONES MONTADAS: " + strings.Join(ids, ", ")
	fmt.Println(respuesta)
	return respuesta
}
