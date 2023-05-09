package latexToImage

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"github.com/google/uuid"
)

func ConvertLatexToImage(latexCode string) (outputFilePath string, err error) {

	template := `
	\documentclass[12pt]{article}
	\usepackage[margin=2.5 cm]{geometry}
	\usepackage{hyperref,enumerate,float,amsfonts,amssymb,amsmath,graphics,graphicx,mathtools}
	\thispagestyle{empty}
	\setlength{\parindent}{0pt}
	\begin{document}
	\fontsize{14pt}{20pt}\selectfont
	%s
	\end{document}
`

	completeCode := fmt.Sprintf(template, latexCode)
	fmt.Println(completeCode)

	// generate a random filename
	id := uuid.New()
	inputFilePath := "pdfs/" + id.String() + ".pdf"
	outputFilePath = "imgs/" + id.String() + ".png"

	fmt.Println("Generated UUID:" + inputFilePath)
	fmt.Println("Generated UUID:" + outputFilePath)

	// Create a new command that runs pdflatex.
	cmd := exec.Command("pdflatex", "--jobname="+inputFilePath[:len(inputFilePath)-4]) // without the extension

	// Set the input to the LaTeX code.
	cmd.Stdin = bytes.NewBufferString(completeCode)

	// Run the command and check for errors.
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error from pdflatex: %v", err)
		return
	}

	// Create a new command that runs magick .
	cmd = exec.Command("magick", "-density", "350", inputFilePath, "-quality", "180", "-trim", "-border", "15%", outputFilePath)



	// Run the command and check for errors.
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error from magick: %v", err)
		return
	}

	// delete the pdf file
	if err = os.Remove(inputFilePath); err != nil {
		fmt.Printf("Error deleting the pdf file: %v", err)
		return
	}

	// return the output
	return outputFilePath, err

}
