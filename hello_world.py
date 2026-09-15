"""Print a friendly greeting.

Usage:
    python hello_world.py
"""

from __future__ import annotations

import sys

GREETING = "Hello, World!"


def greeting() -> str:
    """Return the greeting text."""
    return GREETING


def main() -> int:
    """Print the greeting and return the process exit code."""
    print(greeting())
    return 0


if __name__ == "__main__":
    sys.exit(main())
