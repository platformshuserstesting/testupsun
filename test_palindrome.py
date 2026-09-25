"""Tests for :mod:`palindrome`.

Covers the pure :func:`palindrome.is_palindrome` and
:func:`palindrome.is_palindrome_strict` functions (including the rule that
empty input is not a palindrome) and the :func:`palindrome.main` CLI.

Run with ``python -m unittest discover -v`` or ``pytest -q``.
"""

from __future__ import annotations

import contextlib
import io
import sys
import unittest

import palindrome
from palindrome import is_palindrome, is_palindrome_strict, main


class IsPalindromeTests(unittest.TestCase):
    def test_simple_palindromes(self) -> None:
        for text in ("racecar", "level", "a", "abba"):
            with self.subTest(text=text):
                self.assertTrue(is_palindrome(text))

    def test_non_palindromes(self) -> None:
        for text in ("hello", "python", "ab"):
            with self.subTest(text=text):
                self.assertFalse(is_palindrome(text))

    def test_phrase_ignores_case_and_punctuation(self) -> None:
        self.assertTrue(is_palindrome("A man, a plan, a canal: Panama"))
        self.assertTrue(is_palindrome("Never odd or even"))
        self.assertFalse(is_palindrome("This is not one."))

    def test_strict_mode(self) -> None:
        self.assertTrue(is_palindrome("racecar", strict=True))
        self.assertFalse(is_palindrome("Racecar", strict=True))
        self.assertFalse(is_palindrome("never odd or even", strict=True))

    def test_empty_string_is_not_a_palindrome(self) -> None:
        self.assertFalse(is_palindrome(""))

    def test_punctuation_only_is_not_a_palindrome(self) -> None:
        self.assertFalse(is_palindrome("!?.,"))
        self.assertFalse(is_palindrome("  ,.!? "))

    def test_empty_string_strict_is_not_a_palindrome(self) -> None:
        self.assertFalse(is_palindrome("", strict=True))

    def test_single_space_strict_is_a_palindrome(self) -> None:
        self.assertTrue(is_palindrome(" ", strict=True))

    def test_non_string_input_is_coerced(self) -> None:
        self.assertTrue(is_palindrome(12321))
        self.assertFalse(is_palindrome(12345))

    def test_unicode_casefold(self) -> None:
        self.assertTrue(is_palindrome("Été"))
        self.assertFalse(is_palindrome("Été", strict=True))


class IsPalindromeStrictTests(unittest.TestCase):
    def test_exact_palindromes(self) -> None:
        for text in ("racecar", "abba", "a", " "):
            with self.subTest(text=text):
                self.assertTrue(is_palindrome_strict(text))

    def test_case_and_spaces_are_significant(self) -> None:
        self.assertFalse(is_palindrome_strict("Racecar"))
        self.assertFalse(is_palindrome_strict("never odd or even"))
        self.assertFalse(is_palindrome_strict("A man, a plan, a canal: Panama"))

    def test_empty_string_is_not_a_palindrome(self) -> None:
        self.assertFalse(is_palindrome_strict(""))

    def test_non_string_input_is_coerced(self) -> None:
        self.assertTrue(is_palindrome_strict(12321))
        self.assertFalse(is_palindrome_strict(12345))

    def test_strict_flag_matches_strict_function(self) -> None:
        for text in ("Abba", "abba", "racecar", "Racecar", "", " ", "!?.,", 12321):
            with self.subTest(text=text):
                self.assertEqual(
                    is_palindrome(text, strict=True), is_palindrome_strict(text)
                )


