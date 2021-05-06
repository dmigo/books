# Books

## Examples of requests

1. Check if the service is up and running

    ```zsh
    curl --location --request GET 'http://localhost:8080/health'
    ```

1. Add a book

    ```zsh
    curl --location --request POST 'http://localhost:8080/books' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "title": "The Catcher in the Rye",
    "author": "Jerome David Salinger",
    "publish_year": 1951 
    }'
    ```

1. Get all books

    ```zsh
    curl --location --request GET 'http://localhost:8080/books'
    ```

1. Get one book (Don't forget to replace {BOOK_ID} with a real id of a book)

    ```zsh
    curl --location --request GET 'http://localhost:8080/books/{BOOK_ID}'
    ```
