#!/usr/bin/env python3
"""Check whether a piece of text is a palindrome.

A palindrome reads the same forwards and backwards. By default the check
ignores case and any non-alphanumeric characters, so ``"A man, a plan, a
canal: Panama"`` and ``12321`` both count as palindromes. Empty input -- and
input that contains no alphanumeric characters at all, such as ``"!!!"`` --
normalizes to the empty string and is treated as a palindrome (vacuously
true).

Usage::

    python palindrome.py racecar
    python palindrome.py "A man, a plan, a canal: Panama"
    python palindrome.py --strict "Race car"
    python palindrome.py --case-sensitive Racecar
    python palindrome.py            # prompts for the text

Exit codes: ``0`` palindrome, ``1`` not a palindrome, ``2`` usage/input error.
"""

from __future__ import annotations

import argparse
import sys
from typing import List, Optional, Sequence, Union

__all__ = ["is_palindrome", "main"]

EXIT_PALINDROME = 0
EXIT_NOT_PALINDROME = 1
EXIT_ERROR = 2


def _normalize(
    value: str, *, ignore_case: bool, ignore_non_alphanumeric: bool
) -> str:
    """Return ``value`` reduced to the characters that take part in the check.

    ``str.casefold()`` and ``str.isalnum()`` are Unicode-aware, so non-ASCII
    text is handled correctly.
    """
    characters: List[str] = list(value)
    if ignore_non_alphanumeric:
        characters = [char for char in characters if char.isalnum()]
    normalized = "".join(characters)
    if ignore_case:
        normalized = normalized.casefold()
    return normalized


def is_palindrome(
    value: Union[str, int],
    *,
    ignore_case: bool = True,
    ignore_non_alphanumeric: bool = True,
) -> bool:
    """Return ``True`` when ``value`` reads the same in both directions.

    Args:
        value: The text to check. Integers are accepted and compared as their
            string form; any other non-string type raises ``TypeError``.
        ignore_case: Compare case-insensitively (default ``True``).
        ignore_non_alphanumeric: Drop punctuation and whitespace before
            comparing (default ``True``).

    Returns:
        ``True`` for a palindrome, ``False`` otherwise. Empty input, and input
        whose normalized form is empty, is considered a palindrome.

    Raises:
        TypeError: If ``value`` is neither ``str`` nor ``int``.

    Complexity: O(n) time and O(n) extra space for the normalized copy.
    """
    if isinstance(value, bool) or not isinstance(value, (str, int)):
        raise TypeError(
            "is_palindrome() expects a string (or an int), got "
            f"{type(value).__name__}"
        )
    text = value if isinstance(value, str) else str(value)
    normalized = _normalize(
        text,
        ignore_case=ignore_case,
        ignore_non_alphanumeric=ignore_non_alphanumeric,
    )
    return normalized == normalized[::-1]


def _build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="palindrome.py",
        description="Check whether the given text is a palindrome.",
    )
    parser.add_argument(
        "text",
        nargs="?",
        help="text to check; if omitted, you are prompted for it",
    )
    parser.add_argument(
        "--case-sensitive",
        action="store_true",
        help="treat upper and lower case as different characters",
    )
    parser.add_argument(
        "--strict",
        action="store_true",
        help="keep punctuation and whitespace instead of ignoring them",
    )
    return parser


def main(argv: Optional[Sequence[str]] = None) -> int:
    """Run the command line interface and return the process exit code."""
    args = _build_parser().parse_args(argv)

    text = args.text
    if text is None:
        try:
            text = input("Enter text to check: ")
        except (EOFError, KeyboardInterrupt):
            print("No text provided.", file=sys.stderr)
            return EXIT_ERROR

    result = is_palindrome(
        text,
        ignore_case=not args.case_sensitive,
        ignore_non_alphanumeric=not args.strict,
    )
    verdict = "is a palindrome" if result else "is not a palindrome"
    print(f"{text!r} {verdict}")
    return EXIT_PALINDROME if result else EXIT_NOT_PALINDROME


if __name__ == "__main__":
    raise SystemExit(main())
