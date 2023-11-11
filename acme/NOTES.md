### Mouse Chording

Different actions can be performed by combining mouse button clicks
with keyboard key presses.

**Middle-Click (Button 2) + Right-Click (Button 3)**
Select a range of text in an acme window.

**Middle-Click (Button 2) + Left-Click (Button 1)**
Execute text as a command such as opening a file or sorting a list.

**Middle-Click (Button 2) + Shift + Right-Click (Button 3)**
Copy a text selection.

**Middle-Click (Button 2) + Shift + Left-Click (Button 1)**
Paste text at the current cursor position.

These mouse chording options provide a convenient way to perform various actions
in the acme editor by combining mouse clicks with modifier keys. They enhance
the editing experience and streamline common tasks within the acme environment.

`~/.acme/lookup`:

```
addr {
	select $+
	| plumber
}
```

/lib/plumb/rules:

```
$ext 'txt' {
    data | lookupper
}
```

```python
lookupper:
#!/usr/bin/env python3

import sys

word = sys.stdin.read().strip()
# Perform dictionary lookup using the word
# Display the results or perform any desired actions
```

## Testing

1. Open a text file in Acme that contains some words.
2. Select a word by clicking and dragging the mouse over it.
3. Press the right mouse button to pop up the menu.
4. From the menu, select the "Look Up" option or any other name you provided in
   the plumber rule.
5. Acme will use the plumber to invoke the lookupper script, passing the
   selected word as input.
6. The lookupper script will receive the word, perform the dictionary lookup,
   and display the results.

In `acme` a plumbfile is a special file that triggers an action or command based
on user interactions with file contents. It "listens" for events performs
actions in the editor.


The syntax and structure of plumbfiles are specific to Acme and follow a
particular format. Plumbfiles are typically stored in the plumb/ directory
within the user's home directory.

## Plumbfile

A plumbfile consists of one or more lines of text, where each line
represents a rule or configuration for handling specific events. Each
rule is defined using the following syntax:

Components of a plumbfile rule:

- `name`: A unique identifier for the rule. It is used to reference
  the rule in other parts of the Acme environment.
- `pattern`: A pattern that matches the content of the file being
  plumbed. It can be a regular expression or a literal string.
- `type`: The type of event that triggers the rule. It can be one of
  the following:
  - `put`: Triggered when text is put into the file.
  - `delete`: Triggered when text is deleted from the file.
  - `get`: Triggered when the file is read.
  - `look`: Triggered when a particular sequence of characters is
	typed while the file is being read.
- `command`: The command to execute when the rule is triggered.  It
  can be any valid Acme command or script.

Plumbing Rules Syntax:

1. Rule Syntax: `type:pattern -> action`
2. Action Syntax: `shell_command`
3. Variables: `varname=value`
4. Comments: `# This is a comment`
