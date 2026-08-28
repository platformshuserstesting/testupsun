"""Tests for the hello module."""

import contextlib
import io
import unittest

import hello


class MainTest(unittest.TestCase):
    def test_main_prints_greeting(self):
        buffer = io.StringIO()
        with contextlib.redirect_stdout(buffer):
            hello.main()
        self.assertEqual(buffer.getvalue(), "Hello, World!\n")


if __name__ == "__main__":
    unittest.main()
