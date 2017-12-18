# Upgrade Scripts for 0.11


## Scope

This tool accodomates one change introduced in the 0.11 release: variable name for functions used to be prefixed with `&`, but are now suffixed with `~`:

*   Variable references are changed: `$a:&x` ⟹ `$a:x~`;

*   Assignments are changed: `'&x' = { }` ⟹ `x~ = { }` (quoting is no longer
    required, since `~` is a valid bareword if it does not appear at the
    beginning of a word;

*   Use of `~` in compound nodes get a preceding `''` so that it will not be
    parsed as part of a previous variable: `$x~foo` ⟹ `$x''~foo`.

See `before.elv` and `after.elv` for an example.

The rune `&` is now forbidden in variable names. If a variable contains `&`
after rewriting, a warning is printed.

This tool does not address other compatibility breaks.

## Invocation

This tool can be invoked in one of two ways:

*   Without arguments, it reads stdin and writes stdout.

*   With filename arguments, it rewrites each given file.

It does not accept any flags.
