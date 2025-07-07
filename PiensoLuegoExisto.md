# Explorador de archivos

WillyFE (File Explorer)

Limitación: No usar las estructuras que ya conozco.
Sorpresa: usar B+ Tree.?? jajaja
La base es el B+ Tree.
La idea es crear un explorador de archivos que permita gestionar la estructura de archivos y directorios utilizando un árbol B+ para almacenar sus nombres de forma ordenada.
Debe poder realizar operaciones como crear, eliminar, renombrar archivos y buscarlos. La búsqueda se realizará en dos etapas:
primero buscar coincidencias exactas con el nombre ingresado,
luego, adicionalmente, buscar archivos cuyos nombres contengan la cadena ingresada como subcadena.
La concurrencia se aplicará principalmente en el proceso de carga o indexado de archivos (pensando también en dónde más podría ser útil).

la clave esta en la búsqueda, tiene que ser tan rápida que va a saber lo que queres buscar mañana