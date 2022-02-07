package main

import (
	"fmt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/log"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func run() {
	errCounter := 0
	complete := false
	for i := 2; !complete; i++ {
		if getAndSaveFile(i) {
			errCounter += 1
			if errCounter > 2 {
				complete = true
			}
		}
	}
}

func main() {

	run()
	//url := "https://www.biorxiv.org/content/10.1101/703801v4.full.pdf"
	//url := "https://www.thelancet.com/action/showPdf?pii=S2589-5370%2820%2930228-5"
	//url := ""
	////saveFilePath := "/covid-19/cerebral_micro-structure_changes_3_month_follow-up_60-patients.pdf"
	//saveFilePath := ""
	//err := downloadFile(url, saveFilePath)
	//if err != nil {
	//	fmt.Printf("Error downloading file: %v\n", err)
	//}

}

func getAndSaveFile(i int) bool {
	numberString := strconv.Itoa(i)
	//if i < 10 {
	//	numberString = "0" + strconv.Itoa(i)
	//}
	//filePath := "MIT-8_05-QM-II-" + numberString + ".pdf"
	//filePath := "MATH10111-Exercise-Solutions-" + numberString + ".pdf"
	filePath := "GR/HW" + numberString + "Qsolutions.pdf"

	fileUrl := "https://inside.mines.edu/~aflourno/GR/HW" + numberString + "Qsolutions.pdf"
	//"https://inside.mines.edu/~aflourno/GR/HW2Q.pdf"

	//fileUrl := "https://ocw.mit.edu/courses/physics/8-04-quantum-physics-i-spring-2013/lecture-notes/MIT8_04S13_Lec" + numberString + ".pdf"
	//fileUrl := "https://ocw.mit.edu/courses/physics/8-05-quantum-physics-ii-fall-2013/lecture-notes/MIT8_05F13_Chap_" + numberString + ".pdf"

	//fileUrl := "https://personalpages.manchester.ac.uk/staff/marianne.johnson/lecture" + numberString + "(2019)nopause.pdf"
	//fileUrl := "https://personalpages.manchester.ac.uk/staff/marianne.johnson/ExSol" + numberString + ".pdf"


	err := downloadFile(fileUrl, filePath)
	if err != nil {
		fmt.Println("Error")
		return true
	}
	return false
}

func downloadFile(url string, filepath string) error {
	// Fetch File data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}
	if len(dataBytes) < 50000 {
		return fmt.Errorf("no file found")
	}

		writeToFile(filepath, dataBytes)

	return err
}

func writeToFile(filepath string, bytes []byte) {
	fileCreated := createFile(filepath)
	if !fileCreated {
		fmt.Println("ERROR: File not created")
		return
	}

	file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		log.Error(err)
		return
	}

	_, err = file.Write(bytes)
	if err != nil {
		log.Error(err)
		return
	}

	err = file.Sync()
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Println("File successfully uploaded")
}

/*
Create file
 - If already created, return true.
 - If successfully created, return true.
 - If unsuccessful, return false

 */
func createFile(filepath string) bool {
	if fileDoesExist(filepath) {
		fmt.Println("File", filepath, "does exist.")
		return true
	}

	file, err := os.Create(filepath)
	defer file.Close()
	if err != nil {
		fmt.Println("File has not been created")
		return false
	}

	fmt.Println("File", filepath, "has been created")
	return true
}

func fileDoesExist(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func deleteFile(filepath string) {

}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
