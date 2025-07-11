package core

import (
	"strings"
)

/*
Implementación simplificada de un buscador de archivos utilizando el B+ Tree.
Solo dos tipos de búsqueda:
1. Búsqueda exacta: se busca el nombre del archivo completo.
2. Búsqueda parcial: se busca el nombre del archivo que contenga una subcadena.
*/

// Buscador que opera sobre un B+ Tree
type Buscador struct {
	tree *BPlusTree
}

// Constructor del buscador
func NuevoBuscador(tree *BPlusTree) *Buscador {
	return &Buscador{
		tree: tree,
	}
}

// Buscar archivos por nombre exacto
func (b *Buscador) BuscarExacto(nombreArchivo string) []Archivo {
	clave := strings.ToLower(nombreArchivo)
	nodo := b.tree.EncontrarHoja(clave)

	for _, entrada := range nodo.Entradas {
		if entrada.Clave == clave {
			return entrada.Archivos
		}
	}

	return nil
}

// Buscar archivos que contengan una subcadena
func (b *Buscador) BuscarParcial(subcadena string) []Archivo {
	var resultados []Archivo
	subcadena = strings.ToLower(subcadena)

	// Recorrer todas las hojas usando los enlaces
	nodoActual := b.tree.EncontrarPrimeraHoja()
	for nodoActual != nil {
		for _, entrada := range nodoActual.Entradas {
			if strings.Contains(entrada.Clave, subcadena) {
				resultados = append(resultados, entrada.Archivos...)
			}
		}
		nodoActual = nodoActual.Siguiente
	}

	return resultados
}

// Obtener todos los archivos ordenados por nombre
func (b *Buscador) ObtenerTodosLosArchivos() []Archivo {
	var todosLosArchivos []Archivo

	nodoActual := b.tree.EncontrarPrimeraHoja()
	for nodoActual != nil {
		for _, entrada := range nodoActual.Entradas {
			todosLosArchivos = append(todosLosArchivos, entrada.Archivos...)
		}
		nodoActual = nodoActual.Siguiente
	}

	return todosLosArchivos
}

// Verificar si existe un archivo con el nombre especificado
func (b *Buscador) ExisteArchivo(nombreArchivo string) bool {
	archivos := b.BuscarExacto(nombreArchivo)
	return len(archivos) > 0
}
