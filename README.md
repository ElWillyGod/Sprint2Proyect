# Sprint2Proyect



├── main.go             # Lógica de arranque del explorador
├── core/               # Lógica principal del sistema
│   ├── tree.go         # Implementación del B+ Tree
│   └── searcher.go     # Lógica de búsqueda (exacta, prefijo, subcadena)
├── concurrent/         # Funciones relacionadas a concurrencia
│   └── workers.go      # Control de goroutines, canales, sincronización
├── tools/              # Funciones auxiliares (logging, paths, etc.)
│   └── tools.go        # Herramientas auxiliares para el explorador
└── README.md           # Documentación del proyecto
