package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"net/http"
	"github.com/guillermo/go.procmeminfo"
	"encoding/json"
	"github.com/shirou/gopsutil/cpu"
	"math"

	"github.com/shirou/gopsutil/process"
	"regexp"
	"strconv"
)


/**
*	Función que me muestra la página principal
*/
func indexPageHandler(response http.ResponseWriter, request *http.Request){

	http.ServeFile(response, request, "index.html") //Muestro la página principal

}

/**
*	Función que me devuelve el String completo del estado de un proceso
*	R: Running S: Sleep T: Stop I: Idle Z: Zombie W: Wait L: Lock
*/
func obtenerEstado(caracter string)(estado string){

	if caracter == "R"{
		cantidadRunning+=1
		return "Running"
	}else if caracter == "S"{
		cantidadSleeping+=1
		return "Sleep"
	}else if caracter == "T"{
		cantidadStoped+=1
		return "Stop"
	}else if caracter == "I"{
		return "Idle"
	}else if caracter == "Z"{
		cantidadZombie+=1
		return "Zombie"
	}else if caracter == "W"{
		return "Wait"
	}else if caracter == "L"{
		return "Lock"
	}

	//Retorno uno de error por si no entrara a alguno arriba
	return "Estado indefinido"
}

/**
*	Función que sirve para mandar la información de los procesos que están corriendo actualmente.
*	Esta ruta se llama desde la vista de index.html
*/

type PROCESO struct{
		PID int32
		Usuario string
		Estado string
		Memoria float32
		Nombre string
		Proceso *process.Process
}
type struct_datos struct{
		TotalProcesos int
		TotalEjecucion int
		TotalSuspendidos int
		TotalDetenidos int
		TotalZombie int
		Procesos []PROCESO
	}

var cantidadRunning, cantidadSleeping, cantidadStoped, cantidadZombie int


var arreglo_procesos [] PROCESO


func datosProcesosHandler(response http.ResponseWriter, request *http.Request) {

	//Reinicio mis contadores
	cantidadRunning = 0
	cantidadSleeping = 0
	cantidadStoped = 0
	cantidadZombie = 0


	//var arreglo_procesos [] PROCESO //Defino un arreglo donde voy a poner a todos los procesos
	arreglo_procesos = nil //Vacío el arreglo


	lista_procesos,_ := process.Processes()
	//fmt.Println(lista_procesos)


	for _ , p2 := range lista_procesos{
		usuario, _ := p2.Username() //Return value could be one of these. R: Running S: Sleep T: Stop I: Idle Z: Zombie W: Wait L: Lock
		estado, _ := p2.Status() //Return value could be one of these. R: Running S: Sleep T: Stop I: Idle Z: Zombie W: Wait L: Lock
		memoria, _ := p2.MemoryPercent()
		nombre , _ := p2.Name()
		
		//Agrego el nuevo proceso al arreglo
		arreglo_procesos = append(arreglo_procesos, PROCESO{PID : p2.Pid, Usuario : usuario, Estado : obtenerEstado(estado), Memoria : memoria, Nombre : nombre, Proceso : p2})
	}

	//fmt.Println(len(arreglo_procesos))

	response.Header().Set("Content-Type","application/json")
	response.WriteHeader(http.StatusOK)


	datos := struct_datos {
		TotalProcesos : len(arreglo_procesos),
		TotalEjecucion : cantidadRunning,
		TotalSuspendidos : cantidadSleeping,
		TotalDetenidos : cantidadStoped,
		TotalZombie : cantidadZombie,
		Procesos : arreglo_procesos,
	}
	


	datos_json , _ := json.Marshal(datos)

	response.Write(datos_json)


}


func killProcesoHandler(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "index.html") //Cargo nuevamente la página principal
}



/**
*	Función que sirve para mostrar el monitor de la memoria
*/
func memoriaHandler(response http.ResponseWriter, request *http.Request) {
	
	http.ServeFile(response, request, "memoria.html") //Muestro el monitor de la memoria

}

/**
*	Función que sirve para mandar la información actual de la memoria.
*	Esta ruta se llama desde la vista de memoria.html
*/
func datosmemoriaHandler(response http.ResponseWriter, request *http.Request) {

	meminfo := &procmeminfo.MemInfo{}
	meminfo.Update()

	var total, consumida, porcentaje_consumo, megabytes float64

	megabytes = 1024 * 1024


	total = (float64) (meminfo.Total()) / megabytes //Tamaño en MB
	consumida = (float64) (meminfo.Used()) / megabytes //Tamaño en MB
	porcentaje_consumo = ((consumida * 100) / total)



	response.Header().Set("Content-Type","application/json")
	response.WriteHeader(http.StatusOK)


	type MEMORIA struct {
		Total float64
		Consumida float64
		Porcentaje float64
	}


	datos := MEMORIA{Total : total, Consumida : consumida, Porcentaje : porcentaje_consumo}

	datos_json , _ := json.Marshal(datos)

	response.Write(datos_json)


}


/**
*	Función que sirve para mostrar el monitor de CPU
*/
func CPUHandler(response http.ResponseWriter, request *http.Request) {

	http.ServeFile(response, request, "cpu.html") //Muestro el monitor de CPU

}

/**
*	Función que sirve para mandar la información actual del CPU.
*	Esta ruta se llama desde la vista de memoria.html
*/
func datosCPUHandler(response http.ResponseWriter, request *http.Request) {

	
	vmStat,_ := cpu.Percent(0,false);
	promedio_uso := math.Floor(vmStat[0]*100)/100
	//promedio_uso := math.Floor(vmStat[0])

	response.Header().Set("Content-Type","application/json")
	response.WriteHeader(http.StatusOK)


	type CPU struct {
		Porcentaje float64
	}


	datos := CPU{Porcentaje : promedio_uso}

	datos_json , _ := json.Marshal(datos)

	response.Write(datos_json)


}














var validPath = regexp.MustCompile("^/(kill|save|view)/([a-zA-Z0-9]+)$")

func killHandler(w http.ResponseWriter, r *http.Request, pid string) {

	fmt.Println("Se quiere matar el proceso: ", pid)

	for i := 0; i < len(arreglo_procesos); i++ {
		if string(strconv.Itoa(int(arreglo_procesos[i].PID))) == pid{ //Si encuentro el proceso que se quiere matar
			fmt.Println("Se encontró el proceso: ", pid)
			arreglo_procesos[i].Proceso.Kill()
			break
		}	
	}


	http.Redirect(w, r, "/", http.StatusFound)

}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		//fmt.Println(validPath.FindStringSubmatch(r.URL.Path))
		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {

			http.NotFound(w, r)

			return

		}

		fn(w, r, m[2])

	}

}




var router = mux.NewRouter()

func main(){


	router.HandleFunc("/", indexPageHandler) //Página principal de la aplicación
	router.HandleFunc("/datosprocesos", datosProcesosHandler) //Página principal de la aplicación

	router.HandleFunc("/memoria", memoriaHandler)
	router.HandleFunc("/datosmemoria", datosmemoriaHandler)

	router.HandleFunc("/cpu", CPUHandler)
	router.HandleFunc("/datoscpu", datosCPUHandler)

	//router.HandleFunc("/kill/", makeHandler(killHandler))
	http.HandleFunc("/kill/", makeHandler(killHandler))


	http.Handle("/", router)
	fmt.Println("Servidor corriendo en http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
	

}