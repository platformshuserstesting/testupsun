# testupsun

## Palindrome check

`palindrome.py` is a dependency-free Python program (standard library only)
that checks whether text reads the same forwards and backwards.

### Usage

```bash
python palindrome.py racecar
# "racecar" -> palindrome

python palindrome.py "A man, a plan, a canal: Panama"
# "A man, a plan, a canal: Panama" -> palindrome

python palindrome.py never odd or even
# "never odd or even" -> palindrome

python palindrome.py hello
# "hello" -> not a palindrome
```

By default the comparison ignores case, spaces and punctuation. Use `--strict`
to compare the text exactly:

```bash
python palindrome.py --strict Racecar
# "Racecar" -> not a palindrome
```

With no arguments the program reads standard input and checks one line at a
time:

```bash
printf 'racecar\nhello\n' | python palindrome.py
```

Print the script version and exit (exit code 0, no text is checked and
standard input is not read):

```bash
python palindrome.py --version
# palindrome.py 0.1.0
```

Empty input is **not** a palindrome. This covers an empty argument, blank
lines and empty standard input, and (outside `--strict` mode) text that is
empty once spaces and punctuation are removed, such as `"!?.,"`:

```bash
python palindrome.py ""
# "" -> not a palindrome (empty input)
```

In `--strict` mode nothing is removed, so a single space is a real character
and counts as a palindrome.

### Exit codes

- `0` — every evaluated input is a palindrome.
- `1` — at least one input is not a palindrome (including empty input).

### Library use

```python
from palindrome import is_palindrome, is_palindrome_strict

is_palindrome("Never odd or even")   # True
is_palindrome(12321)                 # True
is_palindrome("")                    # False

is_palindrome_strict("racecar")      # True
is_palindrome_strict("Racecar")      # False (case counts)
is_palindrome_strict("")             # False
```

`is_palindrome` ignores case, spaces and punctuation; `is_palindrome_strict`
compares the text exactly. `is_palindrome(value, strict=True)` still works and
delegates to `is_palindrome_strict`, but calling `is_palindrome_strict`
directly makes the intent explicit.

### Tests

```bash
python -m unittest discover -v
# or, if pytest is available:
pytest -q
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for more on running the tests and
contributing changes.
