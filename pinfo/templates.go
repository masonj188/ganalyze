package pinfo

var bindata string = `
<!DOCTYPE html>
<html>
<head>
	<title>{{ .Name }}</title>
</head>
<body>
	<h1>{{ .Name }}</h1>
	<!--
	{{ if .ModelRes }}
	<h2>Determined Malicious: Yes</h2>
	{{ else }}
	<h2>Determined Malicious: No</h2>
	{{ end }}
	-->
	<h2>Basic Properties</h2>
		<table>
			<tr>
				<td><b>MD5:</b></td>
				<td>{{ .MD5 }}</td>
			</tr>
			<tr>
				<td><b>SHA1:</b></td>
				<td>{{ .SHA1 }}</td>
			</tr>
			<tr>
				<td><b>SHA256:</b></td>
				<td>{{ .SHA256 }}</td>
			</tr>
			<tr>
				<td><b>File Type:</b></td>
				<td>{{ .FileType }}</td>
			</tr>
			<tr>
				<td><b>Magic:</b></td>
				<td>{{ .Magic }}</td>
			</tr>
			<tr>
				<td><b>File Size:</b></td>
				<td>{{ .FSize }} bytes</td>
			</tr>
		</table>

	<h2>Sections</h2>
		<table>
			<tr>
				<th>Name</th>
				<th>Virtual Address</th>
				<th>Virtual Size</th>
				<th>Raw Size</th>
			</tr>
			{{ range .Sections }}
			<tr>
				<td>{{ .SectionHeader.Name }}</td>
				<td>{{ .SectionHeader.VirtualAddress }}</td>
				<td>{{ .SectionHeader.VirtualSize }}</td>
				<td>{{ .SectionHeader.Size }}</td>
			</tr>
			{{ end }}
		</table>
	<h2>Symbols</h2>
		<ul>
			{{ range .Symbols }}
			<li>{{.}}</li>
			{{ end }}
		</ul>
			
		
</body>
</html>
`

var Mainpage string = `
<!DOCTYPE html>
<html>
	<head>
		<title>Ganalyze Report</title>
	</head>
	<body>
		<h1>Files Examined:</h1>
		<ul>
		{{ range .LinkNames }}
			<li><a href={{ .Link }}>{{ .Name }}</a></li>
		{{ end }}
		</ul>
	</body>
</html>
`
