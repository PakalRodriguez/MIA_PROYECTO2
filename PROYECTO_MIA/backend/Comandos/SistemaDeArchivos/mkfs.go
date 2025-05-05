package Comandos

import (
	"PROYECTO_MIA/backend/Herramientas"
	"PROYECTO_MIA/backend/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

func Mkfs(entrada []string) {
	var id string //obligatorio
	paramC := true
	var pathDico string

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro, " ")
		valores := strings.Split(tmp, "=")

		if len(valores) != 2 {
			fmt.Println("ERROR MKFS, valor desconocido de parametros ", valores[1])
			break
		}

		if strings.ToLower(valores[0]) == "id" {
			id = strings.ToUpper(valores[1])
		} else if strings.ToLower(valores[0]) == "type" {
			if strings.ToLower(valores[1]) != "full" {
				fmt.Println("MKFS Error. Valor de -type desconocido")
				paramC = false
				break
			}

			//ERROR EN LOS PARAMETROS LEIDO
		} else {
			fmt.Println("MKFS Error: Parametro desconocido: ", valores[0])
			paramC = false
			break //por si en el camino reconoce algo invalido de una vez se sale
		}

	}

	//obtener la particion correspondiente al id
	if id != "" {
		//BUsca en struck de particiones montadas el id ingresado
		for _, montado := range Structs.Montadas {
			if montado.Id == id {
				pathDico = montado.PathM
			}
		}
		if pathDico == "" {
			fmt.Println("ERROR MKFS NO SE ENCONTRA EL ID")
			paramC = false
		}
	} else {
		fmt.Println("ERROR MKFS NO SE INGRESO ID")
		paramC = false
	}

	if paramC {
		//Abrir el Disco de la particion
		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return
		}

		//Cargar el mbr
		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return
		}

		// Close bin file
		defer file.Close()

		//Buscar particion con el id solicitado
		formatear := true
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				formatear = false //Si encontro la particion

				//Crear el super bloque que contiene los datos del sistema de archivos. Es similar al mbr en los discos
				var newSuperBloque Structs.Superblock
				Herramientas.ReadObject(file, &newSuperBloque, int64(mbr.Partitions[i].Start))

				//Calcular el numero de inodos que caben en la particion. El numero de bloques es el triple de inodos
				//(formula a partir del tamaño de la particion, esta en el enunciado pag. 10)
				//tamaños fisicos: SuperBloque = 92; Inodo = 124; Bloque = 64
				numerador := int(mbr.Partitions[i].Size) - binary.Size(Structs.Superblock{})
				denominador := 4 + binary.Size(Structs.Inode{}) + 3*binary.Size(Structs.Fileblock{})

				n := int32(numerador / denominador) //numero de inodos

				//inicializar atributos generales del superbloque
				newSuperBloque.S_blocks_count = int32(3 * n)      //Total de bloques creados (pueden usarse)
				newSuperBloque.S_free_blocks_count = int32(3 * n) //Numero de bloques libre (Todos estan libres por ahora)

				newSuperBloque.S_inodes_count = n      //Total de inodos creados (pueden usarse)
				newSuperBloque.S_free_inodes_count = n //numero de inodos libres (todos estan libres por ahora)

				newSuperBloque.S_inode_size = int32(binary.Size(Structs.Inode{}))
				newSuperBloque.S_block_size = int32(binary.Size(Structs.Fileblock{}))

				//obtener hora de montaje del sistema de archivos
				ahora := time.Now()
				copy(newSuperBloque.S_mtime[:], ahora.Format("02/01/2006 15:04"))
				//Si fecha de desmontaje coincide con montaje es porque aun no se monta
				copy(newSuperBloque.S_umtime[:], ahora.Format("02/01/2006 15:04"))
				newSuperBloque.S_mnt_count += 1 //Se esta montando por primera vez
				newSuperBloque.S_magic = 0xEF53

				crearEXT2(n, mbr.Partitions[i], newSuperBloque, ahora.Format("02/01/2006 15:04"), file)

				//Fin del formateo
				fmt.Println("Particion con id ", id, " formateada correctamente")

				//Si hubiera una sesion iniciada eliminarla
				break //para que ya no siga recorriendo las demas particiones
			}
		}

		if formatear {
			fmt.Println("MKFS Error. No se pudo formatear la particion con id ", id)
			fmt.Println("MKFS Error. No existe el id")
		}
	}
}

