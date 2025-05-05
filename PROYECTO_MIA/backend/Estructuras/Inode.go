package Estructuras

type Inodo struct {
	I_uid   int64
	I_gid   int64
	I_s     int64
	I_atime [16]byte
	I_ctime [16]byte
	I_mtime [16]byte
	I_block [16]int64
	I_type  int64
	I_perm  int64
}

func NuevoInodo() Inodo {
	var inodo Inodo
	inodo.I_uid = -1
	inodo.I_gid = -1
	inodo.I_s = -1
	inodo.I_atime = [16]byte{}
	inodo.I_ctime = [16]byte{}
	inodo.I_mtime = [16]byte{}
	for i := 0; i < len(inodo.I_block); i++ {
		inodo.I_block[i] = -1
	}
	inodo.I_type = -1
	inodo.I_perm = -1
	return inodo
}
