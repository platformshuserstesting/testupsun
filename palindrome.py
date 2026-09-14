#!/usr/bin/env python3
"""Check whether a piece of text is a palindrome.

Usage examples::

    $ python palindrome.py racecar
    "racecar" is a palindrome.

    $ python palindrome.py "A man, a plan, a canal: Panama"
    "A man, a plan, a canal: Panama" is a palindrome.

    $ python palindrome.py --strict "A man, a plan, a canal: Panama"
    "A man, a plan, a canal: Panama" is not a palindrome.

    $ echo racecar | python palindrome.py

Exit codes: ``0`` when the input is a palindrome, ``1`` when it is not and
``2`` for usage errors.
"""

from __future__ import annotations

import argparse
import sys
from typing import List, Optional

__all__ = ["normalize", "is_palindrome", "main"]


def normalize(text: str, ignore_case: bool = True, alphanumeric_only: bool = True) -> str:
    """Return ``text`` prepared for comparison.

    When ``alphanumeric_only`` is true every character that is not
    alphanumeric (punctuation, whitespace, ...) is dropped. When
    ``ignore_case`` is true the result is case folded, which also handles
    non-ASCII casing such as ``"Été"``.
    """
    if alphanumeric_only:
        text = "".join(character for character in text if character.isalnum())
    if ignore_case:
        text = text.casefold()
    return text


def is_palindrome(text, ignore_case: bool = True, alphanumeric_only: bool = True) -> bool:
    """Return ``True`` when ``text`` reads the same forwards and backwards.

    Non-string values (for example the integer ``121``) are coerced with
    ``str()``. ``None`` is rejected with a ``TypeError``. The empty string
    -- and any input that normalizes to the empty string, such as ``"!!!"``
    -- is considered a palindrome by definition.

    The comparison reverses the normalized text with slicing, which is O(n)
    and more idiomatic than an explicit two-pointer loop.
    """
    if text is None:
        raise TypeError("is_palindrome() requires a string, got None")
    if not isinstance(text, str):
        text = str(text)

    normalized = normalize(text, ignore_case=ignore_case, alphanumeric_only=alphanumeric_only)
    return normalized == normalized[::-1]


def _build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="palindrome.py",
        description="Check whether the given text is a palindrome.",
    )
    parser.add_argument(
        "text",
        nargs="?",
        help="text to check; read from stdin or prompted for when omitted",
    )
    parser.add_argument(
        "--case-sensitive",
        action="store_true",
        help="treat upper and lower case letters as different characters",
    )
    parser.add_argument(
        "--strict",
        action="store_true",
        help="compare the text as-is, keeping punctuation and whitespace",
    )
    return parser


def main(argv: Optional[List[str]] = None) -> int:
    """Run the command line interface and return the process exit code."""
    args = _build_parser().parse_args(argv)

    text = args.text
    if text is None:
        if sys.stdin.isatty():
            text = input("Enter a string: ")
        else:
            text = sys.stdin.read().rstrip("\n")

    result = is_palindrome(
        text,
        ignore_case=not args.case_sensitive,
        alphanumeric_only=not args.strict,
    )

    if result:
        print('"{0}" is a palindrome.'.format(text))
        return 0

    print('"{0}" is not a palindrome.'.format(text))
    return 1


if __name__ == "__main__":
    raise SystemExit(main())
