#!/usr/bin/env python3
"""Palindrome checker.

Provides a reusable :func:`is_palindrome` function and a small command line
interface.

Text is normalized before comparison: it is lowercased and every character
that is not alphanumeric (letters, digits, including non-ASCII ones) is
removed. As a consequence, an empty string -- or a string made only of
punctuation such as ``"!!!"`` -- normalizes to ``""`` and is considered a
palindrome.

Usage::

    $ python palindrome.py "A man, a plan, a canal: Panama"
    "A man, a plan, a canal: Panama" is a palindrome.

    $ echo "level" | python palindrome.py
    "level" is a palindrome.

    >>> from palindrome import is_palindrome
    >>> is_palindrome("racecar")
    True
"""

from __future__ import annotations

import argparse
import sys
from typing import Optional, Sequence


def normalize(text: str) -> str:
    """Return ``text`` lowercased with all non-alphanumeric characters removed."""
    if not isinstance(text, str):
        raise TypeError(f"expected a string, got {type(text).__name__}")
    return "".join(ch for ch in text.lower() if ch.isalnum())


def is_palindrome(text: str) -> bool:
    """Return ``True`` if ``text`` reads the same forwards and backwards.

    Comparison ignores case, punctuation and whitespace. Raises ``TypeError``
    if ``text`` is not a string.
    """
    normalized = normalize(text)
    return normalized == normalized[::-1]


def main(argv: Optional[Sequence[str]] = None) -> int:
    """Run the command line interface.

    Returns ``0`` when the input is a palindrome and ``1`` when it is not.
    """
    parser = argparse.ArgumentParser(
        description="Check whether a string (or number) is a palindrome."
    )
    parser.add_argument(
        "text",
        nargs="?",
        help="text to check; if omitted, a line is read from standard input",
    )
    args = parser.parse_args(argv)

    text = args.text
    if text is None:
        if sys.stdin.isatty():
            text = input("Enter text to check: ")
        else:
            text = sys.stdin.readline().rstrip("\n")

    if is_palindrome(text):
        print(f'"{text}" is a palindrome.')
        return 0

    print(f'"{text}" is not a palindrome.')
    return 1


if __name__ == "__main__":
    raise SystemExit(main())
