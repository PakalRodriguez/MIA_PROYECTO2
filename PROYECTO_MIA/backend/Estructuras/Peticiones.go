package Estructuras

type Comando struct {
	Comand string `json:"text"`
}

type DiscoPeticion struct {
	Driveletter string `json:"disco"`
}

type ParticionPeticion struct {
	Disco  string `json:"disco"`
	Nombre string `json:"particion"`
}

type InicioSesion struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Id       string `json:"id"`
	Nombre   string `json:"nombre"`
}

type Reporte struct {
	Nombre string `json:"reporte"`
}

type ReporteEnviar struct {
	Extension string
	Nombre    string
	Contenido string
}
