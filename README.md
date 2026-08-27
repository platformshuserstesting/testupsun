# testupsun

## Printing program (SCRUM-16)

`print_value.py` is a small, dependency-free interactive CLI script that
prompts you for a value and prints it back.

### Usage

Run it and type a value at the prompt:

```
python print_value.py
```

You can also pipe input into it:

```
echo "hello" | python print_value.py
```

### Behavior

- Enter a value and it is echoed back as `You entered: <value>`.
- Press Enter without typing anything and it prints `No value entered.`
- Sending EOF (Ctrl+D) or interrupting (Ctrl+C) exits cleanly without a
  traceback.
