# testupsun

## Palindrome checker

`palindrome.py` checks whether a string is a palindrome, ignoring case,
punctuation and whitespace. On the command line numbers work too, since they
arrive as text; the library function accepts only strings (`is_palindrome(121)`
raises `TypeError`, use `is_palindrome(str(121))`).

Run it:

```bash
python palindrome.py "A man, a plan, a canal: Panama"
echo "level" | python palindrome.py
python palindrome.py            # prompts for input
```

The script exits with code `0` when the input is a palindrome, `1` when it
is not, and `2` when no input could be read (empty stdin or a blank line, Ctrl-D or Ctrl-C).

Use it as a library:

```python
from palindrome import is_palindrome

is_palindrome("racecar")  # True
```
