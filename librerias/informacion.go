package librerias

import (
	"os"
	"bufio"
	"io/ioutil"
	// "fmt"

	//Para conversiones
	"strconv"
	"strings"

	//Para ejecutar comandos de consola
	"os/exec"
)

//Variables a utilizar
var NumeroRun, NumeroSleep, NumeroStop, NumeroZombie int

func Lectura_archivo(ruta string, tipo int) [6]string {
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
	var texto2 [6]string
	//Itera cada linea
	for scanner.Scan() {
		if tipo == 1 && i ==2 {
			break
		}
   		i++
   		linea := scanner.Text()
		if tipo == 1{
			texto2[i-1] = linea
		} else {
			nombre_aux := strings.Split(linea, ":")
			if nombre_aux[0] == "Pid" {
				texto2[0] = linea
			} else if nombre_aux[0] == "Name" {
				texto2[1] = linea
			} else if nombre_aux[0] == "Uid" {
				texto2[2] = linea
			} else if nombre_aux[0] == "State" {
				texto2[3] = linea
			} else if nombre_aux[0] == "PPid" {
				texto2[4] = linea
			} else if nombre_aux[0] == "VmSize" {
				texto2[5] = linea
			}
		}
	}
	return texto2
}

func Get_directorios(ruta string) []string {
	files, err := ioutil.ReadDir(ruta)
	if err != nil {
		panic(err)
	}

	var texto []string
	for _, archivo := range files {
		if archivo.IsDir() {
			nombre := archivo.Name()
			_,error := strconv.Atoi(nombre)
			if error == nil{
				texto = append(texto, ruta + "/" + nombre + "/status")
			}
		}
	}
	return texto
}

/**
Status returns the process status. Return value could be one of these. 
R: Running S: Sleep T: Stop I: Idle Z: Zombie W: Wait L: Lock 
The character is same within all supported platforms.
*/

func GetStatus(caracter string) string{
	if strings.Contains(caracter, "R") {
		NumeroRun++
		return "Running"
	} else if strings.Contains(caracter, "S") {
		NumeroSleep++
		return "Sleep"
	} else if strings.Contains(caracter, "T") {
		NumeroStop++
		return "Stop"
	} else if strings.Contains(caracter, "I") {
		return "Idle"
	} else if strings.Contains(caracter, "Z") {
		NumeroZombie++
		return "Zombie"
	} else if strings.Contains(caracter, "W") {
		return "Wait"
	} else if strings.Contains(caracter, "L") {
		return "Lock"
	} else {
		return "Error status"
	}
}

func GetNombreUsuario(uid string) string {
	var usuario string
	cmd,error := exec.Command("grep", "x:"+uid, "/etc/passwd").Output()
	if error != nil {
		usuario = "---"
		return usuario
	}
	usuario = strings.Split(string(cmd), ":")[0]
	return usuario
}

func MatarProceso(key string)  {
	_,error := exec.Command("kill", "-15", key).Output()
	if error != nil {
		panic(error)
	}
}