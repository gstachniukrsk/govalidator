package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gstachniukrsk/govalidator"
)

// Employee represents an employee record from CSV
type Employee struct {
	ID         string
	FirstName  string
	LastName   string
	Email      string
	Department string
	Salary     string
	Status     string
}

// Define validation schema for employee records
var (
	emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	employeeSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("id").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(1),
			),
		govalidator.NewField("firstName").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(2),
				govalidator.MaxLengthValidator(50),
			),
		govalidator.NewField("lastName").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(2),
				govalidator.MaxLengthValidator(50),
			),
		govalidator.NewField("email").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.RegexpValidator(*emailPattern),
			),
		govalidator.NewField("department").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.OneOfValidator("Engineering", "Sales", "Marketing", "HR", "Finance"),
			),
		govalidator.NewField("salary").
			Required().
			WithValidators(
				govalidator.NumberValidator,
				govalidator.MinFloatValidator(0),
				govalidator.MaxFloatValidator(1000000),
			),
		govalidator.NewField("status").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.OneOfValidator("active", "inactive", "on-leave"),
			),
	).WithExtra(govalidator.ExtraForbid)
)

// convertEmployeeToMap converts an employee record to a map for validation
func convertEmployeeToMap(emp Employee) map[string]any {
	data := map[string]any{
		"id":         emp.ID,
		"firstName":  emp.FirstName,
		"lastName":   emp.LastName,
		"email":      emp.Email,
		"department": emp.Department,
		"status":     emp.Status,
	}

	// Convert salary string to float for validation
	if salary, err := strconv.ParseFloat(emp.Salary, 64); err == nil {
		data["salary"] = salary
	} else {
		data["salary"] = emp.Salary // Keep as string so validator can catch it
	}

	return data
}

// ValidationResult holds the validation result for a single row
type ValidationResult struct {
	RowNumber int
	Employee  Employee
	Valid     bool
	Errors    []string
}

// validateCSV reads and validates a CSV file
func validateCSV(filename string) ([]ValidationResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file must have a header row and at least one data row")
	}

	// Skip header row
	results := make([]ValidationResult, 0, len(records)-1)

	for i, record := range records[1:] {
		if len(record) != 7 {
			results = append(results, ValidationResult{
				RowNumber: i + 2, // +2 because we skip header and arrays are 0-indexed
				Valid:     false,
				Errors:    []string{"Row must have exactly 7 columns"},
			})
			continue
		}

		emp := Employee{
			ID:         record[0],
			FirstName:  record[1],
			LastName:   record[2],
			Email:      record[3],
			Department: record[4],
			Salary:     record[5],
			Status:     record[6],
		}

		// Convert to map for validation
		data := convertEmployeeToMap(emp)

		// Validate
		valid, errs := employeeSchema.ValidateFlat(
			context.Background(),
			data,
			govalidator.CombinedPresenter(".", ": "),
		)

		results = append(results, ValidationResult{
			RowNumber: i + 2,
			Employee:  emp,
			Valid:     valid,
			Errors:    errs,
		})
	}

	return results, nil
}

// createSampleCSV creates a sample CSV file for testing
func createSampleCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "FirstName", "LastName", "Email", "Department", "Salary", "Status"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write sample data (mix of valid and invalid)
	records := [][]string{
		{"1", "John", "Doe", "john.doe@example.com", "Engineering", "75000", "active"},
		{"2", "Jane", "Smith", "jane.smith@example.com", "Sales", "65000", "active"},
		{"3", "Bob", "J", "bob@invalid", "IT", "55000", "active"}, // Invalid: last name too short, invalid email, invalid department
		{"4", "Alice", "Johnson", "alice@example.com", "Marketing", "70000", "active"},
		{"5", "Charlie", "Brown", "charlie.brown@example.com", "Finance", "-1000", "active"}, // Invalid: negative salary
		{"6", "D", "Williams", "david@example.com", "HR", "60000", "retired"},                // Invalid: first name too short, invalid status
		{"7", "Emma", "Davis", "emma.davis@example.com", "Engineering", "85000", "on-leave"},
		{"8", "", "Wilson", "wilson@example.com", "Sales", "62000", "inactive"}, // Invalid: missing first name
	}

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// printResults prints validation results in a readable format
func printResults(results []ValidationResult) {
	validCount := 0
	invalidCount := 0

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("CSV VALIDATION RESULTS")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	for _, result := range results {
		if result.Valid {
			validCount++
			fmt.Printf("✓ Row %d: VALID - %s %s (%s)\n",
				result.RowNumber,
				result.Employee.FirstName,
				result.Employee.LastName,
				result.Employee.Email,
			)
		} else {
			invalidCount++
			fmt.Printf("✗ Row %d: INVALID - %s %s\n",
				result.RowNumber,
				result.Employee.FirstName,
				result.Employee.LastName,
			)
			for _, err := range result.Errors {
				fmt.Printf("    - %s\n", err)
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("Summary: %d valid, %d invalid out of %d total rows\n",
		validCount, invalidCount, len(results))
	fmt.Println(strings.Repeat("=", 80) + "\n")
}

// exportInvalidRowsJSON exports invalid rows to a JSON file for further processing
func exportInvalidRowsJSON(results []ValidationResult, filename string) error {
	invalidRows := make([]ValidationResult, 0)
	for _, result := range results {
		if !result.Valid {
			invalidRows = append(invalidRows, result)
		}
	}

	if len(invalidRows) == 0 {
		fmt.Println("No invalid rows to export")
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(invalidRows); err != nil {
		return err
	}

	fmt.Printf("Exported %d invalid rows to %s\n", len(invalidRows), filename)
	return nil
}

func main() {
	csvFile := "employees.csv"
	jsonFile := "invalid_rows.json"

	// Create sample CSV file
	fmt.Println("Creating sample CSV file...")
	if err := createSampleCSV(csvFile); err != nil {
		log.Fatalf("Failed to create sample CSV: %v", err)
	}
	fmt.Printf("Created sample file: %s\n", csvFile)

	// Validate CSV
	fmt.Println("\nValidating CSV file...")
	results, err := validateCSV(csvFile)
	if err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	// Print results
	printResults(results)

	// Export invalid rows to JSON
	if err := exportInvalidRowsJSON(results, jsonFile); err != nil {
		log.Fatalf("Failed to export invalid rows: %v", err)
	}

	fmt.Println("\nValidation complete!")
	fmt.Printf("Check '%s' for the sample data\n", csvFile)
	fmt.Printf("Check '%s' for invalid rows details\n", jsonFile)
}
