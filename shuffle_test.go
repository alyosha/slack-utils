package utils

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/nlopes/slack"
)

var (
	testUserList     = []string{"UABCDEFG", "UWXYZ123", "ULMNOP345", "U0000111"}
	testUserListCopy = []string{"UABCDEFG", "UWXYZ123", "ULMNOP345", "U0000111"}
	seed             = int64(21)
	groupSize        = 2
)

func TestShuffle(t *testing.T) {
	var shuffled bool

	client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", "fakeurl.com")))

	shuffle := Shuffle{
		Client:    client,
		GroupSize: groupSize,
		Rand:      rand.New(rand.NewSource(seed)),
	}

	shuffledUsers := shuffle.Shuffle(testUserList)

	for i, user := range shuffledUsers {
		if testUserListCopy[i] != user {
			shuffled = true
		}
	}

	if !shuffled {
		t.Fatal("failed to shuffle user list, order is identical")
	}
}

func TestSplit(t *testing.T) {
	client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", "fakeurl.com")))

	shuffle := Shuffle{
		Client:    client,
		GroupSize: groupSize,
		Rand:      rand.New(rand.NewSource(seed)),
	}

	groups := shuffle.Split(testUserListCopy)

	if len(groups) != 2 {
		t.Fatalf("expected two groups, got: %v", len(groups))
	}

	if groups[0][0] != "UABCDEFG" {
		t.Fatalf("expected to receive user UABCDEFG, but received: %v", groups[0][0])
	}

	if groups[1][0] != "ULMNOP345" {
		t.Fatalf("expected to receive user ULMNOP345, but received: %v", groups[1][0])
	}
}
