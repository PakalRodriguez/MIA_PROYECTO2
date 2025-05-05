# Manual Técnico: Sistema de Archivos EXT2 Simulado

## 1. Descripción de la Arquitectura del Sistema

El sistema de archivos EXT2 simulado en esta aplicación web está compuesto por un backend desarrollado en **Go**, utilizando los frameworks **Fiber** y **CORS**, y un frontend implementado en **React**. La arquitectura general se compone de los siguientes módulos:

### 1.1 Módulos Principales

- **Frontend (React):** Interfaz de usuario donde se ejecutan los comandos del sistema de archivos.
- **Backend (Go + Fiber):** Procesamiento de comandos y simulación del sistema de archivos EXT2.
- **Sistema de Archivos (.mia):** Archivo binario que almacena la estructura del sistema de archivos.

### 1.2 Diagrama de Arquitectura

```
[React UI] <--> [Go Backend (Fiber)] <--> [.mia File System]
```

## 2. Explicación de las Estructuras de Datos

El sistema de archivos se almacena en un archivo binario `.mia`, donde se representan las siguientes estructuras:

### 2.1 MBR (Master Boot Record)
El MBR es la primera estructura dentro del archivo binario y almacena información sobre las particiones del disco.

```go
 type MBR struct {
     MbrSize    int32        // Tamaño total del disco
     FechaC     [16]byte     // Fecha de creación
     Id         int32        // Identificador único
     Fit        [1]byte      // Tipo de ajuste (First, Best, Worst)
     Partitions [4]Partition // Particiones del disco
 }
```

### 2.2 Particiones
Cada partición puede ser primaria o extendida, y almacena información relevante para su gestión.

```go
 type Partition struct {
     Status      [1]byte  // Estado de la partición (Activa/Inactiva)
     Type        [1]byte  // Tipo (Primaria o Extendida)
     Fit         [1]byte  // Ajuste (F, B, W)
     Start       int32    // Inicio de la partición
     Size        int32    // Tamaño de la partición
     Name        [16]byte // Nombre de la partición
 }
```

### 2.3 EBR (Extended Boot Record)
Utilizado en particiones extendidas para manejar particiones lógicas.

```go
 type EBR struct {
     Status [1]byte  // Estado de la partición
     Type   [1]byte  // Tipo (Lógica)
     Fit    [1]byte  // Ajuste
     Start  int32    // Inicio de la partición
     Size   int32    // Tamaño de la partición
     Name   [16]byte // Nombre
     Next   int32    // Apuntador a la siguiente partición lógica
 }
```

### 2.4 Superbloque
Contiene la información general del sistema de archivos y gestiona inodos y bloques.

```go
 type Superblock struct {
     S_filesystem_type   int32    // Tipo de sistema de archivos (EXT2)
     S_inodes_count      int32    // Total de inodos
     S_blocks_count      int32    // Total de bloques
     S_free_blocks_count int32    // Bloques libres
     S_free_inodes_count int32    // Inodos libres
     S_mnt_count         int32    // Veces montado
     S_magic             int32    // Identificador (0xEF53)
 }
```

### 2.5 Inodos
Estructura fundamental que describe archivos y directorios.

```go
 type Inode struct {
     I_uid   int32    // ID de usuario
     I_gid   int32    // ID de grupo
     I_size  int32    // Tamaño del archivo
     I_block [15]int32 // Bloques de datos asociados
     I_type  [1]byte  // Tipo (Archivo/Carpeta)
 }
```

## 3. Descripción de los Comandos Implementados

El sistema soporta varios comandos para la gestión del sistema de archivos. A continuación, se describen algunos de los más relevantes:

### 3.1 `MKDISK`
**Crea un nuevo disco virtual en formato `.mia`.**

- **Ejemplo:**
  ```sh
  mkdisk -Size=3000 -unit=K -path=/home/user/Disco1.mia
  ```

- **Parámetros:**
  - `-Size`: Tamaño del disco en KB, MB o GB.
  - `-unit`: Unidad de medida (`K`, `M`, `G`).
  - `-path`: Ubicación del archivo del disco.

### 3.2 `FDISK`
**Gestiona particiones dentro de un disco virtual.**

- **Ejemplo:**
  ```sh
  fdisk -type=P -name=Part1 -size=500 -unit=M -path=/home/user/Disco1.mia
  ```

- **Parámetros:**
  - `-type`: Tipo de partición (`P`: Primaria, `E`: Extendida, `L`: Lógica).
  - `-name`: Nombre de la partición.
  - `-size`: Tamaño de la partición.

### 3.3 `MOUNT`
**Monta una partición para su uso.**

- **Ejemplo:**
  ```sh
  mount -path=/home/user/Disco1.mia -name=Part1
  ```

- **Parámetros:**
  - `-path`: Ubicación del disco.
  - `-name`: Nombre de la partición.

### 3.4 `REP`
**Genera reportes sobre el estado del sistema de archivos.**

- **Ejemplo:**
  ```sh
  rep -name=mbr -path=/home/user/mbr.png -id=vd1
  ```

- **Parámetros:**
  - `-name`: Tipo de reporte (`mbr`, `disk`, `inode`, etc.).
  - `-path`: Ruta del archivo de salida.
  - `-id`: Identificador del disco montado.

## 4. Conclusión

Este Manual Técnico proporciona una guía completa sobre la arquitectura, estructuras de datos y comandos del sistema de archivos EXT2 simulado. La aplicación combina un backend en **Go con Fiber** y un frontend en **React**, permitiendo la emulación de operaciones sobre discos y particiones dentro de un entorno controlado.

