package concurrent

import (
	"Sprint2Proyect/core"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Worker
type Worker struct {
	tree       *core.BPlusTree
	numWorkers int
	mutex      sync.Mutex
}

type Estadisticas struct {
	TotalArchivos int
	TiempoTotal   time.Duration
}

func NuevoWorker(numWorkers int) *Worker {
	return &Worker{
		tree:       core.NuevoBPlusTree(),
		numWorkers: numWorkers,
	}
}

func (w *Worker) CargarArchivos(rutaDirectorio string) (*core.BPlusTree, *Estadisticas, error) {
	inicio := time.Now()
	stats := &Estadisticas{}
	var statsMutex sync.Mutex

	archivosChan := make(chan core.Archivo, 1000)
	var wg sync.WaitGroup

	const TAMANIO_LOTE = 50

	for i := 0; i < w.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lote := make([]core.Archivo, 0, TAMANIO_LOTE)

			for archivo := range archivosChan {
				lote = append(lote, archivo)

				if len(lote) >= TAMANIO_LOTE {
					w.insertarLote(lote, stats, &statsMutex)
					lote = lote[:0]
				}
			}

			if len(lote) > 0 {
				w.insertarLote(lote, stats, &statsMutex)
			}
		}()
	}

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

func (w *Worker) insertarLote(lote []core.Archivo, stats *Estadisticas, statsMutex *sync.Mutex) {
	w.mutex.Lock()
	for _, archivo := range lote {
		w.tree.Insertar(archivo)
	}
	w.mutex.Unlock()

	statsMutex.Lock()
	stats.TotalArchivos += len(lote)
	statsMutex.Unlock()
}

func RecorrerDirectorio(rutaDir string, procesarArchivo func(core.Archivo)) error {
	entradas, err := os.ReadDir(rutaDir)
	if err != nil {
		return err
	}

	for _, entrada := range entradas {
		rutaCompleta := filepath.Join(rutaDir, entrada.Name())

		if entrada.IsDir() {
			if err := RecorrerDirectorio(rutaCompleta, procesarArchivo); err != nil {
				return err
			}
		} else {
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
func CargarArchivosSimple(rutaDirectorio string, numWorkers int) (*core.BPlusTree, *Estadisticas, error) {
	worker := NuevoWorker(numWorkers)
	return worker.CargarArchivos(rutaDirectorio)
}
