package httpui

const pageTemplate = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>mabata</title>
  <style>
    body { font-family: sans-serif; margin: 2rem auto; max-width: 900px; padding: 0 1rem; }
    table { border-collapse: collapse; width: 100%; margin-top: 1rem; }
    th, td { border: 1px solid #ccc; padding: 0.5rem; text-align: left; }
    .top { display: flex; justify-content: space-between; align-items: center; }
    .muted { color: #666; }
  </style>
</head>
<body>
  <div class="top">
    <div>
      <h1>mabata</h1>
      <p class="muted">OIDC + DuckDB starter</p>
    </div>
    {{if .User}}
      <div>
        <div>{{.User.Name}} {{if .User.Email}}&lt;{{.User.Email}}&gt;{{end}}</div>
        <div><a href="/logout">Logout</a></div>
      </div>
    {{else}}
      <div><a href="/login">Login with Entra ID</a></div>
    {{end}}
  </div>

  {{if .User}}
  <form method="get" action="/">
    <label for="query">Query</label>
    <select name="query" id="query">
      {{range .Queries}}
      <option value="{{.}}" {{if eq $.Selected .}}selected{{end}}>{{.}}</option>
      {{end}}
    </select>
    <button type="submit">Run</button>
  </form>

  {{if .Error}}
    <p>{{.Error}}</p>
  {{end}}

  {{if .Headers}}
  <table>
    <thead>
      <tr>
        {{range .Headers}}<th>{{.}}</th>{{end}}
      </tr>
    </thead>
    <tbody>
      {{range .Rows}}
      <tr>
        {{range .}}<td>{{.}}</td>{{end}}
      </tr>
      {{end}}
    </tbody>
  </table>
  {{end}}
  {{else}}
    <p>Sign in to run approved queries.</p>
  {{end}}
</body>
</html>`
