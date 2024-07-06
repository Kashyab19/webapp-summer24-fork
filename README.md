Before you begin, ensure you meet the following requirements, feel free to skip if you have this installed already:

- Go (version 1.16 or newer) installed on your machine. Download Go.
- Database  (PostgreSQL) installed and running on your machine. The application will bootstrap the database schema automatically upon startup using the GORM.



- Create an environment variable(will be automatically configured in future versions)
  - Configure the following environment variables in a .env file or in your environment before running the application:
  - DB_HOST 
  - DB_USER 
  - DB_PASSWORD 
  - DB_NAME (e.g., mywebappdb) - Ensure this database exists; the application will bootstrap or migrate the schema automatically.
  - JWT_SECRET_KEY - Used for signing JWT tokens.

  - Clone the repository:
    git clone https://github.com/yourusername/yourrepositoryname.git
    cd webapp 
    (Ensure that you have created the env file as it would have been ignored by git)

-   Install the Go dependencies:
    go mod tidy

-   Build the application:
    go build
-   Run the application:
    go run main.go

For testing purposes: https://app.swaggerhub.com/apis-docs/csye6225-webapp/cloud-native-webapp/2024.spring.02

Testing status

