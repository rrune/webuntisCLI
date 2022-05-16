# WebUntisCLI
This is a really simple tool to get your timetable from webuntis.

# Usage

Get your timetable for tomorrow
```bash
webuntis
```

Get your timetable for today + n
```bash
webuntis n
```
for example, today would be
```bash
webuntis 0
```
The day after tomorrow would be
```bash
webuntis 2
```

# Installation/Building
The config gets embedded into the binary, so you have to build it yourself. Copy `config.json.sample` to `config.json` and edit it. Then run 
```bash
go build
```
You need to have at least Go 1.18 installed. Then just put the binary somewhere in your PATH
