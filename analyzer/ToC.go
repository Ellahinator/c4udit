package analyzer

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func ToC_Convertor(fileName string) (string, error) {
	fmt.Println(fileName)
	b, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	str_content := string(b)
	new_str := ""

	// fmt.Println(str_content)

	file_top := regexp.MustCompile(`## QA Issues found`)
	// fmt.Println(file_top.Find([]byte(str_content)))

	buf_l_nc := strings.Builder{}

	low := regexp.MustCompile(`\[L-[0-9][0-9]\].*`)
	nc := regexp.MustCompile(`\[N-[0-9][0-9]\].*`)

	// buf_l_nc.WriteString("## QA Issues found\n")
	buf_l_nc.WriteString("# Table of Contents\n")
	buf_l_nc.WriteString("Low\n")
	for i, x := range low.FindAllString(str_content, -1) {
		new_x := fmt.Sprintf("%d.%s", i+1, strings.Split(x, "]")[1])
		buf_l_nc.WriteString(fmt.Sprintf("- [%s](%s)\n", new_x, "#"+slugify(new_x)))
		// fmt.Println(x, " >>", new_x)

		new_str = strings.ReplaceAll(str_content, x, new_x)
		str_content = new_str

	}

	// fmt.Println(new_str)

	buf_l_nc.WriteString("\nNon-Critical\n")
	for i, x := range nc.FindAllString(str_content, -1) {
		new_x := fmt.Sprintf("%d.%s", i+1, strings.Split(x, "]")[1])
		buf_l_nc.WriteString(fmt.Sprintf("- [%s](%s)\n", new_x, "#"+slugify(new_x)))

		new_str = strings.Replace(new_str, x, new_x, 1)

	}

	// fmt.Println("🎯")
	// fmt.Println(buf_l_nc.String())

	new_str = string(file_top.ReplaceAll([]byte(new_str), []byte(buf_l_nc.String())))

	gas_head := regexp.MustCompile(`## Gas Findings`)
	buf_gas := strings.Builder{}

	buf_gas.WriteString("\n")
	buf_gas.WriteString("# Table of Contents\n")
	buf_gas.WriteString("Gas\n")
	gas := regexp.MustCompile(`\[G-[0-9][0-9]\].*`)
	for i, x := range gas.FindAllString(str_content, -1) {
		new_x := fmt.Sprintf("%d.%s", i+1, strings.Split(x, "]")[1])

		buf_gas.WriteString(fmt.Sprintf("- [%s](%s)\n", new_x, "#"+slugify(new_x)))

		new_str = strings.Replace(new_str, x, new_x, 1)
	}

	buf_gas.WriteString("\n## Gas Findings\n")

	new_str = string(gas_head.ReplaceAll([]byte(new_str), []byte(buf_gas.String())))

	// fmt.Println(new_str)

	return new_str, nil

}
