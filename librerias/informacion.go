package librerias

import (
	"os"
	"bufio"
)

func Lectura_archivo(ruta string, tipo int) []string {
	archivo, error := os.Open(ruta)
	defer func(){
	  archivo.Close()
	  recover()
	}()
   
	if error != nil {
	 panic(error)
	}
   
	scanner := bufio.NewScanner(archivo)
	var i int
	var texto []string
	//Itera cada linea
	for scanner.Scan() {
	 if tipo == 1 && i ==2 {
	  break
	 }
	 i++
	 linea := scanner.Text()
	 texto = append(texto, linea)
	}
	return texto
}
