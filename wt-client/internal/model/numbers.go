package model

import (
	"fmt"
)

type Heading uint16

func (h Heading) Narrate() []string {
	str := fmt.Sprintf("%03d", h)
	sequence := []string{}

	// fmt.Println("checking", str)
	hundreeds := int(str[0] - '0')
	tens := int(str[1] - '0')
	rest := int(str[2] - '0')
	// fmt.Println(hundreeds, tens, rest)
	if hundreeds > 0 {
		sequence = append(sequence, fmt.Sprintf("%d", hundreeds), "100")
	}
	if tens > 1 {
		sequence = append(sequence, fmt.Sprintf("%d", tens)+"0")
	}
	if rest != 0 {
		if tens == 1 {
			sequence = append(sequence, fmt.Sprintf("1%d", rest))
		} else {
			sequence = append(sequence, fmt.Sprintf("%d", rest))
		}
	}

	return sequence
}
