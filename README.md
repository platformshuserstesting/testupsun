# testupsun

## Palindrome checker

`palindrome.py` checks whether a piece of text is a palindrome. It uses only
the Python standard library and runs on Python 3.8+.

### Usage

```
python palindrome.py racecar                          # is a palindrome  (exit 0)
python palindrome.py hello                            # not a palindrome (exit 1)
python palindrome.py "A man, a plan, a canal: Panama" # exit 0
python palindrome.py 12321                            # exit 0
python palindrome.py --strict "Race car"              # keeps spaces   -> exit 1
python palindrome.py --case-sensitive Racecar         # keeps case     -> exit 1
python palindrome.py                                  # prompts for the text
python palindrome.py -h                               # help
```

By default the comparison ignores case and non-alphanumeric characters. Empty
input, or input with no alphanumeric characters (for example `"!!!"`), is
treated as a palindrome. Exit codes: `0` palindrome, `1` not a palindrome, `2`
usage or input error.

The logic is also importable:

```python
from palindrome import is_palindrome

is_palindrome("No 'x' in Nixon")                    # True
is_palindrome("Racecar", ignore_case=False)         # False
is_palindrome("Race car", ignore_non_alphanumeric=False)  # False
```

### Tests

```
python -m unittest discover -v
```
