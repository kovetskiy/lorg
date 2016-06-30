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
```
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
```
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
```
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
```
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
```
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
```
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
