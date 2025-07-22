package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"Sprint2Proyect/concurrent"
	"Sprint2Proyect/core"
)

/*
Logica Principal del proyecto - SIMPLIFICADA
B+ Tree para almacenar archivos con concurrencia básica.
Solo búsqueda exacta y parcial.
*/

func main() {

	//////////////////////////////////////////////////////////////////
	rutaAIndexar := "/home/willy"

	///////////////////////////////////////////////////////////////////

	fmt.Printf("archivos desde: %s\n", rutaAIndexar)
	fmt.Println("...")

	// Usar número optimizado para I/O (no CPU bound)
	numWorkers := runtime.NumCPU() / 3 // Para I/O intensivo, menos workers es mejor
	if numWorkers < 2 {
		numWorkers = 2
	}
	if numWorkers > 6 {
		numWorkers = 6 // Sweet spot para operaciones de disco
	}

	fmt.Printf("%d workers\n", numWorkers)

	// Cargar archivos con versión simple
	tree, stats, err := concurrent.CargarArchivosSimple(rutaAIndexar, numWorkers)
	if err != nil {
		log.Fatal("Error cargando directorio:", err)
	}

	buscador := core.NuevoBuscador(tree)

	fmt.Printf("se cargaron %d archivos en %v\n", stats.TotalArchivos, stats.TiempoTotal)
	fmt.Printf("%.0f archivos/seg\n\n", float64(stats.TotalArchivos)/stats.TiempoTotal.Seconds())

	scanner := bufio.NewScanner(os.Stdin)

	for {
		mostrarMenuSimple()

		fmt.Print("Selecciona una opción (1-3): ")
		if !scanner.Scan() {
			break
		}

		opcion := strings.TrimSpace(scanner.Text())

		switch opcion {
		case "1":
			busquedaExacta(buscador, scanner)
		case "2":
			busquedaParcial(buscador, scanner)
		case "3":
			fmt.Println("bye")
			return
		default:
			fmt.Println("a pero sos bobo")
		}

		fmt.Println("\n" + strings.Repeat("-", 40))
	}
}

func mostrarMenuSimple() {
	fmt.Println("1. Búsqueda exacta (nombre completo)")
	fmt.Println("2. Búsqueda parcial (contiene texto)")
	fmt.Println("3. Salir")
	fmt.Println()
}

func busquedaExacta(buscador *core.Buscador, scanner *bufio.Scanner) {
	fmt.Print("Ingresa el nombre exacto del archivo: ")
	if !scanner.Scan() {
		return
	}

	termino := strings.TrimSpace(scanner.Text())
	if termino == "" {
		fmt.Println("sos chistoso")
		return
	}

	fmt.Printf("Buscando '%s'...\n", termino)
	rutas, nombre := buscador.BuscarExacto(termino)

	if len(rutas) == 0 {
		fmt.Printf("No se encontró el archivo '%s'\n", termino)
	} else {
		fmt.Printf("Se encontraron %d archivo(s) con el nombre '%s':\n", len(rutas), termino)
		for i, ruta := range rutas {
			fmt.Printf("  %d. %s\n     %s\n", i+1, nombre, ruta)
		}
	}
}

func busquedaParcial(buscador *core.Buscador, scanner *bufio.Scanner) {
	fmt.Print("Ingresa el texto a buscar (puede ser parte del nombre): ")
	if !scanner.Scan() {
		return
	}

	termino := strings.TrimSpace(scanner.Text())
	if termino == "" {
		fmt.Println("No ingresaste nada.")
		return
	}

	fmt.Printf("Buscando archivos que contengan '%s'...\n", termino)
	rutas := buscador.BuscarParcial(termino)

	if len(rutas) == 0 {
		fmt.Printf("No se encontraron archivos que contengan '%s'\n", termino)
	} else {
		fmt.Printf("Se encontraron %d archivo(s) que contienen '%s':\n", len(rutas), termino)
		mostrarResultados(rutas, 10) // Mostrar máximo 10 resultados
	}
}

func mostrarResultados(rutas []string, limite int) {
	for i, ruta := range rutas {
		if i >= limite {
			fmt.Printf("... y %d archivo(s) más\n", len(rutas)-limite)
			break
		}
		nombre := filepath.Base(ruta)
		fmt.Printf("  %d. 📄 %s\n     📍 %s\n", i+1, nombre, ruta)
	}
}
