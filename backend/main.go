package main

import (
	"MIA_P1/Analyzer"
	"MIA_P1/Structs"
	"MIA_P1/Utilities"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type ObjIterableData struct {
	ObjIterable string `json:"objIterable"`
}
type responseList struct {
	Status int64    `json:"Status"`
	List   []string `json:"List"`
}

type responseString struct {
	Status int64  `json:"Status"`
	Value  string `json:"Value"`
}

type loginValues struct {
	User     string `json:"User"`
	Password string `json:"Password"`
}
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Discos struct {
	Discos []string `json:"Discos"`
}
type Particiones struct {
	Particiones []string `json:"Particiones"`
}

var particiones_creados []string

func main() {

	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/tasks", obtener_discos)
	http.HandleFunc("/nom_disco", Nombre_Discos)
	http.HandleFunc("/particiones", mandar_particiones)
	http.HandleFunc("/login", loginHandler)
	fmt.Println("Server is running on port 8081")

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")

	// Manejar la solicitud OPTIONS para preflight CORS
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == "POST" {
		r.ParseForm()
		value := r.FormValue("inputValue")
		fmt.Println("Received value from the form:", value)

		//manda el comando al proyecto 1 en la consola
		command, params := getCommandAndParams(value)

		fmt.Println("Command: ", command, "Params: ", params)

		Analyzer.AnalyzeCommnad(command, params)

	}
}
func obtener_discos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Manejar la solicitud OPTIONS para preflight CORS
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a Valid Data")
	}
	var newLoginValues loginValues
	json.Unmarshal(reqBody, &newLoginValues)

	fmt.Println(newLoginValues.User)
	fmt.Println(newLoginValues.Password)

	filepath := "/home/fernando/go/src/MIA_P2/backend/MIA_P1/" + "discos.json"

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Leer el contenido del archivo
	contenido, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return
	}

	// Crear una variable para almacenar los datos decodificados
	var discos Discos

	// Decodificar el contenido JSON en la estructura de datos correspondiente
	err = json.Unmarshal(contenido, &discos)
	if err != nil {
		fmt.Println("Error al decodificar el archivo JSON:", err)
		return
	}
	fmt.Println(discos.Discos)

	var newResponseList responseList
	newResponseList.Status = 200
	newResponseList.List = discos.Discos

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newResponseList)
}

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

func getCommandAndParams(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		params := strings.Join(parts[1:], " ")
		return command, params
	}
	return "", input
}
func Nombre_Discos(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")

	// Manejar la solicitud OPTIONS para preflight CORS
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == "POST" {
		r.ParseForm()
		value := r.FormValue("inputValue")
		fmt.Println("Received value from the form:", value)

		filepath := "./MIA/P1/" + strings.ToUpper(string(value[0])) + ".dsk"
		file, err := Utilities.AbrirFile(filepath)
		if err != nil {
			return
		}
		defer file.Close()

		var TempMBR Structs.MRB
		if err := Utilities.Leer_Object(file, &TempMBR, 0); err != nil {
			return

		}

		for i := 0; i < 4; i++ {
			if TempMBR.M_Partitions[i].Size != 0 {
				temp := string(TempMBR.M_Partitions[i].Name[:])
				particiones_creados = append(particiones_creados, temp[:5])

			} else if TempMBR.M_Partitions[i].Size == 0 {
				break
			}

		}
		fmt.Println(particiones_creados)
		generarJson_disco(particiones_creados)
	}

}

func mandar_particiones(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Manejar la solicitud OPTIONS para preflight CORS
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a Valid Data")
	}
	var newLoginValues loginValues
	json.Unmarshal(reqBody, &newLoginValues)

	fmt.Println(newLoginValues.User)
	fmt.Println(newLoginValues.Password)

	filepath := "/home/fernando/go/src/MIA_P2/backend/MIA_P1/" + "particiones.json"

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Leer el contenido del archivo
	contenido, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return
	}

	// Crear una variable para almacenar los datos decodificados
	var particion Particiones

	// Decodificar el contenido JSON en la estructura de datos correspondiente
	err = json.Unmarshal(contenido, &particion)
	if err != nil {
		fmt.Println("Error al decodificar el archivo JSON:", err)
		return
	}
	fmt.Println(particion.Particiones)

	var newResponseList responseList
	newResponseList.Status = 200
	newResponseList.List = particion.Particiones

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newResponseList)
}

func generarJson_disco(discos []string) error {

	// Crear un mapa para almacenar los discos bajo el identificador "Discos"
	discosMap := map[string][]string{"Particiones": discos}

	// Serializar el mapa a formato JSON
	jsonData, err := json.Marshal(discosMap)
	if err != nil {
		return err
	}
	// Escribir los datos serializados en un archivo
	file, err := os.Create("particiones.json")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Leer el cuerpo de la solicitud
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Decodificar el cuerpo JSON en la estructura User
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Leer los usuarios y contraseñas desde el archivo JSON
	users, err := readUsersFromFile("users.json")
	if err != nil {
		http.Error(w, "Error reading users file", http.StatusInternalServerError)
		return
	}

	// Verificar las credenciales del usuario
	authenticated := false
	for _, u := range users {
		if u.Username == user.Username && u.Password == user.Password {
			authenticated = true
			break
		}
	}

	// Responder con el resultado de la autenticación
	if authenticated {
		fmt.Fprintln(w, "Login successful")
	} else {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	}
}

func readUsersFromFile(filename string) ([]User, error) {
	// Leer el contenido del archivo JSON
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Decodificar el contenido JSON en una lista de usuarios
	var users []User
	if err := json.Unmarshal(content, &users); err != nil {
		return nil, err
	}

	return users, nil
}
