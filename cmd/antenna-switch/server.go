package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Server struct {
	currentPort int
	Ports       []string
}

func mustParse(name string, tmpl string) *template.Template {
	t, err := template.New(name).Parse(tmpl)
	if err != nil {
		log.Fatalf("error parsing template %s: %s", name, err)
	}
	return t
}

var indexTmpl = mustParse("index", `<html>
<head>
<title>Antenna Switch</title>
<style>
body {
	background-color: #03045e;
}
form {
	display: inline;
}
button {
	background-color: #0077b6;
	font-size: 1.25em;
	font-family: monospace;
	margin-right: 1em;
	margin-bottom: 1em;
	border-radius: 0.5em;
	padding: 1em;
	border: 0px;

	-webkit-touch-callout: none;
		-webkit-user-select: none;
		 -khtml-user-select: none;
		   -moz-user-select: none;
			-ms-user-select: none;
				user-select: none;
}
button:hover {
	background-color: #00b4d8 !important;
}
button.selected {
	background-color: #90e0ef;
}
button.calibrate {
	background-color: #005785;
	padding: 0.25em;
}
.centered {
	text-align: center;
}
</style>
</head>
<div class='centered'>
{{range .Ports}}
	<form action='/switch' method='GET'>
		<input type='hidden' name='port' value='{{- .Name -}}' />
		<button type='submit' {{- if .Selected }} class='selected' {{- end -}}>{{ .Name -}}</button>
	</form>
{{end}}
</div>
<div class='centered'>
	<form action='/calibrate' method='GET'>
	<button class='calibrate' type='submit'>Calibrate</button>
	</form>
</div>
</html>`)

func (s *Server) ServeIndex(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	type Port struct {
		Name     string
		Selected bool
	}
	ports := []Port{}
	for i, p := range s.Ports {
		ports = append(ports, Port{
			Name:     p,
			Selected: i == s.currentPort,
		})
	}
	data["Ports"] = ports
	data["CurrentPort"] = s.Ports[s.currentPort]
	indexTmpl.Execute(w, data)
}

func (s *Server) SwitchPorts(w http.ResponseWriter, r *http.Request) {
	portValues := r.URL.Query()["port"]
	if len(portValues) != 0 {
		port := portValues[0]
		for i, p := range s.Ports {
			if p == port {
				s.currentPort = i
				fmt.Println("SWITCH PORTS to", port)
				break
			}
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
