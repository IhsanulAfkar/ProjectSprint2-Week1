package helper

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(email)
}

func Includes(target string, array []string)bool{
	for _, value := range array {
        if value == target {
            return true
        }
    }
    return false
}
func IsURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}
func FormatToIso860(s string)string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return ""
	}

	// Format the time object into ISO 8601 format
	return t.Format("2006-01-02T15:04:05Z07:00")
}

func ConvertToPGList(list []string)string{
	var escaped []string
	for _, s := range list {
		escaped = append(escaped,fmt.Sprintf("\"%s\"", strings.ReplaceAll(s, "\"", "\\\"")))
	}
	return fmt.Sprintf("{%s}", strings.Join(escaped, ","))
}

func IsUUID(s string) bool{
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
    return uuidRegex.MatchString(s)
}