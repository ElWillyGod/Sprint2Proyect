package concurrent

import (
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

// CargarArchivos - carga archivos reutilizando tree.CargarPorLotes
func (w *WorkerBasico) CargarArchivos(rutaDirectorio string) (*core.BPlusTree, *EstadisticasBasicas, error) {
	inicio := time.Now()
	stats := &EstadisticasBasicas{}
	var statsMutex sync.Mutex

	// Canal para recibir archivos del recorrido
	archivosChan := make(chan core.Archivo, 1000)
	var wg sync.WaitGroup

	// Workers para procesar archivos
	for i := 0; i < w.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for archivo := range archivosChan {
				w.mutex.Lock()
				w.tree.Insertar(archivo)
				w.mutex.Unlock()

				statsMutex.Lock()
				stats.TotalArchivos++
				statsMutex.Unlock()
			}
		}()
	}

	// Usar la función genérica de recorrido
	go func() {
		defer close(archivosChan)
		core.RecorrerDirectorio(rutaDirectorio, func(archivo core.Archivo) {
			archivosChan <- archivo
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
