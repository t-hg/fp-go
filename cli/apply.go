package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	C "github.com/urfave/cli/v2"
)

func generateSequenceT(f *os.File, i int) {
	fmt.Fprintf(f, "\n// SequenceT%d is a utility function used to implement the sequence operation for higher kinded types based only on map and ap.\n", i)
	fmt.Fprintf(f, "// The function takes %d higher higher kinded types and returns a higher kinded type of a [Tuple%d] with the resolved values.\n", i, i)
	fmt.Fprintf(f, "func SequenceT%d[\n", i)
	// map as the starting point
	fmt.Fprintf(f, "  MAP ~func(HKT_T1,")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  func(T%d)", j+1)
	}
	fmt.Fprintf(f, " ")
	fmt.Fprintf(f, "T.")
	writeTupleType(f, i)
	fmt.Fprintf(f, ")")
	if i > 1 {
		fmt.Fprintf(f, " HKT_F")
		for k := 1; k < i; k++ {
			fmt.Fprintf(f, "_T%d", k+1)
		}
	} else {
		fmt.Fprintf(f, " HKT_TUPLE%d", i)
	}
	fmt.Fprintf(f, ",\n")
	// the applicatives
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "  AP%d ~func(", j)
		fmt.Fprintf(f, "HKT_F")
		for k := j; k < i; k++ {
			fmt.Fprintf(f, "_T%d", k+1)
		}
		fmt.Fprintf(f, ", HKT_T%d)", j+1)
		if j+1 < i {
			fmt.Fprintf(f, " HKT_F")
			for k := j + 1; k < i; k++ {
				fmt.Fprintf(f, "_T%d", k+1)
			}
		} else {
			fmt.Fprintf(f, " HKT_TUPLE%d", i)
		}
		fmt.Fprintf(f, ",\n")
	}

	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  T%d,\n", j+1)
	}
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  HKT_T%d, // HKT[T%d]\n", j+1, j+1)
	}
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "  HKT_F")
		for k := j; k < i; k++ {
			fmt.Fprintf(f, "_T%d", k+1)
		}
		fmt.Fprintf(f, ", // HKT[")
		for k := j; k < i; k++ {
			fmt.Fprintf(f, "func(T%d) ", k+1)
		}
		fmt.Fprintf(f, "T.")
		writeTupleType(f, i)
		fmt.Fprintf(f, "]\n")
	}
	fmt.Fprintf(f, "  HKT_TUPLE%d any, // HKT[", i)
	writeTupleType(f, i)
	fmt.Fprintf(f, "]\n")
	fmt.Fprintf(f, "](\n")

	// the callbacks
	fmt.Fprintf(f, "  fmap MAP,\n")
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "  fap%d AP%d,\n", j, j)
	}
	// the parameters
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  t%d HKT_T%d,\n", j+1, j+1)
	}
	fmt.Fprintf(f, ") HKT_TUPLE%d {\n", i)

	fmt.Fprintf(f, "  r1 := fmap(t1, tupleConstructor%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, "]())\n")

	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "  r%d := fap%d(r%d, t%d)\n", j+1, j, j, j+1)
	}
	fmt.Fprintf(f, "  return r%d\n", i)

	fmt.Fprintf(f, "}\n")
}

func generateTupleConstructor(f *os.File, i int) {
	// Create the optionize version
	fmt.Fprintf(f, "\n// tupleConstructor%d returns a curried version of [T.MakeTuple%d]\n", i, i)
	fmt.Fprintf(f, "func tupleConstructor%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, " any]()")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, " func(T%d)", j+1)
	}
	fmt.Fprintf(f, " T.Tuple%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, "] {\n")

	fmt.Fprintf(f, "  return F.Curry%d(T.MakeTuple%d[", i, i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, "])\n")

	fmt.Fprintf(f, "}\n")
}

func generateApplyHelpers(filename string, count int) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	pkg := filepath.Base(absDir)
	f, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return err
	}
	defer f.Close()
	// log
	log.Printf("Generating code in [%s] for package [%s] with [%d] repetitions ...", filename, pkg, count)

	// some header
	fmt.Fprintln(f, "// Code generated by go generate; DO NOT EDIT.")
	fmt.Fprintln(f, "// This file was generated by robots at")
	fmt.Fprintf(f, "// %s\n", time.Now())

	fmt.Fprintf(f, "package %s\n\n", pkg)

	// print out some helpers
	fmt.Fprintf(f, `
import (
	F "github.com/ibm/fp-go/function"
	T "github.com/ibm/fp-go/tuple"
)
`)

	for i := 1; i <= count; i++ {
		// tuple constructor
		generateTupleConstructor(f, i)
		// sequenceT
		generateSequenceT(f, i)
	}

	return nil
}

func ApplyCommand() *C.Command {
	return &C.Command{
		Name:  "apply",
		Usage: "generate code for the sequence operations of apply",
		Flags: []C.Flag{
			flagCount,
			flagFilename,
		},
		Action: func(ctx *C.Context) error {
			return generateApplyHelpers(
				ctx.String(keyFilename),
				ctx.Int(keyCount),
			)
		},
	}
}
