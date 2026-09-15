# testupsun

## Palindrome checker

`palindrome.py` checks whether a string (or number) is a palindrome, ignoring
case, punctuation and whitespace.

Run it:

```bash
python palindrome.py "A man, a plan, a canal: Panama"
echo "level" | python palindrome.py
python palindrome.py            # prompts for input
```

The script exits with code `0` when the input is a palindrome and `1` when it
is not.

Use it as a library:

```python
from palindrome import is_palindrome

is_palindrome("racecar")  # True
```

Run the tests:

```bash
python -m unittest test_palindrome -v
```
