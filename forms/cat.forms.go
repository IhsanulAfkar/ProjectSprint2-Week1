package forms

type CreateCatForms struct {
	Name        string
	Race        string
	Sex         string
	AgeInMonth  int
	Description string
	ImageUrls   []string
}