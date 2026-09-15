"""Unit tests for :mod:`palindrome`."""

import unittest

from palindrome import is_palindrome, main, normalize


class NormalizeTests(unittest.TestCase):
    def test_lowercases_and_strips_non_alphanumeric(self):
        self.assertEqual(normalize("A man, a plan!"), "amanaplan")

    def test_keeps_digits(self):
        self.assertEqual(normalize("12 321"), "12321")

    def test_rejects_non_string(self):
        with self.assertRaises(TypeError):
            normalize(None)


class IsPalindromeTests(unittest.TestCase):
    def test_simple_word(self):
        self.assertTrue(is_palindrome("racecar"))

    def test_mixed_case_and_punctuation(self):
        self.assertTrue(is_palindrome("A man, a plan, a canal: Panama"))

    def test_numeric_string(self):
        self.assertTrue(is_palindrome("12321"))

    def test_empty_string(self):
        self.assertTrue(is_palindrome(""))

    def test_punctuation_only(self):
        self.assertTrue(is_palindrome("!!!"))

    def test_non_palindrome(self):
        self.assertFalse(is_palindrome("hello"))

    def test_non_palindrome_phrase(self):
        self.assertFalse(is_palindrome("A man, a plan, a canoe"))

    def test_rejects_non_string(self):
        with self.assertRaises(TypeError):
            is_palindrome(None)


class MainTests(unittest.TestCase):
    def test_exit_code_zero_for_palindrome(self):
        self.assertEqual(main(["level"]), 0)

    def test_exit_code_one_for_non_palindrome(self):
        self.assertEqual(main(["hello"]), 1)


if __name__ == "__main__":
    unittest.main()
