"""Unit tests for :mod:`palindrome`."""

import contextlib
import io
import unittest
from unittest import mock

from palindrome import is_palindrome, main, normalize


class NormalizeTests(unittest.TestCase):
    def test_default_strips_punctuation_and_folds_case(self):
        self.assertEqual(normalize("A man, a plan!"), "amanaplan")

    def test_keep_case(self):
        self.assertEqual(normalize("A man!", ignore_case=False), "Aman")

    def test_keep_non_alnum(self):
        self.assertEqual(normalize("A man!", ignore_non_alnum=False), "a man!")

    def test_non_string_raises_type_error(self):
        with self.assertRaises(TypeError):
            normalize(121)


class IsPalindromeTests(unittest.TestCase):
    def test_simple_palindrome(self):
        self.assertTrue(is_palindrome("racecar"))

    def test_sentence_with_punctuation(self):
        self.assertTrue(is_palindrome("A man, a plan, a canal: Panama"))

    def test_not_a_palindrome(self):
        self.assertFalse(is_palindrome("hello"))

    def test_keeping_punctuation_breaks_sentence(self):
        self.assertFalse(is_palindrome("A man, a plan", ignore_non_alnum=False))

    def test_strict_comparison_is_case_sensitive(self):
        self.assertFalse(
            is_palindrome("RaceCar", ignore_case=False, ignore_non_alnum=False)
        )

    def test_numeric_string(self):
        self.assertTrue(is_palindrome(str(121)))
        self.assertFalse(is_palindrome(str(123)))

    def test_empty_and_single_character(self):
        self.assertTrue(is_palindrome(""))
        self.assertTrue(is_palindrome("a"))

    def test_punctuation_only_normalizes_to_empty(self):
        self.assertTrue(is_palindrome("!!!"))

    def test_non_string_raises_type_error(self):
        with self.assertRaises(TypeError):
            is_palindrome(121)


class CliTests(unittest.TestCase):
    def _run(self, argv, stdin_value=None):
        buffer = io.StringIO()
        with contextlib.redirect_stdout(buffer):
            if stdin_value is None:
                code = main(argv)
            else:
                with mock.patch("builtins.input", side_effect=stdin_value):
                    code = main(argv)
        return code, buffer.getvalue()

    def test_single_palindrome_argument(self):
        code, output = self._run(["racecar"])
        self.assertEqual(code, 0)
        self.assertEqual(output, "'racecar' is a palindrome\n")

    def test_multiple_arguments_report_each_and_fail(self):
        code, output = self._run(["racecar", "hello"])
        self.assertEqual(code, 1)
        self.assertEqual(
            output,
            "'racecar' is a palindrome\n'hello' is not a palindrome\n",
        )

    def test_strict_flag(self):
        code, output = self._run(["--strict", "Racecar"])
        self.assertEqual(code, 1)
        self.assertIn("is not a palindrome", output)

    def test_quiet_flag(self):
        code, output = self._run(["--quiet", "racecar"])
        self.assertEqual(code, 0)
        self.assertEqual(output, "true\n")

    def test_quiet_flag_for_non_palindrome(self):
        code, output = self._run(["--quiet", "hello"])
        self.assertEqual(code, 1)
        self.assertEqual(output, "false\n")

    def test_reads_from_stdin_when_no_arguments(self):
        code, output = self._run([], stdin_value=["racecar"])
        self.assertEqual(code, 0)
        self.assertEqual(output, "'racecar' is a palindrome\n")

    def test_eof_on_stdin_is_friendly(self):
        code, output = self._run([], stdin_value=EOFError())
        self.assertEqual(code, 1)
        self.assertEqual(output, "No input provided.\n")

    def test_help_exits_zero(self):
        buffer = io.StringIO()
        with contextlib.redirect_stdout(buffer):
            with self.assertRaises(SystemExit) as ctx:
                main(["--help"])
        self.assertEqual(ctx.exception.code, 0)
        self.assertIn("usage:", buffer.getvalue())


if __name__ == "__main__":
    unittest.main()
