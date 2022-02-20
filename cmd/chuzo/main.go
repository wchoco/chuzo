package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/wchoco/chuzo"
)

func main() {
	cmd.Execute()
}

var cmd = &cobra.Command{
	Use:   "chuzo template_path src_path",
	Short: "Fill template",
	Args:  cobra.ExactArgs(2),
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	outfile, err := cmd.Flags().GetString("outfile")
	if err != nil {
		log.Fatal(err)
	}

	var out io.Writer
	if outfile == "stdout" {
		out = os.Stdout
	} else {
		out, err = os.Create(outfile)
		if err != nil {
			log.Fatal(err)
		}
	}

	m, err := chuzo.BuildMold(args[0])
	if err != nil {
		log.Fatal(err)
	}

	r, err := os.Open(args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	flush, err := cmd.Flags().GetBool("flush")
	if err != nil {
		log.Fatal(err)
	}

	mat, err := prepareMaterial(cmd, args)
	if err != nil {
		log.Fatal(err)
	}

	err = cast(out, m, mat, flush)
	if err != nil {
		log.Fatal(err)
	}
}

func cast(out io.Writer, m chuzo.Mold, mat chuzo.Material, flush bool) error {
	if flush {
		if err := m.Cast(out, mat); err != nil {
			return err
		}
	} else {
		tmp := bytes.NewBuffer([]byte{})
		if err := m.Cast(tmp, mat); err != nil {
			return err
		}
		_, err := io.Copy(out, tmp)
		if err != nil {
			return err
		}
	}
	return nil
}

func prepareMaterial(cmd *cobra.Command, args []string) (chuzo.Material, error) {
	srcType, err := cmd.Flags().GetString("src-type")
	if err != nil {
		return nil, err
	}

	var mat chuzo.Material
	switch srcType {
	case "yaml":
		mat = chuzo.YAMLMaterial{Path: args[1]}
	default:
		return nil, fmt.Errorf("unknown source type: %s", srcType)
	}

	return mat, nil
}

func init() {
	cmd.Flags().StringP("src-type", "t", "yaml", "source file type")
	cmd.Flags().StringP("outfile", "o", "stdout", "output filename")
	cmd.Flags().Bool("flush", false, "output sequential")
}
