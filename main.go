package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	//"math"
	//"regexp"
	//"strconv"

	"github.com/gorilla/mux"
	"github.com/shirou/gopsutil/process" //Con esto voy a hacer el kill

	"io/ioutil"
)


/**
*	Función que me muestra la página principal
*/
func indexPageHandler(response http.ResponseWriter, request *http.Request){

	http.ServeFile(response, request, "index.html") //Muestro la página principal

}


/**
*	Función que me muestra la página de los procesos
*/
func procesosPageHandler(response http.ResponseWriter, request *http.Request){

	http.ServeFile(response, request, "procesos.html") //Muestro la página de los procesos

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
*	Esta ruta se llama desde la vista de procesos.html
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


/**
*	Función que sirve para mostrar el monitor de la memoria
*/
func memoriaHandler(response http.ResponseWriter, request *http.Request) {
	
	http.ServeFile(response, request, "memoria.html") //Muestro el monitor de la memoria

}

/**
*	Función que lee el archivo que contiene la información de la memoria RAM
*/
func leerRAM(ruta string)(cadena_contenido string){

	bytesLeidos, err := ioutil.ReadFile(ruta)
	if err != nil {
		fmt.Printf("Error leyendo archivo de RAM: %v", err)
		//Devuelvo un json con valores por defecto
		return "{Total : -1.0, Consumida : -1.0, Porcentaje : -1.0}"
	}

	contenido := string(bytesLeidos)
	//fmt.Printf("El contenido del archivo es: %s", contenido)

	return contenido //Retorno el contenido del archivo
}

/**
*	Función que sirve para mandar la información actual de la memoria.
*	Esta ruta se llama desde la vista de memoria.html
*/

type MEMORIA struct {
		Total float64 `json:"Total"`
		Consumida float64 `json:"Consumida"`
		Porcentaje float64 `json:"Porcentaje"`
}

func datosmemoriaHandler(response http.ResponseWriter, request *http.Request) {

	//Voy a leer el archivo que creó el módulo
	string_archivo := leerRAM("/proc/memo_201403525")



	response.Header().Set("Content-Type","application/json")
	response.WriteHeader(http.StatusOK)


	/*type MEMORIA struct {
		Total float64
		Consumida float64
		Porcentaje float64
	}*/


	//datos := MEMORIA{Total : total, Consumida : consumida, Porcentaje : porcentaje_consumo}

	//datos_json , _ := json.Marshal(datos)

	//datos_json , _ := json.Marshal(string_archivo)
	
	

	//Ver si no le tengo que poner comillas al json en el módulo
	//in := `{"firstName":"John","lastName":"Dow"}`
	bytes := []byte(string_archivo)

	var m MEMORIA
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%+v", m)

	datos_json , _ := json.Marshal(m)

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

	response.Header().Set("Content-Type","application/json")
	response.WriteHeader(http.StatusOK)


	type CPU struct {
		Porcentaje float64
	}


	datos := CPU{Porcentaje : 30}

	datos_json , _ := json.Marshal(datos)

	response.Write(datos_json)


}



var router = mux.NewRouter()

func main(){


	router.HandleFunc("/", indexPageHandler) //Página principal de la aplicación
	
	router.HandleFunc("/procesos", procesosPageHandler)
	router.HandleFunc("/datosprocesos", datosProcesosHandler)

	router.HandleFunc("/memoria", memoriaHandler)
	router.HandleFunc("/datosmemoria", datosmemoriaHandler)

	router.HandleFunc("/cpu", CPUHandler)
	router.HandleFunc("/datoscpu", datosCPUHandler)

	////router.HandleFunc("/kill/", makeHandler(killHandler))
	//http.HandleFunc("/kill/", makeHandler(killHandler))


	http.Handle("/", router)
	fmt.Println("Servidor corriendo en http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
	

}
