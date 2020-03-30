package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"math"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/shirou/gopsutil/process"
)


/**
*	Función que me muestra la página principal
*/
func indexPageHandler(response http.ResponseWriter, request *http.Request){

	http.ServeFile(response, request, "index.html") //Muestro la página principal

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

	response.Header().Set("Content-Type","application/json")
	response.WriteHeader(http.StatusOK)


	type MEMORIA struct {
		Total float64
		Consumida float64
		Porcentaje float64
	}


	datos := MEMORIA{Total : 100, Consumida : 80, Porcentaje : 24.3}

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