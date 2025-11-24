# CSV Validator Example

This example demonstrates how to use govalidator to validate CSV file data, such as employee records or data imports.

## Features

- Reads and validates CSV files row by row
- Converts CSV rows to structured data for validation
- Provides detailed error reports for each invalid row
- Exports invalid rows to JSON for further processing
- Handles type conversions (e.g., string to number)
- Demonstrates practical data import validation patterns

## Running the Example

```bash
cd examples/csv-validator
go run main.go
```

The program will:
1. Create a sample `employees.csv` file with valid and invalid records
2. Validate each row against the schema
3. Print validation results to console
4. Export invalid rows to `invalid_rows.json`

## Sample Output

```
Creating sample CSV file...
Created sample file: employees.csv

Validating CSV file...

================================================================================
CSV VALIDATION RESULTS
================================================================================

✓ Row 2: VALID - John Doe (john.doe@example.com)
✓ Row 3: VALID - Jane Smith (jane.smith@example.com)
✗ Row 4: INVALID - Bob J
    - $.lastName: text must be at least 2 characters long
    - $.email: does not match required pattern
    - $.department: value must be one of the allowed values
✓ Row 5: VALID - Alice Johnson (alice@example.com)
✗ Row 6: INVALID - Charlie Brown
    - $.salary: number must be >= 0
✗ Row 7: INVALID - D Williams
    - $.firstName: text must be at least 2 characters long
    - $.status: value must be one of the allowed values
✓ Row 8: VALID - Emma Davis (emma.davis@example.com)
✗ Row 9: INVALID -  Wilson
    - $.firstName: text must be at least 2 characters long

================================================================================
Summary: 4 valid, 4 invalid out of 8 total rows
================================================================================
```

## CSV Format

The example validates employee CSV files with the following columns:

| Column     | Description                                    |
|------------|------------------------------------------------|
| ID         | Employee ID (required, non-empty string)       |
| FirstName  | First name (required, 2-50 characters)         |
| LastName   | Last name (required, 2-50 characters)          |
| Email      | Email address (required, valid email format)   |
| Department | Department (required, one of predefined list)  |
| Salary     | Salary (required, number 0-1,000,000)          |
| Status     | Employment status (required, predefined values)|

## Validation Rules

- **ID**: Required, non-empty string
- **FirstName**: Required, 2-50 characters
- **LastName**: Required, 2-50 characters
- **Email**: Required, valid email format (regex validation)
- **Department**: Required, must be one of: "Engineering", "Sales", "Marketing", "HR", "Finance"
- **Salary**: Required, number between 0 and 1,000,000
- **Status**: Required, must be one of: "active", "inactive", "on-leave"

## Key Patterns Demonstrated

1. **CSV Processing**: Reading and parsing CSV files with the standard library
2. **Data Conversion**: Converting CSV strings to appropriate types for validation
3. **Row-by-Row Validation**: Validating each row independently with error tracking
4. **Error Reporting**: Detailed validation errors with row numbers
5. **Data Export**: Exporting validation results to JSON for further processing
6. **Enum Validation**: Using `OneOfValidator` for restricted value sets
7. **Pattern Matching**: Email validation with regular expressions

## Use Cases

This pattern is useful for:
- Importing employee data from HR systems
- Validating bulk data uploads
- Processing customer records
- Validating financial transaction files
- Data migration quality checks
- ETL pipeline validation

## Extending the Example

You can modify this example for your own CSV validation needs:

1. Change the `Employee` struct to match your data structure
2. Update the `employeeSchema` with your validation rules
3. Modify `convertEmployeeToMap` to handle your data types
4. Customize the error reporting format in `printResults`
