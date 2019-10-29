package types

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
)

type HTTPRequest struct {
	URI    string            `json:"uri, omitempty"`
	Method string            `json:"method, omitempty"`
	Header map[string]string `json:"header, omitempty"`
	Data   string            `json:"data, omitempty"`
}

// AttackConfig consists of information to launch an attack
type AttackConfig struct {
	Name           string `json:"name"`
	LHOST          string `json:"local_host, omitempty"`
	LPORT          string `json:"local_port, omitempty"`
	RHOST          string `json:"remote_host, omitempty"`
	RPORT          string `json:"remote_port, omitempty"`
	SRVHOST        string `json:"srv_host, omitempty"`
	TargetURI      string `json:"target_uri, omitempty"`
	Exploit        string `json:"exploit, omitempty"`
	Payload        string `json:"payload, omitempty"`
	IsReverseShell bool   `json:"reverse_shell, omitempty"`
	Database       string `json:"database, omitempty"`
	Password       string `json:"password, omitempty"`
	Username       string `json:"username, omitempty"`

	HTTPExploit []HTTPRequest `json:"http_exploit, omitempty"`
	HTTPPayload []HTTPRequest `json:"http_payload, omitempty"`
}

func (a AttackConfig) Validate() error {
	if len(a.Name) == 0 {
		return errors.New("missing attack name")
	}

	if len(a.RHOST) == 0 {
		return errors.New("missing victim host IP")
	}

	if len(a.Exploit) == 0 && len(a.HTTPExploit) == 0 {
		return errors.New("missing exploit info")
	}

	for _, e := range a.HTTPExploit {
		if !validHTTPMethod(e.Method) {
			return errors.New("Invalid HTTP request method: " + e.Method)
		}
	}

	for _, e := range a.HTTPExploit {
		if !validHTTPMethod(e.Method) {
			return errors.New("Invalid HTTP request method: " + e.Method)
		}
	}

	return nil
}

func validHTTPMethod(m string) bool {
	method := strings.ToUpper(m)
	if method != "" && method != "GET" && method != "POST" &&
		method != "PUT" && method != "DELETE" && method != "PATCH" {
		return false
	}
	return true
}

// Header returns the header string
func (r HTTPRequest) ConstructHeader() string {
	headers := []string{}

	for k, v := range r.Header {
		headers = append(headers, fmt.Sprintf("%s: %s", k, v))
	}
	return strings.Join(headers, "; ")
}

func (r HTTPRequest) ConstructMethod() string {
	method := strings.ToUpper(r.Method)
	if method == "" {
		method = "GET"
	}
	return "-X " + method
}

func (r HTTPRequest) ConstructData() string {
	if len(r.Data) == 0 {
		return ""
	}
	return fmt.Sprintf("-d '%s'", r.Data)
}

func (r HTTPRequest) ConstructURI(ac AttackConfig) string {
	var uri string
	if len(r.URI) == 0 {
		uri = "/"
	} else if r.URI[0] != '/' {
		t := template.New("http request")
		t, _ = t.Parse(r.URI)

		var tpl bytes.Buffer
		if err := t.Execute(&tpl, ac); err == nil {
			uri = tpl.String()
			return uri
		}
	} else {
		uri = r.URI
	}

	return ac.RHOST + ":" + ac.RPORT + uri
}
