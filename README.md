# 🚀 WillyFE

<div align="center">

![Go](https://img.shields.io/badge/Go-1.24.5-00ADD8?style=for-the-badge&logo=go)
![Status](https://img.shields.io/badge/Status-En%20Desarrollo-yellow?style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

**Explorador de Archivos Concurrente con B+ Tree**

*Un sistema de indexación y búsqueda de archivos de alto rendimiento construido en Go*

</div>

---

## 📋 Descripción

WillyFE es un explorador de archivos concurrente que utiliza una estructura de datos B+ Tree para indexar y buscar archivos de manera eficiente. El sistema aprovecha las capacidades de concurrencia de Go mediante goroutines y canales para lograr un rendimiento óptimo.

## ✨ Características

- 🌳 **B+ Tree**: Estructura de datos optimizada para búsquedas eficientes
- 🔍 **Búsqueda Múltiple**: Soporte para búsqueda exacta, por prefijo y subcadena
- ⚡ **Concurrencia**: Procesamiento paralelo con goroutines
- 🔄 **Sincronización**: Control avanzado de canales y sincronización
- 📊 **Alto Rendimiento**: Optimizado para manejar grandes volúmenes de archivos

## 🏗️ Estructura del Proyecto

```
Sprint2Proyect/
├── 📄 main.go             # Lógica de arranque del explorador
├── 🎯 core/               # Lógica principal del sistema
│   ├── tree.go            # Implementación del B+ Tree
│   └── searcher.go        # Lógica de búsqueda (exacta, prefijo, subcadena)
├── 🔄 concurrent/         # Funciones relacionadas a concurrencia
│   └── workers.go         # Control de goroutines, canales, sincronización
├── � go.mod              # Configuración del módulo Go
├── 📝 PiensoLuegoExisto.md # Documentación adicional
└── 📚 README.md           # Documentación del proyecto
```

## 🚀 Instalación y Uso

### Prerrequisitos
- Go 1.24.5 o superior
- Sistema operativo compatible con Go

### Instalación

1. **Clona el repositorio**
   ```bash
   git clone https://github.com/ElWillyGod/Sprint2Proyect.git
   cd Sprint2Proyect
   ```

2. **Instala las dependencias**
   ```bash
   go mod tidy
   ```

3. **Ejecuta el proyecto**
   ```bash
   go run main.go
   ```

## 💻 Uso

El sistema indexará automáticamente los archivos desde la ruta configurada y proporcionará capacidades de búsqueda interactiva:

```bash
# Ejemplo de ejecución
archivos desde: /ruta/especificada
...
```

## 🏗️ Arquitectura

### Componentes Principales

- **🎯 Core**: Contiene la lógica principal del B+ Tree y los algoritmos de búsqueda
- **🔄 Concurrent**: Maneja la concurrencia y paralelización de tareas

### Patrones de Diseño Utilizados

- **Concurrencia**: Uso de goroutines y canales para procesamiento paralelo
- **Separación de responsabilidades**: Modularización clara de funcionalidades
- **Estructura de datos eficiente**: B+ Tree para búsquedas optimizadas

## 🤝 Contribución

Las contribuciones son bienvenidas. Por favor:

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## 👨‍💻 Autor

**ElWillyGod** - [@ElWillyGod](https://github.com/ElWillyGod)

---

<div align="center">

⭐ ¡Dale una estrella si te gusta este proyecto!

</div>
