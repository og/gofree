package f_test

//
// func TestBasic(t *testing.T) {
// 	db := f.NewDatabase(f.DataSourceName{
// 		DriverName: "mysql",
// 		User: "root",
// 		Password: "somepass",
// 		Host: "localhost",
// 		Port: "3306",
// 		DB: "test_gofree",
// 	})
// 	type IDAndName struct {
// 		ID string
// 		Name string
// 	}
// 	data := struct {
// 		User struct {
// 			A IDAndName
// 			B IDAndName
// 		}
// 		Book struct {
// 			A1 IDAndName
// 			A2 IDAndName
// 			B1 IDAndName
// 			B2 IDAndName
// 		}
// 	}{}
// 	mock, err := f.ResetMockData(db, nil, "./relation.mock.json") ; ge.Check(err)
// 	{
// 		data.User.A.ID = mock.Local.String("userAID")
// 		data.User.A.Name = mock.Local.String("userAName")
// 		data.User.B.ID = mock.Local.String("userBID")
// 		data.User.B.Name = mock.Local.String("userBName")
//
// 		data.Book.A1.ID = mock.Local.String("bookA1ID")
// 		data.Book.A1.Name = mock.Local.String("bookA1Name")
// 		data.Book.A2.ID = mock.Local.String("bookA2ID")
// 		data.Book.A2.Name = mock.Local.String("bookA2Name")
// 		data.Book.B1.ID = mock.Local.String("bookB1ID")
// 		data.Book.B1.Name = mock.Local.String("bookB1Name")
// 		data.Book.B2.ID = mock.Local.String("bookB2ID")
// 		data.Book.B2.Name = mock.Local.String("bookB2Name")
// 	}
// 	userList := []User{} ; _=userList
// 	db.ListQB(&userList, f.QB{})
// 	l.V(userList)
// }
