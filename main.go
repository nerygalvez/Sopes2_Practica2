package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	//"math"
	//"regexp"
	//"strconv"

	"github.com/gorilla/mux"
	//"github.com/shirou/gopsutil/process" //Con esto voy a hacer el kill

	"io/ioutil"
)


/**
*	Función que me muestra la página principal
*/
func indexPageHandler(response http.ResponseWriter, request *http.Request){

	http.ServeFile(response, request, "index.html") //Muestro la página principal

}

/**
*	Función que sirve para mandar la información de los procesos que están corriendo actualmente.
*	Esta ruta se llama desde la vista de index.html
*/
func datosProcesosHandler(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type","application/json")
	response.WriteHeader(http.StatusOK)

	type struct_proceso struct {
		Total float64
		Consumida float64
		Porcentaje float64
	}

	datos := struct_proceso{
		Total : 4.0,
		Consumida : 10.0,
		Porcentaje : 5.0,
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

	datos_json , _ := json.Marshal(string_archivo)

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
	router.HandleFunc("/datosprocesos", datosProcesosHandler) //Página principal de la aplicación

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