class MainTests(unittest.TestCase):
    def _run(self, argv: list[str], stdin: str | None = None) -> tuple[int, str]:
        buffer = io.StringIO()
        original_stdin = sys.stdin
        if stdin is not None:
            sys.stdin = io.StringIO(stdin)
        try:
            with contextlib.redirect_stdout(buffer):
                code = main(argv)
        finally:
            sys.stdin = original_stdin
        return code, buffer.getvalue()

    def test_palindrome_argument(self) -> None:
        code, output = self._run(["racecar"])
        self.assertEqual(code, 0)
        self.assertEqual(output, '"racecar" -> palindrome\n')

    def test_non_palindrome_argument(self) -> None:
        code, output = self._run(["hello"])
        self.assertEqual(code, 1)
        self.assertEqual(output, '"hello" -> not a palindrome\n')

    def test_multiple_positionals_are_joined(self) -> None:
        code, output = self._run(["never", "odd", "or", "even"])
        self.assertEqual(code, 0)
        self.assertEqual(output, '"never odd or even" -> palindrome\n')

    def test_strict_flag(self) -> None:
        code, output = self._run(["--strict", "Racecar"])
        self.assertEqual(code, 1)
        self.assertEqual(output, '"Racecar" -> not a palindrome\n')

    def test_strict_flag_exact_palindrome(self) -> None:
        code, output = self._run(["--strict", "racecar"])
        self.assertEqual(code, 0)
        self.assertEqual(output, '"racecar" -> palindrome\n')

    def test_strict_flag_empty_argument(self) -> None:
        code, output = self._run(["--strict", ""])
        self.assertEqual(code, 1)
        self.assertEqual(output, '"" -> not a palindrome (empty input)\n')

    def test_stdin_lines(self) -> None:
        code, output = self._run([], stdin="racecar\nhello\n")
        self.assertEqual(code, 1)
        self.assertEqual(
            output,
            '"racecar" -> palindrome\n"hello" -> not a palindrome\n',
        )

    def test_empty_argument(self) -> None:
        code, output = self._run([""])
        self.assertEqual(code, 1)
        self.assertEqual(output, '"" -> not a palindrome (empty input)\n')

    def test_blank_stdin_line(self) -> None:
        code, output = self._run([], stdin="racecar\n\n")
        self.assertEqual(code, 1)
        self.assertEqual(
            output,
            '"racecar" -> palindrome\n"" -> not a palindrome (empty input)\n',
        )

    def test_empty_stdin(self) -> None:
        code, output = self._run([], stdin="")
        self.assertEqual(code, 1)
        self.assertEqual(output, '"" -> not a palindrome (empty input)\n')


class VersionTests(unittest.TestCase):
    def _run_exit(self, argv: list[str], stdin: str | None = None) -> tuple[int, str]:
        buffer = io.StringIO()
        original_stdin = sys.stdin
        if stdin is not None:
            sys.stdin = io.StringIO(stdin)
        try:
            with contextlib.redirect_stdout(buffer):
                with self.assertRaises(SystemExit) as ctx:
                    main(argv)
        finally:
            sys.stdin = original_stdin
        return ctx.exception.code, buffer.getvalue()

    def test_version_constant_format(self) -> None:
        self.assertIsInstance(palindrome.__version__, str)
        self.assertRegex(palindrome.__version__, r"^\d+\.\d+\.\d+$")

    def test_version_flag_prints_version_and_exits_zero(self) -> None:
        code, output = self._run_exit(["--version"])
        self.assertEqual(code, 0)
        self.assertIn(palindrome.__version__, output)

    def test_version_flag_skips_check(self) -> None:
        code, output = self._run_exit(["--version", "racecar"])
        self.assertEqual(code, 0)
        self.assertIn(palindrome.__version__, output)
        self.assertNotIn("->", output)

    def test_version_flag_does_not_read_stdin(self) -> None:
        fake_stdin = io.StringIO("hi\n")
        original_stdin = sys.stdin
        sys.stdin = fake_stdin
        try:
            with contextlib.redirect_stdout(io.StringIO()):
                with self.assertRaises(SystemExit) as ctx:
                    main(["--version"])
        finally:
            sys.stdin = original_stdin
        self.assertEqual(ctx.exception.code, 0)
        self.assertEqual(fake_stdin.tell(), 0)

    def test_help_lists_version(self) -> None:
        code, output = self._run_exit(["--help"])
        self.assertEqual(code, 0)
        self.assertIn("--version", output)


if __name__ == "__main__":
    unittest.main()
