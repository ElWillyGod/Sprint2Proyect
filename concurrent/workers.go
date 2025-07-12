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

// CargarArchivos - carga archivos usando workers concurrentes optimizado
func (w *WorkerBasico) CargarArchivos(rutaDirectorio string) (*core.BPlusTree, *EstadisticasBasicas, error) {
	inicio := time.Now()
	stats := &EstadisticasBasicas{}
	var statsMutex sync.Mutex

	// Canal para recibir archivos del recorrido
	archivosChan := make(chan core.Archivo, 1000)
	var wg sync.WaitGroup

	// Buffer para acumular archivos antes de insertar en lotes
	const TAMAÑO_LOTE = 50
	
	// Workers para procesar archivos en lotes
	for i := 0; i < w.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lote := make([]core.Archivo, 0, TAMAÑO_LOTE)
			
			for archivo := range archivosChan {
				lote = append(lote, archivo)
				
				// Procesar lote cuando esté lleno
				if len(lote) >= TAMAÑO_LOTE {
					w.insertarLote(lote, stats, &statsMutex)
					lote = lote[:0] // Limpiar sin reallocar
				}
			}
			
			// Procesar lote final si no está vacío
			if len(lote) > 0 {
				w.insertarLote(lote, stats, &statsMutex)
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

// Insertar un lote de archivos con una sola adquisición de mutex
func (w *WorkerBasico) insertarLote(lote []core.Archivo, stats *EstadisticasBasicas, statsMutex *sync.Mutex) {
	w.mutex.Lock()
	for _, archivo := range lote {
		w.tree.Insertar(archivo)
	}
	w.mutex.Unlock()
	
	// Actualizar estadísticas
	statsMutex.Lock()
	stats.TotalArchivos += len(lote)
	statsMutex.Unlock()
}

// Mantener la función simple para compatibilidad
func CargarArchivosSimple(rutaDirectorio string, numWorkers int) (*core.BPlusTree, *EstadisticasBasicas, error) {
	worker := NuevoWorkerBasico(numWorkers)
	return worker.CargarArchivos(rutaDirectorio)
}
