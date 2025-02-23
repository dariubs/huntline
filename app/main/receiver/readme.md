# Huntline Scheduler

The Huntline Scheduler is a command-line application developed in Go for retrieving top products from the ProductHunt API and persisting them to a database. Designed with operational flexibility in mind, the application features a dynamic scheduling mechanism that allows for both repeatable (default) and one-off execution of tasks. The architecture and parameterization have been developed in accordance with rigorous management science methodologies, ensuring that model sensitivity and execution control are maintained at the highest standards.

## Features

- **Dynamic Scheduling:** Executes tasks at a user-specified time (default: 00:30 PST).
- **Repeatability Control:** Toggle between continuous (repeatable) and single execution modes.
- **Immediate Execution Option:** An optional parameter to run the task immediately prior to scheduling.
- **Date Customization:** Override the default “yesterday” date with a custom date (format: YYYY-MM-DD).
- **Database Integration:** Persists fetched product data using a configurable database connection.

## Requirements

- **Go:** Version 1.16 or higher is recommended.
- **Environment Variables:** A `.env` file containing:
  - `PH_API_KEY` – Your ProductHunt API key.
- **Database:** A properly configured database, as defined in the `db.ConnectToDB()` implementation.

## Installation

1. **Clone the Repository:**

   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. **Install Dependencies:**

   Run the following command to ensure that all required packages are available:

   ```bash
   go mod tidy
   ```

3. **Configure Environment Variables:**

   Create a `.env` file in the root directory and add the following:

   ```env
   PH_API_KEY=your_producthunt_api_key
   ```

## Usage

The application supports several command-line parameters, which facilitate precise control over task execution. Below is a detailed guide to each parameter.

### Command-Line Parameters

- **`-run-now`**  
  **Description:** Executes the task immediately before commencing the scheduled runs.  
  **Type:** Boolean flag  
  **Default:** `true`  
  **Usage Example:**

  ```bash
  go run main.go -run-now=false
  ```

- **`-date`**  
  **Description:** Specifies the date (in the format `YYYY-MM-DD`) for which product data should be fetched, overriding the default “yesterday” parameter.  
  **Type:** String flag  
  **Default:** Empty (defaults to yesterday’s date based on PST)  
  **Usage Example:**

  ```bash
  go run main.go -date 2025-02-22
  ```

- **`-repeat`**  
  **Description:** Determines whether the task should execute repeatedly at the scheduled time. Setting this flag to `false` results in a one-off execution.  
  **Type:** Boolean flag  
  **Default:** `false`  
  **Usage Example:**

  ```bash
  go run main.go -repeat=true
  ```

- **`-schedule`**  
  **Description:** Specifies the scheduled time for task execution in 24-hour format (HH:MM), interpreted in the PST timezone.  
  **Type:** String flag  
  **Default:** `"00:30"`  
  **Usage Example:**

  ```bash
  go run main.go -schedule 13:45 -repeat=true
  ```

## License

This project is licensed under the MIT License. For further details, please refer to the [LICENSE](LICENSE) file.
