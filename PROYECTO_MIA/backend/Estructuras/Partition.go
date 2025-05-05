package Estructuras

type Partition struct {
	Part_status      byte
	Part_type        byte
	Part_fit         byte
	Part_start       int64
	Part_s           int64
	Part_name        [16]byte
	Part_correlative int64
	Part_id          [4]byte
}

func NuevaPartition() Partition {
	var partition Partition
	partition.Part_status = '0'
	partition.Part_type = 'P'
	partition.Part_fit = 'W'
	partition.Part_start = -1
	partition.Part_s = 0
	partition.Part_name = [16]byte{}
	partition.Part_correlative = -1
	partition.Part_id = [4]byte{}
	return partition
}
