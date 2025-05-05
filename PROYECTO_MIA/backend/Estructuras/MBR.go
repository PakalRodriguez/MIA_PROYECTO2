package Estructuras

type MBR struct {
	Mbr_tamano         int64
	Mbr_fecha_creacion [16]byte
	Mbr_disk_signature int64
	Dsk_fit            [1]byte
	Mbr_partition_1    Partition
	Mbr_partition_2    Partition
	Mbr_partition_3    Partition
	Mbr_partition_4    Partition
}

func NuevoMBR() MBR {
	var mbr MBR
	return mbr
}
