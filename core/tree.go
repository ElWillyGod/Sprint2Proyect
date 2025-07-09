package core

/*
implementacion de B+ Tree para almacenar los archivos y sus rutas.
*/

type archivo struct {
	nombre       string // No es necesario
	ruta         string
	esDirectorio bool
}

type listaArchivos struct {
	clave    string
	archivos []archivo
}

type nodo struct {
	hojas      []listaArchivos
	siguientes *nodo
}
