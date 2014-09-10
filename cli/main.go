package cli

import (
	"code.google.com/p/go.crypto/ssh/terminal"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	grypt "polydawn.net/grypt"
)

func Run(args []string) {
	app := cli.NewApp()
	app.Name = "grypt"
	app.Usage = "grypt is a tool that allows you to store secrets in a git repository."
	app.Version = "v0.1"

	app.Commands = []cli.Command{
		{
			Name:  "generate-key",
			Usage: "generate a key used to lock and unlock secrets you ask grypt to keep.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "keyring",
					Value: "default",
					Usage: "the keyring to use.  for example, you may wish to keep multiple secrets for one project, and lock some of them with a keyring for prod, and another with a keyring for dev.",
				},
				cli.BoolFlag{
					Name:  "random-key",
					Usage: "if random-key is enabled, a completely random key will be generated.  when using this kind of key, **you will not be able to regenerate your key from the password alone**, and you will need to implement some kind of out-of-band management strategy for sharing this key with your collaborators.",
				},
				cli.StringFlag{
					Name:   "password",
					Usage:  "provide a password for generating the key.  if not provided, it will be prompted for interactively -- and this is the preferred mechanism, as this prevents your password from showing up trivially to others on the same machine in the output of commands like `ps`!  in grypt's default behavior, this password will be used to derive the symmetric key which will then encrypt your secrets; if combined with the '--random-key' option, this will instead cause the random key to be itself protected in another layer of symmetric encryption with a password-derived key (much like adding a password to ssh keys).",
					EnvVar: "GRYPT_PASSWORD",
				},
				cli.StringFlag{
					Name: "scheme",
					Usage: `selects the encryption schema to use to keep your secrets.

						Valid encryption schemes are:

						 * AES-256/SHA-256          (default, aes256sha256)
						 * AES-256/Keccak-256       (keccak, aes256keccak256)
						 * AES-256/BLAKE2-256       (blake2, aes256blake2256)
						 * Blowfish-448/SHA-256     (blowfish, blowfish448sha256)
						 * Blowfish-448/BLAKE2-512  (blakefish, blowfish448blake2512)
						`,
					Value:  "aes256sha256",
					EnvVar: "GRYPT_PASSWORD",
				},
			},
			Action: func(c *cli.Context) {
				encryptionScheme, err := grypt.ParseScheme(c.String("scheme"))
				if err != nil {
					panic(fmt.Sprintf("Unable to determine encryption scheme: %v", err))
				}

				password := []byte(c.String("password"))
				if len(password) == 0 {
					// interactive prompt
					fmt.Fprintf(os.Stderr, "passphrase: ")
					var err error
					password, err = terminal.ReadPassword(0)
					if err != nil {
						panic(err)
					}
				}

				GenerateKey(
					c.String("keyring"),
					c.Bool("random-key"),
					password,
					encryptionScheme,
				)
			},
		},
		{
			Name:   "keep-secret",
			Usage:  "tells grypt to keep this file a secret for you.",
			Action: func(c *cli.Context) { /* TODO */ },
			// TODO: not sure how to get this cli library to generate help that hints there's a nonoptional positional argument.  might have to replace large swaths of their helptext template.
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "keyring",
					Value: "default",
					Usage: "the keyring to use.  for example, you may wish to keep multiple secrets for one project, and lock some of them with a keyring for prod, and another with a keyring for dev.",
				},
			},
		},
		{
			Name:   "unlock",
			Usage:  "tells grypt that this repo has secrets, and it's time to open them now.",
			Action: func(c *cli.Context) { /* TODO */ },
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "keyring",
					Value: "default",
					Usage: "the keyring to use.  for example, you may wish to keep multiple secrets for one project, and lock some of them with a keyring for prod, and another with a keyring for dev.  if some secrets have already been unlocked, this does not lock them; it just unlocks the new secrets you asked for (in other words, you are not limited to using one key to open a set of secrets at a time; you may open several sets of secrets under their respective keys at the same time).",
				},
			},
			// note that there really is no option for only doing unlock on specific files.  we can't think of a use for that -- if you've already entered credentials/keys to unlock any file with those keys, you clearly trust this hardware with all files under those same keys, so we might as well open them all.
		},
		{
			Name:   "lock",
			Usage:  "tells grypt to lock down all secrets again.",
			Action: func(c *cli.Context) { /* TODO */ },
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "keyring",
					Value: "",
					Usage: "the keyring to use.  if not specified, grypt will lock *all* secrets.",
				},
			},
			// note that this is really kind of a strange operation.  unless you have full-disk encryption and know the (many) practical limitations of securing hardware, be advised that this command alone is unlikely to wipe all traces of your secrets from recoverable reach if you're about to hand your device over to a sufficiently capable adversary.
		},
	}

	app.Run(args)
}
