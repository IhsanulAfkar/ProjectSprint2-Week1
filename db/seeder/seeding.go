// package seeder

// import (
// 	"week1/db"
// 	"week1/helper/hash"
// )

// func main() {
// 	db.Init()
// 	conn := db.CreateConn()
// 	usersQuery := "INSERT INTO public.user (name, email, password)VALUES ('Ihsanul', 'ihsanul2001@gmail.com', ?), ('Ihsanul2', 'ihsanulafkar01@gmail.com',?)"
// 	hashPass1,_ := hash.HashPassword("12345678")
// 	hashPass2,_ := hash.HashPassword("12345678")
// 	conn.Exec("")
// }