package internal_test

import (
	"testing"

	"github.com/moreirathomas/golastic/internal"
)

func TestValidate(t *testing.T) {
	valid := internal.Book{
		Title: "Harry Potter and the Philosopher's Stone",
		Author: internal.Author{
			Firstname: "J. K.",
			Lastname:  "Rowling",
		},
		Abstract: "Harry Potter's life is miserable. His parents are dead and he's stuck with his heartless relatives, who force him to live in a tiny closet under the stairs. But his fortune changes when he receives a letter that tells him the truth about himself: he's a wizard. A mysterious visitor rescues him from his relatives and takes him to his new home, Hogwarts School of Witchcraft and Wizardry.",
	}

	if err := valid.Validate(false); err != nil {
		t.Errorf("unexpected error: want nil, got %s", err)
	}

	invalid := internal.Book{
		Title: "",
		Author: internal.Author{
			Firstname: "J. K. \u200d",
			Lastname:  "Rowling",
		},
		Abstract: "Harry Potter's life is miserable. His parents are dead and he's stuck with his heartless relatives, who force him to live in a tiny closet under the stairs. But his fortune changes when he receives a letter that tells him the truth about himself: he's a wizard. A mysterious visitor rescues him from his relatives and takes him to his new home, Hogwarts School of Witchcraft and Wizardry.",
	}

	if err := invalid.Validate(false); err == nil {
		t.Errorf("unexpected nil error")
	}

	partial := internal.Book{
		Title: "Harry Potter and the Chamber of Secrets",
	}
	if err := partial.Validate(true); err != nil {
		t.Errorf("unexpected error: want nil, got %s", err)
	}
}
