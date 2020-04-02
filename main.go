package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	//"math"
	"regexp"
	"strconv"

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
*	Función que lee el archivo que contiene la información de la memoria RAM
*/
func leerProcesos(ruta string)(cadena_contenido string){

	bytesLeidos, err := ioutil.ReadFile(ruta)
	if err != nil {
		fmt.Printf("Error leyendo archivo de Procesos: %v", err)
		//Devuelvo| un json con valores por defecto
		return "{\"Procesos\" : [], \"Total\" : -1.0, \"Ejecucion\" : -1.0, \"Suspendidos\" : -1.0, \"Detenidos\" : -1.0, \"Zombies\" : -1.0}"
	}

	contenido := string(bytesLeidos)
	//fmt.Printf("El contenido del archivo es: %s", contenido)

	return contenido //Retorno el contenido del archivo
}

/**
*	Función que sirve para mandar la información de los procesos que están corriendo actualmente.
*	Esta ruta se llama desde la vista de procesos.html
*/

func datosProcesosHandler(response http.ResponseWriter, request *http.Request) {

	type PROCESO struct{
		PID int32 `json:"PID"`
		Nombre string `json:"Nombre"`
		Usuario float64 `json:"Usuario"`
		Ram float64 `json:"Ram"`
		Cpu float64 `json:"Cpu"`
		Estado string `json:"Estado"`
		Hijos [] PROCESO `json:"Hijos"`
	}

	type struct_datos struct{
		Procesos []PROCESO `json:"Procesos"`
		Total int `json:"Total"`
		Ejecucion int `json:"Ejecucion"`
		Suspendidos int `json:"Suspendidos"`
		Detenidos int `json:"Detenidos"`
		Zombies int `json:"Zombies"`
	}

	//Voy a leer el archivo que creó el módulo
	string_archivo := leerProcesos("/proc/201403525_procesos")



	response.Header().Set("Content-Type","application/json")
	response.WriteHeader(http.StatusOK)

	//in := `{"firstName":"John","lastName":"Dow"}`
	bytes := []byte(string_archivo)

	var m struct_datos
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%+v", m)

	datos_json , _ := json.Marshal(m)

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
		return "{\"Total\" : -1.0, \"Consumida\" : -1.0, \"Porcentaje\" : -1.0}"
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
	string_archivo := leerRAM("/proc/201403525_memo")



	response.Header().Set("Content-Type","application/json")
	response.WriteHeader(http.StatusOK)

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
*	Función que lee el archivo que contiene la información del CPU
*/
func leerCPU(ruta string)(cadena_contenido string){

	bytesLeidos, err := ioutil.ReadFile(ruta)
	if err != nil {
		fmt.Printf("Error leyendo archivo del CPU: %v", err)
		//Devuelvo un json con valores por defecto
		return "{\"Total\" : -1.0, \"Utilizado\" : -1.0, \"Libre\" : -1.0, \"Promedio\" : -1.0}"
	}

	contenido := string(bytesLeidos)
	//fmt.Printf("El contenido del archivo es: %s", contenido)

	return contenido //Retorno el contenido del archivo
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
type CPU struct {
	Total float64 `json:"Total"`
	Utilizado float64 `json:"Utilizado"`
	Libre float64 `json:"Libre"`
	Promedio float64 `json:"Promedio"`
}
func datosCPUHandler(response http.ResponseWriter, request *http.Request) {

	//Voy a leer el archivo que creó el módulo
	string_archivo := leerRAM("/proc/201403525_cpu")



	response.Header().Set("Content-Type","application/json")
	response.WriteHeader(http.StatusOK)

	//in := `{"firstName":"John","lastName":"Dow"}`
	bytes := []byte(string_archivo)

	var m CPU
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%+v", m)

	datos_json , _ := json.Marshal(m)

	response.Write(datos_json)


}


/**
*	Con esto voy a hacer el kill del proceso
*/
var validPath = regexp.MustCompile("^/(kill|save|view)/([a-zA-Z0-9]+)$")

func killHandler(w http.ResponseWriter, r *http.Request, pid string) {

	fmt.Println("Se quiere matar el proceso: ", pid)

	//Obtengo los procesos
	lista_procesos,_ := process.Processes()

	for _ , p2 := range lista_procesos{

		if string(strconv.Itoa(int(p2.Pid))) == pid{ //Si encuentro el proceso que se quiere matar
			fmt.Println("Se encontró el proceso: ", pid)
			p2.Kill()
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
	
	router.HandleFunc("/procesos", procesosPageHandler)
	router.HandleFunc("/datosprocesos", datosProcesosHandler)

	router.HandleFunc("/memoria", memoriaHandler)
	router.HandleFunc("/datosmemoria", datosmemoriaHandler)

	router.HandleFunc("/cpu", CPUHandler)
	router.HandleFunc("/datoscpu", datosCPUHandler)

	////router.HandleFunc("/kill/", makeHandler(killHandler))
	http.HandleFunc("/kill/", makeHandler(killHandler))


	http.Handle("/", router)
	fmt.Println("Servidor corriendo en http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
	

}
