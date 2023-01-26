package personRepository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LayssonENS/go-genealogy-api/domain"
	"testing"
	"time"
)

func TestCreatePerson(t *testing.T) {
	// Create a mocks for the sql.DB struct
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new postgresPersonRepo
	repo := &postgresPersonRepo{DB: db}

	// Create a test person request
	person := domain.PersonRequest{
		Name:      "John Doe",
		Email:     "johndoe@example.com",
		BirthDate: time.Now().Format("2006-01-02"),
	}

	// Define the expected query and its response
	birthdate, _ := time.Parse("2006-01-02", person.BirthDate)
	mock.ExpectExec("INSERT INTO person").WithArgs(
		"John Doe",
		"johndoe@example.com",
		birthdate).WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the CreatePerson function
	err = repo.CreatePerson(person)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure that the correct query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllPerson(t *testing.T) {
	// Create a mocks for the sql.DB struct
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new postgresPersonRepo
	repo := &postgresPersonRepo{DB: db}

	// Define the expected query and its response
	rows := sqlmock.NewRows([]string{"id", "name", "email", "birth_date", "created_at"}).
		AddRow(1, "John Doe", "johndoe@example.com", time.Now(), time.Now()).
		AddRow(2, "Jane Smith", "janesmith@example.com", time.Now(), time.Now())
	mock.ExpectQuery("SELECT id, name, email, birth_date, created_at FROM person").WillReturnRows(rows)

	// Call the GetAllPerson function
	people, err := repo.GetAllPerson()
	if err != nil {
		t.Fatal(err)
	}

	// Check the response
	if len(people) != 2 {
		t.Errorf("expected 2 people but got %d", len(people))
	}

	// Ensure that the correct query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByID(t *testing.T) {
	// Create a mocks for the sql.DB struct
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new postgresPersonRepo
	repo := &postgresPersonRepo{DB: db}

	// Define the expected query and its response
	rows := sqlmock.NewRows([]string{"id", "name", "email", "birth_date", "created_at"}).
		AddRow(1, "John Doe", "johndoe@example.com", time.Now(), time.Now())
	mock.ExpectQuery("SELECT id, name, email, birth_date, created_at FROM person WHERE id = \\$1").WithArgs(1).WillReturnRows(rows)

	// Call the GetByID function
	person, err := repo.GetByID(1)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response
	if person.ID != 1 {
		t.Errorf("expected ID 1 but got %d", person.ID)
	}

	// Ensure that the correct query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
