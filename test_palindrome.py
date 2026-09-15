"""Tests for the palindrome checker."""

import contextlib
import io
import unittest

from palindrome import (
    EXIT_INVALID_INPUT,
    EXIT_NOT_PALINDROME,
    EXIT_PALINDROME,
    is_palindrome,
    main,
)


class IsPalindromeTests(unittest.TestCase):
    def test_simple_palindrome(self):
        self.assertTrue(is_palindrome("racecar"))

    def test_simple_non_palindrome(self):
        self.assertFalse(is_palindrome("hello"))

    def test_single_character(self):
        self.assertTrue(is_palindrome("a"))

    def test_empty_after_normalization(self):
        self.assertTrue(is_palindrome(""))
        self.assertTrue(is_palindrome("!!!"))

    def test_mixed_case_ignored_by_default(self):
        self.assertTrue(is_palindrome("RaceCar"))

    def test_punctuation_and_spaces_ignored_by_default(self):
        self.assertTrue(is_palindrome("A man, a plan, a canal: Panama"))

    def test_numeric_input(self):
        self.assertTrue(is_palindrome(12321))
        self.assertFalse(is_palindrome(12345))

    def test_case_sensitive_comparison(self):
        self.assertFalse(is_palindrome("RaceCar", ignore_case=False))
        self.assertTrue(is_palindrome("racecar", ignore_case=False))

    def test_strict_comparison_keeps_spaces(self):
        self.assertFalse(is_palindrome("never odd or even", alphanumeric_only=False))
        self.assertTrue(is_palindrome("neveroddoreven", alphanumeric_only=False))


class MainTests(unittest.TestCase):
    def _run(self, argv):
        stream = io.StringIO()
        with contextlib.redirect_stdout(stream):
            code = main(argv)
        return code, stream.getvalue()

    def test_palindrome_exit_code(self):
        code, output = self._run(["racecar"])
        self.assertEqual(code, EXIT_PALINDROME)
        self.assertEqual(output.strip(), '"racecar" is a palindrome.')

    def test_non_palindrome_exit_code(self):
        code, output = self._run(["hello"])
        self.assertEqual(code, EXIT_NOT_PALINDROME)
        self.assertEqual(output.strip(), '"hello" is not a palindrome.')

    def test_multiple_arguments_are_joined(self):
        code, _ = self._run(["A", "man", "a", "plan", "a", "canal", "Panama"])
        self.assertEqual(code, EXIT_PALINDROME)

    def test_case_sensitive_flag(self):
        code, _ = self._run(["--case-sensitive", "Racecar"])
        self.assertEqual(code, EXIT_NOT_PALINDROME)

    def test_strict_flag(self):
        code, _ = self._run(["--strict", "never odd or even"])
        self.assertEqual(code, EXIT_NOT_PALINDROME)

    def test_no_comparable_characters(self):
        code, output = self._run(["!!!"])
        self.assertEqual(code, EXIT_INVALID_INPUT)
        self.assertIn("No comparable characters", output)


if __name__ == "__main__":
    unittest.main()
