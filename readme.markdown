# LORG

## Formatting and placeholders

### Synopsis

When you set new format for a logger, lorg will examine given string and look
for placeholders' signatures which must be used with the following syntax:

```
${placeholder}
${placeholder[:option]}
```

- `placeholder` - unique placeholder name like `level` or `time`.
- `option` - positional argument which will be passed to placeholder function,
    this argument usually should be used for controlling placeholder behavior.
    For example, `time` placeholder use `option` value instead of time layout.

### Level

Level placeholder returns level of current logging entry. `level` can take 3
positional options:

`${level[:format[:alignment[:short]]]}`

- `format` - will be used to format a string representation of logging level;
- `alignment` - describes how to align resulting string, can be `left` or
    `right`.
- `short` - use short string representation of logging level, `WARNING -> WARN`.

Example:
```go
lorg.SetFormat(
    lorg.NewFormat(`[${level}] %s`),
)
lorg.Info("info")
lorg.Warning("warning")
```

Output:
```
[INFO] info
[WARNING] warning
```

Example:
```go
lorg.SetFormat(
    lorg.NewFormat(`${level:[%s]} %s`),
)
lorg.Info("info")
lorg.Warning("warning")
```

Output:
```
[INFO] info
[WARNING] warning
````

Example:
```go
lorg.SetFormat(
    lorg.NewFormat(`${level:[%s]:left} %s`),
)
lorg.Info("info")
lorg.Warning("warning")
```

Output:
```
[INFO]    info
[WARNING] warning
```

Example:
```go
lorg.SetFormat(
    lorg.NewFormat(`${level:[%s]:right} %s`),
)
lorg.Info("info")
lorg.Warning("warning")
```

Output:
```
   [INFO] info
[WARNING] warning
```

Example:
```go
lorg.SetFormat(
    lorg.NewFormat(`${level:[%s]:left:short} %s`),
)
lorg.Info("info")
lorg.Warning("warning")
```

Output:
```
[INFO]  info
[WARN]  warning
```

Example:
```go
lorg.SetFormat(
    lorg.NewFormat(`${level:[%s]:right:short} %s`),
)
lorg.Info("info")
lorg.Warning("warning")
lorg.Debug("debug")
```

Output:
```
 [INFO] info
 [WARN] warning
[DEBUG] debug
```

### Time

Time placeholder returns current time. `time` can take 1 positional option:

`${time[:layout]}`

- `layout` - will be to format current time;

Example:
```go
lorg.SetFormat(
    lorg.NewFormat(`${time:15:04:05} %s`),
)
lorg.Info("info")
lorg.Warning("warning")
```
Output:
```
09:21:44 info
09:21:44 warning
```

Example:
```go
lorg.SetFormat(
    lorg.NewFormat(`${time:15:04} %s`),
)
lorg.Info("info")
lorg.Warning("warning")
```

Output:
```
09:21 info
09:21 warning
```

### File

File placeholder returns a filename where has been called logging function.

```
${file[:mode]}
```

Placeholder accept two modes:

- `short` - default behavior, return base name of resulting filename.
- `long` - return full resulting filename

Example:

```go
lorg.SetFormat(
    lorg.NewFormat(`${file} %s`),
)
lorg.Info("info")
lorg.Warning("warning")
```

Output:
```
a.go info
a.go warning
```

Example:
```go
lorg.SetFormat(
    lorg.NewFormat(`${file:short} %s`),
)

lorg.Info("info")
lorg.Warning("warning")
```

Output:
```
a.go info
a.go warning
```

Example:
```go
lorg.SetFormat(
    lorg.NewFormat(`${file:long} %s`),
)
lorg.Info("info")
lorg.Warning("warning")
```

Output
```
/home/operator/a.go info
/home/operator/a.go warning
```

### Line

Line placeholder returns a number of line where has been called logging function.

```
${line}
```

Example:
```go
lorg.SetFormat(
    lorg.NewFormat(`${line} %s`),
)
lorg.Info("info")
lorg.Warning("warning")
```

Output:
```
11 info
12 warning
```
