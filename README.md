# The Blog API

The Blog API is a robust and efficient solution for managing a blog platform. Developed using Golang and Gin, it embodies modern practices and technologies for seamless performance and scalability. Leveraging PostgreSQL as its database engine ensures reliability and flexibility in data management, while JWT (JSON Web Tokens) authentication enhances security by providing a stateless authentication mechanism.

This API serves as a foundational component for building and managing a dynamic blogging platform, offering a comprehensive set of endpoints for user authentication, post management, comment handling, user profile management, and more. With clear and concise documentation and a user-friendly architecture, integrating this API into your project is straightforward and hassle-free.

Whether you're developing a personal blog, a collaborative writing platform, or an enterprise-level content management system, The Blog API provides the necessary tools and functionality to streamline your development process and deliver a seamless user experience. Unlock the power of modern web development with The Blog API.

## Requirements

- Go 1.16 or higher
- PostgreSQL 12 or higher

## Getting Started

1. Clone this repository:

```bash
git clone https://github.com/khunaungpaing/the-blog-api.git
```

2. Navigate to the project directory:

```bash
cd the-blog-api
```

3. Ensure that PostgreSQL is running:

```bash
make postgres
```

4. Initialize Swagger documentation:
```bash
make doc
```

5. Compile and run the application using the Makefile:

```bash
make run
```
6. Visit the API documentation at http://localhost:8080/swagger/index.html.

Note: You may need to adjust the PORT configuration in your `.env` file if your server runs on a different port.

## Makefile Commands
1. `make run`: Compile and run the application.
2. `make postgres`: Start PostgreSQL database using Docker Compose.
3. `make doc`: Initialize Swagger documentation.

## Endpoints

### Auth

| Method | Endpoint               | Description                                                                  |
| ------ | ---------------------- | ---------------------------------------------------------------------------- |
| POST   | /users/signup          | Sign up a new user                                                           |
| POST   | /users/login           | Login an existing user                                                      |

### Posts

| Method | Endpoint               | Description                                                                  |
| ------ | ---------------------- | ---------------------------------------------------------------------------- |
| POST   | /posts                 | Create a new post                                                           |
| GET    | /posts                 | Get all the posts                                                           |
| GET    | /posts/:postId         | Get a specific post                                                         |
| DELETE | /posts/:postId         | Delete a specific post                                                      |
| PATCH  | /posts/:postId         | Update a specific post                                                      |
| GET    | /posts/:postId/comments | Get all the comments for a specific post                                     |
| POST   | /posts/:postId/comments | Create a new comment for a specific post                                     |
| DELETE | /posts/:postId/comments/:commentId | Delete a specific comment for a specific post                              |
| PATCH  | /posts/:postId/comments/:commentId | Update a specific comment for a specific post                              |
| GET    | /posts/:postId/likes    | Get all the likes for a specific post                                        |
| POST   | /posts/:postId/likes    | Like a specific post                                                        |
| DELETE | /posts/:postId/likes    | Unlike a specific post                                                      |

### Users

| Method | Endpoint               | Description                                                                  |
| ------ | ---------------------- | ---------------------------------------------------------------------------- |
| GET    | /users/profile         | Get the profile of the logged-in user                                       |
| PATCH  | /users/profile         | Update the profile of the logged-in user                                    |

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.