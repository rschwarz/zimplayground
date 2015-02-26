package main

import (
	"html/template"
)

const inputTemplateStr string = `
<!DOCTYPE HTML>
<html>
  <head>
    <title>Model Input</title>
    <link rel="stylesheet" href="http://yui.yahooapis.com/pure/0.6.0/pure-min.css">
  </head>
  <body>
    <h2>Model Input</h2>
    <form action="/solve/" method="POST" class="pure-form">
      <div><label>Input your Zimpl model here:</label></div>
      <div><textarea name="model" rows="24" cols="80"></textarea></div>
      <div><input type="submit" value="Solve" class="pure-button"></div>
    </form>
  </body>
</html>
`

var inputTemplate = template.Must(template.New("input").Parse(inputTemplateStr))

const resultTemplateStr string = `
<!DOCTYPE HTML>
<html>
  <head>
    <title>Solver Output</title>
    <link rel="stylesheet" href="http://yui.yahooapis.com/pure/0.6.0/pure-min.css">
  </head>
  <body>
    <h2>Model</h2>
    <div><pre>{{.Model}}</pre></div>

    {{if .Output}}
    <h2>Solution Values</h2>
    <div><pre>{{.Solution}}</pre></div>

    <h2>Solver Output</h2>
    <div><pre>{{.Output}}</pre></div>
    {{else}}
    <p>Solving not complete yet. Please retry later.</p>
    {{end}}
  </body>
</html>
`

var resultTemplate = template.Must(template.New("input").Parse(resultTemplateStr))
