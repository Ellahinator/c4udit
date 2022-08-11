package analyzer

import (
	"fmt"
	"testing"
)

func TestToc(t *testing.T) {
	file := "../c4udit-report.md"
	fmt.Println("🧪")
	ToC_Convertor(file)
	fmt.Println("🧪")

}
