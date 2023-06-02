package main

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/stealthrocket/timecraft/format"
)

const exportUsage = `
Usage:	timecraft export <resource type> <resource id> <output file> [options]

   The export command reads resources from the time machine registry and writes
   them to local files. This command is useful to extract data generated by
   timecraft and visualize it with other tools.

   The last argument is the location where the resource is exported, typically
   a path on the file system. The special value "-" may be set to write the
   resource to stdout.

Options:
   -c, --config  Path to the timecraft configuration file (overrides TIMECRAFTCONFIG)
   -h, --help    Show this usage information
`

func export(ctx context.Context, args []string) error {
	flagSet := newFlagSet("timecraft export", exportUsage)

	args, err := parseFlags(flagSet, args)
	if err != nil {
		return err
	}
	if len(args) != 3 {
		perrorf(`Expected resource type, id, and output file as argument` + useCmd("export"))
		return exitCode(2)
	}
	resource, err := findResource("describe", args[0])
	if err != nil {
		perror(err)
		return exitCode(2)
	}
	config, err := loadConfig()
	if err != nil {
		return err
	}
	registry, err := config.openRegistry()
	if err != nil {
		return err
	}

	if resource.typ == "log" {
		// How should we handle logs?
		// - write the manifest.json + segments to a tar archive?
		// - combines the segments into a single log file?
		// - ???
		return errors.New(`TODO`)
	}

	var hash format.Hash
	if resource.typ == "process" {
		processID, err := uuid.Parse(args[1])
		if err != nil {
			return errors.New(`malformed process id (not a UUID)`)
		}
		manifest, err := registry.LookupLogManifest(ctx, processID)
		if err != nil {
			return err
		}
		hash = manifest.Process.Digest
	} else {
		desc, err := registry.LookupDescriptor(ctx, format.ParseHash(args[1]))
		if err != nil {
			return err
		}
		hash = desc.Digest
	}

	r, err := registry.LookupResource(ctx, hash)
	if err != nil {
		return err
	}
	defer r.Close()

	w := io.Writer(os.Stdout)
	if outputFile := args[2]; outputFile != "-" {
		f, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	_, err = io.Copy(w, r)
	return err
}