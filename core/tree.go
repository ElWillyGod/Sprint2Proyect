package core

import (
	"os"
	"path/filepath"
	"strings"
)

/*
Implementación completa de B+ Tree para almacenar los archivos y sus rutas.
Con división de nodos y balanceo automático.
*/

const (
	// Orden del B+ Tree - máximo número de hijos por nodo interno
	ORDEN = 4
	// Máximo número de entradas por hoja
	MAX_ENTRADAS_HOJA = ORDEN - 1
	// Máximo número de claves por nodo interno
	MAX_CLAVES_INTERNO = ORDEN - 1
)

type Archivo struct {
	NombreArchivo string // Ej: "nota.txt"
	RutaCompleta  string // Ej: "/home/willy/docs/nota.txt"

}

type EntradaHoja struct {
	Clave string
	Rutas []string
}

// dividir en dos nodos, uno hoja y otro interno
type Nodo struct {
	EsHoja bool

	// Solo para nodos hoja
	Entradas  []EntradaHoja
	Siguiente *Nodo

	// Solo para nodos internos
	Claves []string
	Hijos  []*Nodo
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

func (tree *BPlusTree) Insertar(archivo Archivo) {
	clave := strings.ToLower(archivo.NombreArchivo)

	// Si la raíz se divide, necesitamos crear una nueva raíz
	nuevaClave, nuevoNodo := tree.insertarEnNodo(tree.Raiz, clave, archivo)

	if nuevoNodo != nil {
		// La raíz se dividió, crear nueva raíz
		nuevaRaiz := &Nodo{
			EsHoja: false,
			Claves: []string{nuevaClave},
			Hijos:  []*Nodo{tree.Raiz, nuevoNodo},
		}
		tree.Raiz = nuevaRaiz
	}
}

// Retorna (clave promovida, nuevo nodo) si hay división, (nil, nil) si no
func (tree *BPlusTree) insertarEnNodo(nodo *Nodo, clave string, archivo Archivo) (string, *Nodo) {
	if nodo.EsHoja {
		return tree.insertarEnHoja(nodo, clave, archivo)
	} else {
		// Encontrar el hijo correcto
		indice := tree.encontrarIndiceHijo(nodo, clave)
		clavePromovida, nuevoHijo := tree.insertarEnNodo(nodo.Hijos[indice], clave, archivo)

		if nuevoHijo != nil {
			// El hijo se dividió, insertar la clave promovida en este nodo
			return tree.insertarClaveEnInterno(nodo, clavePromovida, nuevoHijo)
		}

		return "", nil
	}
}

// Retorna (clave promovida, nuevo nodo) si hay división, ("", nil) si no
func (tree *BPlusTree) insertarEnHoja(nodo *Nodo, clave string, archivo Archivo) (string, *Nodo) {
	// Buscar si la clave ya existe
	for i, entrada := range nodo.Entradas {
		if entrada.Clave == clave {
			nodo.Entradas[i].Rutas = append(nodo.Entradas[i].Rutas, archivo.RutaCompleta)
			return "", nil
		}
	}

	// Crear nueva entrada
	nuevaEntrada := EntradaHoja{
		Clave: clave,
		Rutas: []string{archivo.RutaCompleta},
	}

	// Insertar manteniendo orden
	posicion := tree.encontrarPosicionInsercion(nodo.Entradas, clave)
	nodo.Entradas = append(nodo.Entradas, EntradaHoja{})
	copy(nodo.Entradas[posicion+1:], nodo.Entradas[posicion:])
	nodo.Entradas[posicion] = nuevaEntrada

	// Verificar si necesita división
	if len(nodo.Entradas) > MAX_ENTRADAS_HOJA {
		return tree.dividirHoja(nodo)
	}

	return "", nil
}

func (tree *BPlusTree) encontrarPosicionInsercion(entradas []EntradaHoja, clave string) int {
	for i, entrada := range entradas {
		if clave < entrada.Clave {
			return i
		}
	}
	return len(entradas)
}

// Dividir una hoja cuando excede el máximo de entradas
func (tree *BPlusTree) dividirHoja(nodo *Nodo) (string, *Nodo) {
	medio := (len(nodo.Entradas) + 1) / 2

	// Crear nuevo nodo hermano
	nuevoNodo := &Nodo{
		EsHoja:    true,
		Entradas:  make([]EntradaHoja, len(nodo.Entradas)-medio),
		Siguiente: nodo.Siguiente,
	}

	// Mover la mitad de las entradas al nuevo nodo
	copy(nuevoNodo.Entradas, nodo.Entradas[medio:])

	// Mantener enlaces entre hojas
	nodo.Siguiente = nuevoNodo
	nodo.Entradas = nodo.Entradas[:medio]

	// La clave promovida es la primera clave del nuevo nodo
	clavePromovida := nuevoNodo.Entradas[0].Clave

	return clavePromovida, nuevoNodo
}

// Insertar clave en nodo interno después de división de hijo
func (tree *BPlusTree) insertarClaveEnInterno(nodo *Nodo, clave string, nuevoHijo *Nodo) (string, *Nodo) {
	// Insertar la clave en la posición correcta
	posicion := tree.encontrarPosicionClaveInterna(nodo.Claves, clave)

	// Insertar clave
	nodo.Claves = append(nodo.Claves, "")
	copy(nodo.Claves[posicion+1:], nodo.Claves[posicion:])
	nodo.Claves[posicion] = clave

	// Insertar hijo (después de la clave insertada)
	nodo.Hijos = append(nodo.Hijos, nil)
	copy(nodo.Hijos[posicion+2:], nodo.Hijos[posicion+1:])
	nodo.Hijos[posicion+1] = nuevoHijo

	// Verificar si necesita división
	if len(nodo.Claves) > MAX_CLAVES_INTERNO {
		return tree.dividirInterno(nodo)
	}

	return "", nil
}

// Dividir un nodo interno cuando excede el máximo de claves
func (tree *BPlusTree) dividirInterno(nodo *Nodo) (string, *Nodo) {
	medio := len(nodo.Claves) / 2
	clavePromovida := nodo.Claves[medio]

	// Crear nuevo nodo hermano
	nuevoNodo := &Nodo{
		EsHoja: false,
		Claves: make([]string, len(nodo.Claves)-medio-1),
		Hijos:  make([]*Nodo, len(nodo.Hijos)-medio-1),
	}

	// Mover claves e hijos al nuevo nodo
	copy(nuevoNodo.Claves, nodo.Claves[medio+1:])
	copy(nuevoNodo.Hijos, nodo.Hijos[medio+1:])

	// Truncar el nodo original
	nodo.Claves = nodo.Claves[:medio]
	nodo.Hijos = nodo.Hijos[:medio+1]

	return clavePromovida, nuevoNodo
}

// Encontrar posición para insertar clave en nodo interno
func (tree *BPlusTree) encontrarPosicionClaveInterna(claves []string, clave string) int {
	for i, c := range claves {
		if clave < c {
			return i
		}
	}
	return len(claves)
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
