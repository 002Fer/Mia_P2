package Analyzer

import (
	"MIA_P1/DiskManagement"
	"MIA_P1/FileSystem"
	"MIA_P1/User"
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

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

var disco int = 0

var Logicas_extendidas int = 0

func Analyze() {

	for true {
		var input string
		fmt.Println("Enter command: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input = scanner.Text()

		command, params := getCommandAndParams(input)

		fmt.Println("Command: ", command, "Params: ", params)

		AnalyzeCommnad(command, params)

	}
}

func AnalyzeCommnad(command string, params string) {

	if strings.Contains(command, "exec") {
		fn_exec(params)

	} else if strings.Contains(command, "mkdisk") {
		fn_mkdisk(params, int32(disco))
		disco++
	} else if strings.Contains(command, "unmount") {
		Unmount(params)
	} else if strings.Contains(command, "fdisk") {
		fn_fdisk(params)
	} else if strings.Contains(command, "mount") {
		fn_mount(params)
	} else if strings.Contains(command, "mkfs") {
		fn_mkfs(params)
	} else if strings.Contains(command, "rmdisk") {
		fn_rmdisk(params)
	} else if strings.Contains(command, "login") {
		fn_login(params)
	} else if strings.Contains(command, "logout") {
		fn_logout()
	} else if strings.Contains(command, "rep") {
		fn_reporte(params)
	} else {
		fmt.Println("Error: Command not found")
	}

}
func fn_logout() {
	User.Logout()
}
func correr_programa(correr string) {

	file, err := os.Open(correr)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		line := scanner.Text()

		command, params := getCommandAndParams(line)

		fmt.Println("Command: ", command, "Params: ", params)

		AnalyzeCommnad(command, params)
	}
}

func fn_exec(input string) {
	fs := flag.NewFlagSet("exec", flag.ExitOnError)
	path_ := fs.String("path", "", "direccion archivo")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}
	correr_programa(*path_)
}

func fn_login(input string) {
	// Define flags
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	user := fs.String("user", "", "Usuario")
	pass := fs.String("pass", "", "Contraseña")
	id := fs.String("id", "", "Id")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "user", "pass", "id":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// Call the function
	User.Login(*user, *pass, *id)

}
func fn_mkfs(input string) {

	fs := flag.NewFlagSet("mkfs", flag.ExitOnError)
	id := fs.String("id", "", "Id")
	type_ := fs.String("type", "", "Tipo")
	fs_ := fs.String("fs", "2fs", "Fs")

	fs.Parse(os.Args[1:])

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id", "type", "fs":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	FileSystem.Mkfs(*id, *type_, *fs_)

}

func fn_mount(input string) {

	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	driveletter := fs.String("driveletter", "", "Letra")
	name := fs.String("name", "", "Nombre")

	fs.Parse(os.Args[1:])

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "driveletter", "name":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	DiskManagement.Mount(*driveletter, *name)
}

func fn_fdisk(input string) {

	fs := flag.NewFlagSet("fdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	driveletter := fs.String("driveletter", "", "Letra")
	name := fs.String("name", "", "Nombre")
	unit := fs.String("unit", "m", "Unidad")
	type_ := fs.String("type", "p", "Tipo")
	fit := fs.String("fit", "bf", "Ajuste")

	fs.Parse(os.Args[1:])

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "driveletter", "name", "type":
			if flagValue == "e" {
				fs.Set(flagName, flagValue)
				fmt.Println("Particion logiva no: ", Logicas_extendidas)
				Logicas_extendidas++
			} else {
				fs.Set(flagName, flagValue)
				fmt.Println(Logicas_extendidas)
			}

		default:
			fmt.Println("Error: Flag not found")
		}
	}

	DiskManagement.Fdisk(*size, *driveletter, *name, *unit, *type_, *fit)
}

func fn_mkdisk(params string, disco int32) {

	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	fit := fs.String("fit", "bf", "Ajuste")
	unit := fs.String("unit", "m", "Unidad")
	num_disco := disco
	fmt.Println(num_disco)

	fs.Parse(os.Args[1:])

	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	DiskManagement.Mkdisk(*size, *fit, *unit, num_disco)

}
func fn_rmdisk(input string) {

	fs := flag.NewFlagSet("rmdisk", flag.ExitOnError)
	driveletter := fs.String("driveletter", "", "Letra")

	fs.Parse(os.Args[1:])

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "driveletter":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}
	DiskManagement.RMdisk(*driveletter)

}
func Unmount(input string) {

	fs := flag.NewFlagSet("unmount", flag.ExitOnError)
	id := fs.String("id", "", "Id")

	fs.Parse(os.Args[1:])

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}
	FileSystem.Unmount_1(*id)

}
func fn_reporte(input string) {

	fs := flag.NewFlagSet("rep", flag.ExitOnError)
	id := fs.String("id", "", "Id")
	path := fs.String("path", "", "Carpeta")
	name := fs.String("name", "", "Name")

	fs.Parse(os.Args[1:])

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id", "path", "name":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	DiskManagement.Imagen_disco(*id, *path, *name)

}

func correr(input string) {

	fs := flag.NewFlagSet("rep", flag.ExitOnError)
	id := fs.String("id", "", "Id")
	path := fs.String("path", "", "Carpeta")
	name := fs.String("name", "", "Name")

	fs.Parse(os.Args[1:])

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id", "path", "name":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	DiskManagement.Imagen_disco(*id, *path, *name)

}
