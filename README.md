# testupsun

Small standard-library Python programs. No third-party dependencies are required.

## Palindrome checker

`palindrome.py` checks whether the given text reads the same forwards and backwards.
By default the comparison ignores case, spaces and punctuation.

```
python palindrome.py racecar                            # "racecar" is a palindrome.
python palindrome.py hello                              # "hello" is not a palindrome.
python palindrome.py "A man, a plan, a canal: Panama"   # palindrome
python palindrome.py A man a plan a canal Panama        # arguments are joined with spaces
python palindrome.py 12321                              # palindrome
python palindrome.py --case-sensitive Racecar           # not a palindrome
python palindrome.py --strict "never odd or even"       # not a palindrome (spaces kept)
```

Running it with no arguments from an interactive terminal prompts for the text to check.

Exit-code convention (note that a non-zero code is used as a verdict, not only for failure):

| Code | Meaning |
|---|---|
| 0 | The input is a palindrome |
| 1 | The input is not a palindrome |
| 2 | Empty input, or no comparable characters |

The comparison is also importable as a plain function:

```python
from palindrome import is_palindrome

is_palindrome("Was it a car or a cat I saw?")           # True
is_palindrome("Racecar", ignore_case=False)             # False
is_palindrome("never odd or even", alphanumeric_only=False)  # False
```

## Hello world

```
python hello_world.py   # Hello, World!
```

## Tests

```
python -m unittest discover -v
```
