package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/stealthrocket/timecraft/format"
	"github.com/stealthrocket/timecraft/internal/object"
	"github.com/stealthrocket/timecraft/internal/print/human"
	"github.com/stealthrocket/timecraft/internal/timemachine"
	"github.com/stealthrocket/timecraft/internal/timemachine/wasicall"
	"github.com/stealthrocket/wasi-go"
	"github.com/stealthrocket/wasi-go/imports"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/sys"
)

const runUsage = `
Usage:	timecraft run [options] [--] <module> [args...]

Options:
   -c, --config                   Path to the timecraft configuration file (overrides TIMECRAFTCONFIG)
   -D, --dial addr                Expose a socket connected to the specified address
   -e, --env name=value           Pass an environment variable to the guest module
       --fly-blind                Disable recording of the guest module execution
   -h, --help                     Show this usage information
   -L, --listen addr              Expose a socket listening on the specified address
   -S, --sockets extension        Enable a sockets extension, one of none, auto, path_open, wasmedgev1, wasmedgev2 (default to auto)
       --record-batch-size size   Number of records written per batch (default to 4096)
       --record-compression type  Compression to use when writing records, either snappy or zstd (default to zstd)
   -T, --trace                    Enable strace-like logging of host function calls
`

func run(ctx context.Context, args []string) error {
	var (
		envs        stringList
		listens     stringList
		dials       stringList
		batchSize   = human.Count(4096)
		compression = compression("zstd")
		sockets     = sockets("auto")
		flyBlind    = false
		trace       = false
	)

	flagSet := newFlagSet("timecraft run", runUsage)
	customVar(flagSet, &envs, "e", "env")
	customVar(flagSet, &listens, "L", "listen")
	customVar(flagSet, &dials, "D", "dial")
	customVar(flagSet, &sockets, "S", "sockets")
	boolVar(flagSet, &trace, "T", "trace")
	boolVar(flagSet, &flyBlind, "fly-blind")
	customVar(flagSet, &batchSize, "record-batch-size")
	customVar(flagSet, &compression, "record-compression")

	if err := flagSet.Parse(args); err != nil {
		return err
	}

	envs = append(os.Environ(), envs...)
	args = flagSet.Args()

	if len(args) == 0 {
		return errors.New(`missing "--" separator before the module path`)
	}

	config, err := loadConfig()
	if err != nil {
		return err
	}
	registry, err := config.createRegistry()
	if err != nil {
		return err
	}

	wasmPath := args[0]
	wasmName := filepath.Base(wasmPath)
	wasmCode, err := os.ReadFile(wasmPath)
	if err != nil {
		return fmt.Errorf("could not read wasm file '%s': %w", wasmPath, err)
	}

	runtime := config.newRuntime(ctx)
	defer runtime.Close(ctx)

	wasmModule, err := runtime.CompileModule(ctx, wasmCode)
	if err != nil {
		return err
	}
	defer wasmModule.Close(ctx)

	// When running cmd.Root from testable examples, the standard streams are
	// not set to alternative files and the fd numbers are not 0, 1, 2.
	stdin := int(os.Stdin.Fd())
	stdout := int(os.Stdout.Fd())
	stderr := int(os.Stderr.Fd())

	builder := imports.NewBuilder().
		WithName(wasmName).
		WithArgs(args[1:]...).
		WithEnv(envs...).
		WithDirs("/").
		WithListens(listens...).
		WithDials(dials...).
		WithStdio(stdin, stdout, stderr).
		WithSocketsExtension(string(sockets), wasmModule).
		WithTracer(trace, os.Stderr)

	if !flyBlind {
		var c timemachine.Compression
		switch compression {
		case "snappy":
			c = timemachine.Snappy
		case "zstd":
			c = timemachine.Zstd
		case "none", "":
			c = timemachine.Uncompressed
		default:
			return fmt.Errorf("invalid compression type %q", compression)
		}

		processID := uuid.New()
		startTime := time.Now().UTC()

		module, err := registry.CreateModule(ctx, &format.Module{
			Code: wasmCode,
		}, object.Tag{
			Name:  "timecraft.module.name",
			Value: wasmModule.Name(),
		})
		if err != nil {
			return err
		}

		runtime, err := registry.CreateRuntime(ctx, &format.Runtime{
			Runtime: "timecraft",
			Version: currentVersion(),
		})
		if err != nil {
			return err
		}

		config, err := registry.CreateConfig(ctx, &format.Config{
			Runtime: runtime,
			Modules: []*format.Descriptor{module},
			Args:    append([]string{wasmName}, args...),
			Env:     envs,
		})
		if err != nil {
			return err
		}

		process, err := registry.CreateProcess(ctx, &format.Process{
			ID:        processID,
			StartTime: startTime,
			Config:    config,
		})
		if err != nil {
			return err
		}

		if err := registry.CreateLogManifest(ctx, processID, &format.Manifest{
			Process:   process,
			StartTime: startTime,
		}); err != nil {
			return err
		}

		logSegment, err := registry.CreateLogSegment(ctx, processID, 0)
		if err != nil {
			return err
		}
		defer logSegment.Close()
		logWriter := timemachine.NewLogWriter(logSegment)

		recordWriter := timemachine.NewLogRecordWriter(logWriter, int(batchSize), c)
		defer recordWriter.Flush()

		builder = builder.WithWrappers(func(s wasi.System) wasi.System {
			return wasicall.NewRecorder(s, startTime, func(record *timemachine.RecordBuilder) {
				if err := recordWriter.WriteRecord(record); err != nil {
					panic(err)
				}
			})
		})

		fmt.Fprintf(os.Stderr, "%s\n", processID)
	}

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var system wasi.System
	ctx, system, err = builder.Instantiate(ctx, runtime)
	if err != nil {
		return err
	}
	defer system.Close(ctx)

	return instantiate(ctx, runtime, wasmModule)
}

func instantiate(ctx context.Context, runtime wazero.Runtime, compiledModule wazero.CompiledModule) error {
	module, err := runtime.InstantiateModule(ctx, compiledModule, wazero.NewModuleConfig().
		WithStartFunctions())
	if err != nil {
		return err
	}
	defer module.Close(ctx)

	ctx, cancel := context.WithCancelCause(ctx)
	go func() {
		_, err := module.ExportedFunction("_start").Call(ctx)
		module.Close(ctx)
		cancel(err)
	}()

	<-ctx.Done()

	err = context.Cause(ctx)
	switch err {
	case context.Canceled, context.DeadlineExceeded:
		err = nil
	}

	switch e := err.(type) {
	case *sys.ExitError:
		if rc := e.ExitCode(); rc != 0 {
			return exitCode(rc)
		}
		err = nil
	}

	return err
}
