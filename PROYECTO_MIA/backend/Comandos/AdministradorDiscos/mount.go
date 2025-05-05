package Comandos

import (
	"PROYECTO_MIA/backend/Herramientas"
	"PROYECTO_MIA/backend/Structs"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Para su proyecto, la letra es el disco y el numero es la particion
// var Pmontaje []Structs.Mount//GUarda en Ram las particones montadas
func Mount(entrada []string) {
	var name string //Nobre de la particion a montar
	var path string //Path del Disco
	paramC := true

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro, " ")
		valores := strings.Split(tmp, "=")

		if len(valores) != 2 {
			fmt.Println("ERROR MOUNT, valor desconocido de parametros ", valores[1])
			return //Finaliza comando
		}

		//******************* PATH *************
		if strings.ToLower(valores[0]) == "path" {
			path = strings.ReplaceAll(valores[1], "\"", "")
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				fmt.Println("ERROR MOUNT: El disco no existe")
				paramC = false
				break // Terminar el bucle porque encontramos un nombre único
			}
			//********************  NAME *****************
		} else if strings.ToLower(valores[0]) == "name" {
			// Eliminar comillas
			name = strings.ReplaceAll(valores[1], "\"", "")
			// Eliminar espacios en blanco al final
			name = strings.TrimSpace(name)

			//******************* ERROR EN LOS PARAMETROS *************
		} else {
			fmt.Println("ERROR MOUNT: Parametro desconocido: ", valores[0])
			paramC = false
			break //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	if paramC {
		if path != "" {
			if name != "" {
				// Abrir y cargar el disco
				disco, err := Herramientas.OpenFile(path)
				if err != nil {
					fmt.Println("ERROR NO SE PUEDE LEER EL DISCO ")
					return
				}

				//Se crea un mbr para cargar el mbr del disco
				var mbr Structs.MBR
				//Guardo el mbr leido
				if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
					return
				}

				// cerrar el archivo del disco
				defer disco.Close()

				montar := true // para guardar error si no se puede montar
				for i := 0; i < 4; i++ {
					nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					if nombre == name {
						montar = false
						if string(mbr.Partitions[i].Type[:]) != "E" {
							if string(mbr.Partitions[i].Status[:]) != "A" {
								var id string
								var nuevaLetra byte = 'A' // A
								contador := 1
								modificada := false //para saber si ya hay una particion montada en el disco

								//Verifica si el path existe dentro de las particiones montadas para calcular la nueva letra
								for k := 0; k < len(Structs.Pmontaje); k++ {
									if Structs.Pmontaje[k].MPath == path {
										//MOdifica el struct
										Structs.Pmontaje[k].Cont = Structs.Pmontaje[k].Cont + 1
										contador = int(Structs.Pmontaje[k].Cont)
										nuevaLetra = Structs.Pmontaje[k].Letter
										modificada = true
										break
									}
								}

								if !modificada {
									if len(Structs.Pmontaje) > 0 {
										nuevaLetra = Structs.Pmontaje[len(Structs.Pmontaje)-1].Letter + 1
									}
									Structs.AddPathM(path, nuevaLetra, 1)
								}

								id = "57" + strconv.Itoa(contador) + string(nuevaLetra) //Id de particion
								fmt.Println("ID:  Letra ", string(nuevaLetra), " cont ", contador)
								//Agregar al struct de montadas
								Structs.AddMontadas(id, path)

								//TODO modificar la particion que se va a montar
								copy(mbr.Partitions[i].Status[:], "A")
								copy(mbr.Partitions[i].Id[:], id)
								mbr.Partitions[i].Correlative = int32(contador)

								//sobreescribir el mbr para guardar los cambios
								if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
									return
								}
								fmt.Println("Particion con nombre ", name, " montada correctamente. ID: ", id)
							} else {
								fmt.Println("ERROR MOUNT. ESTA PARTICION YA FUE MONTADA PREVIAMENTE")
								return
							}
						} else {
							fmt.Println("ERROR MOUNT. No se puede montar una particion extendida")
							return
						}
					}
				}

				if montar {
					fmt.Println("ERROR MOUNT. No se pudo montar la particion ", name)
					fmt.Println("ERROR MOUNT. No se encontro la particion")
					return
				}

			} else {
				fmt.Println("ERROR: FALTA NAME  EN MOUNT")
			}
		} else {
			fmt.Println("ERROR: FALTA PARAMETRO PATH EN MOUNT")
		}
	}
}
