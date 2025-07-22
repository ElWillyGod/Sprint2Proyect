package core

import (
	"strings"
)

/*
* implementación completa de B+ Tree para almacenar los archivos y sus rutas.
* con división de nodos.
 */

const (
	ORDEN        = 4
	MAX_ENTRADAS = ORDEN - 1
)

type Archivo struct {
	NombreArchivo string
	RutaCompleta  string
}

type EntradaHoja struct {
	Clave string
	Rutas []string
}

// ver de usar dos tipos de nodos
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

	nuevaClave, nuevoNodo := tree.insertarEnNodo(tree.Raiz, clave, archivo)

	if nuevoNodo != nil {
		nuevaRaiz := &Nodo{
			EsHoja: false,
			Claves: []string{nuevaClave},
			Hijos:  []*Nodo{tree.Raiz, nuevoNodo},
		}
		tree.Raiz = nuevaRaiz
	}
}

/*
* Retorna (clave promovida, nuevo nodo) si hay división
* (nil, nil) si no
 */
func (tree *BPlusTree) insertarEnNodo(nodo *Nodo, clave string, archivo Archivo) (string, *Nodo) {
	if nodo.EsHoja {
		return tree.insertarEnHoja(nodo, clave, archivo)
	} else {
		indice := tree.encontrarIndiceHijo(nodo, clave)
		clavePromovida, nuevoHijo := tree.insertarEnNodo(nodo.Hijos[indice], clave, archivo)

		if nuevoHijo != nil {
			return tree.insertarClaveEnInterno(nodo, clavePromovida, nuevoHijo)
		}

		return "", nil
	}
}

/*
* retorna (clave promovida, nuevo nodo) si hay división
* ("", nil) si no
 */
func (tree *BPlusTree) insertarEnHoja(nodo *Nodo, clave string, archivo Archivo) (string, *Nodo) {
	for i, entrada := range nodo.Entradas {
		if entrada.Clave == clave {
			nodo.Entradas[i].Rutas = append(nodo.Entradas[i].Rutas, archivo.RutaCompleta)
			return "", nil
		}
	}
	nuevaEntrada := EntradaHoja{
		Clave: clave,
		Rutas: []string{archivo.RutaCompleta},
	}

	posicion := tree.encontrarPosicionInsercion(nodo.Entradas, clave)
	nodo.Entradas = append(nodo.Entradas, EntradaHoja{})
	copy(nodo.Entradas[posicion+1:], nodo.Entradas[posicion:])
	nodo.Entradas[posicion] = nuevaEntrada

	if len(nodo.Entradas) > MAX_ENTRADAS {
		return tree.dividirHoja(nodo)
	}

	return "", nil
}

/*
* dividir una hoja cuando excede el máximo de entradas
 */
func (tree *BPlusTree) dividirHoja(nodo *Nodo) (string, *Nodo) {
	medio := (len(nodo.Entradas) + 1) / 2

	// Crear nuevo nodo hermano
	nuevoNodo := &Nodo{
		EsHoja:    true,
		Entradas:  make([]EntradaHoja, len(nodo.Entradas)-medio),
		Siguiente: nodo.Siguiente,
	}

	copy(nuevoNodo.Entradas, nodo.Entradas[medio:])

	nodo.Siguiente = nuevoNodo
	nodo.Entradas = nodo.Entradas[:medio]

	clavePromovida := nuevoNodo.Entradas[0].Clave

	return clavePromovida, nuevoNodo
}

/*
* insertar clave en nodo interno después de división de hijo
 */
func (tree *BPlusTree) insertarClaveEnInterno(nodo *Nodo, clave string, nuevoHijo *Nodo) (string, *Nodo) {
	posicion := tree.encontrarPosicionClaveInterna(nodo.Claves, clave)

	nodo.Claves = append(nodo.Claves, "")
	copy(nodo.Claves[posicion+1:], nodo.Claves[posicion:])
	nodo.Claves[posicion] = clave

	nodo.Hijos = append(nodo.Hijos, nil)
	copy(nodo.Hijos[posicion+2:], nodo.Hijos[posicion+1:])
	nodo.Hijos[posicion+1] = nuevoHijo

	if len(nodo.Claves) > MAX_ENTRADAS {
		return tree.dividirInterno(nodo)
	}

	return "", nil
}

func (tree *BPlusTree) dividirInterno(nodo *Nodo) (string, *Nodo) {
	medio := len(nodo.Claves) / 2
	clavePromovida := nodo.Claves[medio]

	nuevoNodo := &Nodo{
		EsHoja: false,
		Claves: make([]string, len(nodo.Claves)-medio-1),
		Hijos:  make([]*Nodo, len(nodo.Hijos)-medio-1),
	}

	copy(nuevoNodo.Claves, nodo.Claves[medio+1:])
	copy(nuevoNodo.Hijos, nodo.Hijos[medio+1:])

	nodo.Claves = nodo.Claves[:medio]
	nodo.Hijos = nodo.Hijos[:medio+1]

	return clavePromovida, nuevoNodo
}

func (tree *BPlusTree) encontrarPosicionInsercion(entradas []EntradaHoja, clave string) int {
	izquierda, derecha := 0, len(entradas)

	for izquierda < derecha {
		medio := (izquierda + derecha) / 2
		if clave < entradas[medio].Clave {
			derecha = medio
		} else {
			izquierda = medio + 1
		}
	}

	return izquierda
}

func (tree *BPlusTree) encontrarPosicionClaveInterna(claves []string, clave string) int {
	izquierda, derecha := 0, len(claves)

	for izquierda < derecha {
		medio := (izquierda + derecha) / 2
		if clave < claves[medio] {
			derecha = medio
		} else {
			izquierda = medio + 1
		}
	}

	return izquierda
}

func (tree *BPlusTree) encontrarIndiceHijo(nodo *Nodo, clave string) int {
	izquierda, derecha := 0, len(nodo.Claves)

	for izquierda < derecha {
		medio := (izquierda + derecha) / 2
		if clave < nodo.Claves[medio] {
			derecha = medio
		} else {
			izquierda = medio + 1
		}
	}

	return izquierda
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

// para avanzar
func (tree *BPlusTree) EncontrarPrimeraHoja() *Nodo {
	nodo := tree.Raiz
	for !nodo.EsHoja {
		nodo = nodo.Hijos[0]
	}
	return nodo
}

func (tree *BPlusTree) buscarEnHoja(nodo *Nodo, clave string) []string {
	izquierda, derecha := 0, len(nodo.Entradas)

	for izquierda < derecha {
		medio := (izquierda + derecha) / 2
		if nodo.Entradas[medio].Clave < clave {
			izquierda = medio + 1
		} else {
			derecha = medio
		}
	}

	if izquierda < len(nodo.Entradas) && nodo.Entradas[izquierda].Clave == clave {
		return nodo.Entradas[izquierda].Rutas
	}

	return nil
}
