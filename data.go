package chuzo

import (
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

type Mold struct {
	tmpl *template.Template
}

func (m *Mold) Cast(wr io.Writer, mt Material) error {
	data, err := mt.Melt()
	if err != nil {
		return err
	}

	if err := m.tmpl.Execute(wr, data); err != nil {
		return err
	}
	return nil
}

func BuildMold(filename string) (Mold, error) {
	tmpl, err := template.ParseFiles(filename)
	tmpl = tmpl.Option("missingkey=error")
	if err != nil {
		return Mold{}, err
	}
	return Mold{tmpl: tmpl}, nil
}

type Material interface {
	Melt() (interface{}, error)
}

type YAMLMaterial struct {
	Path string
}

func (ym YAMLMaterial) Melt() (interface{}, error) {
	m := make(map[string]interface{})
	r, err := os.Open(ym.Path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	d, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	yaml.Unmarshal(d, &m)
	return m, nil
}
