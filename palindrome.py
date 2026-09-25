#!/usr/bin/env python3
"""Check whether text is a palindrome.

Usage:
    python palindrome.py racecar
    python palindrome.py never odd or even
    python palindrome.py --strict Racecar
    printf 'racecar\\nhello\\n' | python palindrome.py
    python palindrome.py --version

By default the comparison ignores case, spaces and punctuation: the text is
casefolded and every non-alphanumeric character is dropped before it is
compared with its reverse. With ``--strict`` the raw text is compared as-is,
so case, spaces and punctuation are all significant.

Library callers should use :func:`is_palindrome` for the forgiving check and
:func:`is_palindrome_strict` for the exact-match check.
``is_palindrome(value, strict=True)`` is kept for compatibility and simply
delegates to :func:`is_palindrome_strict`.

Empty input is *not* a palindrome. This covers the empty string, blank lines
read from standard input, empty standard input, and (outside ``--strict``
mode) text that is empty once spaces and punctuation are removed, such as
``"!?.,"``. The CLI reports such input as ``not a palindrome (empty input)``.

Options:
    --strict   compare the raw text exactly (case, spaces, punctuation matter).
    --version  print the script version (``__version__``) and exit with 0
               without checking any text or reading standard input.

Exit codes:
    0  every evaluated input is a palindrome, or ``--version``/``--help``
       was requested.
    1  at least one input is not a palindrome (including empty input).
"""

from __future__ import annotations

import argparse
import sys

__version__ = "0.1.0"


def _normalize(text: str) -> str:
    """Return ``text`` casefolded and stripped of non-alphanumeric characters.

    Args:
        text: The text to normalize.

    Returns:
        The casefolded text containing only alphanumeric characters.

    Examples:
        >>> _normalize("A man, a plan!")
        'amanaplan'
        >>> _normalize("  ,.!? ")
        ''
    """
    return "".join(ch for ch in text.casefold() if ch.isalnum())


def _candidate(value: object, *, strict: bool = False) -> str:
    """Return the string that ``is_palindrome`` compares with its reverse.

    Args:
        value: The value to check; non-strings are coerced with ``str()``.
        strict: When ``True`` the raw text is returned unchanged; otherwise
            it is normalized with :func:`_normalize`.

    Returns:
        The text to compare.
    """
    text = value if isinstance(value, str) else str(value)
    return text if strict else _normalize(text)


def is_palindrome_strict(value: object) -> bool:
    """Return ``True`` if ``value`` is exactly equal to its reverse.

    Unlike :func:`is_palindrome`, nothing is normalized: case, spaces and
    punctuation are all significant.

    Args:
        value: The value to check. Non-string values are coerced with
            ``str()``, so ``is_palindrome_strict(12321)`` is ``True``.

    Returns:
        ``True`` if the raw text is non-empty and equal to its reverse,
        ``False`` otherwise. Empty input returns ``False``.

    Examples:
        >>> is_palindrome_strict("racecar")
        True
        >>> is_palindrome_strict("Racecar")
        False
        >>> is_palindrome_strict("never odd or even")
        False
        >>> is_palindrome_strict("")
        False
        >>> is_palindrome_strict(" ")
        True
    """
    candidate = _candidate(value, strict=True)
    # Nothing to compare: empty input is deliberately not a palindrome.
    if not candidate:
        return False
    return candidate == candidate[::-1]


def is_palindrome(value: object, *, strict: bool = False) -> bool:
    """Return ``True`` if ``value`` reads the same forwards and backwards.

    Case, spaces and punctuation are ignored. For an exact comparison call
    :func:`is_palindrome_strict` directly.

    Args:
        value: The value to check. Non-string values are coerced with
            ``str()``, so ``is_palindrome(12321)`` is ``True``.
        strict: When ``False`` (the default) case, spaces and punctuation are
            ignored. When ``True`` the call delegates to
            :func:`is_palindrome_strict`; prefer calling that function
            directly so the intent is explicit.

    Returns:
        ``True`` if the (normalized, unless ``strict``) text is non-empty and
        equal to its reverse, ``False`` otherwise. Empty input, or input that
        is empty after normalization, returns ``False``.

    Examples:
        >>> is_palindrome("A man, a plan, a canal: Panama")
        True
        >>> is_palindrome("hello")
        False
        >>> is_palindrome("")
        False
        >>> is_palindrome("!?.,")
        False
        >>> is_palindrome(" ", strict=True)
        True
    """
    if strict:
        return is_palindrome_strict(value)
    candidate = _candidate(value)
    # Nothing left to compare: empty input is deliberately not a palindrome.
    if not candidate:
        return False
    return candidate == candidate[::-1]


def _build_parser() -> argparse.ArgumentParser:
    """Return the argument parser for the command line interface."""
    parser = argparse.ArgumentParser(
        description="Check whether text is a palindrome.",
    )
    parser.add_argument(
        "text",
        nargs="*",
        help="text to check; multiple words are joined with a single space. "
        "If omitted, lines are read from standard input.",
    )
    parser.add_argument(
        "--strict",
        action="store_true",
        help="compare the text exactly, so case, spaces and punctuation matter",
    )
    parser.add_argument(
        "--version",
        action="version",
        version=f"%(prog)s {__version__}",
        help="print the script version and exit",
    )
    return parser


def main(argv: list[str] | None = None) -> int:
    """Run the command line interface and return the process exit code.

    Positional arguments are joined with a single space and checked as one
    input. With no positional arguments each line of standard input is
    checked separately; empty standard input is checked as a single empty
    input. Empty inputs are reported as ``not a palindrome (empty input)``.

    Args:
        argv: Command line arguments, excluding the program name. Defaults to
            ``sys.argv[1:]`` when ``None``.

    Returns:
        ``0`` if every input is a palindrome, ``1`` if at least one is not.

    Examples:
        >>> main(["racecar"])
        "racecar" -> palindrome
        0
        >>> main([""])
        "" -> not a palindrome (empty input)
        1
    """
    args = _build_parser().parse_args(argv)

    if args.text:
        inputs = [" ".join(args.text)]
    else:
        inputs = [line.rstrip("\n") for line in sys.stdin] or [""]

    check = is_palindrome_strict if args.strict else is_palindrome

    all_palindromes = True
    for text in inputs:
        if check(text):
            print(f'"{text}" -> palindrome')
        elif not _candidate(text, strict=args.strict):
            print(f'"{text}" -> not a palindrome (empty input)')
            all_palindromes = False
        else:
            print(f'"{text}" -> not a palindrome')
            all_palindromes = False

    return 0 if all_palindromes else 1


if __name__ == "__main__":
    raise SystemExit(main())
