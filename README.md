grypt is a tool that allows one to store secrets in a git repository.

Getting Started
===============

Here's an example to start a repository using grypt, assuming you're inside a
repository.

If you want a random key:
	% grypt keygen .git/key

If you want to derive a key from a passphrase (perhaps for easy sharing later):
	% grypt phrase .git/key

Note: run `keygen` or `phrase`, not both.

	% grypt init .git/key

grypt will print out a suggestion on what to enter in the repository's
`.gitattributes` file. For more information, see gitattributes(5).

`grypt help` will display some online help.

How It Works
============

grypt uses deterministic encryption and enciphers/deciphers data as it is
written to the git object store. If a repository is not configured to use grypt,
the encrypted blob is displayed. git's filter support is used for this, see
git-config(1) for more information.
