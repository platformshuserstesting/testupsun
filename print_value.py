#!/usr/bin/env python3
"""Interactive printing program.

Prompts the user to enter a value and prints it back. Uses only the
Python standard library, so it runs on a stock Python 3 interpreter.

Usage:
    python print_value.py

    # or with piped input:
    echo "hello" | python print_value.py
"""

import sys

PROMPT = "Enter a value: "
EMPTY_MESSAGE = "No value entered."


def format_value(value):
    """Return the output string for a given input value.

    An empty (or whitespace-only) value yields a friendly message; any
    other value is echoed back to the caller.
    """
    if value is None or value.strip() == "":
        return EMPTY_MESSAGE
    return "You entered: {}".format(value)


def main():
    """Prompt for a value and print it back.

    Returns an exit code: 0 on success, 1 when input is aborted via
    EOF (Ctrl+D) or interrupt (Ctrl+C).
    """
    try:
        value = input(PROMPT)
    except (EOFError, KeyboardInterrupt):
        # Exit cleanly without dumping a raw traceback.
        print()
        return 1

    print(format_value(value))
    return 0


if __name__ == "__main__":
    sys.exit(main())