func crearEXT2(n int32, particion Structs.Partition, newSuperBloque Structs.Superblock, date string, file *os.File) {
	fmt.Println("Superbloque: ", newSuperBloque)
	fmt.Println("Fecha: ", date)

	//completar los atributos del super bloque. La estructura de la particion formateada es:
	// | Superbloque | Bitmap Inodos | Bitmap Bloques | Inodos | Bloques |

	//tipo del sistema de archivos
	newSuperBloque.S_filesystem_type = 2 //2 -> EXT2; 3 -> EXT3
	//Bitmap Inodos inicia donde termina el superbloque fisicamente (y el superbloque esta al inicio de la particion)
	newSuperBloque.S_bm_inode_start = particion.Start + int32(binary.Size(Structs.Superblock{}))
	//Bitmap bloques inicia donde termina el de inodos. Se suma n que es el numero de inodos maximo
	newSuperBloque.S_bm_block_start = newSuperBloque.S_bm_inode_start + n
	//Se crea el primer Inodo. Esta al final de los bloques que son 3 veces el numero de inodos
	newSuperBloque.S_inode_start = newSuperBloque.S_bm_block_start + 3*n
	//Se crea el primer bloque, este esta al final de los inodos fisicos
	newSuperBloque.S_block_start = newSuperBloque.S_inode_start + n*int32(binary.Size(Structs.Inode{}))

	//Se restan 2 bloques y dos inodos. uno para la carpeta raiz y otro para el archivo users.txt
	//lo que se crea al formatear es /users.txt (la carpeta usa un inodo y el archivo otro)
	newSuperBloque.S_free_inodes_count -= 2
	newSuperBloque.S_free_blocks_count -= 2

	//primer inodo libre
	//newSuperBloque.S_first_ino = newSuperBloque.S_inode_start + 2*int32(binary.Size(Structs.Inode{})) //multiplico por 2 porque hay 2 inodos creados
	newSuperBloque.S_first_ino = int32(2)
	//primer bloque libre
	//newSuperBloque.S_first_blo = newSuperBloque.S_block_start + 2*int32(binary.Size(Structs.Fileblock{})) //multiplicar por 2 porque hay 2 bloques creados
	newSuperBloque.S_first_blo = int32(2)

	//limpio (formateo) el espacio del bitmap de inodos para evitar inconsistencias
	bmInodeData := make([]byte, n)
	bmInodeErr := Herramientas.WriteObject(file, bmInodeData, int64(newSuperBloque.S_bm_inode_start))
	if bmInodeErr != nil {
		fmt.Println("MKFS Error: ", bmInodeErr)
		return
	}

	//limpiar (formatear) el espacio del bitmap de bloques para evitar inconsistencias
	bmBlockData := make([]byte, 3*n)
	bmBlockErr := Herramientas.WriteObject(file, bmBlockData, int64(newSuperBloque.S_bm_block_start))
	if bmBlockErr != nil {
		fmt.Println("MKFS Error: ", bmInodeErr)
		return
	}

	//creo un inodo y lleno el arreglo de bloques con -1
	var newInode Structs.Inode
	for i := 0; i < 15; i++ {
		newInode.I_block[i] = -1
	}

	//creo todos los inodos del sistema de archivos
	for i := int32(0); i < n; i++ {
		err := Herramientas.WriteObject(file, newInode, int64(newSuperBloque.S_inode_start+i*int32(binary.Size(Structs.Inode{}))))
		if err != nil {
			fmt.Println("MKFS Error: ", err)
			return
		}
	}

	//Crear todos los bloques de carpeta que se pueden crear
	fileBlocks := make([]Structs.Fileblock, 3*n) //lo puedo trabajar asi porque son instancias de la estructura, el inode llevaban valores
	fileBlocksErr := Herramientas.WriteObject(file, fileBlocks, int64(newSuperBloque.S_bm_block_start))
	if fileBlocksErr != nil {
		fmt.Println("MKFS Error: ", fileBlocksErr)
		return
	}

	//Crear el Inode 0
	var Inode0 Structs.Inode
	Inode0.I_uid = 1
	Inode0.I_gid = 1
	Inode0.I_size = 0 //por ser carpeta no tiene tamaño como tal. para saber si existe basarse en I_ui/I_gid
	//unica vez que las 3 fechas son iguales
	copy(Inode0.I_atime[:], date)
	copy(Inode0.I_ctime[:], date)
	copy(Inode0.I_mtime[:], date)
	copy(Inode0.I_type[:], "0") //como es raiz es de tipo carpeta
	copy(Inode0.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode0.I_block[i] = -1
	}

	Inode0.I_block[0] = 0 //apunta al bloque 0

	//Crear el folder con la estructura
	// 	. 		| 0   -> actual (a si mismo)
	// 	..      | 0   -> el padre
	//users.txt | 1
	//			|-1

	var folderBlock0 Structs.Folderblock //Bloque0 -> carpetas
	folderBlock0.B_content[0].B_inodo = 0
	copy(folderBlock0.B_content[0].B_name[:], ".")
	folderBlock0.B_content[1].B_inodo = 0
	copy(folderBlock0.B_content[1].B_name[:], "..")
	folderBlock0.B_content[2].B_inodo = 1
	copy(folderBlock0.B_content[2].B_name[:], "users.txt")
	folderBlock0.B_content[3].B_inodo = -1

	//Inode1 que es el que contiene el archivo (Bloque 0 apunta a este nuevo inodo)
	var Inode1 Structs.Inode
	Inode1.I_uid = 1
	Inode1.I_gid = 1
	Inode1.I_size = int32(binary.Size(Structs.Folderblock{}))
	copy(Inode1.I_atime[:], date)
	copy(Inode1.I_ctime[:], date)
	copy(Inode1.I_mtime[:], date)
	copy(Inode1.I_type[:], "1") //es del archivo
	copy(Inode0.I_perm[:], "664")
	for i := int32(0); i < 15; i++ {
		Inode1.I_block[i] = -1
	}
	//Inode1 apunta al bloque1 (en este caso el bloque1 contiene el archivo)
	Inode1.I_block[0] = 1
	data := "1,G,root\n1,U,root,root,123\n"
	var fileBlock1 Structs.Fileblock //Bloque1 -> archivo
	copy(fileBlock1.B_content[:], []byte(data))
	fmt.Println("Creado users.txt con los datos : \n", data)

	//resumen
	//Inodo 0 -> Bloque 0 -> Inodo1 -> bloque1 (archivo)

	//Crear la carpeta raiz /
	//crear el archivo users.txt

	//fmt.Println("Superbloque: ", newSuperBloque)

	// Escribir el superbloque
	Herramientas.WriteObject(file, newSuperBloque, int64(particion.Start))

	//escribir el bitmap de inodos
	Herramientas.WriteObject(file, byte(1), int64(newSuperBloque.S_bm_inode_start))
	Herramientas.WriteObject(file, byte(1), int64(newSuperBloque.S_bm_inode_start+1)) //Se escribieron dos inode

	//escribir el bitmap de bloques (se usaron dos bloques)
	Herramientas.WriteObject(file, byte(1), int64(newSuperBloque.S_bm_block_start))
	Herramientas.WriteObject(file, byte(1), int64(newSuperBloque.S_bm_block_start+1))

	//escribir inodes
	//Inode0
	Herramientas.WriteObject(file, Inode0, int64(newSuperBloque.S_inode_start))
	//Inode1
	Herramientas.WriteObject(file, Inode1, int64(newSuperBloque.S_inode_start+int32(binary.Size(Structs.Inode{}))))

	//Escribir bloques
	//bloque0
	Herramientas.WriteObject(file, folderBlock0, int64(newSuperBloque.S_block_start))
	//bloque1
	Herramientas.WriteObject(file, fileBlock1, int64(newSuperBloque.S_block_start+int32(binary.Size(Structs.Fileblock{}))))
	// Fin crear EXT2
}
