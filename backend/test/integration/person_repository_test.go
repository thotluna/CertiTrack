package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"certitrack/backend/feature/person"
	"certitrack/backend/shared/database"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

type PersonRepositoryTestSuite struct {
	suite.Suite
	container testcontainers.Container
	db        *gorm.DB
	repo      person.PersonRepository
}

func (s *PersonRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(30*time.Second),
			wait.ForListeningPort("5432/tcp").WithStartupTimeout(10*time.Second),
		),
	}

	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            false,
	})
	if err != nil {
		s.T().Fatal("Failed to start container:", err)
	}
	s.container = pgContainer

	port, err := pgContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		s.T().Fatal("Failed to get mapped port:", err)
	}
	host, err := pgContainer.Host(ctx)
	if err != nil {
		s.T().Fatal("Failed to get container host:", err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=testuser password=testpass dbname=testdb sslmode=disable",
		host, port.Int())

	var db *gorm.DB
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = database.ConnectForTesting(connStr)
		if err == nil {
			s.db = db
			break
		}
		if i == maxRetries-1 {
			s.T().Fatal("Failed to connect to test database after", maxRetries, "attempts:", err)
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	s.repo = person.NewGormRepository(s.db)
}

func (s *PersonRepositoryTestSuite) TearDownSuite() {
	sqlDB, err := s.db.DB()
	if err == nil {
		sqlDB.Close()
	}

	if s.container != nil {
		s.container.Terminate(context.Background())
	}
}

func (s *PersonRepositoryTestSuite) SetupTest() {
	s.db.Exec("TRUNCATE TABLE persons CASCADE")
}

func (s *PersonRepositoryTestSuite) TestCreateAndFindPerson() {
	newPerson := &person.Person{
		ID:        uuid.New(),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "+1234567890",
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.Create(context.Background(), newPerson)
	s.NoError(err)

	found, err := s.repo.FindByID(context.Background(), newPerson.ID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(newPerson.ID, found.ID)
	s.Equal(newPerson.Email, found.Email)
	s.Equal(newPerson.FirstName, found.FirstName)
	s.Equal(newPerson.LastName, found.LastName)
	s.Equal(newPerson.Phone, found.Phone)
	s.Equal(newPerson.Status, found.Status)

	persons, err := s.repo.FindAll(context.Background(), map[string]interface{}{"limit": 10, "offset": 0})
	s.NoError(err)
	s.Len(persons, 1)
	s.Equal(newPerson.ID, persons[0].ID)
	count, err := s.repo.Count(context.Background(), map[string]interface{}{})
	s.NoError(err)
	s.Equal(int64(1), count)
	persons, err = s.repo.FindAll(context.Background(), map[string]interface{}{"name": "John"})
	s.NoError(err)
	s.Len(persons, 1)
	s.Equal(newPerson.ID, persons[0].ID)
}

func (s *PersonRepositoryTestSuite) TestFindByID_NotFound() {
	_, err := s.repo.FindByID(context.Background(), uuid.New())
	s.Error(err)
	s.Equal(person.ErrPersonNotFound, err)
}

func (s *PersonRepositoryTestSuite) TestCreate_DuplicateEmail() {
	p1 := &person.Person{
		ID:        uuid.New(),
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Status:    "active",
	}
	err := s.repo.Create(context.Background(), p1)
	s.NoError(err)

	p2 := &person.Person{
		ID:        uuid.New(),
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane.doe@example.com",
		Status:    "active",
	}
	err = s.repo.Create(context.Background(), p2)
	s.Error(err)
	s.Equal(person.ErrEmailAlreadyExists, err)
}

func (s *PersonRepositoryTestSuite) TestUpdatePerson() {
	p := &person.Person{
		ID:        uuid.New(),
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Status:    "active",
	}
	s.NoError(s.repo.Create(context.Background(), p))

	p.FirstName = "Jane Updated"
	p.Phone = "+9876543210"
	err := s.repo.Update(context.Background(), p)
	s.NoError(err)
	updated, err := s.repo.FindByID(context.Background(), p.ID)
	s.NoError(err)
	s.Equal("Jane Updated", updated.FirstName)
	s.Equal("+9876543210", updated.Phone)
	nonExistent := &person.Person{
		ID:        uuid.New(),
		FirstName: "Non",
		LastName:  "Existent",
		Email:     "non.existent@example.com",
	}
	err = s.repo.Update(context.Background(), nonExistent)
	s.Error(err)
	s.Equal(person.ErrPersonNotFound, err)
}

func (s *PersonRepositoryTestSuite) TestUpdate_DuplicateEmail() {
	p1 := &person.Person{
		ID:        uuid.New(),
		FirstName: "Alice",
		LastName:  "Smith",
		Email:     "alice.smith@example.com",
		Status:    "active",
	}
	s.NoError(s.repo.Create(context.Background(), p1))

	p2 := &person.Person{
		ID:        uuid.New(),
		FirstName: "Bob",
		LastName:  "Johnson",
		Email:     "bob.johnson@example.com",
		Status:    "active",
	}
	s.NoError(s.repo.Create(context.Background(), p2))

	p2.Email = p1.Email
	err := s.repo.Update(context.Background(), p2)
	s.Error(err)
	s.Equal(person.ErrEmailAlreadyExists, err)
}

func (s *PersonRepositoryTestSuite) TestDeletePerson() {
	p := &person.Person{
		ID:        uuid.New(),
		FirstName: "To Delete",
		LastName:  "Person",
		Email:     "delete.me@example.com",
		Status:    "active",
	}
	s.NoError(s.repo.Create(context.Background(), p))

	err := s.repo.Delete(context.Background(), p.ID)
	s.NoError(err)

	_, err = s.repo.FindByID(context.Background(), p.ID)
	s.Equal(person.ErrPersonNotFound, err)

	err = s.repo.Delete(context.Background(), uuid.New())
	s.Error(err)
	s.Equal(person.ErrPersonNotFound, err)
}

func (s *PersonRepositoryTestSuite) TestFindAllWithPagination() {
	for i := 0; i < 15; i++ {
		p := &person.Person{
			ID:        uuid.New(),
			FirstName: "User",
			LastName:  "Test",
			Email:     fmt.Sprintf("user.test%d@example.com", i),
			Status:    "active",
		}
		s.NoError(s.repo.Create(context.Background(), p))
	}

	persons, err := s.repo.FindAll(context.Background(), map[string]interface{}{
		"limit":  5,
		"offset": 0,
	})
	s.NoError(err)
	s.Len(persons, 5)

	persons, err = s.repo.FindAll(context.Background(), map[string]interface{}{
		"limit":  5,
		"offset": 5,
	})
	s.NoError(err)
	s.Len(persons, 5)

	persons, err = s.repo.FindAll(context.Background(), map[string]interface{}{
		"status": "active",
	})
	s.NoError(err)
	s.True(len(persons) >= 15)

	count, err := s.repo.Count(context.Background(), map[string]interface{}{
		"status": "active",
	})
	s.NoError(err)
	s.True(count >= 15)
}

func TestPersonRepositorySuite(t *testing.T) {
	suite.Run(t, new(PersonRepositoryTestSuite))
}
