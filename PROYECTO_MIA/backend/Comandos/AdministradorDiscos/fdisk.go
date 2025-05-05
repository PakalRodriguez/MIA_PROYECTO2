package Comandos

import (
	"PROYECTO_MIA/backend/Herramientas"
	"PROYECTO_MIA/backend/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Fdisk(parametros []string) {
	fmt.Println("FDISK")
	//PARAMETROS: -size -unit -path -type -fit -name
	var size int    //obligatorio si es creacion
	var path string //obligatorio (es el "path", es una letra nombre de la particion, path ya esta fijado)
	var name string //obligatorio Nombre de la particion
	unit := 1024    //opcional /valor por defecto en KB por eso es 1024
	typee := "P"    //opcional Valores: P, E, L
	fit := "W"      //opcional valores para fit: f, w, b

	var opcion int        // 0 -> crear; 1 -> add; 2 -> delete (por defecto es 0 = CREAR)
	paramC := true        //Para validar que los parametros cumplen con los requisitos
	sizeInit := false     //Sirve para saber si se inicializo size (por si no viniera el parametro por ser opcional) false -> no inicializado
	var sizeValErr string //Para reportar el error si no se pudo convertir a entero el size

	//mismo proceso que el mkdisk para manejar parametros
	for _, parametro := range parametros[1:] {
		//quito los espacios en blano despues de cada parametro
		tmp2 := strings.TrimRight(parametro, " ")
		tmp := strings.Split(tmp2, "=")

		//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
		if len(tmp) != 2 {
			fmt.Println("FDISK Error: Valor desconocido del parametro ", tmp[0])
			paramC = false
			break
		}

		//SIZE
		if strings.ToLower(tmp[0]) == "size" {
			sizeInit = true
			var err error
			size, err = strconv.Atoi(tmp[1]) //se convierte el valor en un entero
			if err != nil {
				sizeValErr = tmp[1] //guarda para el reporte del error si es necesario validar size
			}

			//PATH
		} else if strings.ToLower(tmp[0]) == "path" {
			//homonimo al path
			path = tmp[1]
			nombreDisco := strings.Split(path, "/")
			disco := nombreDisco[len(nombreDisco)-1]
			//Se valida si existe el disco ingresado
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				fmt.Println("FDISK Error: El disco ", disco, " no existe")
				paramC = false
				break // Terminar el bucle porque encontramos un nombre único
			}

			//NAME
		} else if strings.ToLower(tmp[0]) == "name" {
			// Eliminar comillas
			name = strings.ReplaceAll(tmp[1], "\"", "")
			// Eliminar espacios en blanco al final
			name = strings.TrimSpace(name)

			//UNIT
		} else if strings.ToLower(tmp[0]) == "unit" {
			//k ya esta predeterminado
			if strings.ToLower(tmp[1]) == "b" {
				//asigno el valor del parametro en su respectiva variable
				unit = 1
			} else if strings.ToLower(tmp[1]) == "m" {
				unit = 1048576 //1024*1024
			} else if strings.ToLower(tmp[1]) != "k" {
				fmt.Println("FDISK Error en -unit. Valores aceptados: b, k, m. ingreso: ", tmp[1])
				paramC = false
				break
			}

			//TYPE
		} else if strings.ToLower(tmp[0]) == "type" {
			//p esta predeterminado
			if strings.ToLower(tmp[1]) == "e" {
				typee = "E"
			} else if strings.ToLower(tmp[1]) == "l" {
				typee = "L"
			} else if strings.ToLower(tmp[1]) != "p" {
				fmt.Println("FDISK Error en -type. Valores aceptados: e, l, p. ingreso: ", tmp[1])
				paramC = false
				break
			}

			//FIT
		} else if strings.ToLower(tmp[0]) == "fit" {
			//Si el ajuste es BF (best fit)
			if strings.ToLower(tmp[1]) == "bf" {
				//asigno el valor del parametro en su respectiva variable
				fit = "B"
				//Si el ajuste es WF (worst fit)
			} else if strings.ToLower(tmp[1]) == "ff" {
				//asigno el valor del parametro en su respectiva variable
				fit = "F"
				//Si el ajuste es ff ya esta definido por lo que si es distinto es un error
			} else if strings.ToLower(tmp[1]) != "wf" {
				fmt.Println("FDISK Error en -fit. Valores aceptados: BF, FF o WF. ingreso: ", tmp[1])
				paramC = false
				break
			}

			//ERROR EN LOS PARAMETROS LEIDOS
		} else {
			fmt.Println("FDISK Error: Parametro desconocido ", tmp[0])
			paramC = false
			break //por si en el camino reconoce algo invalido de una vez se sale
		}

	}

	//Si va a crear una particion verificar el size
	if opcion == 0 && paramC {
		if sizeInit { //Si viene el parametro size
			if sizeValErr == "" { //Si es un numero (si es numero la variable sizeValErr sera una cadena vacia)
				if size <= 0 { //se valida que sea mayor a 0 (positivo)
					fmt.Println("FDISK Error: -size debe ser un valor positivo mayor a cero (0). se leyo ", size)
					paramC = false
				}
			} else { //Si sizeValErr es una cadena (por lo que no se pudo dar valor a size)
				fmt.Println("FDISK Error: -size debe ser un valor numerico. se leyo ", sizeValErr)
				paramC = false
			}
		} else { //Si no viene el parametro size
			fmt.Println("FDISK Error: No se encuentra el parametro -size")
			paramC = false
		}
	}

	//si todos los parametros son correctos
	if paramC {
		if path != "" && name != "" {
			// Abrir y cargar el disco
			filepath := path
			disco, err := Herramientas.OpenFile(filepath)
			if err != nil {
				fmt.Println("FDisk Error: No se pudo leer el disco")
				return
			}

			//Se crea un mbr para cargar el mbr del disco
			var mbr Structs.MBR
			//Guardo el mbr leido
			if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
				return
			}

			//CREAR (opcion: 0 -> crear; 1 -> add; 2 -> delete)
			if opcion == 0 {

				//Si la particion es tipo extendida validar que no exista alguna extendida
				isPartExtend := false //Indica si se puede usar la particion extendida
				isName := true        //Valida si el nombre no se repite (true no se repite)
				if typee == "E" {
					for i := 0; i < 4; i++ {
						tipo := string(mbr.Partitions[i].Type[:])
						//fmt.Println("tipo ", tipo)
						if tipo != "E" {
							isPartExtend = true
						} else {
							isPartExtend = false
							isName = false //Para que ya no evalue el nombre ni intente hacer nada mas
							fmt.Println("FDISK Error. Ya existe una particion extendida")
							fmt.Println("FDISK Error. No se puede crear la nueva particion con nombre: ", name)
							break
						}
					}
				}

				//verificar si  el nombre existe en las particiones primarias o extendida
				if isName {
					for i := 0; i < 4; i++ {
						nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
						if nombre == name {
							isName = false
							fmt.Println("FDISK Error. Ya existe la particion : ", name)
							fmt.Println("FDISK Error. No se puede crear la nueva particion con nombre: ", name)
							break
						}
					}
				}

				//verificar si existe en las logicas

				//INGRESO DE PARTICIONES PRIMARIAS Y/O EXTENDIDA (SIN LOGICAS)
				sizeNewPart := size * unit //Tamaño de la nueva particion (tamaño * unidades)
				guardar := false           //Indica si se debe guardar la particion, es decir, escribir en el disco
				var newPart Structs.Partition
				if (typee == "P" || isPartExtend) && isName { //para que  isPartExtend sea true, typee tendra que ser "E"
					sizeMBR := int32(binary.Size(mbr)) //obtener el tamaño del mbr (el que ocupa fisicamente: 165)
					//Para manejar los demas ajustes hacer un if del fit para llamar a la funcion adecuada
					//F = primer ajuste; B = mejor ajuste; else -> peor ajuste

					//INSERTAR PARTICION (Primer ajuste)
					mbr, newPart = primerAjuste(mbr, typee, sizeMBR, int32(sizeNewPart), name, fit) //int32(sizeNewPart) es para castear el int a int32 que es el tipo que tiene el atributo en el struct Partition
					guardar = newPart.Size != 0

					//escribimos el MBR en el archivo. Lo que no se llegue a escribir en el archivo (aqui) se pierde, es decir, los cambios no se guardan
					if guardar {
						//sobreescribir el mbr
						if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {
							return
						}

						//Se agrega el ebr de la particion extendida en el disco
						if isPartExtend {
							var ebr Structs.EBR
							ebr.Start = newPart.Start
							ebr.Next = -1
							if err := Herramientas.WriteObject(disco, ebr, int64(ebr.Start)); err != nil {
								return
							}
						}
						//para verificar que lo guardo
						var TempMBR2 Structs.MBR
						// Read object from bin file
						if err := Herramientas.ReadObject(disco, &TempMBR2, 0); err != nil {
							return
						}
						fmt.Println("\nParticion con nombre " + name + " creada exitosamente")
						Structs.PrintMBR(TempMBR2)
					} else {
						//Lo podría eliminar pero tendria que modificar en el metodo del ajuste todos los errores para que aparezca el nombre que se intento ingresar como nueva particion
						fmt.Println("FDISK Error. No se puede crear la nueva particion con nombre: ", name)
					}

					//else if para ingreso de particiones logicas
					//****************
				} else if typee == "L" && isName {
					var partExtend Structs.Partition
					if string(mbr.Partitions[0].Type[:]) == "E" {
						partExtend = mbr.Partitions[0]
					} else if string(mbr.Partitions[1].Type[:]) == "E" {
						partExtend = mbr.Partitions[1]
					} else if string(mbr.Partitions[2].Type[:]) == "E" {
						partExtend = mbr.Partitions[2]
					} else if string(mbr.Partitions[3].Type[:]) == "E" {
						partExtend = mbr.Partitions[3]
					} else {
						fmt.Println("FDISK Error. No existe una particion extendida en la cual crear un particion logica")
					}

					//valido que la particion extendida si exista (podría haber entrado al error que no existe extendida)
					if partExtend.Size != 0 {
						//si tuviera los demas ajustes con un if del fit y uso el metodo segun ajuste
						primerAjusteLogicas(disco, partExtend, int32(sizeNewPart), name, fit) //int32(sizeNewPart) es para castear el int a int32 que es el tipo que tiene el atributo en el struct Partition
						//repLogicas(partExtend, disco)
					}
				}
				//a esta altura sigue abierto el archivo

				//------------------------------ADD---------------------

				//--------------------- Eliminar particiones -----------------------------------------------------

			} else {
				//Probablemente nunca entre aqui (se podría quitar)
				fmt.Println("FDISK Error. Operación desconocida (operaciones aceptadas: crear, modificar o eliminar)")
			}
			//Fin operaciones crear, modificar (add) y eliminar

			// Cierro el disco
			defer disco.Close()
			fmt.Println("======End FDISK======")
		} else {
			fmt.Println("FDISK Error. No se encontro parametro letter y/o name")
		}
	} //Fin if paramC
} //Fin FDisk

