package main

//Importaciones
import (
	"net/http"
	"fmt"

	//Para lectura de los archivos
	"./librerias"
	"strings"

	//Para la lectura de los procesos
	"github.com/shirou/gopsutil/cpu"

	//Para usar json
	"encoding/json"

	//Para conversiones
	"strconv"
	"math"

	//Para hacer el api rest
	"github.com/gorilla/mux"
)

//=======================================================================

//Funcion Principal
func main() {
	router := mux.NewRouter().StrictSlash(true)
	// Para los archivos staticos (css,js)
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	//Rutas de API-REST
	router.HandleFunc("/PROCESS", lista_procesos)
	router.HandleFunc("/RAM", memoria_proceso)
	router.HandleFunc("/CPU", cpu_proceso)
   
	//Rutas para cliente
	router.HandleFunc("/Principal.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w,r, "./public/Principal.html")
	})

	router.HandleFunc("/CPU.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w,r, "./public/CPU.html")
	})

	router.HandleFunc("/RAM.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w,r, "./public/RAM.html")
	})

	router.HandleFunc("/kill/{id}", kill_proceso)

	//Servidor levantado
	fmt.Println("Servidor levantado en el puerto: 3000")
	http.ListenAndServe(":3000", router)
}

func memoria_proceso(w http.ResponseWriter, r *http.Request){
	informacion := librerias.Lectura_archivo("/proc/meminfo", 1)
	MemTotal := informacion[0] 
	MemFree := informacion[1]

	//Haciendo Reemplazos para obtener los datos
	MemTotal = strings.Replace(MemTotal, "MemTotal:", "", -1)
	MemTotal = strings.Replace(MemTotal, " ", "", -1)
	MemTotal = strings.Replace(MemTotal, "kB", "", -1)

	MemFree = strings.Replace(MemFree, "MemFree:", "", -1)
	MemFree = strings.Replace(MemFree, " ", "", -1)
	MemFree = strings.Replace(MemFree, "kB", "", -1)

	//Conversiones y calculos
	MemTotal_,_:= strconv.Atoi(MemTotal)
	MemTotal_ = MemTotal_ / 1000

	MemFree_,_ := strconv.Atoi(MemFree)
	MemFree_ = MemFree_ / 1000

	MemConsumida := MemTotal_ - MemFree_
	PorcentajeConsumo := (float32(MemConsumida) / float32(MemTotal_) )*100

	info_ram := RAM{
		Total_Ram_Servidor: MemTotal_,
		Total_Ram_Consumida: MemConsumida,
		Porcentaje_Consumo_Ram: PorcentajeConsumo,
	}

	JSON_Data,_ := json.Marshal(info_ram)
	w.Write(JSON_Data)
}

func cpu_proceso(w http.ResponseWriter, r *http.Request) {
  porcent,_ := cpu.Percent(0,false);
  promedio_uso := math.Floor(porcent[0]*100)/100

  info_cpu := CPU{
    Porcentaje : promedio_uso,
  }

  JSON_Data , _ := json.Marshal(info_cpu)
  w.Write(JSON_Data)
}

func lista_procesos(w http.ResponseWriter, r *http.Request){
	var arr_process []PROCESS

	//Obteniendo lista de directorios
	lista_directorios := librerias.Get_directorios("/proc")

	//Recorriendo cada directorio
	for _,dir := range lista_directorios {
		informacion := librerias.Lectura_archivo(dir,2)

		//Obteniendo cada atributo
		Pid_ := strings.Split(informacion[0],":")[1]
		Pid_ = strings.Replace(Pid_, "\t", "", -1)

		Nombre_ := strings.Split(informacion[1], ":")[1]
		Nombre_ = strings.Replace(Nombre_, "\t", "", -1)

		Usuario_ := strings.Split(informacion[2], ":")[1]
		Usuario_ = strings.Replace(Usuario_, "\t", " ", -1)
		Usuario_ = strings.Split(Usuario_, " ")[1]  //En la posicion 0 hay un espacio en blanco " "

		Estado_ := strings.Split(informacion[3], ":")[1]
		Estado_ = strings.Replace(Estado_, " ", "", -1)

		Ppid_ := strings.Split(informacion[4],":")[1]
		Ppid_ = strings.Replace(Ppid_, "\t", "", -1)

		Porcentaje_ := ""

		if informacion[5] != "" {
			Porcentaje_ = strings.Split(informacion[5],":")[1]
			Porcentaje_ = strings.Replace(Porcentaje_, "\t", "", -1)
			Porcentaje_ = strings.Replace(Porcentaje_, " ", "", -1)
		} else {
			Porcentaje_ = "---"
		}

		info_process := PROCESS {
			PID: Pid_,
			Nombre: Nombre_,
			Usuario: librerias.GetNombreUsuario(Usuario_),
			Estado: librerias.GetStatus(Estado_),
			PorcentajeRAM: Porcentaje_,
			Proceso_padre: Ppid_,
		}

		arr_process = append(arr_process, info_process)
	}

	//Agregando informacion general
	info_general := Info_general{
		Procesos_en_ejecucion: librerias.NumeroRun,
		Procesos_suspendidos: librerias.NumeroSleep,
		Procesos_detenidos: librerias.NumeroStop,
		Procesos_zombie: librerias.NumeroZombie,
		Total_procesos: len(arr_process),
		List_Procesos: arr_process,
	}

	JSON_Data , _ := json.Marshal(info_general)
	w.Write(JSON_Data)
}

func kill_proceso(w http.ResponseWriter, r *http.Request){
	key := mux.Vars(r)["id"]
	librerias.MatarProceso(key)
	http.Redirect(w,r,"/public/Principal.html", http.StatusFound)
}
//=======================================================================

//Estructuras a utilizar
type RAM struct {
	Total_Ram_Servidor int
	Total_Ram_Consumida int
	Porcentaje_Consumo_Ram float32
}

type CPU struct {
  Porcentaje float64
}

type PROCESS struct {
	PID string
	Nombre string
	Usuario string
	Estado string
	PorcentajeRAM string
	Proceso_padre string
}

type Info_general struct {
	Procesos_en_ejecucion int
	Procesos_suspendidos int
	Procesos_detenidos int
	Procesos_zombie int
	Total_procesos int
	List_Procesos []PROCESS
}