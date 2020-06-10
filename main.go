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

 //Para el socket
//  socketio "github.com/googollee/go-socket.io"
)

//Variables a utilizar
var numeroRun, numeroSleep, numeroStop, numeroZombie int32
// var arr_process []Proceso
//=======================================================================

//Funcion Principal
func main() {
 //Rutas
//  http.HandleFunc("/", lista_procesos)
 http.HandleFunc("/RAM", memoria_proceso)
 http.HandleFunc("/CPU", cpu_proceso)

//  server, err := socketio.NewServer(nil)
//  if err != nil {
//    fmt.Println(err)
//  }

// server.OnConnect("/", func(s socketio.Conn) error {
//   s.SetContext("")
//   fmt.Println("connected:", s.ID())
//   return nil
// })
// server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
//   fmt.Println("notice:", msg)
//   s.Emit("reply", "have "+msg)
// })

// server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
//   s.SetContext(msg)
//   return "recv " + msg
// })

// server.OnError("/", func(s socketio.Conn, e error) {
//   fmt.Println("meet error:", e)
// })
// server.OnDisconnect("/", func(s socketio.Conn, reason string) {
//   fmt.Println("closed", reason)
// })

 fs := http.FileServer(http.Dir("./public"))
 http.Handle("/",fs)

//  http.Handle("/socket.io/", server)

//  go server.Serve()
//  defer server.Close()

 //Servidor levantado
 http.ListenAndServe(":3000", nil)
 fmt.Println("Servidor levantado en el puerto: 3000")
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
	fmt.Println(string(JSON_Data))
	w.Write(JSON_Data)
}

func cpu_proceso(w http.ResponseWriter, r *http.Request) {
  porcent,_ := cpu.Percent(0,false);
  promedio_uso := math.Floor(porcent[0]*100)/100

  info_cpu := CPU{
    Porcentaje : promedio_uso,
  }

  JSON_Data , _ := json.Marshal(info_cpu)
  fmt.Println(string(JSON_Data))
  w.Write(JSON_Data)
}

type RAM struct {
	Total_Ram_Servidor int
	Total_Ram_Consumida int
	Porcentaje_Consumo_Ram float32
}

type CPU struct {
  Porcentaje float64
}