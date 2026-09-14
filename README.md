# testupsun

## Palindrome checker

`palindrome.py` is a standard-library-only Python program (3.8+) that checks
whether a string — or a number written as text — is a palindrome.

By default case, spaces and punctuation are ignored; pass `--strict` for an
exact character-by-character comparison.

### Run

```sh
python palindrome.py racecar
python palindrome.py "A man, a plan, a canal: Panama"
python palindrome.py racecar hello      # one line per input, exit code 1
python palindrome.py --strict "Racecar"
python palindrome.py --quiet racecar    # prints only "true" / "false"
echo racecar | python palindrome.py     # reads one line from stdin
python palindrome.py --help
```

The program exits with `0` when every checked input is a palindrome and `1`
otherwise.

### Use as a library

```python
from palindrome import is_palindrome

is_palindrome("A man, a plan, a canal: Panama")  # True
is_palindrome("RaceCar", ignore_case=False, ignore_non_alnum=False)  # False
is_palindrome(str(121))  # True — non-string input raises TypeError
```

### Test

```sh
python -m unittest test_palindrome -v
```
