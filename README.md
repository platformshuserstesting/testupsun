# testupsun

A minimal "Hello, World!" Python 3 program.

## Usage

Run the program directly with Python 3:

```sh
python hello.py
```

Expected output:

```
Hello, World!
```

The script is dependency-free (standard library only) and uses a guarded
`if __name__ == "__main__"` entry point, so it can also be imported without
printing anything—call `hello.main()` explicitly to run it.
