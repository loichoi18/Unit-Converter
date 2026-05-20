package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// ------------------- Conversion Logic -------------------

// Length: all stored as meters
var lengthToMeter = map[string]float64{
	"millimeter": 0.001,
	"centimeter": 0.01,
	"meter":      1,
	"kilometer":  1000,
	"inch":       0.0254,
	"foot":       0.3048,
	"yard":       0.9144,
	"mile":       1609.344,
}

// Weight: all stored as grams
var weightToGram = map[string]float64{
	"milligram": 0.001,
	"gram":      1,
	"kilogram":  1000,
	"ounce":     28.3495,
	"pound":     453.592,
}

func convertLength(value float64, from, to string) (float64, error) {
	fromFactor, ok1 := lengthToMeter[from]
	toFactor, ok2 := lengthToMeter[to]
	if !ok1 || !ok2 {
		return 0, fmt.Errorf("unknown length unit")
	}
	meters := value * fromFactor
	return meters / toFactor, nil
}

func convertWeight(value float64, from, to string) (float64, error) {
	fromFactor, ok1 := weightToGram[from]
	toFactor, ok2 := weightToGram[to]
	if !ok1 || !ok2 {
		return 0, fmt.Errorf("unknown weight unit")
	}
	grams := value * fromFactor
	return grams / toFactor, nil
}

func convertTemperature(value float64, from, to string) (float64, error) {
	// First convert to Celsius as intermediate
	var celsius float64
	switch from {
	case "celsius":
		celsius = value
	case "fahrenheit":
		celsius = (value - 32) * 5 / 9
	case "kelvin":
		celsius = value - 273.15
	default:
		return 0, fmt.Errorf("unknown temperature unit: %s", from)
	}

	// Then convert from Celsius to target
	switch to {
	case "celsius":
		return celsius, nil
	case "fahrenheit":
		return celsius*9/5 + 32, nil
	case "kelvin":
		return celsius + 273.15, nil
	default:
		return 0, fmt.Errorf("unknown temperature unit: %s", to)
	}
}

// Round to a sensible number of decimal places
func roundResult(val float64) string {
	if val == math.Trunc(val) {
		return strconv.FormatFloat(val, 'f', 0, 64)
	}
	// Up to 6 decimals, trim trailing zeros
	s := strconv.FormatFloat(val, 'f', 6, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

// ------------------- Template Data -------------------

type PageData struct {
	Category    string
	Units       []string
	Value       string
	FromUnit    string
	ToUnit      string
	Result      string
	ResultFrom  string
	ResultTo    string
	ShowResult  bool
	Error       string
}

// ------------------- Handlers -------------------

func lengthHandler(w http.ResponseWriter, r *http.Request) {
	units := []string{"millimeter", "centimeter", "meter", "kilometer", "inch", "foot", "yard", "mile"}
	data := PageData{Category: "Length", Units: units}

	if r.Method == http.MethodPost {
		r.ParseForm()
		valStr := r.FormValue("value")
		from := r.FormValue("from")
		to := r.FormValue("to")
		data.Value = valStr
		data.FromUnit = from
		data.ToUnit = to

		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			data.Error = "Please enter a valid number."
		} else {
			result, err := convertLength(val, from, to)
			if err != nil {
				data.Error = err.Error()
			} else {
				data.ShowResult = true
				data.Result = roundResult(result)
				data.ResultFrom = fmt.Sprintf("%s %s", valStr, from)
				data.ResultTo = fmt.Sprintf("%s %s", data.Result, to)
			}
		}
	}

	tmpl.Execute(w, data)
}

func weightHandler(w http.ResponseWriter, r *http.Request) {
	units := []string{"milligram", "gram", "kilogram", "ounce", "pound"}
	data := PageData{Category: "Weight", Units: units}

	if r.Method == http.MethodPost {
		r.ParseForm()
		valStr := r.FormValue("value")
		from := r.FormValue("from")
		to := r.FormValue("to")
		data.Value = valStr
		data.FromUnit = from
		data.ToUnit = to

		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			data.Error = "Please enter a valid number."
		} else {
			result, err := convertWeight(val, from, to)
			if err != nil {
				data.Error = err.Error()
			} else {
				data.ShowResult = true
				data.Result = roundResult(result)
				data.ResultFrom = fmt.Sprintf("%s %s", valStr, from)
				data.ResultTo = fmt.Sprintf("%s %s", data.Result, to)
			}
		}
	}

	tmpl.Execute(w, data)
}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	units := []string{"celsius", "fahrenheit", "kelvin"}
	data := PageData{Category: "Temperature", Units: units}

	if r.Method == http.MethodPost {
		r.ParseForm()
		valStr := r.FormValue("value")
		from := r.FormValue("from")
		to := r.FormValue("to")
		data.Value = valStr
		data.FromUnit = from
		data.ToUnit = to

		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			data.Error = "Please enter a valid number."
		} else {
			result, err := convertTemperature(val, from, to)
			if err != nil {
				data.Error = err.Error()
			} else {
				data.ShowResult = true
				data.Result = roundResult(result)
				data.ResultFrom = fmt.Sprintf("%s %s", valStr, from)
				data.ResultTo = fmt.Sprintf("%s %s", data.Result, to)
			}
		}
	}

	tmpl.Execute(w, data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/length", http.StatusFound)
}

