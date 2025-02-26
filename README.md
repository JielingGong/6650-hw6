# CS6650HW6_GO

This system stores album information, including an image, in a **MySQL database (AWS RDS MySQL)**. 

## Data Model

The database consists of a single `albums` table with the following schema:

| Column  | Type                     | Description                          |
|---------|--------------------------|--------------------------------------|
| `id`    | INT (PK, AUTO_INCREMENT) | Unique identifier for each album    |
| `title` | VARCHAR(255)             | The title of the album              |
| `year`  | INT                      | The release year of the album       |
| `artist`| VARCHAR(255)             | The artist associated with the album |
| `image` | LONGBLOB                 | The album cover image (binary data) |

Images are stored in the database using the **LONGBLOB** data type to accommodate binary image data. The test image used is **25 KB** in size.

## API Endpoints

The system provides the following API endpoints:

- **POST `/album`** – Uploads a new album with metadata and an image.
- **GET `/album/{id}`** – Retrieves album details, including the stored image.

## Technology Stack

- **Backend:** Go
- **Router:** Gorilla Mux
- **Database:** MySQL (AWS RDS)
- **Hosting:** AWS EC2

This setup allows efficient album storage and retrieval with a simple yet scalable design.
