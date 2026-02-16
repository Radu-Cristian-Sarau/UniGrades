# UniGrades

**A terminal user interface (TUI) application for managing and monitoring university course grades, ECTS credits, and academic progress.**

Created because I deeply dislike having to use two-factor authentication (2FA) anytime I want to view this data on my university's internal websites.
## Features

- **Interactive Dashboard** – View your courses in a table with real-time statistics
- **Grade Analytics** – Track average grades, grade distribution per year, and ECTS progress
- **Course Management** – Add, edit, and delete course records with simple commands

## Prerequisites

- **Go 1.25.6** or later
- **MongoDB** (this project uses a cloud instance via MongoDB Atlas)
- **Environment Configuration** – MongoDB connection URI

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/Radu-Cristian-Sarau/UniGrades.git
cd UniGrades
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Set Up MongoDB

Create a `.env` file in the project root with your MongoDB connection URI:

```env
MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority
```

### 4. Build and Run

```bash
go run main.go
```

## Usage

### Starting the Application

When you launch UniGrades, you'll see a university picker screen. Select a university to proceed to the grades view.

### Commands

Once in the grades view, interact with your courses using these terminal commands:

| Command | Description | Example |
|---------|-------------|---------|
| `/add` | Add a new course | `/add Applied_Math 1 8 5` |
| `/edit` | Modify course information | `/edit Applied_Math Grade 9` |
| `/delete` | Remove a course | `/delete Applied_Math` |

### Navigation

- **Arrow Keys or J and K Keys** – Navigate the list of universities
- **Enter** – Select an option or confirm input
- **Ctrl + Q** – Go back to the university picker screen
- **Ctrl + C** – Quit the application

## Project Structure

```
UniGrades/
├── main.go                               # Application entry point
├── go.mod / go.sum                       # Dependency management
├── internal/
│   ├── api/                              # MongoDB API layer
│   │   ├── api.go                        # Data fetching operations
│   │   └── server.go
│   ├── computations/                     # Business logic calculations
│   │   ├── averages.go                   # Grade average calculations
│   │   ├── total_ects.go                 # ECTS credit computation
│   ├── screens/                          # UI Screen definitions
│   │   ├── grades/                       # Grades dashboard screen
│   │   │   ├── data_screen.go
│   │   │   └── data_screen_style.go
│   │   └── picker/                       # University selection screen
│   │       └── model.go
│   ├── tui/                              # Terminal UI components
│   │   ├── title.go                      # App title rendering
│   │   ├── table_renderer.go             # Course table display
│   │   ├── average_grades_renderer.go    # Grade statistics
│   │   ├── average_grades_per_year_renderer.go # Grade stats per year
│   │   ├── total_ects_renderer.go        # ECTS statistics
│   │   ├── total_ects_per_year_renderer.go # ECTS stats per year
│   │   └── *_style.go                    # Styling and colors
│   └── university/                       # University data models
│       └── university.go
└── README.md
```

## Architecture

UniGrades follows the architecture pattern below:

- **Models** – Located in `internal/screens/` (Bubble Tea models for each screen)
- **API Layer** – `internal/api/` handles MongoDB operations
- **Computations** – `internal/computations/` contains pure business logic
- **UI Layer** – `internal/tui/` manages all rendering and styling
- **Domain Models** – `internal/university/` contains data structures

### Technology Stack

- **Language** – Go 1.25.6
- **TUI Framework** – [Charmbracelet Bubbles](https://github.com/charmbracelet/bubbles) & [Bubbletea](https://github.com/charmbracelet/bubbletea)
- **Styling** – [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Database** – [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
- **Charts** – [NimbleMarkets ntcharts](https://github.com/NimbleMarkets/ntcharts)
- **Environment** – [godotenv](https://github.com/joho/godotenv)

## Development

### Running Locally

```bash
go run main.go
```

### Building an Executable

```bash
go build -o unigrades
./unigrades
```

---

**Created by:** Radu-Cristian Sarău  
**Repository:** [UniGrades on GitHub](https://github.com/Radu-Cristian-Sarau/UniGrades)
