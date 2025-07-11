package main

/*
Logica Principal del proyecto.
La idea es crear un B+ Tree que almacene los archivos y sus rutas.
la busqueda se realizará de forma concurrente, utilizando goroutines para cargar los archivos en el B+ Tree y para realizar búsquedas en el árbol.
dos busquedas:
exacta, se realizará buscando el nombre del archivo completo.
parcial, se realizará buscando el nombre del archivo que contenga una subcadena.
concurrencia en la carga de archivos y en la búsqueda.
*/