/*func primerAjuste(mbr Structs.MBR, typee string, sizeMBR int32, sizeNewPart int32, name string, fit string) (Structs.MBR, Structs.Partition) {
	var newPart Structs.Partition
	var noPart Structs.Partition //para revertir el set info (simula volverla null)

	//PARTICION 1 (libre) - (size = 0 no se ha creado)
	if mbr.Partitions[0].Size == 0 {
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if mbr.Partitions[1].Size == 0 {
			if mbr.Partitions[2].Size == 0 {
				//caso particion 4 (no existe)
				if mbr.Partitions[3].Size == 0 {
					//859 <= 1024 - 165
					if sizeNewPart <= mbr.MbrSize-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
					}
				}
			}
		}
		//Fin de 1 no existe

		//PARTICION 2 (no existe)
	} else if mbr.Partitions[1].Size == 0 {
		//Si no hay espacio antes de particion 1
		newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2) //el nuevo inicio es donde termina 1
		if mbr.Partitions[2].Size == 0 {
			if mbr.Partitions[3].Size == 0 {
				if sizeNewPart <= mbr.MbrSize-newPart.Start {
					mbr.Partitions[1] = newPart
				} else {
					newPart = noPart
					fmt.Println("FDISK Error. Espacio insuficiente")
				}
			}
		}
		//Fin particion 2 no existe

		//PARTICION 3
	} else if mbr.Partitions[2].Size == 0 {
		//despues de 2
		newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
		if mbr.Partitions[3].Size == 0 {
			if sizeNewPart <= mbr.MbrSize-newPart.Start {
				mbr.Partitions[2] = newPart
			} else {
				newPart = noPart
				fmt.Println("FDISK Error. Espacio insuficiente")
			}
		}
		//Fin particion 3

		//PARTICION 4
	} else if mbr.Partitions[3].Size == 0 {
		if sizeNewPart <= mbr.MbrSize-mbr.Partitions[2].GetEnd() {
			//despues de 3
			newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
			mbr.Partitions[3] = newPart
		} else {
			newPart = noPart
			fmt.Println("FDISK Error. Espacio insuficiente")
		}
		//Fin particion 4
	} else {
		newPart = noPart
		fmt.Println("FDISK Error. Particiones primarias y/o extendidas ya no disponibles")
	}

	return mbr, newPart
}
*/

