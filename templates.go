package main

import (
	"html/template"
)

const inputTemplateStr string = `
<!DOCTYPE HTML>
<html>
  <head>
    <title>Model Input</title>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap.min.css">
  </head>
  <body>
    <div class="container-fluid">
      <h2>Model Input</h2>
      <form action="/solve/" method="POST">
        <div><label>Input your Zimpl model here:</label></div>
        <div style="font-family:monospace;">
          <textarea name="model" rows="24" class="form-control">{{.Model}}</textarea>
        </div>
        <button type="submit" class="btn btn-default">Solve</button>
      </form>
    </div>
  </body>
</html>
`

var inputTemplate = template.Must(template.New("input").Parse(inputTemplateStr))

const resultTemplateStr string = `
<!DOCTYPE HTML>
<html>
  <head>
    <title>Solver Output</title>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap.min.css">
  </head>
  <body>
    <div class="container-fluid">
      <h2>Model</h2>
      <div><pre>{{.Model}}</pre></div>
      <form action="/input/" method="POST">
        <input type="hidden" name="prefilled" value="{{.Model}}">
        <button type="submit" class="btn btn-default">Edit</button>
      </form>

      {{if .Output}}
      <h2>Solution Values</h2>
      <div><pre>{{.Solution}}</pre></div>
      
      <h2>Solver Output</h2>
      <div><pre>{{.Output}}</pre></div>
      {{else}}
      <p>Solving not complete yet. Please retry later.</p>
      {{end}}
    </div>
  </body>
</html>
`

var resultTemplate = template.Must(template.New("input").Parse(resultTemplateStr))
