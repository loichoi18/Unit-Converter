# Unit Converter

A simple server-side web application built with **Go** and **HTML/CSS** that converts between different units of measurement. No database, no JavaScript frameworks — just a clean Go backend serving HTML templates with traditional form submissions.

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)

## Demo

**Input Screen** — Select a category tab, enter a value, choose units, and hit Convert.

**Result Screen** — The converted value is displayed below the form with a Reset option.

## Features

- **Three conversion categories** with tab-based navigation
  - **Length**: millimeter, centimeter, meter, kilometer, inch, foot, yard, mile
  - **Weight**: milligram, gram, kilogram, ounce, pound
  - **Temperature**: Celsius, Fahrenheit, Kelvin
- **Server-side rendering** — form POSTs to the same page, Go processes the conversion and re-renders with results
- **Single-file architecture** — all logic, handlers, and templates live in one `main.go`
- **Smart rounding** — results display up to 6 decimal places with trailing zeros trimmed
- **Input validation** — handles non-numeric input gracefully with error messages
- **Form state preservation** — selected units and input value persist after conversion

## How It Works

The application uses a base-unit intermediary pattern for conversions:

| Category    | Base Unit | Method                                              |
|-------------|-----------|-----------------------------------------------------|
| Length      | Meter     | `input × fromFactor / toFactor`                     |
| Weight      | Gram      | `input × fromFactor / toFactor`                     |
| Temperature | Celsius   | Explicit formula chain (e.g. `°F = °C × 9/5 + 32`) |

Each category has its own HTTP handler that:

1. Serves the form on `GET` requests
2. Parses form data, runs the conversion, and re-renders with results on `POST` requests

All three categories share a single Go `html/template` that conditionally renders the correct units and results.

## Getting Started

### Prerequisites

- [Go 1.22+](https://go.dev/dl/) installed on your machine

### Run

```bash
git clone https://github.com/loichoi18/Unit-Coverter
go run main.go
```

The server starts at **http://localhost:8080**. It automatically redirects `/` to `/length`.

### Build (optional)

```bash
go build -o unit-converter main.go
./unit-converter
```

## Project Structure

```
unit-converter/
├── main.go       # All application code (handlers, conversion logic, HTML template)
├── go.mod        # Go module definition
└── README.md
```

### Code Layout inside `main.go`

| Section              | Description                                                  |
|----------------------|--------------------------------------------------------------|
| Conversion Logic     | `convertLength`, `convertWeight`, `convertTemperature` funcs |
| Template Data        | `PageData` struct passed to the HTML template                |
| Handlers             | `lengthHandler`, `weightHandler`, `temperatureHandler`       |
| Template             | Inline HTML/CSS template using Go's `html/template`          |
| Main                 | Route registration and server startup                        |

## Routes

| Route          | Method | Description                          |
|----------------|--------|--------------------------------------|
| `/`            | GET    | Redirects to `/length`               |
| `/length`      | GET    | Length conversion form                |
| `/length`      | POST   | Processes length conversion           |
| `/weight`      | GET    | Weight conversion form                |
| `/weight`      | POST   | Processes weight conversion           |
| `/temperature` | GET    | Temperature conversion form           |
| `/temperature` | POST   | Processes temperature conversion      |

## Example Conversions

| Input          | From       | To         | Result          |
|----------------|------------|------------|-----------------|
| 20             | Foot       | Centimeter | 609.6           |
| 1              | Mile       | Kilometer  | 1.609344        |
| 1              | Kilogram   | Pound      | 2.204624        |
| 100            | Celsius    | Fahrenheit | 212             |
| 0              | Celsius    | Kelvin     | 273.15          |

## Technologies Used

- **Go** — HTTP server, routing, form parsing, template rendering (standard library only, zero dependencies)
- **HTML5** — semantic markup with accessible form labels
- **CSS3** — custom styling with CSS variables, animations, and responsive design
- **Google Fonts** — DM Sans + DM Serif Display for typography
