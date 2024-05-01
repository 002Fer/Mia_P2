package DiskManagement

import (
	"MIA_P1/Structs"
	"MIA_P1/Utilities"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var aux_tamaño int
var discos_creados []string
var nombre_disk string

func Mount(driveletter string, name string) {
	fmt.Println("======Start ======")
	fmt.Println("Driveletter:", driveletter)
	fmt.Println("Name:", name)

	filepath := "./MIA/P1/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := Utilities.AbrirFile(filepath)
	if err != nil {
		return
	}

	var TempMBR Structs.MRB

	if err := Utilities.Leer_Object(file, &TempMBR, 0); err != nil {
		return

	}

	fmt.Println("-------------")

	var index int = -1
	var count = 0

	for i := 0; i < 4; i++ {
		if TempMBR.M_Partitions[i].Size != 0 {
			count++
			if strings.Contains(string(TempMBR.M_Partitions[i].Name[:]), name) {

				if strings.Contains(string(TempMBR.M_Partitions[i].Status[:]), "1") {
					fmt.Println("Partition is mounted")
					break

				} else if strings.Contains(string(TempMBR.M_Partitions[i].Type[:]), "e") {
					fmt.Println("NO se puede montar una extendida")
					break
				}
				index = i
				break
			}
		}
	}

	if index != -1 {
		fmt.Println("Partition found")

	} else {
		fmt.Println("Partition not found")
		return
	}

	id := strings.ToUpper(driveletter) + strconv.Itoa(count) + "50"

	copy(TempMBR.M_Partitions[index].Status[:], "1")
	copy(TempMBR.M_Partitions[index].Id[:], id)

	if err := Utilities.Escribir_Object(file, TempMBR, 0); err != nil {
		return
	}

	var TempMBR2 Structs.MRB

	if err := Utilities.Leer_Object(file, &TempMBR2, 0); err != nil {
		return
	}

	Structs.PrintMBR(TempMBR2)

	defer file.Close()

	fmt.Println("======End MOUNT======")
}

func Fdisk(size int, driveletter string, name string, unit string, type_ string, fit string) {
	fmt.Println("======Start FDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Driveletter:", driveletter)
	fmt.Println("Name:", name)
	fmt.Println("Unit:", unit)
	fmt.Println("Type:", type_)
	fmt.Println("Fit:", fit)
	fit = string(fit[1])
	aux_tamaño = len(name)
	aux_tamaño = aux_tamaño - 1

	if fit != "b" && fit != "w" && fit != "f" {
		fmt.Println("Error: Fit must be b, w or f")
		return
	}

	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: Unit must be b, k or m")
		return
	}

	if type_ != "p" && type_ != "e" && type_ != "l" {
		fmt.Println("Error: Type must be p, e or l")
		return
	}

	if (unit == "B") || (unit == "") {
		size = size * 1
	} else if unit == "k" {
		size = size * 1024
	} else if unit == "M" {
		size = size * 1024 * 1024
	}

	filepath := "./MIA/P1/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := Utilities.AbrirFile(filepath)
	if err != nil {
		return
	}
	var TempMBR Structs.MRB

	if err := Utilities.Leer_Object(file, &TempMBR, 0); err != nil {
		return
	}

	Structs.PrintMBR(TempMBR)

	fmt.Println("-------------")

	if type_ == "l" {

		var temp_next int32
		var index int = -1
		var count = 0
		var tamaño int32 = 0
		for i := 0; i < 4; i++ {
			if TempMBR.M_Partitions[i].Size != 0 {
				count++
				tamaño += TempMBR.M_Partitions[i].Size
				if strings.Contains(string(TempMBR.M_Partitions[i].Type[:]), "e") {
					index = i

					break
				}
			}
		}
		fmt.Println("inicio de la particion logica")
		inicio := int32(tamaño + int32(binary.Size(Structs.MRB{})))
		fmt.Println(temp_next)
		if index != -1 {
			fmt.Println("Partition found")
			fmt.Println("Creando Particion logica")
			var TempEBR Structs.EBR

			if err := Utilities.Leer_Object(file, &TempEBR, 0); err != nil {
				return
			}

			var count = 0
			var gap = int32(0)
			var gap2 = int32(0)

			for i := 0; i < 15; i++ {
				if TempMBR.M_EBR[i].E_part_s != 0 {
					count++
					gap = TempMBR.M_EBR[i].E_part_start + TempMBR.M_EBR[i].E_part_s + inicio
				}
			}

			for i := 0; i < 15; i++ {
				if TempMBR.M_EBR[i].E_part_s == 0 {
					TempMBR.M_EBR[i].E_part_s = int32(size)

					if count == 0 {
						TempMBR.M_EBR[i].E_part_start = inicio
						TempMBR.M_EBR[i].E_part_next = 0
					} else {
						TempMBR.M_EBR[i].E_part_start = gap
						TempMBR.M_EBR[i].E_part_next = gap + int32(size)

					}

					copy(TempMBR.M_EBR[i].E_part_name[:], name)
					copy(TempMBR.M_EBR[i].E_part_fit[:], fit)
					TempMBR.M_EBR[i].E_part_next = gap2
					copy(TempMBR.M_EBR[i].E_part_mounth[:], "0")
					fmt.Println(fmt.Sprintf("Name: %s, mount: %s, start: %d, size: %d, next: %s", string(TempMBR.M_EBR[i].E_part_name[:]), string(TempMBR.M_EBR[i].E_part_mounth[:]), TempMBR.M_EBR[i].E_part_start, TempMBR.M_EBR[i].E_part_s, TempMBR.M_EBR[i].E_part_next))
					break
				}
			}
			if err := Utilities.Escribir_Object(file, TempMBR, 0); err != nil {
				return

			}

			var TempEBR2 Structs.EBR

			if err := Utilities.Leer_Object(file, &TempEBR2, 0); err != nil {
				return
			}

			defer file.Close()
		} else {
			fmt.Println("No se ha creado nunguna particion extendida")
			return
		}
	} else if type_ == "e" {
		var index int = -1
		var count = 0

		for i := 0; i < 4; i++ {
			if TempMBR.M_Partitions[i].Size != 0 {
				count++
				if strings.Contains(string(TempMBR.M_Partitions[i].Type[:]), "e") {
					index = i
					break
				}
			}
		}

		if index != -1 {
			fmt.Println("No se pueden tener 2 particiones extendidas en este disco")

		} else {
			fmt.Println("Creando particion extendida")
			var count = 0
			var gap = int32(0)

			for i := 0; i < 4; i++ {
				if TempMBR.M_Partitions[i].Size != 0 {
					count++
					gap = TempMBR.M_Partitions[i].Start + TempMBR.M_Partitions[i].Size
				}
			}

			for i := 0; i < 4; i++ {
				if TempMBR.M_Partitions[i].Size == 0 {
					TempMBR.M_Partitions[i].Size = int32(size)

					if count == 0 {
						TempMBR.M_Partitions[i].Start = int32(binary.Size(TempMBR))
					} else {
						TempMBR.M_Partitions[i].Start = gap
					}

					copy(TempMBR.M_Partitions[i].Name[:], name)
					copy(TempMBR.M_Partitions[i].Fit[:], fit)
					copy(TempMBR.M_Partitions[i].Status[:], "0")
					copy(TempMBR.M_Partitions[i].Type[:], type_)
					TempMBR.M_Partitions[i].Correlative = int32(count + 1)
					break
				}
			}

			if err := Utilities.Escribir_Object(file, TempMBR, 0); err != nil {
				return

			}

			var TempMBR2 Structs.MRB

			if err := Utilities.Leer_Object(file, &TempMBR2, 0); err != nil {
				return
			}

			Structs.PrintMBR(TempMBR2)

			defer file.Close()
		}
	} else {

		var count = 0
		var gap = int32(0)

		for i := 0; i < 4; i++ {
			if TempMBR.M_Partitions[i].Size != 0 {
				count++
				gap = TempMBR.M_Partitions[i].Start + TempMBR.M_Partitions[i].Size
			}
		}

		for i := 0; i < 4; i++ {
			if TempMBR.M_Partitions[i].Size == 0 {
				TempMBR.M_Partitions[i].Size = int32(size)

				if count == 0 {
					TempMBR.M_Partitions[i].Start = int32(binary.Size(TempMBR))
				} else {
					TempMBR.M_Partitions[i].Start = gap
				}

				copy(TempMBR.M_Partitions[i].Name[:], name)
				copy(TempMBR.M_Partitions[i].Fit[:], fit)
				copy(TempMBR.M_Partitions[i].Status[:], "0")
				copy(TempMBR.M_Partitions[i].Type[:], type_)
				TempMBR.M_Partitions[i].Correlative = int32(count + 1)
				break
			}
		}

		if err := Utilities.Escribir_Object(file, TempMBR, 0); err != nil {
			return

		}

		var TempMBR2 Structs.MRB

		if err := Utilities.Leer_Object(file, &TempMBR2, 0); err != nil {
			return
		}

		Structs.PrintMBR(TempMBR2)

		defer file.Close()
	}

	fmt.Println("======End FDISK======")
}

