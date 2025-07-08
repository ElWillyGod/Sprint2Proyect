# Sprint2Proyect



├── main.go             # Lógica de arranque del explorador <br>
├── core/               # Lógica principal del sistema<br>
│   ├── tree.go         # Implementación del B+ Tree<br>
│   └── searcher.go     # Lógica de búsqueda (exacta, prefijo, subcadena)<br>
├── concurrent/         # Funciones relacionadas a concurrencia<br>
│   └── workers.go      # Control de goroutines, canales, sincronización<br>
├── tools/              # Funciones auxiliares (logging, paths, etc.)<br>
│   └── tools.go        # Herramientas auxiliares para el explorador<br>
└── README.md           # Documentación del proyecto<br>
