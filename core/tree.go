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
	// EsDirectorio eliminado ya que solo almacenamos archivos
}

type EntradaHoja struct {
	Clave    string
	Archivos []Archivo
}

// dividir en dos nodos, uno hoja y otro interno
type Nodo struct {
	Claves []string
	EsHoja bool

	// Solo para nodos hoja
	Entradas  []EntradaHoja
	Siguiente *Nodo

	// Solo para nodos internos
	Hijos []*Nodo
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

// Carga solo archivos (no directorios) del directorio especificado
func (tree *BPlusTree) CargarDirectorio(rutaDirectorio string) error {
	return filepath.Walk(rutaDirectorio, func(ruta string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Solo insertar archivos, ignorar directorios
		if !info.IsDir() {
			archivo := Archivo{
				NombreArchivo: info.Name(),
				RutaCompleta:  ruta,
			}

			tree.Insertar(archivo)
		}
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
	for i, entrada := range nodo.Entradas {
		if entrada.Clave == clave {
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

// Carga de archivos por lotes usando os.ReadDir
func (tree *BPlusTree) CargarPorLotes(rutaDirectorio string, tamañoLote int) error {
	var lote []Archivo
	return tree.procesarDirectorioLotes(rutaDirectorio, tamañoLote, &lote)
}

func (tree *BPlusTree) procesarDirectorioLotes(rutaDir string, tamañoLote int, lote *[]Archivo) error {
	return RecorrerDirectorio(rutaDir, func(archivo Archivo) {
		*lote = append(*lote, archivo)
		
		// Si el lote está lleno, procesarlo
		if len(*lote) >= tamañoLote {
			tree.insertarLote(*lote)
			*lote = (*lote)[:0] // Limpiar lote sin reallocar
		}
	})
}

// Insertar un lote completo de archivos
func (tree *BPlusTree) insertarLote(archivos []Archivo) {
	for _, archivo := range archivos {
		tree.Insertar(archivo)
	}
}

// Función genérica de recorrido que pueden usar ambos módulos
func RecorrerDirectorio(rutaDir string, procesarArchivo func(Archivo)) error {
	entradas, err := os.ReadDir(rutaDir)
	if err != nil {
		return err
	}

	for _, entrada := range entradas {
		rutaCompleta := filepath.Join(rutaDir, entrada.Name())
		
		if entrada.IsDir() {
			// Recursión
			if err := RecorrerDirectorio(rutaCompleta, procesarArchivo); err != nil {
				return err
			}
		} else {
			// Procesar archivo usando la función callback
			archivo := Archivo{
				NombreArchivo: entrada.Name(),
				RutaCompleta:  rutaCompleta,
			}
			procesarArchivo(archivo)
		}
	}
	return nil
}
