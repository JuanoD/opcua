package main

import (
	"bytes"
	"encoding/csv"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	f, err := os.Open(filepath.Join("..", "..", "schema", "NodeIds.csv"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rows := make([][]string, 0)
	reader := csv.NewReader(f)

	binary := "_Encoding_DefaultBinary"
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if strings.Contains(record[0], binary) {
			// rewrite name to match our type definitions
			record[0] = strings.ReplaceAll(record[0], binary, "")
			rows = append(rows, record)
		}
		// skip all json and xml to make final bundle smaller
		// if strings.Contains(record[0], "Encoding_DefaultJson") || strings.Contains(record[0], "Encoding_DefaultXml") {
		// 	continue
		// }
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, rows); err != nil {
		panic(err)
	}

	out := filepath.Join("..", "..", "src", "id", "id.ts")
	if err := ioutil.WriteFile(out, b.Bytes(), 0644); err != nil {
		panic(err)
	}

	log.Printf("Wrote %s", out)
}

var tmpl = template.Must(template.New("").Parse(
	`// Code generated by cmd/id. DO NOT EDIT!

{{range .}}export const Id{{index . 0}} = {{index . 1}}
{{end}}

export const mapIdToName = new Map([
  {{ range . }}[{{ index . 1 }}, '{{ index . 0 }}'],
  {{ end }}
])

export const mapNameToId = new Map([
  {{ range . }}['{{ index . 0 }}', {{ index . 1 }}],
  {{ end }}
])`))
