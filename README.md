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

2. Install dependencies:

```bash
cd the-blog-api
go mod download
```

3. Create a.env file and set the required environment variables:

```bash
cp.env.example.env
```

4. Run the migrations to create the database tables:

```bash
go run main.go migrate
```

5. Start the server:

```bash
go run main.go
```

## Endpoints

### Auth

| Method | Endpoint               | Description                                                                  |
| ------ | ---------------------- | ---------------------------------------------------------------------------- |
| POST   | /signup                | Sign up a new user                                                           |
| POST   | /login                 | Login an existing user                                                      |

### Posts

| Method | Endpoint               | Description                                                                  |
| ------ | ---------------------- | ---------------------------------------------------------------------------- |
| POST   | /posts                 | Create a new post                                                           |
| GET    | /posts                 | Get all the posts                                                           |
| GET    | /posts/:postId         | Get a specific post                                                         |
| DELETE | /posts/:postId         | Delete a specific post                                                      |
| PUT    | /posts/:postId         | Update a specific post                                                      |
| GET    | /posts/:postId/comments | Get all the comments for a specific post                                     |
| POST   | /posts/:postId/comments | Create a new comment for a specific post                                     |
| DELETE | /posts/:postId/comments/:commentId | Delete a specific comment for a specific post                              |
| PUT    | /posts/:postId/comments/:commentId | Update a specific comment for a specific post                              |
| GET    | /posts/:postId/likes    | Get all the likes for a specific post                                        |
| POST   | /posts/:postId/likes    | Like a specific post                                                        |
| DELETE | /posts/:postId/likes    | Unlike a specific post                                                      |

### Users

| Method | Endpoint               | Description                                                                  |
| ------ | ---------------------- | ---------------------------------------------------------------------------- |
| GET    | /users                 | Get all the users                                                           |
| GET    | /users/:userId         | Get a specific user                                                         |
| DELETE | /users/:userId         | Delete a specific user                                                      |
| PUT    | /users/:userId         | Update a specific user                                                      |

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.