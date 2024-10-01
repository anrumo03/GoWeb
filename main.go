package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func makeHandler(imageBase64 []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		lista := escogerImagenes(imageBase64)
		hostname, erro := os.Hostname()
		if erro != nil {
			fmt.Println("Error al obtener el nombre del host:", erro)
			return
		}

		data := struct {
			Title    string
			Hostname string
			Imagen1  string
			Imagen2  string
			Imagen3  string
		}{
			Title:    "Servidor de Imágenes",
			Hostname: hostname,
			Imagen1:  lista[0],
			Imagen2:  lista[1],
			Imagen3:  lista[2],
		}

		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error al renderizar la plantilla: %v", err), http.StatusInternalServerError)
		}
	}
}

func obtenerImagenes(directory string) []string {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	var imageFiles []string

	// Iterar y obtener los nombres de los archivos
	for _, file := range files {
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if imageExtensions[ext] {
				imageFiles = append(imageFiles, file.Name())
			}
		}
	}

	return imageFiles
}

func convertir64(arr []string, ruta string) []string {

	var images64 []string

	for _, file := range arr {
		path := ruta + "\\" + file
		imageData, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		encodedImage := base64.StdEncoding.EncodeToString(imageData)
		images64 = append(images64, encodedImage)
	}

	return images64
}

func escogerImagenes(imagenes []string) []string {

	var lista []string
	selected := make(map[int]bool)

	for i := 0; i < 3; i++ {
		var randomIndex int
		for {
			// Obtén un índice aleatorio basado en el tamaño del slice
			randomIndex = rand.Intn(len(imagenes))
			// Comprueba si el índice ya ha sido seleccionado
			if !selected[randomIndex] {
				break
			}
		}
		// Marca el índice como seleccionado
		selected[randomIndex] = true

		// Accede al elemento aleatorio
		randomElement := imagenes[randomIndex]
		lista = append(lista, randomElement)
	}

	return lista
}

func main() {

	// Definir un flag para el directorio de imágenes
	dirFlag := flag.String("dir", "", "Ruta del directorio que contiene las imágenes")
	flag.Parse()

	// Verificar si se proporcionó el directorio
	if *dirFlag == "" {
		fmt.Println("Debe especificar la ruta del directorio usando el flag -dir")
		return
	}

	imageFiles := obtenerImagenes(*dirFlag)
	images64 := convertir64(imageFiles, *dirFlag)

	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))

	// Manejar la ruta principal "/"
	http.HandleFunc("/", makeHandler(images64))

	// Iniciar el servidor en el puerto 8080
	fmt.Println("Servidor ejecutándose en http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
