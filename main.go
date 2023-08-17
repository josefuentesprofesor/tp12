package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Record struct {
	Timestamp   time.Time
	PKID        int
	Source      string
	Measurement float64
	Event       string
}

var records []Record
var lastID int

func saveToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		row := []string{
			record.Timestamp.Format(time.RFC3339),
			strconv.Itoa(record.PKID),
			record.Source,
			strconv.FormatFloat(record.Measurement, 'f', -1, 64),
			record.Event,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func loadFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Permitir cantidad variable de campos

	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	records = []Record{}
	lastID = 0

	for _, row := range rows {
		timestamp, _ := time.Parse(time.RFC3339, row[0])
		pkid, _ := strconv.Atoi(row[1])
		measurement, _ := strconv.ParseFloat(row[3], 64)

		record := Record{
			Timestamp:   timestamp,
			PKID:        pkid,
			Source:      row[2],
			Measurement: measurement,
			Event:       row[4],
		}
		records = append(records, record)

		if pkid > lastID {
			lastID = pkid
		}
	}

	return nil
}

func createRecord(source string, measurement float64, event string) {
	lastID++
	record := Record{
		Timestamp:   time.Now(),
		PKID:        lastID,
		Source:      source,
		Measurement: measurement,
		Event:       event,
	}
	records = append(records, record)
}

func readAllRecords() {
	//TODO Recorrer la coleccion records, acceder a
	//todos los campos de cada record
	// (PKID, Timestamp, Source, etc)
	//e imprimirlos por consola
	/*	for ???? range ???? {
			//fmt.Printf("ID: %d\n", record.PKID)
		}
	*/
}

func updateRecord(pkid int, source string, measurement float64, event string) error {
	for i, record := range records {
		if record.PKID == pkid {
			records[i].Source = source
			records[i].Measurement = measurement
			records[i].Event = event
			return nil
		}
	}
	return fmt.Errorf("Record with PKID %d not found", pkid)
}

func deleteRecord(pkid int) error {
	for i, record := range records {
		if record.PKID == pkid {
			records = append(records[:i], records[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Record with PKID %d not found", pkid)
}

func deleteAllRecords() {
	records = []Record{}
	lastID = 0
}

func main() {
	// Cargar datos desde el archivo CSV si existe
	filename := "data.csv"
	if _, err := os.Stat(filename); err == nil {
		err := loadFromCSV(filename)
		if err != nil {
			fmt.Println("Error loading data from CSV:", err)
			return
		}
	}

	var choice int
	for {
		fmt.Println("1. Crear registro")
		fmt.Println("2. Leer todos los registros")
		fmt.Println("3. Actualizar registro")
		fmt.Println("4. Borrar registro")
		fmt.Println("5. Borrar todos los registros")
		fmt.Println("6. Salir")
		fmt.Print("Seleccione una opción: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var source, event string
			var measurement float64
			fmt.Print("Ingrese la fuente: ")
			fmt.Scanln(&source)
			fmt.Print("Ingrese la medida: ")
			fmt.Scanln(&measurement)
			fmt.Print("Ingrese el evento: ")
			fmt.Scanln(&event)
			createRecord(source, measurement, event)
		case 2:
			readAllRecords()
		case 3:
			var pkid int
			fmt.Print("Ingrese el ID del registro a actualizar: ")
			fmt.Scanln(&pkid)
			var source, event string
			var measurement float64
			fmt.Print("Ingrese la nueva fuente: ")
			fmt.Scanln(&source)
			fmt.Print("Ingrese la nueva medida: ")
			fmt.Scanln(&measurement)
			fmt.Print("Ingrese el nuevo evento: ")
			fmt.Scanln(&event)
			err := updateRecord(pkid, source, measurement, event)
			if err != nil {
				fmt.Println(err)
			}
		case 4:
			var pkid int
			fmt.Print("Ingrese el ID del registro a borrar: ")
			fmt.Scanln(&pkid)
			err := deleteRecord(pkid)
			if err != nil {
				fmt.Println(err)
			}
		case 5:
			deleteAllRecords()
			fmt.Println("Todos los registros han sido borrados.")
		case 6:
			// Guardar los datos en el archivo CSV antes de salir
			err := saveToCSV(filename)
			if err != nil {
				fmt.Println("Error saving data to CSV:", err)
			}
			fmt.Println("Saliendo...")
			return
		default:
			fmt.Println("Opción no válida")
		}
	}
}
