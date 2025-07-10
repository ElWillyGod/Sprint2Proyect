package core

import (
	"os"
	"path/filepath"
	"strings"
)

/*
Implementación de B+ Tree para almacenar los archivos y sus rutas.
*/

type Archivo struct {
	NombreArchivo string // Ej: "nota.txt"
	RutaCompleta  string // Ej: "/home/willy/docs/nota.txt"
	EsDirectorio  bool   // true si es un directorio
}

type EntradaHoja struct {
	Clave    string  
	Archivos []Archivo 
}

type Nodo struct {
	Claves    []string 
	EsHoja    bool
	
	// Solo para nodos hoja
	Entradas  []EntradaHoja
	Siguiente *Nodo
	
	// Solo para nodos internos
	Hijos     []*Nodo
}

type BPlusTree struct {
	Raiz *Nodo
}


func NuevoBPlusTree() *BPlusTree {
	raiz := &Nodo{
		EsHoja:   true,
		Entradas: make([]EntradaHoja, 0),
	}
	return &BPlusTree{
		Raiz: raiz,
	}
}

// sera dios
func (tree *BPlusTree) CargarDirectorio(rutaDirectorio string) error {
	return filepath.Walk(rutaDirectorio, func(ruta string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		archivo := Archivo{
			NombreArchivo: info.Name(),
			RutaCompleta:  ruta,
			EsDirectorio:  info.IsDir(),
		}
		
		tree.Insertar(archivo)
		return nil
	})
}


func (tree *BPlusTree) Insertar(archivo Archivo) {
	clave := strings.ToLower(archivo.NombreArchivo) // puros problemas
	tree.insertarEnNodo(tree.Raiz, clave, archivo)
}


func (tree *BPlusTree) insertarEnNodo(nodo *Nodo, clave string, archivo Archivo) {
	if nodo.EsHoja {
		tree.insertarEnHoja(nodo, clave, archivo)
	} else {
		// el hijo correcto
		indice := tree.encontrarIndiceHijo(nodo, clave)
		tree.insertarEnNodo(nodo.Hijos[indice], clave, archivo)
	}
}


func (tree *BPlusTree) insertarEnHoja(nodo *Nodo, clave string, archivo Archivo) {
	// Buscar si la clave ya existe
	for i := range nodo.Entradas {
		if nodo.Entradas[i].Clave == clave {
			// Agregar archivo a la entrada existente
			nodo.Entradas[i].Archivos = append(nodo.Entradas[i].Archivos, archivo)
			return
		}
	}
	
	// Crear nueva entrada
	nuevaEntrada := EntradaHoja{
		Clave:    clave,
		Archivos: []Archivo{archivo},
	}
	
	// Insertar manteniendo orden
	posicion := tree.encontrarPosicionInsercion(nodo.Entradas, clave)
	nodo.Entradas = append(nodo.Entradas, EntradaHoja{})
	copy(nodo.Entradas[posicion+1:], nodo.Entradas[posicion:])
	nodo.Entradas[posicion] = nuevaEntrada
}


func (tree *BPlusTree) encontrarPosicionInsercion(entradas []EntradaHoja, clave string) int {
	for i, entrada := range entradas {
		if clave < entrada.Clave {
			return i
		}
	}
	return len(entradas)
}


func (tree *BPlusTree) encontrarIndiceHijo(nodo *Nodo, clave string) int {
	for i, c := range nodo.Claves {
		if clave < c {
			return i
		}
	}
	return len(nodo.Claves) // último hijo
}


func (tree *BPlusTree) EncontrarHoja(clave string) *Nodo {
	return tree.encontrarHoja(tree.Raiz, clave)
}

func (tree *BPlusTree) encontrarHoja(nodo *Nodo, clave string) *Nodo {
	if nodo.EsHoja {
		return nodo
	}
	
	indice := tree.encontrarIndiceHijo(nodo, clave)
	return tree.encontrarHoja(nodo.Hijos[indice], clave)
}


func (tree *BPlusTree) EncontrarPrimeraHoja() *Nodo {
	nodo := tree.Raiz
	for !nodo.EsHoja {
		nodo = nodo.Hijos[0]
	}
	return nodo
}