func Mkdisk(size int, fit string, unit string, numDisco int32) {
	a_disco := numDisco
	lista := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L"}
	NUMERO := lista[a_disco]
	filepath := "./MIA/P1/" + strings.ToUpper(NUMERO) + ".dsk"
	fmt.Println("======Start MKDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Fit:", fit)
	fmt.Println("Unit:", unit)
	fit = string(fit[1])
	identificador := ".disk"
	nombre_disk = NUMERO + identificador

	discos_creados = append(discos_creados, nombre_disk)
	generarJson_disco(discos_creados)

	if fit != "b" && fit != "w" && fit != "f" {
		fmt.Println("Error: Fit must be b, w or f")
		return
	}

	if size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	if unit != "k" && unit != "m" {
		fmt.Println("Error: Unit must be k or m")
		return
	}

	err := Utilities.Crear_File(filepath)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	if unit == "k" {
		size = size * 1024
	} else {
		size = size * 1024 * 1024
	}

	file, err := Utilities.AbrirFile(filepath)
	if err != nil {
		return
	}

	arreglo := make([]byte, 1024)
	// create array of byte(0)
	for i := 0; i <= size/1024; i++ {
		err := Utilities.Escribir_Object(file, arreglo, int64(i*1024))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
	fecha := Utilities.FechaActual()

	var newMRB Structs.MRB
	newMRB.M_Size = int32(size)
	newMRB.M_Signature = 10
	copy(newMRB.M_Fit[:], fit)
	copy(newMRB.M_CreationDate[:], []byte(fecha))

	if err := Utilities.Escribir_Object(file, newMRB, 0); err != nil {
		return
	}

	var TempMBR Structs.MRB

	if err := Utilities.Leer_Object(file, &TempMBR, 0); err != nil {
		return
	}

	Structs.PrintMBR(TempMBR)

	defer file.Close()

	fmt.Println("======End MKDISK======")

}
func RMdisk(driveletter string) {
	var op int32
	fmt.Println("Driveletter:", driveletter)

	filepath := "./MIA/P1/" + strings.ToUpper(driveletter) + ".dsk"
	fmt.Println(filepath)
	fmt.Println("¿Desea eliminar el archivo?")
	fmt.Println("1. Si")
	fmt.Println("2. NO")
	fmt.Println("Elija una opcion")
	fmt.Scanln(&op)
	switch op {
	case 1:
		err := os.Remove(filepath)
		if err != nil {
			fmt.Println("Error al eliminar el archivo,archivo no encontrado", err)
			return
		}
		fmt.Println("Se elimino el disco")
		fmt.Println("")
		return
	case 2:
		fmt.Println("Se mantiene el disco")
	default:
		fmt.Println("Opcion invalida")
		return
	}
}
func Imagen_disco(id string, path string, name string) {
	fmt.Println("Id: ", id)
	fmt.Println("Carpeta: ", path)
	fmt.Println("Nombre: ", name)
	driveletter := string(id[0])

	var nombre string
	nombre = name

	filepath := "./MIA/P1/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := Utilities.AbrirFile(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	var TempMBR Structs.MRB
	var temp_nom string
	if err := Utilities.Leer_Object(file, &TempMBR, 0); err != nil {
		return
	}

	if name == "disk" {
		for i := 0; i < 4; i++ {
			if TempMBR.M_Partitions[i].Size != 0 {
				temp := string(TempMBR.M_Partitions[i].Name[:])
				temp2 := string(TempMBR.M_Partitions[i].Type[:])

				temp_nom = string(temp[:aux_tamaño])
				fmt.Println(temp_nom)

				if temp2 == "e" {

					nombre += "|" + "{Extendida|{ "
					for i := 0; i < 15; i++ {
						if TempMBR.M_EBR[i].E_part_s != 0 {

							nombre += "|" + "EBR" + "|" + "Logica"
						}

					}
					nombre += "}}"

				} else {
					nombre += "|" + "Primaria"
				}
			} else {
				nombre += "|" + "Libre"
			}

		}

		grafo(nombre, path)

		fmt.Println("Grafo generado exitosamente en", path)
	} else if name == "mbr" {
		/////////////////////////////////////////////////////////////////
		fmt.Println("generando reporte de mbr")
		var temp string
		for i := 0; i < len(path); i++ {
			if path[i] == '.' {
				break
			}

			temp += string(path[i])
		}

		dotFilePath := temp

		// Solicitar la información de la tabla interactivamente
		attributes := make(map[string]string)
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Println("Ingresa el nombre del atributo (o escribe 'fin' para terminar):")
			scanner.Scan()
			attributeName := scanner.Text()
			if attributeName == "fin" {
				break
			}
			fmt.Println("Ingresa el valor del atributo:")
			scanner.Scan()
			attributeValue := scanner.Text()
			attributes[attributeName] = attributeValue
		}

		// Generar el contenido del archivo DOT
		dotContent := generateDOT(attributes)

		// Escribir el contenido en el archivo DOT
		err := writeDOTFile(dotFilePath, dotContent)
		if err != nil {
			fmt.Println("Error al escribir el archivo DOT:", err)
			return
		}

		fmt.Println("Archivo DOT generado exitosamente:", dotFilePath)

		// Generar la imagen a partir del archivo DOT
		path = path + "disk.png"
		err = generateImageFromDOT(dotFilePath, path)
		if err != nil {
			fmt.Println("Error al generar la imagen:", err)
			return
		}

		fmt.Println("Imagen generada exitosamente:", path)
	}

}
func grafo(datos string, direccion string) error {
	graphCode := fmt.Sprintf(`digraph structs {
        bgcolor="#68d9e2"
        node [shape=record];
        structs [label="%s"];
    }`, datos)

	var temp string
	for i := 0; i < len(direccion); i++ {
		if direccion[i] == '.' {
			break
		}

		temp += string(direccion[i])
	}
	dotPath := temp + ".dot"
	direccion = direccion + "disk.png"

	// Abre el archivo en modo escritura
	file, err := os.Create(dotPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Escribe el código Graphviz en el archivo
	_, err = file.WriteString(graphCode)
	if err != nil {
		return err
	}

	// Ejecuta el comando "dot" de Graphviz para generar la imagen desde el archivo DOT
	cmd := exec.Command("dot", "-Tpng", "-o", direccion, dotPath)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil

}

func generateDOT(attributes map[string]string) string {
	dot := `digraph G {
  a0 [shape=none label=<
  <TABLE cellspacing="10" cellpadding="10" style="rounded" bgcolor="#6495ED">
  
  <TR>
  <TD bgcolor="yellow">REPORTE MBR</TD>
  </TR>
  `
	for key, value := range attributes {
		dot += fmt.Sprintf(`  
  <TR>
  <TD bgcolor="yellow">%s</TD>
  <TD bgcolor="yellow">%s</TD>
  </TR>
  `, key, value)
	}

	dot += `
  <TR>
  <TD bgcolor="#D3D3D3">Particion</TD>
  </TR>
  
</TABLE>>];
}`
	return dot
}

func writeDOTFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func generateImageFromDOT(dotFilePath, imageFilePath string) error {
	cmd := exec.Command("dot", "-Tpng", dotFilePath, "-o", imageFilePath)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
func generarJson_disco(discos []string) error {
	// Crear un mapa para almacenar los discos bajo el identificador "Discos"
	discosMap := map[string][]string{"Discos": discos}

	// Serializar el mapa a formato JSON
	jsonData, err := json.Marshal(discosMap)
	if err != nil {
		return err
	}
	// Escribir los datos serializados en un archivo
	file, err := os.Create("discos.json")
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
