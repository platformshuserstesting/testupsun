#!/usr/bin/env python3
"""Check whether a string (or a number written as text) is a palindrome.

The module exposes two reusable helpers, :func:`normalize` and
:func:`is_palindrome`, plus a small command line interface.

By default comparison ignores case and any non-alphanumeric characters, so
``"A man, a plan, a canal: Panama"`` is considered a palindrome.  Use the
keyword flags (or ``--strict`` on the command line) for an exact,
character-by-character comparison.

Usage examples::

    $ python palindrome.py racecar
    'racecar' is a palindrome
    $ python palindrome.py racecar hello
    'racecar' is a palindrome
    'hello' is not a palindrome
    $ python palindrome.py --strict "Racecar"
    'Racecar' is not a palindrome
    $ python palindrome.py --quiet racecar
    true
    $ echo racecar | python palindrome.py
    Enter text to check: 'racecar' is a palindrome

The process exits with status ``0`` when every checked input is a palindrome
and ``1`` otherwise, which makes the program easy to use from shell scripts.
"""

from __future__ import annotations

import argparse
from typing import List, Optional, Sequence

__all__ = ["normalize", "is_palindrome", "main"]


def normalize(
    text: str,
    *,
    ignore_case: bool = True,
    ignore_non_alnum: bool = True,
) -> str:
    """Return ``text`` prepared for palindrome comparison.

    Args:
        text: The text to normalize.
        ignore_case: Fold case with :meth:`str.casefold` when true.
        ignore_non_alnum: Drop every character that is not alphanumeric.

    Raises:
        TypeError: If ``text`` is not a :class:`str`.
    """
    if not isinstance(text, str):
        raise TypeError(f"expected str, got {type(text).__name__}")

    result = text
    if ignore_non_alnum:
        result = "".join(ch for ch in result if ch.isalnum())
    if ignore_case:
        result = result.casefold()
    return result


def is_palindrome(
    text: str,
    *,
    ignore_case: bool = True,
    ignore_non_alnum: bool = True,
) -> bool:
    """Return ``True`` when ``text`` reads the same forwards and backwards.

    The empty string, a single character and — with the default flags — input
    made up only of punctuation (which normalizes to the empty string) are all
    treated as palindromes.

    Non-string input raises :class:`TypeError`; convert numbers explicitly with
    ``is_palindrome(str(121))``.
    """
    normalized = normalize(
        text,
        ignore_case=ignore_case,
        ignore_non_alnum=ignore_non_alnum,
    )
    return normalized == normalized[::-1]


def _build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="palindrome.py",
        description="Check whether the given text is a palindrome.",
    )
    parser.add_argument(
        "text",
        nargs="*",
        help="text to check; when omitted a single line is read from stdin",
    )
    parser.add_argument(
        "--strict",
        action="store_true",
        help="compare exactly, keeping case, spaces and punctuation",
    )
    parser.add_argument(
        "--quiet",
        action="store_true",
        help="print only 'true' or 'false' for each input",
    )
    return parser


def main(argv: Optional[Sequence[str]] = None) -> int:
    """Run the command line interface and return the process exit code."""
    args = _build_parser().parse_args(argv)

    values: List[str] = [str(value) for value in args.text]
    if not values:
        try:
            values = [input("Enter text to check: ")]
        except EOFError:
            print("No input provided.")
            return 1

    all_palindromes = True
    for value in values:
        result = is_palindrome(
            value,
            ignore_case=not args.strict,
            ignore_non_alnum=not args.strict,
        )
        all_palindromes = all_palindromes and result
        if args.quiet:
            print("true" if result else "false")
        else:
            suffix = "is a palindrome" if result else "is not a palindrome"
            print(f"{value!r} {suffix}")

    return 0 if all_palindromes else 1


if __name__ == "__main__":
    raise SystemExit(main())
