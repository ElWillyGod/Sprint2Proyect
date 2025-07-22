# Explorador de archivos

WillyFE (File Explorer) (que buen nombre)

Limitación: No usar las estructuras que ya conozco.
Sorpresa: usar B+ Tree.?? jajaja

La base es el B+ Tree.

La idea es crear un explorador de archivos que permita gestionar la estructura de archivos y directorios utilizando un árbol B+ para almacenar sus nombres de forma ordenada.
Debe poder realizar operaciones como crear, eliminar, renombrar archivos y buscarlos. La búsqueda se realizará en dos etapas:
primero buscar coincidencias exactas con el nombre ingresado,
luego, buscar archivos cuyos nombres contengan la cadena ingresada como subcadena.

La concurrencia se aplicará principalmente en el proceso de carga o indexado de archivos también en la búsqueda (Pensar donde mas).

la clave esta en la búsqueda, tiene que ser tan rápida que va a saber lo que queres buscar mañana.

cargar los archivos en el B+ Tree con concurrencia.
buscar en el B+ Tree, la coincidencia exacta del nombre, se busca de forma concurrente el nombre parcial.(posible implementacion de Suffix Tree para la busqueda parcial)


implementacion del B+ Tree:
las hojas van a tener todas las rutas de los archivos que tengan ese nombre.
la clave tal vez sera la primera letra del nombre del archivo para la busqueda exacta.

la ultima comparacion de las hojas es la que es parcial y exacta o todo el recorrido sera concurrente, un recorrido entero para buscar el nombre exacto y otro recorrido para buscar el nombre parcial.

solo las rutas, porque hacer distincion entre archivos y directorios

batch insert -> tecnica de inserción por lotes para mejorar la velocidad de carga.
lexograficamente tremenda palabra.(top 3)
funciones anonimas, o lambdas en python


Guion:

bueno, espero que recuerden mi nombre, sino van a tener que adivinar, no se si alguno usa windows 11 de forma diaria o almenos 3 veces por semana, si si lo usan no se si se an dado cuenta de que el explorador de archivos es exageradamente lento, y eso me molesta, por eso decidi hacer un explorador de archivos que sea mas rapido que el de microsoft, 