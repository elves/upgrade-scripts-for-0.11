# Upgrade Scripts for 0.11


## Scope

This tool accodomates one change introduced in the 0.11 release: variable name for functions used to be prefixed with `&`, but are now suffixed with `~`.

Before:

```
'&x' = { echo x }
# equivalent to "fn x { echo x }"
echo $&x
```

After:

```
x~ = { echo x }
echo $x~
```

Since `~` is now allowed as part of variable name, it also fixes code like

```
echo $x~2
```

to

```
echo $x''~2
```

Also, `&` is now forbidden in variable names. If a variable contains `&` after
rewriting, a warning is printed.

It does not address other compatibility breaks.

## Invocation

This tool can be invoked in one of two ways:

*   Without arguments, it reads stdin and writes stdout.

*   With filename arguments, it rewrites each given file.

It does not accept any flags.
