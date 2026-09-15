"""Tests for the hello world program."""

import contextlib
import io
import unittest

from hello_world import greeting, main


class HelloWorldTests(unittest.TestCase):
    def test_greeting_text(self):
        self.assertEqual(greeting(), "Hello, World!")

    def test_main_prints_greeting(self):
        stream = io.StringIO()
        with contextlib.redirect_stdout(stream):
            code = main()
        self.assertEqual(code, 0)
        self.assertEqual(stream.getvalue(), "Hello, World!\n")


if __name__ == "__main__":
    unittest.main()
