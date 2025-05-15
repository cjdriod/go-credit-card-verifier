# GoLang Dummy Credit Card Verification REST Endpoint System

This is a simple GoLang application that provides a RESTful API for credit card verification.
It includes endpoints to perform various operations related to credit card verification,
such as validating a credit card number, checking the card type, checking card history,
and identify black listed card.

## **Features**

- **Credit Card Validation**: Verify the validity of a credit card number using the Luhn algorithm.
- **Credit Card Type Detection**: Determine the type of credit card (e.g., Visa, MasterCard, American Express).
- **Black List Bank Card**: Record card into black list bank card database.
- **Report Bank Card Activity**: Report bank card activity to system and share with public on check.

![Language](https://img.shields.io/github/languages/top/cjdriod/go-credit-card-verifier?style=flat-square)
![Size](https://img.shields.io/github/repo-size/cjdriod/go-credit-card-verifier?style=flat-square)

## 🔨Installation

1. **Clone the Repository**:

    ```bash
    git clone https://github.com/cjdriod/go-credit-card-verifier.git
    ```

2. **Navigate to the Project Directory**:

    ```bash
    cd go-credit-card-verifier
    ```

3. **Build the Application**:

    ```bash
    go build -o main
    ```

## **⛷️ Run application**

### With Docker

 ```bash
 // Run
docker-compose up -d --build
```
 ```bash
// Shutdown
docker-compose down -v --rmi all --remove-orphans   
```

### Without Docker

```bash
go run .\cmd
```

### With binary file

```bash
./main
```

## **⚙️ Config Environment Variables**

| Variable                  | Description                                 | Default Value |
|---------------------------|---------------------------------------------|---------------|
| MYSQL_ACC                 | MySQL account username                      |               |
| MYSQL_PASSWORD            | MySQL account password                      |               |
| MYSQL_HOST                | MySQL host                                  |               |
| MYSQL_PORT                | MySQL port                                  | 3306          |
| MYSQL_DB_NAME             | MySQL database name                         |               |
| ENABLE_PREMIUM_CARD_CHECK | Enable premium card check                   | false         |
| JWT_SECRET                | JWT secret key                              |               |
| HTTPS_MODE                | Enable HTTPS mode                           | true          |
| APP_SERVER_DOMAIN         | Domain for the application server           |               |
| APP_SERVER_PORT           | Port for the application server             | 8080          |
| GIN_MODE                  | GIN mode for the application                | debug         |
| ENV                       | Environment mode (Production / Development) |               |   

## ⚔️ **Contributing**

Contributions are welcome! If you'd like to contribute to this project, please feel free to open a pull request or
submit an issue with your suggestions or changes.

## 📝 **License**

This project is licensed under the MIT License - see the LICENSE file for details.