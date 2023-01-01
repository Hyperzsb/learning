package main

import (
	"html/template"
	"os"
)

func HTMLDemo() error {
	const (
		filename string = ".build/static/demo.html"
	)

	section := Section{Title: "Sec 1",
		Story: []string{"This is the first section.",
			"It is a long section."},
		Options: []Option{{Arc: "Opt 1", Text: "This is the first option"},
			{Arc: "Opt 2", Text: "This is the second option"}},
	}

	htmlTemplate, err := template.New("page").Parse(
		"<html>\n" +
			"<body>\n" +
			"\t<h1>{{.Title}}</h1>\n" +
			"{{range .Story}}" +
			"\t<p>{{.}}</p>\n" +
			"{{end}}" +
			"{{range .Options}}" +
			"\t<h2>Options:</h2>\n" +
			"\t<p><b>{{.Arc}}: </b>{{.Text}}</p>\n" +
			"{{end}}" +
			"</body>\n" +
			"</html>\n")
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	err = htmlTemplate.Execute(file, section)

	return nil
}
