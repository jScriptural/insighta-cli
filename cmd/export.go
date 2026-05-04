package cmd

import (
	"flag"
	"log"
)
func Export(args []string){
	fs := flag.NewFlagSet("export",flag.ExitOnError);
	format := fs.String("format","","format for the exported data-support csv and json")
	fs.Parse(args[0:2]);
	if len(args) > 2 {
		query = ParseArgs(fs.Args());
	}

	log.Println("query: ",query)
	_ = *format
}
