"""Check whether a piece of text is a palindrome.

Usage:
    python palindrome.py racecar
    python palindrome.py "A man, a plan, a canal: Panama"
    python palindrome.py --case-sensitive --strict "never odd or even"
"""

from __future__ import annotations

import argparse
import sys

EXIT_PALINDROME = 0
EXIT_NOT_PALINDROME = 1
EXIT_INVALID_INPUT = 2


def _normalize(value: str, ignore_case: bool, alphanumeric_only: bool) -> str:
    """Return ``value`` reduced to the characters that take part in the comparison."""
    text = value
    if alphanumeric_only:
        text = "".join(character for character in text if character.isalnum())
    if ignore_case:
        text = text.casefold()
    return text


def is_palindrome(value: object, *, ignore_case: bool = True, alphanumeric_only: bool = True) -> bool:
    """Return ``True`` when ``value`` reads the same forwards and backwards.

    Non-string input is converted with :func:`str`, so ``is_palindrome(12321)`` works.
    By default the comparison ignores case and non-alphanumeric characters.
    """
    text = value if isinstance(value, str) else str(value)
    normalized = _normalize(text, ignore_case, alphanumeric_only)
    return normalized == normalized[::-1]


def _build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="palindrome.py",
        description="Check whether the given text is a palindrome.",
    )
    parser.add_argument(
        "text",
        nargs="*",
        help="Text to check. Multiple arguments are joined with spaces.",
    )
    parser.add_argument(
        "--case-sensitive",
        action="store_true",
        help="Treat upper and lower case characters as different.",
    )
    parser.add_argument(
        "--strict",
        action="store_true",
        help="Keep spaces and punctuation instead of ignoring them.",
    )
    return parser


def main(argv: list[str] | None = None) -> int:
    """Run the command line interface and return the process exit code."""
    parser = _build_parser()
    args = parser.parse_args(argv)

    if args.text:
        text = " ".join(args.text)
    elif sys.stdin.isatty():
        text = input("Enter text to check: ")
    else:
        parser.print_usage(sys.stderr)
        print("error: no text supplied", file=sys.stderr)
        return EXIT_INVALID_INPUT

    ignore_case = not args.case_sensitive
    alphanumeric_only = not args.strict

    if not _normalize(text, ignore_case, alphanumeric_only):
        print("No comparable characters in the input.")
        return EXIT_INVALID_INPUT

    if is_palindrome(text, ignore_case=ignore_case, alphanumeric_only=alphanumeric_only):
        print(f'"{text}" is a palindrome.')
        return EXIT_PALINDROME

    print(f'"{text}" is not a palindrome.')
    return EXIT_NOT_PALINDROME


if __name__ == "__main__":
    sys.exit(main())