func primerAjuste(mbr Structs.MBR, typee string, sizeMBR int32, sizeNewPart int32, name string, fit string) (Structs.MBR, Structs.Partition) {
	var newPart Structs.Partition
	var noPart Structs.Partition //para revertir el set info (simula volverla null)

	//PARTICION 1 (libre) - (size = 0 no se ha creado)
	if mbr.Partitions[0].Size == 0 {
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if mbr.Partitions[1].Size == 0 {
			if mbr.Partitions[2].Size == 0 {
				//caso particion 4 (no existe)
				if mbr.Partitions[3].Size == 0 {
					//859 <= 1024 - 165
					if sizeNewPart <= mbr.MbrSize-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
					}
				} else {
					//particion 4 existe
					// 600 < 765 - 165 (600 maximo aceptado)
					if sizeNewPart <= mbr.Partitions[3].Start-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						//Si cabe despues de 4
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//Reordeno el correlativo para que coincida con el numero de particion en que se guardo
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					}
				}
				//Fin no existe particion 4
			} else {
				// 3 existe
				//entre mbr y 3 -> 300 <= 465 -165
				if sizeNewPart <= mbr.Partitions[2].Start-sizeMBR {
					mbr.Partitions[0] = newPart
				} else {
					//si no cabe entre el mbr y 3 debe ser despues de 3, es decir, en 4
					newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					} else {
						//4 existe
						//hay espacio entre 3 y 4
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							//Reordenando los correlativos
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3 //new part traia 4 y quedo en la tercer particion por eso tambien se modifica aqui
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							//Hay espacio despues de 4
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//reconfiguro los correlativos
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					} //fin si hay espacio entre 3 y 4
				} //fin si no cabe antes de 3
			} //fin 3 existe
		} else {
			//2 existe
			//Si la nueva particion se puede guardar antes de 2
			if sizeNewPart <= mbr.Partitions[1].Start-sizeMBR {
				mbr.Partitions[0] = newPart
			} else {
				//Si no cabe entre mbr y 2
				//Validar si existen 3 y 4
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				if mbr.Partitions[2].Size == 0 {
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[2] = newPart
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					} else {
						//4 existe (estamos entre 2 y 4)
						//62 < 69-6 (62 maximo aceptado)
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[2] = newPart
						} else {
							//Si no cabe entre 2 y 4, ver si cabe despues de 4
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							if sizeNewPart <= mbr.MbrSize-newPart.Start { //1 <= 100-99
								mbr.Partitions[2] = mbr.Partitions[3]
								mbr.Partitions[3] = newPart
								//reordeno correlativos
								mbr.Partitions[2].Correlative = 3
							} else {
								newPart = noPart
								fmt.Println("FDISK Error. Espacio insuficiente")
							}
						} //Fin si cabe antes o despues de 4
					} //fin de 4 existe o no existe
				} else {
					//3 existe
					//entre 2 y 3
					if sizeNewPart <= mbr.Partitions[2].Start-newPart.Start {
						mbr.Partitions[0] = mbr.Partitions[1]
						mbr.Partitions[1] = newPart
						//Reordeno correlativos
						mbr.Partitions[0].Correlative = 1
						mbr.Partitions[1].Correlative = 2
					} else if mbr.Partitions[3].Size == 0 {
						//entre 3 y el final
						//cambiamos el inicio de la nueva particion porque 3 existe y no cabe antes de 3
						newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					} else {
						//si 4 existe
						//hay espacio entre 3 y 4
						newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 3)
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[0] = mbr.Partitions[1]
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							//Reordeno correlativos
							mbr.Partitions[0].Correlative = 1
							mbr.Partitions[1].Correlative = 2
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							//entre 4 y el final
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[0] = mbr.Partitions[1]
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//Reordeno correlativos
							mbr.Partitions[0].Correlative = 1
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					} //Fin si 4 existe o no (3 activa)
				} //Fin 3 existe o no existe
			} //Fin entre 2 y final (antes de 2 o depues de 2)
		} //Fin 2 existe o no existe
		//Fin de 1 no existe

		//PARTICION 2 (no existe)
	} else if mbr.Partitions[1].Size == 0 {
		//Si hay espacio entre el mbr y particion 1
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start { //particion 1 ya existe (debe existir para entrar a este bloque)
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			//Reordeno correlativo
			mbr.Partitions[1].Correlative = 2
		} else {
			//Si no hay espacio antes de particion 1
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2) //el nuevo inicio es donde termina 1
			if mbr.Partitions[2].Size == 0 {
				if mbr.Partitions[3].Size == 0 {
					if sizeNewPart <= mbr.MbrSize-newPart.Start {
						mbr.Partitions[1] = newPart
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
					}
				} else {
					//4 existe
					//entre 1 y 4
					if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
						mbr.Partitions[1] = newPart
					} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
						//despues de 4
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						mbr.Partitions[2] = mbr.Partitions[3]
						mbr.Partitions[3] = newPart
						//Reordeno correlativo
						mbr.Partitions[2].Correlative = 3
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
					}
				} //Fin 4 existe o no existe
			} else {
				//3 Activa
				//entre 1 y 3
				if sizeNewPart <= mbr.Partitions[2].Start-newPart.Start {
					mbr.Partitions[1] = newPart
				} else {
					//despues de 3
					newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 3)
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
							//corrijo el correlativo
							mbr.Partitions[3].Correlative = 4
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					} else {
						//4 existe
						//entre 3 y 4
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							//Corrijo el correlativo
							mbr.Partitions[1].Correlative = 2
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							//Despues de 4
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//Corrijo los correlativos
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
						}
					} //fin 4 existe o no existe
				} //Fin para entre 1 y 3, y despues de 3
			} //Fin 3 existe o no existe
		} //Fin antes o despues de particion 1
		//Fin particion 2 no existe

		//PARTICION 3
	} else if mbr.Partitions[2].Size == 0 {
		//antes de 1
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {
			mbr.Partitions[2] = mbr.Partitions[1]
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			//Reordeno los correlativos
			mbr.Partitions[2].Correlative = 3
			mbr.Partitions[1].Correlative = 2
		} else {
			//entre 1 y 2
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2)
			if sizeNewPart <= mbr.Partitions[1].Start-newPart.Start {
				mbr.Partitions[2] = mbr.Partitions[1]
				mbr.Partitions[1] = newPart
				//Reordeno correlativo
				mbr.Partitions[2].Correlative = 3
			} else {
				//despues de 2
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				if mbr.Partitions[3].Size == 0 {
					if sizeNewPart <= mbr.MbrSize-newPart.Start {
						mbr.Partitions[2] = newPart
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
					}
				} else {
					//4 existe
					//entre 2 y 4
					if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
						mbr.Partitions[2] = newPart
					} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
						//despues de 4
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						mbr.Partitions[2] = mbr.Partitions[3]
						mbr.Partitions[3] = newPart
						//Reordeno correlativo
						mbr.Partitions[2].Correlative = 3
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
					}
				} //Fin de 4 existe o no existe
			} //Fin espacio entre 1 y 2 o despues de 2
		} //Fin espacio antes de 1
		//Fin particion 3

		//PARTICION 4
	} else if mbr.Partitions[3].Size == 0 {
		//antes de 1
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {
			mbr.Partitions[3] = mbr.Partitions[2]
			mbr.Partitions[2] = mbr.Partitions[1]
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			//Reordeno los correlativos
			mbr.Partitions[3].Correlative = 4
			mbr.Partitions[2].Correlative = 3
			mbr.Partitions[1].Correlative = 2
		} else {
			//si no cabe antes de 1
			//entre 1 y 2
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2)
			if sizeNewPart <= mbr.Partitions[1].Start-newPart.Start {
				mbr.Partitions[3] = mbr.Partitions[2]
				mbr.Partitions[2] = mbr.Partitions[1]
				mbr.Partitions[1] = newPart
				//Reordeno correlativos
				mbr.Partitions[3].Correlative = 4
				mbr.Partitions[2].Correlative = 3
			} else if sizeNewPart <= mbr.Partitions[2].Start-mbr.Partitions[1].GetEnd() {
				//entre 2 y 3
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				mbr.Partitions[3] = mbr.Partitions[2]
				mbr.Partitions[2] = newPart
				//Reordeno correlativo
				mbr.Partitions[3].Correlative = 4
			} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[2].GetEnd() {
				//despues de 3
				newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
				mbr.Partitions[3] = newPart
			} else {
				newPart = noPart
				fmt.Println("FDISK Error. Espacio insuficiente")
			}
		} //Fin antes y despues de 1
		//Fin particion 4
	} else {
		newPart = noPart
		fmt.Println("FDISK Error. Particiones primarias y/o extendidas ya no disponibles")
	}

	return mbr, newPart
}

