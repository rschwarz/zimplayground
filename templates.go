package main

import (
	"html/template"
)

const inputTemplateStr string = `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN"
  "http://www.w3.org/TR/html4/loose.dtd">
<html>
  <head>
    <title>Model Input</title>
  </head>
  <body>
    <h1>Model Input</h1>

    <form action="/solve/" method="POST">
      <div>
        Input your Zimpl model here:
      </div>
      <div>
        <textarea name="model" rows="24" cols="80"></textarea>
      </div>
      <div>
        <input type="submit" value="Solve">
      </div>
    </form>
  </body>
</html>
`

var inputTemplate = template.Must(template.New("input").Parse(inputTemplateStr))

const resultTemplateStr string = `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN"
  "http://www.w3.org/TR/html4/loose.dtd">
<html>
  <head>
    <title>Solver Output</title>
  </head>
  <body>
    <h1>Model</h1>
    <div>
      <pre>{{.Model}}</pre>
    </div>

    {{if .Output}}
    <h1>Results</h1>
    <h2>Solution Values</h2>
    <div>
      <pre>{{.Solution}}</pre>
    </div>

    <h2>Solver Output</h2>
    <div>
      <pre>{{.Output}}</pre>
    </div>
    {{else}}
    <div>
      Solving not complete yet. Please retry later.
    </div>
    {{end}}
  </body>
</html>
`

var resultTemplate = template.Must(template.New("input").Parse(resultTemplateStr))