// ------------------- Template -------------------

// capitalize first letter for display
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

var funcMap = template.FuncMap{
	"capitalize": capitalize,
	"lower":      strings.ToLower,
}

var tmpl = template.Must(template.New("page").Funcs(funcMap).Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Unit Converter — {{.Category}}</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link href="https://fonts.googleapis.com/css2?family=DM+Sans:ital,wght@0,400;0,500;0,600;0,700&family=DM+Serif+Display&display=swap" rel="stylesheet">
    <style>
        *, *::before, *::after {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        :root {
            --bg: #f5f0eb;
            --card: #ffffff;
            --text: #1a1a1a;
            --text-muted: #6b6560;
            --accent: #4f46e5;
            --accent-hover: #4338ca;
            --accent-light: #eef2ff;
            --border: #e0dbd5;
            --error: #dc2626;
            --success-bg: #f0fdf4;
            --success-border: #86efac;
            --shadow: 0 1px 3px rgba(0,0,0,0.06), 0 4px 12px rgba(0,0,0,0.04);
            --shadow-lg: 0 4px 20px rgba(0,0,0,0.08);
            --radius: 12px;
        }

        body {
            font-family: 'DM Sans', sans-serif;
            background: var(--bg);
            color: var(--text);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: flex-start;
            padding: 40px 20px;
        }

        .container {
            width: 100%;
            max-width: 480px;
        }

        h1 {
            font-family: 'DM Serif Display', serif;
            font-size: 2rem;
            font-weight: 400;
            margin-bottom: 28px;
            letter-spacing: -0.02em;
        }

        /* ---- Navigation Tabs ---- */
        .tabs {
            display: flex;
            gap: 4px;
            margin-bottom: 32px;
            background: var(--border);
            padding: 4px;
            border-radius: 10px;
        }

        .tabs a {
            flex: 1;
            text-align: center;
            text-decoration: none;
            color: var(--text-muted);
            font-weight: 500;
            font-size: 0.9rem;
            padding: 10px 0;
            border-radius: 8px;
            transition: all 0.2s ease;
        }

        .tabs a:hover {
            color: var(--text);
            background: rgba(255,255,255,0.5);
        }

        .tabs a.active {
            background: var(--card);
            color: var(--accent);
            box-shadow: 0 1px 4px rgba(0,0,0,0.06);
            font-weight: 600;
        }

        /* ---- Card ---- */
        .card {
            background: var(--card);
            border-radius: var(--radius);
            padding: 32px;
            box-shadow: var(--shadow);
        }

        /* ---- Form ---- */
        label {
            display: block;
            font-size: 0.82rem;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.06em;
            color: var(--text-muted);
            margin-bottom: 8px;
        }

        input[type="text"],
        select {
            width: 100%;
            padding: 14px 16px;
            font-family: 'DM Sans', sans-serif;
            font-size: 1rem;
            border: 2px solid var(--border);
            border-radius: 10px;
            background: var(--card);
            color: var(--text);
            transition: border-color 0.2s ease, box-shadow 0.2s ease;
            appearance: none;
            -webkit-appearance: none;
        }

        input[type="text"]:focus,
        select:focus {
            outline: none;
            border-color: var(--accent);
            box-shadow: 0 0 0 3px var(--accent-light);
        }

        input[type="text"]::placeholder {
            color: #b5b0aa;
        }

        /* custom select arrow */
        .select-wrapper {
            position: relative;
        }

        .select-wrapper::after {
            content: "";
            position: absolute;
            right: 16px;
            top: 50%;
            transform: translateY(-50%);
            width: 0;
            height: 0;
            border-left: 5px solid transparent;
            border-right: 5px solid transparent;
            border-top: 6px solid var(--text-muted);
            pointer-events: none;
        }

        .field {
            margin-bottom: 20px;
        }

        /* swap icon between selects */
        .swap-icon {
            display: flex;
            justify-content: center;
            margin: 4px 0;
        }

        .swap-icon svg {
            color: var(--text-muted);
            opacity: 0.5;
        }

        /* ---- Button ---- */
        button {
            width: 100%;
            padding: 14px;
            font-family: 'DM Sans', sans-serif;
            font-size: 1rem;
            font-weight: 600;
            color: #fff;
            background: var(--accent);
            border: none;
            border-radius: 10px;
            cursor: pointer;
            transition: background 0.2s ease, transform 0.1s ease;
            margin-top: 8px;
        }

        button:hover {
            background: var(--accent-hover);
        }

        button:active {
            transform: scale(0.98);
        }

        /* ---- Result ---- */
        .result {
            margin-top: 28px;
            padding: 24px;
            background: var(--accent-light);
            border: 2px solid #c7d2fe;
            border-radius: var(--radius);
            text-align: center;
            animation: fadeIn 0.3s ease;
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(6px); }
            to { opacity: 1; transform: translateY(0); }
        }

        .result-label {
            font-size: 0.82rem;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.06em;
            color: var(--accent);
            margin-bottom: 8px;
        }

        .result-value {
            font-family: 'DM Serif Display', serif;
            font-size: 1.6rem;
            color: var(--text);
            word-break: break-all;
        }

        .result-value .equals {
            color: var(--text-muted);
            margin: 0 4px;
        }

        /* ---- Reset link ---- */
        .reset-link {
            display: block;
            text-align: center;
            margin-top: 16px;
            font-size: 0.88rem;
            color: var(--text-muted);
            text-decoration: none;
            font-weight: 500;
            transition: color 0.2s;
        }

        .reset-link:hover {
            color: var(--accent);
        }

        /* ---- Error ---- */
        .error {
            margin-top: 20px;
            padding: 14px 18px;
            background: #fef2f2;
            border: 2px solid #fecaca;
            border-radius: var(--radius);
            color: var(--error);
            font-size: 0.92rem;
            font-weight: 500;
            text-align: center;
        }

        /* ---- Footer ---- */
        .footer {
            text-align: center;
            margin-top: 28px;
            font-size: 0.78rem;
            color: var(--text-muted);
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Unit Converter</h1>

        <nav class="tabs">
            <a href="/length" class="{{if eq .Category "Length"}}active{{end}}">Length</a>
            <a href="/weight" class="{{if eq .Category "Weight"}}active{{end}}">Weight</a>
            <a href="/temperature" class="{{if eq .Category "Temperature"}}active{{end}}">Temperature</a>
        </nav>

        <div class="card">
            <form method="POST" action="">
                <div class="field">
                    <label for="value">Enter the {{lower .Category}} to convert</label>
                    <input type="text" id="value" name="value" placeholder="e.g. 20" value="{{.Value}}" autocomplete="off" />
                </div>

                <div class="field">
                    <label for="from">Convert from</label>
                    <div class="select-wrapper">
                        <select id="from" name="from">
                            {{range .Units}}
                            <option value="{{.}}" {{if eq . $.FromUnit}}selected{{end}}>{{capitalize .}}</option>
                            {{end}}
                        </select>
                    </div>
                </div>

                <div class="swap-icon">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="7 3 7 21"/>
                        <polyline points="3 17 7 21 11 17"/>
                        <polyline points="17 21 17 3"/>
                        <polyline points="13 7 17 3 21 7"/>
                    </svg>
                </div>

                <div class="field">
                    <label for="to">Convert to</label>
                    <div class="select-wrapper">
                        <select id="to" name="to">
                            {{range .Units}}
                            <option value="{{.}}" {{if eq . $.ToUnit}}selected{{end}}>{{capitalize .}}</option>
                            {{end}}
                        </select>
                    </div>
                </div>

                <button type="submit">Convert</button>
            </form>

            {{if .ShowResult}}
            <div class="result">
                <div class="result-label">Result</div>
                <div class="result-value">
                    {{.ResultFrom}} <span class="equals">=</span> {{.ResultTo}}
                </div>
            </div>
            <a href="" class="reset-link">Reset</a>
            {{end}}

            {{if .Error}}
            <div class="error">{{.Error}}</div>
            {{end}}
        </div>

        <div class="footer">Built with Go &amp; HTML</div>
    </div>
</body>
</html>
`))

// ------------------- Main -------------------

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/length", lengthHandler)
	http.HandleFunc("/weight", weightHandler)
	http.HandleFunc("/temperature", temperatureHandler)

	port := ":8080"
	fmt.Printf("Unit Converter running at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
