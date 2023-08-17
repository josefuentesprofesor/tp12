package main

import (
	"os"
	"testing"
)

func TestSaveAndLoadToCSV(t *testing.T) {
	filename := "test_data.csv"
	defer os.Remove(filename)

	deleteAllRecords()
	createRecord("Source1", 10.5, "Event1")
	createRecord("Source2", 20.0, "Event2")

	err := saveToCSV(filename)
	if err != nil {
		t.Errorf("Error al guardar en el archivo CSV: %v", err)
	}

	err = loadFromCSV(filename)
	if err != nil {
		t.Errorf("Error cargando del archivo CSV: %v", err)
	}

	if len(records) != 2 {
		t.Errorf("Se esperaban 2 registros cargados del CSV, se cargaron %d", len(records))
	}

	// Verificar valores de cada registro
	record := records[0]
	if record.PKID != 1 || record.Source != "Source1" || record.Measurement != 10.5 || record.Event != "Event1" {
		t.Errorf("Los registros leidos no coinciden con lo esperado")
	}

	record = records[1]
	if record.PKID != 2 || record.Source != "Source2" || record.Measurement != 20.0 || record.Event != "Event2" {
		t.Errorf("Los registros leidos no coinciden con lo esperado")
	}
}

func TestCRUDOperations(t *testing.T) {
	// Prueba de creación
	deleteAllRecords()
	createRecord("Source1", 10.5, "Event1")
	createRecord("Source2", 20.0, "Event2")

	if len(records) != 2 {
		t.Errorf("Se esperaban 2 registros despues de la creacion, se obtuvieron %d", len(records))
	}

	// Prueba de lectura
	record := records[0]
	if record.PKID != 1 || record.Source != "Source1" || record.Measurement != 10.5 || record.Event != "Event1" {
		t.Errorf("El registro leido no coincide con lo esperado")
	}

	// Prueba de lectura de todos los registros
	readAllRecords()

	// Prueba de actualización
	updateRecord(2, "UpdatedSource", 15.0, "UpdatedEvent")
	record = records[1]
	if record.Source != "UpdatedSource" || record.Measurement != 15.0 || record.Event != "UpdatedEvent" {
		t.Errorf("Falló la actualización")
	}

	// Prueba de eliminación
	deleteRecord(1)
	if len(records) != 1 {
		t.Errorf("Se esperaba 1 registro luego del borrado, se obtuvo %d", len(records))
	}

	// Prueba de borrar todos los registros
	deleteAllRecords()
	if len(records) != 0 {
		t.Errorf("El borrado de todos los registros falló")
	}
}

func TestMain(m *testing.M) {
	filename := "test_data.csv"
	// Crear un archivo de prueba
	file, _ := os.Create(filename)
	file.Close()

	// Ejecutar pruebas
	result := m.Run()

	// Limpiar después de las pruebas
	os.Remove(filename)

	os.Exit(result)
}
