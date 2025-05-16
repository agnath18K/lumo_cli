package create

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// generatePythonProject creates a new Python web framework project
func generatePythonProject(framework string, options map[string]string) (string, error) {
	// Get project name from options or use a default
	projectName := options["name"]
	if projectName == "" {
		// Default project name based on framework
		switch strings.ToLower(framework) {
		case "fastapi":
			projectName = "fastapi_app"
		case "flask":
			projectName = "flask_app"
		default:
			projectName = "python_app"
		}
	}

	// Check if Python is installed
	if err := checkPythonInstalled(); err != nil {
		return "", err
	}

	// Create the project based on the framework
	switch strings.ToLower(framework) {
	case "fastapi":
		return setupFastAPIProject(projectName, options)
	case "flask":
		return setupFlaskProject(projectName, options)
	default:
		return "", fmt.Errorf("unsupported Python framework: %s", framework)
	}
}

// checkPythonInstalled verifies that Python is installed
func checkPythonInstalled() error {
	// Try python3 first
	cmd := exec.Command("python3", "--version")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Try python as fallback
	cmd = exec.Command("python", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Python is not installed or not in PATH. Please install Python first: https://www.python.org/downloads/")
	}

	return nil
}

// getPythonCommand returns the appropriate Python command (python3 or python)
func getPythonCommand() string {
	// Try python3 first
	cmd := exec.Command("python3", "--version")
	if err := cmd.Run(); err == nil {
		return "python3"
	}

	// Use python as fallback
	return "python"
}

