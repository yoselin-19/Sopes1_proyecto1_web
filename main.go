package main

//Importaciones
import (
 "net/http"
 "fmt"
)

//funciones
func main() {

 //Rutas
 fileServer := http.FileServer(http.Dir("public"))
 http.Handle("/", http.StripPrefix("/", fileServer)) 

 //Servidor levantado
 http.ListenAndServe(":3000", nil)
 fmt.Println("Servidor levantado en el puerto: 3000")
}
