package util

import (
	"fmt"
)

type calNumber []number

type number struct {
	num  int
	base int
}

func NewNumber(num int, base int) number {
	return number{num: num, base: base}
}

func NewCalNumber(nums ...number) calNumber {
	calNum := calNumber{}

	for _, num := range nums {
		calNum = append(calNum, num)
	}

	return calNum
}

func (c calNumber) Sub(subtrahend []number) (difference calNumber, err error) {
	minuend := c
	difference = make(calNumber, len(c))
	for i := len(minuend) - 1; i >= 0; i-- {
		if minuend[i].num < subtrahend[i].num {
			if i-1 < 0 {
				return nil, fmt.Errorf("subtrahend greater than minuend failed: %v - %v", minuend, subtrahend)
			}
			minuend[i-1].num--
			minuend[i].num += minuend[i].base
		}
		difference[i].num = minuend[i].num - subtrahend[i].num
		difference[i].base = minuend[i].base
	}
	return
}

func (c calNumber) ToNumbers() []number {
	return []number(c)
}

func (c calNumber) ToInts() []int {
	ints := []int{}

	for _, num := range c {
		ints = append(ints, num.num)
	}

	return ints
}

func (n *number) GetBase() int {
	return n.base
}

func (n *number) GetNum() int {
	return n.num
}
