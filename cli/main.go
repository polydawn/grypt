package cli

import (
	"fmt"
	"os"

	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/codegangsta/cli"

	"polydawn.net/grypt/gitutil"
	"polydawn.net/grypt/schema"
)

func Run(myName string, args ...string) {
	// set up shared bits of context information.  any subcommand will use this.
	ctx := gitutil.DetectContext()
	ctx.GryptName = myName
	ctx.GryptVersion = "v0.1"

	// first sack-n-grab on any subcommands that are the plumbing for git.
	// we don't run these through the main cli system because they're not complex enough to need it and we don't actually really want help text for these.
	if len(args) > 1 {
		// FIXME: don't know how to do keyring setup yet.  probably something with extra tuples tossed into gitattributes at the time of keep-secret.
		ctx.Keyring = "default"

		switch args[1] {
		case "git-clean":
			PlumbingClean(ctx, os.Stdin, os.Stdout)
			return
		case "git-smudge":
			PlumbingSmudge(ctx, os.Stdin, os.Stdout)
			return
		case "git-textconv":
			PlumbingTextconv(ctx, args[2], os.Stdout)
			return
		}
	}

	// construct the main cli args parser and help text generator
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
					Usage:  "provide a password for generating the key.  if not provided, it will be prompted for interactively -- and this is the preferred mechanism, as this prevents your password from showing up trivially to others on the same machine in the output of commands like `ps`!  in grypt's default behavior, this password will be used to derive the symmetric key which will then encrypt your secrets.",
					EnvVar: "GRYPT_PASSWORD",
				},
				cli.StringFlag{
					Name: "scheme",
					Usage: `selects the encryption schema to use to keep your secrets.

						Valid encryption schemes are:

						 * AES-256/SHA-256          (default, aes256sha256ctr)
						 * AES-256/Keccak-256       (keccak, aes256keccak256ctr)
						 * AES-256/BLAKE2-256       (blake2, aes256blake2256ctr)
						 * Blowfish-448/SHA-256     (blowfish, blowfish448sha256ctr)
						 * Blowfish-448/BLAKE2-512  (blakefish, blowfish448blake2512ctr)
						`,
					Value: "aes256sha256ctr",
				},
			},
			Action: func(c *cli.Context) {
				encryptionScheme := schema.ParseSchema(c.String("scheme"))

				password := []byte(c.String("password"))
				if len(password) == 0 {
					// interactive prompt
					fmt.Fprintf(os.Stderr, "passphrase: ")
					var err error
					password, err = terminal.ReadPassword(0) // jesus christ this will leave your terminal fucked if you ctrl-c out of it -.- come on guise
					if err != nil {
						panic(err)
					}
				}

				ctx.Keyring = c.String("keyring")

				GenerateKey(
					ctx,
					c.Bool("random-key"),
					password,
					encryptionScheme,
				)
			},
		},
		{
			Name:  "keep-secret",
			Usage: "tells grypt to keep this file a secret for you.",
			// TODO: not sure how to get this cli library to generate help that hints there's a nonoptional positional argument.  might have to replace large swaths of their helptext template.
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "keyring",
					Value: "default",
					Usage: "the keyring to use.  for example, you may wish to keep multiple secrets for one project, and lock some of them with a keyring for prod, and another with a keyring for dev.",
				},
			},
			Action: func(c *cli.Context) {
				// ick, because handling positional args manually is what i wanted to do
				if len(c.Args()) == 0 {
					panic("which files should we keep secret?")
				}
				files := c.Args()

				ctx.Keyring = c.String("keyring")

				KeepSecret(
					ctx,
					files,
				)
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

	// parse and run.  dispatches control to command implementations via the Action function pointers in the structs above.
	app.Run(args)
}
