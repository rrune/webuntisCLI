# WebUntisCLI
This is a really simple tool to get your timetable from webuntis.

# Usage

Get your timetable for today
```bash
webuntis
```

Get your timetable for today + n
```bash
webuntis n
```
for example, tomorrow would be
```bash
webuntis 1
```

# Installation/Building
The config gets embedded into the binary, so you have to build it yourself. Copy `config.json.sample` to `config.json` and edit it. Then run 
```bash
go build
```
You need to have at least Go 1.18 installed. Then just put the binary somewhere in your PATH