func primerAjusteLogicas(disco *os.File, partExtend Structs.Partition, sizeNewPart int32, name string, fit string) {
	//Se crea un ebr para cargar el ebr desde el disco y la particion extendida
	save := true //false indica que guardo en el primer ebr, true significa que debe seguir buscando
	var actual Structs.EBR
	sizeEBR := int32(binary.Size(actual)) //obtener el tamaño del ebr (el que ocupa fisicamente: 31)
	//fmt.Println("Tamaño fisico del ebr ", sizeEBR)

	//Guardo el ebr leido
	if err := Herramientas.ReadObject(disco, &actual, int64(partExtend.Start)); err != nil {
		return
	}

	//NOTA: debe caber la particion con el tamaño establecido MAS su EBR
	//NOTA2: Recordar que a la hora de escribir (usar) la particion se inicia donde termina fisicamente la estructura del ebr
	//ej: si el ebr ocupa 5 bytes y la particion es de 10 bytes. los primeros 5 son del ebr entonces uso de 5-15 para escribir en el archivo el contenido de la particion

	//si el primer ebr esta vacio o no existe
	if actual.Size == 0 {
		if actual.Next == -1 {
			//validar si el tamaño de la nueva particion junto al ebr es menor al tamaño de la particion extendida
			if sizeNewPart+sizeEBR <= partExtend.Size {
				actual.SetInfo(fit, partExtend.Start, sizeNewPart, name, -1)
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return
				}
				save = false //ya guardo la nueva particion
				fmt.Println("Particion con nombre ", name, " creada correctamente")
			} else {
				fmt.Println("FDISK Error. Espacio insuficiente logicas")
			}
		} else {
			//Para insertar si se elimino la primera particion (primer EBR)
			//Si actual.Next no es -1 significa que hay otra particion despues de la actual y actual.next tiene el inicio de esa particion
			disponible := actual.Next - partExtend.Start //del inicio hasta donde inicia la siguiente
			if sizeNewPart+sizeEBR <= disponible {
				actual.SetInfo(fit, partExtend.Start, sizeNewPart, name, actual.Next)
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return
				}
				save = false //ya guardo la nueva particion
				fmt.Println("Particion con nombre ", name, " creada correctamente")
			} else {
				fmt.Println("FDISK Error. Espacio insuficiente logicas 2")
			}
		}
		//Si esta despues del primer ebr
	}

	if save {
		//siguiente = actual.next //el valor del siguiente es el inicio de la siguiente particion
		for actual.Next != -1 {
			//si el ebr y la particion caben
			if sizeNewPart+sizeEBR <= actual.Next-actual.GetEnd() {
				break
			}
			//paso al siguiente ebr (simula un actual = actual.next)
			if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
				return
			}

		}

		//Despues de la ultima particion
		if actual.Next == -1 {
			//ya no es el tamaño porque ya hay espacio ocupado por lo que tomo donde termina la extendida y se resta donde termina la ultima
			if sizeNewPart+sizeEBR <= (partExtend.GetEnd() - actual.GetEnd()) {
				//guardar cambios en el ebr actual (cambio el Next)
				actual.Next = actual.GetEnd()
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return
				}

				//crea y guarda la nueva particion logica
				newStart := actual.GetEnd()                          //la nueva ebr inicia donde termina la ultima ebr
				actual.SetInfo(fit, newStart, sizeNewPart, name, -1) //cambia actual con los nuevos valores
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return
				}
				fmt.Println("Particion con nombre ", name, " creada correctamente")
			} else {
				fmt.Println("FDISK Error. Espacio insuficiente logicas 3")
			}
		} else {
			//Entre dos particiones
			if sizeNewPart+sizeEBR <= (actual.Next - actual.GetEnd()) {
				siguiente := actual.Next //guardo el siguiente de la actual para ponerlo en el siguiente de la nueva particion
				//guardar cambio de siguiente en la actual
				actual.Next = actual.GetEnd()
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return
				}

				//agrego la nueva particion apuntando a la siguiente de la actual
				newStart := actual.GetEnd()                                 //la nueva ebr inicia donde termina la ultima ebr
				actual.SetInfo(fit, newStart, sizeNewPart, name, siguiente) //cambia actual con los nuevos valores
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return
				}
				fmt.Println("Particion con nombre ", name, " creada correctamente")
			} else {
				fmt.Println("FDISK Error. Espacio insuficiente logicas 4")
			}
		}
	}
}
