"""Unit tests for the palindrome checker."""

import contextlib
import io
import unittest
from unittest import mock

from palindrome import is_palindrome, main


class IsPalindromeTests(unittest.TestCase):
    def test_simple_palindromes(self):
        for value in ("racecar", "madam", "noon", "a", ""):
            with self.subTest(value=value):
                self.assertTrue(is_palindrome(value))

    def test_non_palindromes(self):
        for value in ("hello", "python", "ab"):
            with self.subTest(value=value):
                self.assertFalse(is_palindrome(value))

    def test_mixed_case_is_ignored_by_default(self):
        self.assertTrue(is_palindrome("RaceCar"))

    def test_punctuation_and_whitespace_ignored_by_default(self):
        self.assertTrue(is_palindrome("A man, a plan, a canal: Panama"))
        self.assertTrue(is_palindrome("No 'x' in Nixon"))

    def test_digits(self):
        self.assertTrue(is_palindrome("12321"))
        self.assertFalse(is_palindrome("12345"))

    def test_integer_input_is_coerced(self):
        self.assertTrue(is_palindrome(12321))
        self.assertFalse(is_palindrome(12345))

    def test_punctuation_only_normalizes_to_empty(self):
        self.assertTrue(is_palindrome("!!!"))

    def test_unicode(self):
        self.assertTrue(is_palindrome("Ana"))
        self.assertTrue(is_palindrome("アカサカ"[::-1] + "アカサカ"))
        self.assertFalse(is_palindrome("アカサカ"))

    def test_case_sensitive_flag(self):
        self.assertFalse(is_palindrome("Racecar", ignore_case=False))
        self.assertTrue(is_palindrome("racecar", ignore_case=False))

    def test_strict_flag(self):
        self.assertFalse(
            is_palindrome("Race car", ignore_non_alphanumeric=False)
        )
        self.assertTrue(
            is_palindrome("racecar", ignore_non_alphanumeric=False)
        )

    def test_flag_combination(self):
        self.assertFalse(
            is_palindrome(
                "A man, a plan, a canal: Panama",
                ignore_case=False,
                ignore_non_alphanumeric=False,
            )
        )

    def test_invalid_types_raise_type_error(self):
        for value in (None, 1.5, ["r"], True):
            with self.subTest(value=value):
                with self.assertRaises(TypeError):
                    is_palindrome(value)


class MainTests(unittest.TestCase):
    def run_main(self, argv):
        stdout = io.StringIO()
        with contextlib.redirect_stdout(stdout):
            code = main(argv)
        return code, stdout.getvalue()

    def test_palindrome_argument(self):
        code, output = self.run_main(["racecar"])
        self.assertEqual(code, 0)
        self.assertIn("is a palindrome", output)

    def test_non_palindrome_argument(self):
        code, output = self.run_main(["hello"])
        self.assertEqual(code, 1)
        self.assertIn("is not a palindrome", output)

    def test_sentence_with_punctuation(self):
        code, _ = self.run_main(["A man, a plan, a canal: Panama"])
        self.assertEqual(code, 0)

    def test_numeric_argument(self):
        code, _ = self.run_main(["12321"])
        self.assertEqual(code, 0)

    def test_empty_argument(self):
        code, _ = self.run_main([""])
        self.assertEqual(code, 0)

    def test_strict_flag(self):
        code, _ = self.run_main(["--strict", "Race car"])
        self.assertEqual(code, 1)

    def test_case_sensitive_flag(self):
        code, _ = self.run_main(["--case-sensitive", "Racecar"])
        self.assertEqual(code, 1)

    def test_prompt_used_when_argument_missing(self):
        with mock.patch("builtins.input", return_value="madam"):
            code, output = self.run_main([])
        self.assertEqual(code, 0)
        self.assertIn("is a palindrome", output)

    def test_eof_on_prompt_returns_error_code(self):
        stderr = io.StringIO()
        with mock.patch("builtins.input", side_effect=EOFError):
            with contextlib.redirect_stderr(stderr):
                code, _ = self.run_main([])
        self.assertEqual(code, 2)
        self.assertIn("No text provided.", stderr.getvalue())

    def test_help_exits_zero(self):
        with contextlib.redirect_stdout(io.StringIO()):
            with self.assertRaises(SystemExit) as ctx:
                main(["-h"])
        self.assertEqual(ctx.exception.code, 0)


if __name__ == "__main__":
    unittest.main()
