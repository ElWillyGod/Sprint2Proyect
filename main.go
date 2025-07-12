package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	// Establecer la ruta que quieres indexar (cámbiala aquí)
	rutaAIndexar := "/home" // <-- CAMBIA ESTA RUTA POR LA QUE QUIERAS

	fmt.Printf("📁 Cargando archivos desde: %s\n", rutaAIndexar)
	fmt.Println("⏳ Por favor espera...")

	// Usar número optimizado para I/O (no CPU bound)
	numWorkers := runtime.NumCPU() / 3 // Para I/O intensivo, menos workers es mejor
	if numWorkers < 2 {
		numWorkers = 2
	}
	if numWorkers > 6 {
		numWorkers = 6 // Sweet spot para operaciones de disco
	}

	fmt.Printf("🚀 Usando %d workers\n", numWorkers)

	// Cargar archivos con versión simple
	tree, stats, err := concurrent.CargarArchivosSimple(rutaAIndexar, numWorkers)
	if err != nil {
		log.Fatal("Error cargando directorio:", err)
	}

	// Crear el buscador
	buscador := core.NuevoBuscador(tree)

	// Mostrar estadísticas simples
	fmt.Printf("✅ ¡Listo! Se cargaron %d archivos en %v\n", stats.TotalArchivos, stats.TiempoTotal)
	fmt.Printf("⚡ Velocidad: %.0f archivos/seg\n\n", float64(stats.TotalArchivos)/stats.TiempoTotal.Seconds())

	// Scanner para leer entrada del usuario
	scanner := bufio.NewScanner(os.Stdin)

	// Menú interactivo simple
	for {
		mostrarMenuSimple()

		fmt.Print("Selecciona una opción (1-4): ")
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
			mostrarTodosLosArchivos(buscador)
		case "4":
			fmt.Println("👋 ¡Hasta luego!")
			return
		default:
			fmt.Println("❌ Opción inválida. Intenta de nuevo.")
		}

		fmt.Println("\n" + strings.Repeat("-", 40))
	}
}

func mostrarMenuSimple() {
	fmt.Println("🔍 ¿Qué quieres hacer?")
	fmt.Println("1. 🎯 Búsqueda exacta (nombre completo)")
	fmt.Println("2. 🔎 Búsqueda parcial (contiene texto)")
	fmt.Println("3. 📋 Ver todos los archivos")
	fmt.Println("4. 🚪 Salir")
	fmt.Println()
}

func busquedaExacta(buscador *core.Buscador, scanner *bufio.Scanner) {
	fmt.Print("📝 Ingresa el nombre exacto del archivo: ")
	if !scanner.Scan() {
		return
	}

	termino := strings.TrimSpace(scanner.Text())
	if termino == "" {
		fmt.Println("❌ No ingresaste nada.")
		return
	}

	fmt.Printf("🔍 Buscando '%s'...\n", termino)
	archivos := buscador.BuscarExacto(termino)

	if len(archivos) == 0 {
		fmt.Printf("❌ No se encontró el archivo '%s'\n", termino)
	} else {
		fmt.Printf("✅ Se encontraron %d archivo(s) con el nombre '%s':\n", len(archivos), termino)
		for i, archivo := range archivos {
			fmt.Printf("  %d. 📄 %s\n     📍 %s\n", i+1, archivo.NombreArchivo, archivo.RutaCompleta)
		}
	}
}

func busquedaParcial(buscador *core.Buscador, scanner *bufio.Scanner) {
	fmt.Print("📝 Ingresa el texto a buscar (puede ser parte del nombre): ")
	if !scanner.Scan() {
		return
	}

	termino := strings.TrimSpace(scanner.Text())
	if termino == "" {
		fmt.Println("❌ No ingresaste nada.")
		return
	}

	fmt.Printf("🔍 Buscando archivos que contengan '%s'...\n", termino)
	archivos := buscador.BuscarParcial(termino)

	if len(archivos) == 0 {
		fmt.Printf("❌ No se encontraron archivos que contengan '%s'\n", termino)
	} else {
		fmt.Printf("✅ Se encontraron %d archivo(s) que contienen '%s':\n", len(archivos), termino)
		mostrarResultados(archivos, 10) // Mostrar máximo 10 resultados
	}
}

func mostrarTodosLosArchivos(buscador *core.Buscador) {
	fmt.Println("📋 Obteniendo todos los archivos...")
	archivos := buscador.ObtenerTodosLosArchivos()

	fmt.Printf("📊 Total de archivos: %d\n", len(archivos))
	fmt.Println("📋 Primeros 15 archivos:")
	mostrarResultados(archivos, 15)
}

func mostrarResultados(archivos []core.Archivo, limite int) {
	for i, archivo := range archivos {
		if i >= limite {
			fmt.Printf("... y %d archivo(s) más\n", len(archivos)-limite)
			break
		}
		fmt.Printf("  %d. 📄 %s\n     📍 %s\n", i+1, archivo.NombreArchivo, archivo.RutaCompleta)
	}
}
