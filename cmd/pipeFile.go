/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var (
	fileName string
	delay    int
)

// pipeFileCmd represents the pipeFile command
var pipeFileCmd = &cobra.Command{
	Use:   "pipeFile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runPipeFile()
	},
}

func init() {
	rootCmd.AddCommand(pipeFileCmd)

	pipeFileCmd.PersistentFlags().StringVarP(&fileName, "fileName", "f", "", "Path to file")
	pipeFileCmd.MarkPersistentFlagRequired("fileName")

	pipeFileCmd.PersistentFlags().IntVarP(&delay, "delay", "d", 40, "delay between text lines in milliseconds")
}

func runPipeFile() {
	fileOut, err := os.Create(fileName + ".out")
	if err != nil {
		fmt.Println("Error creating file: ", err)
		return
	}

	writer := bufio.NewWriter(fileOut)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go handleExit(sig, cleanup, fileOut)


    fileIn, err := os.Open(fileName)
    if err != nil {
        fmt.Println("Error opening file: ", err)
        return
    }
	defer fileIn.Close()

	fmt.Println("Piping to: ", fileName + ".out")
	for {
        scanner := bufio.NewScanner(fileIn)
		for scanner.Scan() {
			writer.WriteString(scanner.Text() + "\n")
			writer.Flush()
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file: ", err)
			return
		}
        fileIn.Seek(0,0)
		fmt.Println("End of file reached, restarting...")
	}

}
func cleanup(param interface{}) {
	fileOut, _ := param.(*os.File)
	err := fileOut.Close()
	if err != nil {
		fmt.Println("Error closing file: ", err)
	}
	fmt.Println("Exiting...")
	err = os.Remove(fileName + ".out")
	if err != nil {
		fmt.Println("Error deleting file: ", err)
		return
	}

}
