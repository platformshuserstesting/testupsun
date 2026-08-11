"""A small, self-contained Python test program.

Exposes a trivial, deterministic helper (:func:`add`) and a ``main``
entrypoint. Importing this module has no side effects; execution is
guarded by ``if __name__ == "__main__":``.
"""


def add(a: int, b: int) -> int:
    """Return the sum of two numbers."""
    return a + b


def main() -> None:
    """Print a deterministic greeting and a sample computation."""
    print("Hello from testupsun!")
    print(f"add(2, 3) = {add(2, 3)}")


if __name__ == "__main__":
    main()
