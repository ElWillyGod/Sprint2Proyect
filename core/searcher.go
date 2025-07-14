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
func (b *Buscador) BuscarExacto(nombreArchivo string) ([]string, string) {
	clave := strings.ToLower(nombreArchivo)
	nodo := b.tree.EncontrarHoja(clave)

	for _, entrada := range nodo.Entradas {
		if entrada.Clave == clave {
			return entrada.Rutas, entrada.Clave
		}
	}

	return nil, ""
}

// Buscar archivos que contengan una subcadena
func (b *Buscador) BuscarParcial(subcadena string) []string {
	var rutas []string
	subcadena = strings.ToLower(subcadena)

	// Recorrer todas las hojas usando los enlaces
	nodoActual := b.tree.EncontrarPrimeraHoja()
	for nodoActual != nil {
		for _, entrada := range nodoActual.Entradas {
			if strings.Contains(entrada.Clave, subcadena) {
				rutas = append(rutas, entrada.Rutas...)
			}
		}
		nodoActual = nodoActual.Siguiente
	}

	return rutas
}
