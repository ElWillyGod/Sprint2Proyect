package concurrent

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"Sprint2Proyect/core"
)

// WorkerBasico - Version muy simple
type WorkerBasico struct {
	tree       *core.BPlusTree
	numWorkers int
	mutex      sync.Mutex
}

// EstadisticasBasicas - Solo lo esencial
type EstadisticasBasicas struct {
	TotalArchivos int
	TiempoTotal   time.Duration
}

// NuevoWorkerBasico crea un worker simple
func NuevoWorkerBasico(numWorkers int) *WorkerBasico {
	return &WorkerBasico{
		tree:       core.NuevoBPlusTree(),
		numWorkers: numWorkers,
	}
}

// CargarArchivos - carga archivos de forma simple
func (w *WorkerBasico) CargarArchivos(rutaDirectorio string) (*core.BPlusTree, *EstadisticasBasicas, error) {
	inicio := time.Now()
	stats := &EstadisticasBasicas{}

	// Canal simple
	archivosChan := make(chan core.Archivo, 100)
	var wg sync.WaitGroup

	// Workers simples
	for i := 0; i < w.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for archivo := range archivosChan {
				w.mutex.Lock()
				w.tree.Insertar(archivo)
				w.mutex.Unlock()
			}
		}()
	}

	// Recorrer directorio
	go func() {
		defer close(archivosChan)

		filepath.Walk(rutaDirectorio, func(ruta string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if !info.IsDir() {
				archivo := core.Archivo{
					NombreArchivo: info.Name(),
					RutaCompleta:  ruta,
				}

				stats.TotalArchivos++
				archivosChan <- archivo
			}

			return nil
		})
	}()

	wg.Wait()
	stats.TiempoTotal = time.Since(inicio)
	return w.tree, stats, nil
}

// Mantener la función simple para compatibilidad
func CargarArchivosSimple(rutaDirectorio string, numWorkers int) (*core.BPlusTree, *EstadisticasBasicas, error) {
	worker := NuevoWorkerBasico(numWorkers)
	return worker.CargarArchivos(rutaDirectorio)
}