// setupFastAPIProject creates a new FastAPI project
func setupFastAPIProject(projectName string, options map[string]string) (string, error) {
	// Create project directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return "", fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create virtual environment
	pythonCmd := getPythonCommand()
	cmd := exec.Command(pythonCmd, "-m", "venv", filepath.Join(projectName, "venv"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to create virtual environment: %w", err)
	}

	// Determine pip command based on OS
	var pipCmd string
	if os.PathSeparator == '/' {
		// Unix-like systems
		pipCmd = filepath.Join(projectName, "venv", "bin", "pip")
	} else {
		// Windows
		pipCmd = filepath.Join(projectName, "venv", "Scripts", "pip.exe")
	}

	// Install FastAPI and dependencies
	cmd = exec.Command(pipCmd, "install", "fastapi", "uvicorn[standard]", "pydantic")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to install FastAPI dependencies: %w", err)
	}

	// Create project structure
	dirs := []string{
		filepath.Join(projectName, "app"),
		filepath.Join(projectName, "app", "api"),
		filepath.Join(projectName, "app", "core"),
		filepath.Join(projectName, "app", "models"),
		filepath.Join(projectName, "app", "schemas"),
		filepath.Join(projectName, "tests"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create main.py
	mainPath := filepath.Join(projectName, "app", "main.py")
	mainContent := `from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(
    title="FastAPI App",
    description="FastAPI application with automatic interactive documentation",
    version="0.1.0",
)

# Configure CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
async def root():
    return {"message": "Hello World"}

@app.get("/items/{item_id}")
async def read_item(item_id: int, q: str = None):
    return {"item_id": item_id, "q": q}
`
	if err := os.WriteFile(mainPath, []byte(mainContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create main.py: %w", err)
	}

	// Create __init__.py files
	initFiles := []string{
		filepath.Join(projectName, "app", "__init__.py"),
		filepath.Join(projectName, "app", "api", "__init__.py"),
		filepath.Join(projectName, "app", "core", "__init__.py"),
		filepath.Join(projectName, "app", "models", "__init__.py"),
		filepath.Join(projectName, "app", "schemas", "__init__.py"),
	}

	for _, file := range initFiles {
		if err := os.WriteFile(file, []byte(""), 0644); err != nil {
			return "", fmt.Errorf("failed to create %s: %w", file, err)
		}
	}

	// Create config.py
	configPath := filepath.Join(projectName, "app", "core", "config.py")
	configContent := `from pydantic import BaseSettings
import os

class Settings(BaseSettings):
    API_V1_STR: str = "/api/v1"
    PROJECT_NAME: str = "FastAPI App"

    # CORS
    BACKEND_CORS_ORIGINS: list = ["*"]

    class Config:
        case_sensitive = True
        env_file = ".env"

settings = Settings()
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create config.py: %w", err)
	}

	// Create a sample model
	modelPath := filepath.Join(projectName, "app", "models", "item.py")
	modelContent := `from sqlalchemy import Column, Integer, String
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()

class Item(Base):
    __tablename__ = "items"

    id = Column(Integer, primary_key=True, index=True)
    title = Column(String, index=True)
    description = Column(String)
`
	if err := os.WriteFile(modelPath, []byte(modelContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create item.py: %w", err)
	}

	// Create a sample schema
	schemaPath := filepath.Join(projectName, "app", "schemas", "item.py")
	schemaContent := `from pydantic import BaseModel

class ItemBase(BaseModel):
    title: str
    description: str = None

class ItemCreate(ItemBase):
    pass

class Item(ItemBase):
    id: int

    class Config:
        orm_mode = True
`
	if err := os.WriteFile(schemaPath, []byte(schemaContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create item schema: %w", err)
	}

	// Create a sample API router
	routerPath := filepath.Join(projectName, "app", "api", "items.py")
	routerContent := `from fastapi import APIRouter, HTTPException
from typing import List, Optional

from app.schemas.item import Item, ItemCreate

router = APIRouter()

# Mock database
items_db = {}

@router.get("/items/", response_model=List[Item])
async def read_items(skip: int = 0, limit: int = 100):
    return list(items_db.values())[skip : skip + limit]

@router.post("/items/", response_model=Item)
async def create_item(item: ItemCreate):
    item_id = len(items_db) + 1
    db_item = Item(id=item_id, **item.dict())
    items_db[item_id] = db_item
    return db_item

@router.get("/items/{item_id}", response_model=Item)
async def read_item(item_id: int):
    if item_id not in items_db:
        raise HTTPException(status_code=404, detail="Item not found")
    return items_db[item_id]
`
	if err := os.WriteFile(routerPath, []byte(routerContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create items.py: %w", err)
	}

	// Create requirements.txt
	reqPath := filepath.Join(projectName, "requirements.txt")
	reqContent := `fastapi>=0.68.0,<0.69.0
pydantic>=1.8.0,<2.0.0
uvicorn>=0.15.0,<0.16.0
sqlalchemy>=1.4.23,<1.5.0
`
	if err := os.WriteFile(reqPath, []byte(reqContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create requirements.txt: %w", err)
	}

	// Create README.md
	readmePath := filepath.Join(projectName, "README.md")
	readmeContent := `# FastAPI Application

This is a FastAPI application with a clean architecture.

## Setup

1. Create a virtual environment:
   ` + "```" + `
   python -m venv venv
   ` + "```" + `

2. Activate the virtual environment:
   - On Windows: ` + "`venv\\Scripts\\activate`" + `
   - On Unix or MacOS: ` + "`source venv/bin/activate`" + `

3. Install dependencies:
   ` + "```" + `
   pip install -r requirements.txt
   ` + "```" + `

## Running the Application

Run the application with:

` + "```" + `
uvicorn app.main:app --reload
` + "```" + `

The API will be available at http://localhost:8000

## API Documentation

- Interactive API documentation: http://localhost:8000/docs
- Alternative API documentation: http://localhost:8000/redoc
`
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create README.md: %w", err)
	}

	return fmt.Sprintf("✅ FastAPI project '%s' created successfully!", projectName), nil
}

// setupFlaskProject creates a new Flask project
func setupFlaskProject(projectName string, options map[string]string) (string, error) {
	// Create project directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return "", fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create virtual environment
	pythonCmd := getPythonCommand()
	cmd := exec.Command(pythonCmd, "-m", "venv", filepath.Join(projectName, "venv"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to create virtual environment: %w", err)
	}

	// Determine pip command based on OS
	var pipCmd string
	if os.PathSeparator == '/' {
		// Unix-like systems
		pipCmd = filepath.Join(projectName, "venv", "bin", "pip")
	} else {
		// Windows
		pipCmd = filepath.Join(projectName, "venv", "Scripts", "pip.exe")
	}

	// Install Flask and dependencies
	cmd = exec.Command(pipCmd, "install", "flask", "flask-sqlalchemy", "flask-migrate", "python-dotenv")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to install Flask dependencies: %w", err)
	}

	// Create project structure
	dirs := []string{
		filepath.Join(projectName, "app"),
		filepath.Join(projectName, "app", "static"),
		filepath.Join(projectName, "app", "static", "css"),
		filepath.Join(projectName, "app", "static", "js"),
		filepath.Join(projectName, "app", "templates"),
		filepath.Join(projectName, "app", "models"),
		filepath.Join(projectName, "app", "routes"),
		filepath.Join(projectName, "migrations"),
		filepath.Join(projectName, "tests"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create __init__.py files
	initFiles := []string{
		filepath.Join(projectName, "app", "__init__.py"),
		filepath.Join(projectName, "app", "models", "__init__.py"),
		filepath.Join(projectName, "app", "routes", "__init__.py"),
	}

	for _, file := range initFiles {
		if err := os.WriteFile(file, []byte(""), 0644); err != nil {
			return "", fmt.Errorf("failed to create %s: %w", file, err)
		}
	}

	// Create app/__init__.py
	appInitPath := filepath.Join(projectName, "app", "__init__.py")
	appInitContent := `from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from flask_migrate import Migrate
from config import Config

db = SQLAlchemy()
migrate = Migrate()

def create_app(config_class=Config):
    app = Flask(__name__)
    app.config.from_object(config_class)

    db.init_app(app)
    migrate.init_app(app, db)

    # Register blueprints
    from app.routes import main_bp
    app.register_blueprint(main_bp)

    return app

from app import models
`
	if err := os.WriteFile(appInitPath, []byte(appInitContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create app/__init__.py: %w", err)
	}

	// Create config.py
	configPath := filepath.Join(projectName, "config.py")
	configContent := `import os
from dotenv import load_dotenv

basedir = os.path.abspath(os.path.dirname(__file__))
load_dotenv(os.path.join(basedir, '.env'))

class Config:
    SECRET_KEY = os.environ.get('SECRET_KEY') or 'you-will-never-guess'
    SQLALCHEMY_DATABASE_URI = os.environ.get('DATABASE_URL') or \
        'sqlite:///' + os.path.join(basedir, 'app.db')
    SQLALCHEMY_TRACK_MODIFICATIONS = False
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create config.py: %w", err)
	}

	// Create app/models/user.py
	userModelPath := filepath.Join(projectName, "app", "models", "user.py")
	userModelContent := `from app import db
from datetime import datetime

class User(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(64), index=True, unique=True)
    email = db.Column(db.String(120), index=True, unique=True)
    password_hash = db.Column(db.String(128))
    created_at = db.Column(db.DateTime, default=datetime.utcnow)

    def __repr__(self):
        return f'<User {self.username}>'
`
	if err := os.WriteFile(userModelPath, []byte(userModelContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create user.py: %w", err)
	}

	// Create app/routes/__init__.py with blueprint
	routesInitPath := filepath.Join(projectName, "app", "routes", "__init__.py")
	routesInitContent := `from flask import Blueprint

main_bp = Blueprint('main', __name__)

from app.routes import routes
`
	if err := os.WriteFile(routesInitPath, []byte(routesInitContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create routes/__init__.py: %w", err)
	}

	// Create app/routes/routes.py
	routesPath := filepath.Join(projectName, "app", "routes", "routes.py")
	routesContent := `from flask import render_template, jsonify
from app.routes import main_bp

@main_bp.route('/')
def index():
    return render_template('index.html', title='Home')

@main_bp.route('/api/users')
def get_users():
    return jsonify({
        'users': [
            {'id': 1, 'username': 'user1'},
            {'id': 2, 'username': 'user2'}
        ]
    })
`
	if err := os.WriteFile(routesPath, []byte(routesContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create routes.py: %w", err)
	}

	// Create app/templates/base.html
	baseTemplatePath := filepath.Join(projectName, "app", "templates", "base.html")
	baseTemplateContent := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ title }} - Flask App</title>
    <link rel="stylesheet" href="{{ url_for('static', filename='css/style.css') }}">
</head>
<body>
    <header>
        <nav>
            <ul>
                <li><a href="{{ url_for('main.index') }}">Home</a></li>
            </ul>
        </nav>
    </header>

    <main>
        {% block content %}{% endblock %}
    </main>

    <footer>
        <p>&copy; {{ now.year }} Flask App</p>
    </footer>

    <script src="{{ url_for('static', filename='js/main.js') }}"></script>
</body>
</html>
`
	if err := os.WriteFile(baseTemplatePath, []byte(baseTemplateContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create base.html: %w", err)
	}

	// Create app/templates/index.html
	indexTemplatePath := filepath.Join(projectName, "app", "templates", "index.html")
	indexTemplateContent := `{% extends "base.html" %}

{% block content %}
    <h1>Welcome to Flask App</h1>
    <p>This is a simple Flask application with a clean architecture.</p>
{% endblock %}
`
	if err := os.WriteFile(indexTemplatePath, []byte(indexTemplateContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create index.html: %w", err)
	}

	// Create app/static/css/style.css
	cssPath := filepath.Join(projectName, "app", "static", "css", "style.css")
	cssContent := `body {
    font-family: Arial, sans-serif;
    line-height: 1.6;
    margin: 0;
    padding: 0;
    color: #333;
}

header {
    background-color: #4a69bd;
    color: white;
    padding: 1rem;
}

nav ul {
    display: flex;
    list-style: none;
    padding: 0;
}

nav ul li {
    margin-right: 1rem;
}

nav ul li a {
    color: white;
    text-decoration: none;
}

main {
    padding: 2rem;
    max-width: 1200px;
    margin: 0 auto;
}

footer {
    background-color: #f1f2f6;
    text-align: center;
    padding: 1rem;
    margin-top: 2rem;
}
`
	if err := os.WriteFile(cssPath, []byte(cssContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create style.css: %w", err)
	}

	// Create app/static/js/main.js
	jsPath := filepath.Join(projectName, "app", "static", "js", "main.js")
	jsContent := `// Main JavaScript file
console.log('Flask app loaded');
`
	if err := os.WriteFile(jsPath, []byte(jsContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create main.js: %w", err)
	}

	// Create run.py
	runPath := filepath.Join(projectName, "run.py")
	runContent := `from app import create_app, db
from app.models.user import User

app = create_app()

@app.shell_context_processor
def make_shell_context():
    return {'db': db, 'User': User}

if __name__ == '__main__':
    app.run(debug=True)
`
	if err := os.WriteFile(runPath, []byte(runContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create run.py: %w", err)
	}

	// Create requirements.txt
	reqPath := filepath.Join(projectName, "requirements.txt")
	reqContent := `flask==2.0.1
flask-sqlalchemy==2.5.1
flask-migrate==3.1.0
python-dotenv==0.19.0
`
	if err := os.WriteFile(reqPath, []byte(reqContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create requirements.txt: %w", err)
	}

	// Create README.md
	readmePath := filepath.Join(projectName, "README.md")
	readmeContent := `# Flask Application

This is a Flask application with a clean architecture.

## Setup

1. Create a virtual environment:
   ` + "```" + `
   python -m venv venv
   ` + "```" + `

2. Activate the virtual environment:
   - On Windows: ` + "`venv\\Scripts\\activate`" + `
   - On Unix or MacOS: ` + "`source venv/bin/activate`" + `

3. Install dependencies:
   ` + "```" + `
   pip install -r requirements.txt
   ` + "```" + `

4. Initialize the database:
   ` + "```" + `
   flask db init
   flask db migrate -m "Initial migration"
   flask db upgrade
   ` + "```" + `

## Running the Application

Run the application with:

` + "```" + `
flask run
` + "```" + `

Or:

` + "```" + `
python run.py
` + "```" + `

The application will be available at http://localhost:5000
`
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create README.md: %w", err)
	}

	// Create .env file
	envPath := filepath.Join(projectName, ".env")
	envContent := `SECRET_KEY=dev-key-please-change-in-production
FLASK_APP=run.py
FLASK_ENV=development
`
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create .env file: %w", err)
	}

	// Create .gitignore
	gitignorePath := filepath.Join(projectName, ".gitignore")
	gitignoreContent := `# Python
__pycache__/
*.py[cod]
*$py.class
*.so
.Python
venv/
ENV/
env/
.env

# Flask
instance/
.webassets-cache
app.db

# Migrations
migrations/versions/

# IDE
.idea/
.vscode/
*.swp
*.swo
`
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create .gitignore: %w", err)
	}

	return fmt.Sprintf("✅ Flask project '%s' created successfully!", projectName), nil
}
