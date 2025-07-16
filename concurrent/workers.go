package concurrent

import (
	"Sprint2Proyect/core"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// WorkerBasico
type WorkerBasico struct {
	tree       *core.BPlusTree
	numWorkers int
	mutex      sync.Mutex
}

type EstadisticasBasicas struct {
	TotalArchivos int
	TiempoTotal   time.Duration
}

// NuevoWorkerBasico crea un worker
func NuevoWorkerBasico(numWorkers int) *WorkerBasico {
	return &WorkerBasico{
		tree:       core.NuevoBPlusTree(),
		numWorkers: numWorkers,
	}
}

// CargarArchivos
func (w *WorkerBasico) CargarArchivos(rutaDirectorio string) (*core.BPlusTree, *EstadisticasBasicas, error) {
	inicio := time.Now()
	stats := &EstadisticasBasicas{}
	var statsMutex sync.Mutex

	archivosChan := make(chan core.Archivo, 1000)
	var wg sync.WaitGroup

	const TAMAÑO_LOTE = 50

	for i := 0; i < w.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lote := make([]core.Archivo, 0, TAMAÑO_LOTE)

			for archivo := range archivosChan {
				lote = append(lote, archivo)

				if len(lote) >= TAMAÑO_LOTE {
					w.insertarLote(lote, stats, &statsMutex)
					lote = lote[:0]
				}
			}

			if len(lote) > 0 {
				w.insertarLote(lote, stats, &statsMutex)
			}
		}()
	}

	// funcion fantasma jsj (funciones anonimas, lambda)
	go func() {
		defer close(archivosChan)
		RecorrerDirectorio(rutaDirectorio, func(archivo core.Archivo) {
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

// Función genérica de recorrido que pueden usar ambos módulos
func RecorrerDirectorio(rutaDir string, procesarArchivo func(core.Archivo)) error {
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
			archivo := core.Archivo{
				NombreArchivo: entrada.Name(),
				RutaCompleta:  rutaCompleta,
			}
			procesarArchivo(archivo)
		}
	}
	return nil
}

// Mantener la función simple para compatibilidad
func CargarArchivosSimple(rutaDirectorio string, numWorkers int) (*core.BPlusTree, *EstadisticasBasicas, error) {
	worker := NuevoWorkerBasico(numWorkers)
	return worker.CargarArchivos(rutaDirectorio)
}
