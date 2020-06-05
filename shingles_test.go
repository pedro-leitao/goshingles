package shingles

import (
	"testing"
)

var linearText string = `A B C D E F G A B C H I J K L M`
var loremIpsum string = `
Everything, everything, everything, everything..
In its right place
In its right place
In its right place
Right place

Yesterday I woke up sucking a lemon
Yesterday I woke up sucking a lemon
Yesterday I woke up sucking a lemon
Yesterday I woke up sucking a lemon

Everything, everything, everything..
In its right place
In its right place
Right place

There are two colours in my head
There are two colours in my head
What is that you tried to say?
What was that you tried to say?
Tried to say.. tried to say..
Tried to say.. tried to say..

Everything in its right place 
`

func TestSimpleIncorporate(t *testing.T) {
	var shingles Shingles
	shingles.Initialize(TRIGRAM)
	shingles.Incorporate(linearText, false)
	got := shingles.Count()
	expected := 13
	if got != expected {
		t.Errorf("TestSimpleIncorporate: got %v, expected %v", got, expected)
	}
	shingles.SortedWalk()
}

func TestComplexIncorporate(t *testing.T) {
	var shingles Shingles
	shingles.Initialize(TRIGRAM)
	shingles.Incorporate(loremIpsum, true)
	got := shingles.Count()
	expected := 35
	if got != expected {
		t.Errorf("TestComplexIncorporate: got %v, expected %v", got, expected)
	}
	shingles.SortedWalk()
}
