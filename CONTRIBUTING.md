# Contributing

Thanks for helping improve the palindrome checker. This guide explains how to
run the tests and try the program locally.

## Requirements

- Python 3.8 or newer.
- Nothing to install: the program and its tests use only the standard
  library.

## Running the tests

From the repository root, run the whole suite:

```bash
python -m unittest discover -v
```

To run just the palindrome test module:

```bash
python -m unittest test_palindrome -v
```

If you have [pytest](https://pytest.org) installed, it can run the same tests:

```bash
pytest -q
```

All tests should pass before you open a pull request.

## Trying the CLI by hand

```bash
python palindrome.py racecar
# "racecar" -> palindrome            (exit code 0)

python palindrome.py --strict Racecar
# "Racecar" -> not a palindrome      (exit code 1)

echo hello | python palindrome.py
# "hello" -> not a palindrome        (exit code 1)
```

Exit codes:

- `0` — every evaluated input is a palindrome.
- `1` — at least one input is not a palindrome (including empty input).

Check the exit code with `echo $?` right after running a command.

## Conventions

- Keep the project standard-library only; do not add third-party
  dependencies.
- Add or update tests in `test_palindrome.py` with every change.
- Bump `__version__` in `palindrome.py` (shown by `--version`) whenever you
  change the program's behaviour.
- Give every public function a docstring describing its arguments, return
  value and examples.
