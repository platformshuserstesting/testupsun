# testupsun

## Palindrome checker

`palindrome.py` is a dependency-free Python 3 program that checks whether a
string (or a number) is a palindrome. It can be used as a command line tool or
imported as a module.

### Usage

```bash
python palindrome.py racecar
# "racecar" is a palindrome.

python palindrome.py hello
# "hello" is not a palindrome.

python palindrome.py "A man, a plan, a canal: Panama"
# "A man, a plan, a canal: Panama" is a palindrome.

echo racecar | python palindrome.py
```

When no argument is given and the program runs in a terminal it prompts with
`Enter a string: `; when stdin is piped it reads the text from there.

### Flags

| Flag | Effect |
| --- | --- |
| `--case-sensitive` | Treat upper and lower case letters as different characters. |
| `--strict` | Compare the text as-is, keeping punctuation and whitespace. |

### Exit codes

| Code | Meaning |
| --- | --- |
| `0` | The input is a palindrome. |
| `1` | The input is not a palindrome. |
| `2` | Usage error (unknown flag, too many arguments). |

### Behaviour notes

By default the comparison ignores case (using `str.casefold()`, so accented
text such as `Été` works) and ignores every non-alphanumeric character. The
empty string — and any input that normalizes to it, such as `"!!!"` — is
considered a palindrome by definition.

### Library use

```python
from palindrome import is_palindrome

is_palindrome("A man, a plan, a canal: Panama")            # True
is_palindrome("Racecar", ignore_case=False)                 # False
is_palindrome("A man, a plan", alphanumeric_only=False)     # False
is_palindrome(121)                                          # True
is_palindrome(None)                                         # raises TypeError
```

### Running the tests

```bash
python -m pytest tests/ -q
```
