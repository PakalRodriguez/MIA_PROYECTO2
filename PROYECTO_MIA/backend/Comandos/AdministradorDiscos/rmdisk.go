package Comandos

import (
	"fmt"
	"os"
	"strings"
)

func Rmdisk(parametros []string) {
	fmt.Println("RMDISK")

	// Verificar que se haya recibido el parámetro correcto
	if len(parametros) < 2 {
		fmt.Println("RMDISK Error: Falta el parámetro -path")
		return
	}

	// Limpiar espacios y dividir en clave-valor
	tmp2 := strings.TrimSpace(parametros[1])
	tmp := strings.Split(tmp2, "=")

	if len(tmp) != 2 {
		fmt.Println("RMDISK Error: Parámetro incorrecto", tmp[0])
		return
	}

	// Verificar que el parámetro sea path
	if strings.ToLower(tmp[0]) == "path" {
		path := tmp[1] // Se usa directamente el path ingresado

		// Verificar si el archivo existe
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Println("RMDISK Error: El archivo especificado no existe ->", path)
			return
		}

		// Eliminar el archivo sin confirmación
		if err := os.Remove(path); err != nil {
			fmt.Println("RMDISK Error: No se pudo eliminar el disco en", path)
			return
		}

		fmt.Println("RMDISK: Disco eliminado correctamente ->", path)
	} else {
		fmt.Println("RMDISK Error: Parámetro desconocido", tmp[0])
	}
}
