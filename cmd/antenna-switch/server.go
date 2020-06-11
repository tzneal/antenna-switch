package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/tzneal/antenna-switch/ticcmd"
)

type Server struct {
	currentPort int
	ports       []Port
	messages    []string
	tic         *ticcmd.Client
}

func NewServer(ports []Port) (*Server, error) {
	tic, err := ticcmd.NewClient("")
	if err != nil {
		return nil, err
	}

	st, err := tic.Status()
	if err != nil {
		return nil, err
	}

	// find the closest configured port
	currentPort := 0
	currentPortDistance := 999
	for i, p := range ports {
		delta := p.Position - st.CurrentPosition
		if delta < 0 {
			delta *= -1
		}
		if delta < currentPortDistance {
			currentPortDistance = delta
			currentPort = i
		}
	}
	s := &Server{
		currentPort: currentPort,
		ports:       ports,
		tic:         tic,
	}
	return s, nil
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
<meta http-equiv="refresh" content="5">
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
button:disabled {
	color: #000;
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
.message {
	text-align: left;
	font-family: monospace;
	font-size: 1.25em;
	color: #fff;
}
@media only screen and (max-device-width: 850px) {
	button {
		font-size: 3em;
	}
	.message {
		font-size: 2.5em;
	}
}
</style>
</head>
<div class='centered'>
{{range .ports}}
	<form action='/switch' method='GET'>
		<input type='hidden' name='port' value='{{- .Label -}}' />
		<button type='submit' {{- if .Selected }} class='selected' disabled {{- end -}}>{{ .Label -}}</button>
	</form>
{{end}}
</div>
<div class='centered'>
	<form action='/calibrate' method='GET'>
	<button class='calibrate' type='submit'>Calibrate</button>
	</form>
</div>
<div class='message'>
	<pre>
{{range .Messages}}
{{.}}
{{end}}
	</pre>
</div>
</html>`)

func ReverseSlice(s interface{}) {
	size := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
func (s *Server) ServeIndex(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	type PortOption struct {
		Label    string
		Position int
		Selected bool
	}
	ports := []PortOption{}
	for i, p := range s.ports {
		ports = append(ports, PortOption{
			Label:    p.Label,
			Position: p.Position,
			Selected: i == s.currentPort,
		})
	}
	data["ports"] = ports

	// print logs in reverse order
	var msgs = []string{}
	msgs = append(msgs, s.messages...)
	ReverseSlice(msgs)
	data["Messages"] = msgs

	data["CurrentPort"] = s.ports[s.currentPort]
	indexTmpl.Execute(w, data)
}

func (s *Server) SwitchPorts(w http.ResponseWriter, r *http.Request) {
	portValues := r.URL.Query()["port"]
	if len(portValues) != 0 {
		port := portValues[0]
		for i, p := range s.ports {
			if p.Label == port {
				s.currentPort = i
				if err := s.tic.Energize(); err != nil {
					s.AddMessage("error energizing: " + err.Error())
				}
				defer s.tic.Deenergize()
				s.AddMessage("switching to " + port)
				if err := s.tic.SetPosition(p.Position); err != nil {
					log.Printf("error setting positions: %s", err)
					s.AddMessage(err.Error())
				}
				if err := s.tic.WaitForPosition(p.Position, 5*time.Second); err != nil {
					log.Printf("error waiting for positions: %s", err)
					s.AddMessage(err.Error())
				}
				break
			}
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) Calibrate(w http.ResponseWriter, r *http.Request) {
	s.AddMessage("Calibrating")
	s.tic.Energize()
	defer s.tic.Deenergize()
	maxPos := -999
	currentPort := -1
	for i, p := range s.ports {
		if p.Position > maxPos {
			maxPos = p.Position
			currentPort = i
		}
	}
	s.tic.SetPosition(2 * maxPos)
	s.tic.WaitForPosition(2*maxPos, 5*time.Second)
	s.tic.SetKnownPosition(maxPos)
	s.currentPort = currentPort
	s.AddMessage("Calibration Complete")
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) AddMessage(msg string) {
	s.messages = append(s.messages, fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), msg))
	if len(s.messages) > 10 {
		s.messages = s.messages[1:]
	}
}
