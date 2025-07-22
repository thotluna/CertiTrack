# Class Diagram - People Module (CertiTrack)

```mermaid
classDiagram
    %% Domain Entities
    class Person {
        +ID: string
        +FirstName: string
        +LastName: string
        +Email: string
        +Phone: string
        +Status: string
        +CreatedAt: datetime
        +UpdatedAt: datetime
    }

    %% DTOs
    class CreatePersonRequest {
        +FirstName: string
        +LastName: string
        +Email: string
        +Phone: string
        +Validate()
    }

    class UpdatePersonRequest {
        +FirstName: string
        +LastName: string
        +Email: string
        +Phone: string
        +Validate()
    }

    class PersonResponse {
        +ID: string
        +FullName: string
        +Email: string
        +Phone: string
        +Status: string
        +CreatedAt: string
    }

    %% Data Access Layer
    class PersonRepository {
        +Create(Person) error
        +FindByID(string) (Person, error)
        +FindAll(filters) ([]Person, error)
        +Update(Person) error
        +Delete(string) error
    }

    %% Business Logic Layer
    class PersonService {
        -repository: PersonRepository
        +Create(CreatePersonRequest) (PersonResponse, error)
        +GetByID(string) (PersonResponse, error)
        +List(filters) ([]PersonResponse, error)
        +Update(string, UpdatePersonRequest) (PersonResponse, error)
        +Delete(string) error
        -toResponse(p: Person) PersonResponse
        -toEntity(req: CreatePersonRequest) Person
        -mergePerson(p: Person, req: UpdatePersonRequest) Person
    }

    %% HTTP Handlers
    class PersonHandler {
        -service: PersonService
        +RegisterRoutes(router)
        +Create(ctx)
        +GetByID(ctx)
        +List(ctx)
        +Update(ctx)
        +Delete(ctx)
        -bindAndValidate(data: any) error
        -handleError(err: error) (int, any)
    }

    %% Relationships
    PersonHandler *-- PersonService : contains >
    PersonService *-- PersonRepository : contains >
    
    %% Dependencies
    PersonHandler ..> CreatePersonRequest : uses >
    PersonHandler ..> UpdatePersonRequest : uses >
    
    %% Service Operations
    PersonService ..> Person : manages >
    PersonService ..> PersonResponse : returns >
    
    %% Repository Operations
    PersonRepository "1" -- "*" Person : manages >
    
    %% Conversions
    CreatePersonRequest ..> Person : converts to >
    UpdatePersonRequest ..> Person : updates >
    Person ..> PersonResponse : transforms to >
    
    %% Service Methods
    PersonService --> Person : toResponse()
    PersonService --> CreatePersonRequest : toEntity()
    PersonService --> UpdatePersonRequest : mergePerson()
```

## Explicación del Diagrama (en español)

### Componentes Principales:

1. **Person**: Entidad principal que representa a una persona en el sistema.
2. **DTOs (Data Transfer Objects)**:
   - `CreatePersonRequest`: Para la creación de nuevas personas
   - `UpdatePersonRequest`: Para actualizaciones parciales
   - `PersonResponse`: Formato de respuesta de la API
3. **PersonRepository**: Maneja el acceso a la base de datos
4. **PersonService**: Contiene la lógica de negocio
5. **PersonHandler**: Maneja las peticiones HTTP (antes llamado Controller)

### Relaciones:

1. **Composición (◆)**:
   - `PersonHandler` contiene un `PersonService`
   - `PersonService` contiene un `PersonRepository`

2. **Dependencia (..>)**:
   - El handler usa los DTOs para validación
   - El servicio maneja las entidades y genera respuestas

3. **Asociación (--)**:
   - El repositorio gestiona múltiples instancias de `Person`

4. **Conversiones (..|>)**:
   - Los DTOs se convierten a entidades
   - Las entidades se transforman a respuestas

### Flujo de Datos:
1. Las peticiones HTTP llegan al `PersonHandler`
2. El handler valida los datos usando los DTOs
3. Delega la lógica al `PersonService`
4. El servicio convierte los DTOs a entidades
5. El repositorio persiste las entidades
6. Las entidades se transforman en respuestas
7. Las respuestas se devuelven al cliente
