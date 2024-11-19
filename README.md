# Simple Chat

Simple Chat is a lightweight real-time chat application that allows multiple users to communicate instantly. It uses WebSockets for efficient message delivery and provides a responsive, user-friendly interface.

## Features

- **Real-time Messaging**: Instant message delivery using WebSockets.
- **Go Backend**: The server is built with Go for performance and simplicity.
- **Dynamic Templates**: HTML templates with Go's `template` package for dynamic content rendering.
- **TailwindCSS**: Responsive and modern UI design using TailwindCSS CDN.

## Project Structure

- **Server**: A Go-based WebSocket server that manages connections and broadcasts messages.
- **Frontend**: Built with HTML templates, styled using TailwindCSS, and rendered dynamically.
- **Templates**: Uses Go's `template` package with `{{ template "title" . }}` syntax for embedding dynamic data.

## Installation

1. **Clone the Repository**:
    ```bash
    git clone https://github.com/MementoMori11723/Simple-Chat.git
    cd simple-chat
    ```

2. **Run the Application**:
    Ensure you have `make` and `go` installed. Then, run the following command:
    ```bash
    make run
    ```
    Or, you can use the following commands:
    ```bash
    make
    ```
    Or, you can also run the application manually:
    ```bash
    go run .
    ```

3. **Access the Application**:
    Open your browser and go to `http://localhost:8080`.

## Usage

- Navigate to the homepage to join the chat room.
- Start chatting with other connected users.
- Error handling is implemented for invalid routes.

## File Structure

- **`app.go`**: Main server file.
- **`pages/`**: Contains HTML templates (`index.html`, `layouts.html`, `error.html`).
- **`Makefile`**: Simplifies running the application.

## Future Improvements

- User authentication.
- Private messaging.
- Chat history storage.

---

Feel free to customize further if needed!
