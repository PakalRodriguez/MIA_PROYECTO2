package Estructuras

type EBR struct {
	Part_mount byte
	Part_fit   byte
	Part_start int64
	Part_s     int64
	Part_next  int64
	Part_name  [16]byte
}

func NuevoEBR() EBR {
	var ebr EBR
	ebr.Part_mount = '0'
	ebr.Part_fit = 'W'
	ebr.Part_start = -1
	ebr.Part_s = 0
	ebr.Part_next = -1
	ebr.Part_name = [16]byte{}
	return ebr
}
