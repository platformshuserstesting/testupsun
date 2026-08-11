"""A minimal Python test program.

Exposes a pure helper function (:func:`add`) and a runnable entrypoint
(:func:`main`). Importing this module has no side effects; execution only
happens when run directly via ``python main.py``.
"""


def add(a: int, b: int) -> int:
    """Return the sum of two numbers."""
    return a + b


def main() -> None:
    """Print a deterministic greeting and a sample computation."""
    print("Hello from testupsun!")
    print(f"2 + 3 = {add(2, 3)}")


if __name__ == "__main__":
    main()
