package main

import (
	"strconv"
	"strings"
)

type intList []int64

func (l intList) Len() int {
	return len(l)
}

func (l intList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l intList) Less(i, j int) bool {
	return l[i] < l[j]
}

//generateIntList from string columns specified
func generateIntList(list string) (intList, error) {

	var out []int64

	byComma := strings.Split(list, ",")

	for _, v := range byComma {

		if strings.Contains(v, "-") {

			ssv := strings.Split(v, "-")
			start, err := strconv.ParseInt(ssv[0], 10, 64)
			if err != nil {
				return nil, err
			}
			end, err := strconv.ParseInt(ssv[1], 10, 64)
			if err != nil {
				return nil, err
			}

			for i := start; i <= end; i++ {
				out = append(out, i)
			}

		} else {

			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}

			out = append(out, i)
		}
	}

	return out, nil
}
