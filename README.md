# Mondou REPL

A simple interpreted language for chatting with your computer. The word _mondou_ (jp. 問答) consists of kanji "question" and "answer" and has a meaning of "dialogue".

## Usage

`1 2 + 3 * 4 - 5 /` - evaluates postfix notation expression. The example results in `1`.

`!num 1 2 + 3 * 4 - 5 /` - stores the value in variable. When variable doesn't exist, it's created.

`<var` - prints the value to stdout.

`>var` - reads the value as a string from stdin and stores it in variable. When variable doesn't exist, exception is thrown.

`:/path/to/file.mondou` - loads the specified file and executes mondou script inside it.

`<"Hello, world!"  // hello world` - single line comment.

`@func_name(arg1, arg2=228) { ... }` - function definition. COMING SOON!

## Types

- int (`0`, `-4`, `1048576`)
- float64 (`3.14`, `1.`, `.618`)
- string (`"\"some text\"\nand a new line"`)
