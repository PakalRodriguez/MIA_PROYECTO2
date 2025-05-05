package Estructuras

import "unsafe"

type Superbloque struct {
	S_filesystem_type   int64
	S_inodes_count      int64
	S_blocks_count      int64
	S_free_inodes_count int64
	S_free_blocks_count int64
	S_mtime             [16]byte
	S_umtime            [16]byte
	S_mnt_count         int64
	S_magic             int64
	S_inode_s           int64
	S_block_s           int64
	S_first_ino         int64
	S_first_blo         int64
	S_bm_inode_start    int64
	S_bm_block_start    int64
	S_inode_start       int64
	S_block_start       int64
}

func NuevoSuperbloque() Superbloque {
	var superBloque Superbloque
	superBloque.S_mnt_count = 0
	superBloque.S_magic = 0xEF53
	superBloque.S_inode_s = int64(unsafe.Sizeof(Inodo{}))
	superBloque.S_block_s = int64(unsafe.Sizeof(BloqueCarpeta{}))
	superBloque.S_first_ino = 0
	superBloque.S_first_blo = 0
	superBloque.S_bm_inode_start = -1
	superBloque.S_bm_block_start = -1
	superBloque.S_inode_start = -1
	superBloque.S_block_start = -1
	return superBloque
}
